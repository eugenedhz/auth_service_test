package tokens

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/eugenedhz/auth_service_test/internal/service/auth"
)

type TokenManager struct {
	signingKey     string
	accessTokenTTL time.Duration
}

func NewTokenManager(signingKey string, accessTokenTTL time.Duration) *TokenManager {
	return &TokenManager{
		signingKey:     signingKey,
		accessTokenTTL: accessTokenTTL,
	}
}

func (t *TokenManager) GenerateAccessToken(userID string, userIpAddress string, tokenID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["jti"] = tokenID
	claims["sub"] = userID
	claims["exp"] = t.accessTokenTTL
	claims["userIpAddress"] = userIpAddress

	tokenString, err := token.SignedString([]byte(t.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenManager) GenerateRefreshToken(tokenID string) ([]byte, error) {
	var token []byte

	sha1 := sha1.New()
	io.WriteString(sha1, t.signingKey)

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		fmt.Println(err.Error())

		return token, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return token, err
	}

	token = gcm.Seal(nonce, nonce, []byte(tokenID), nil)
	return token, nil
}

func (t *TokenManager) ParseAccessToken(tokenString string) (*auth.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.signingKey), nil
	})

	claims := &auth.JWTClaims{}
	if err != nil {
		return claims, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		claims.TokenID = payload["jti"].(string)
		claims.UserID = payload["sub"].(string)
		claims.UserIpAddress = payload["userIpAddress"].(string)

		return claims, nil
	}

	return claims, errors.New("Invalid access token")
}

func (t *TokenManager) ParseRefreshToken(token []byte) (string, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, t.signingKey)

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(token) < (nonceSize + 1) {
		return "", errors.New("Invalid token byte length")
	}
	nonce, ciphertext := token[:nonceSize], token[nonceSize:]

	tokenUUID, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(tokenUUID), nil
}
