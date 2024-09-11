package http

import (
	"github.com/gin-gonic/gin"

	"github.com/eugenedhz/auth_service_test/internal/service/auth"
)

func RegisterHTTPEndpoints(router *gin.Engine, authService *auth.AuthService) {
	h := NewAuthHandler(authService)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/signin", h.SignIn)
		authEndpoints.POST("/refresh", h.Refresh)
	}
}
