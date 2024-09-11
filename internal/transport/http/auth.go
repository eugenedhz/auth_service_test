package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eugenedhz/auth_service_test/internal/domain/models"
	"github.com/eugenedhz/auth_service_test/internal/service"
	"github.com/eugenedhz/auth_service_test/internal/service/auth"
)

type AuthService interface {
	GetTokens(ctx context.Context, userID string, userIpAddress string) (*models.Tokens, error)
	RefreshTokens(ctx context.Context, claims *auth.JWTClaims, userIpAddress string) (*models.Tokens, error)
	ParseTokens(ctx context.Context, tokens *models.Tokens) (*auth.JWTClaims, error)
}

type AuthHandler struct {
	auth AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		auth: authService,
	}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	userID := c.DefaultQuery("userID", "")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: userIdNotProvided})
		return
	}

	userIpAddress := c.ClientIP()
	tokens, err := h.auth.GetTokens(c.Request.Context(), userID, userIpAddress)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
			return
		}

		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: internalErrorMessage})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	tokens := &models.Tokens{}

	if err := c.BindJSON(tokens); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: tokensNotProvided})
		return
	}

	claims, err := h.auth.ParseTokens(c.Request.Context(), tokens)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAccessToken) {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
			return
		}

		if errors.Is(err, service.ErrInvalidRefreshToken) {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
			return
		}

		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: internalErrorMessage})
		return
	}

	userIpAddress := c.ClientIP()
	refreshedTokens, err := h.auth.RefreshTokens(c.Request.Context(), claims, userIpAddress)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
			return
		}

		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: internalErrorMessage})
		return
	}

	c.JSON(http.StatusOK, refreshedTokens)
}
