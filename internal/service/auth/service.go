package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"github.com/beevik/guid"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/eugenedhz/auth_service_test/internal/domain/models"
	"github.com/eugenedhz/auth_service_test/internal/repository"
	"github.com/eugenedhz/auth_service_test/internal/service"
)

type AuthService struct {
	sessionRepo  SessionRepository
	userRepo     UserRepository
	tokenManager TokenManager
	email        EmailSender
	log          *slog.Logger
}

type SessionRepository interface {
	Save(ctx context.Context, session *models.Session) error
	Get(ctx context.Context, userGUID string) (*models.Session, error)
}

type UserRepository interface {
	Get(ctx context.Context, userGUID string) (*models.User, error)
}

type TokenManager interface {
	GenerateAccessToken(userID string, userIpAddress string, tokenID string) (string, error)
	GenerateRefreshToken(tokenID string) ([]byte, error)
	ParseAccessToken(token string) (*JWTClaims, error)
	ParseRefreshToken(token []byte) (string, error)
}

type JWTClaims struct {
	TokenID       string
	UserID        string
	UserIpAddress string
}

type EmailSender interface {
	SendEmail(subject string, body string, emailTo string) error
}

func NewAuthService(
	sessionRepository SessionRepository,
	userRepository UserRepository,
	tokenManager TokenManager,
	emailSender EmailSender,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		sessionRepo:  sessionRepository,
		userRepo:     userRepository,
		tokenManager: tokenManager,
		email:        emailSender,
		log:          logger,
	}
}

func (s *AuthService) GetTokens(ctx context.Context, userID string, userIpAddress string) (*models.Tokens, error) {
	const op = "AuthService.GetTokens"

	log := s.log.With(
		slog.String("op", op),
		slog.String("userID", userID),
	)

	log.Info("Attempting to get tokens")

	if !(guid.IsGuid(userID)) {
		log.Warn("Invalid userID")

		return &models.Tokens{}, service.ErrInvalidUserID
	}

	_, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Warn("User not found")

			return &models.Tokens{}, service.ErrUserNotFound
		}

		log.Error("Unable to get user", "error", err.Error())
		return &models.Tokens{}, err
	}

	tokenID := uuid.New().String()
	accessToken, err := s.tokenManager.GenerateAccessToken(userID, userIpAddress, tokenID)
	if err != nil {
		log.Error("Unable to generate access token", "error", err.Error())
		return &models.Tokens{}, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(tokenID)
	if err != nil {
		log.Error("Unable to generate refresh token", "error", err.Error())
		return &models.Tokens{}, err
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword(refreshToken, 0)
	if err != nil {
		log.Error("Unable to generate refresh token hash", "error", err.Error())
		return &models.Tokens{}, err
	}

	session := &models.Session{
		UserID:           userID,
		RefreshTokenHash: refreshTokenHash,
	}
	err = s.sessionRepo.Save(ctx, session)
	if err != nil {
		log.Error("Unable to save session", "error", err.Error())
		return &models.Tokens{}, err
	}

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: base64.URLEncoding.EncodeToString(refreshToken),
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, accessTokenPayload *JWTClaims, userIpAddress string) (*models.Tokens, error) {
	const op = "AuthService.RefreshTokens"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("Attempting to refresh tokens")

	if userIpAddress != accessTokenPayload.UserIpAddress {
		log.Info("Different user ip address from tokens")

		user, err := s.userRepo.Get(ctx, accessTokenPayload.UserID)

		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				log.Warn("User not found")

				return &models.Tokens{}, service.ErrUserNotFound
			}

			log.Error("Unable to get user", "error", err.Error())
			return &models.Tokens{}, err
		}

		subject := "Different IP address auth refresh"
		body := fmt.Sprintf("Your tokens have been refreshed by different ip address %s", userIpAddress)
		err = s.email.SendEmail(subject, body, user.Email)
		if err != nil {
			log.Error("Unable to send email", "error", err.Error())

			return &models.Tokens{}, err
		}
	}

	return s.GetTokens(ctx, accessTokenPayload.UserID, userIpAddress)
}

func (s *AuthService) ParseTokens(ctx context.Context, tokens *models.Tokens) (*JWTClaims, error) {
	const op = "AuthService.ParseTokens"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("Attempting to parse tokens")

	refreshToken, err := base64.URLEncoding.DecodeString(tokens.RefreshToken)
	if err != nil {
		log.Warn("Unable to decode refresh token")

		return &JWTClaims{}, err
	}

	tokenID, err := s.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		log.Warn("Unable to parse refresh token")

		return &JWTClaims{}, service.ErrInvalidRefreshToken
	}

	accessTokenPayload, err := s.tokenManager.ParseAccessToken(tokens.AccessToken)
	if err != nil {
		log.Warn("Unable to parse access token")

		return &JWTClaims{}, service.ErrInvalidAccessToken
	}

	if tokenID != accessTokenPayload.TokenID {
		return &JWTClaims{}, service.ErrInvalidRefreshToken
	}

	session, err := s.sessionRepo.Get(ctx, accessTokenPayload.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			log.Warn("Session not found")

			return &JWTClaims{}, service.ErrSessionNotFound
		}

		log.Error("Unable to get session", "error", err.Error())
		return &JWTClaims{}, err
	}

	err = bcrypt.CompareHashAndPassword(session.RefreshTokenHash, refreshToken)
	if err != nil {
		return &JWTClaims{}, service.ErrInvalidRefreshToken
	}

	return accessTokenPayload, nil
}
