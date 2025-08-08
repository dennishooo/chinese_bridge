package handler

import (
	"chinese-bridge-game/internal/auth/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.GET("/google", h.GoogleLogin)      // GET for initial OAuth redirect
		auth.POST("/google", h.GoogleLogin)     // POST for token exchange
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
	}
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Implementation will be added in later tasks
	c.JSON(200, gin.H{"message": "Google login endpoint"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Implementation will be added in later tasks
	c.JSON(200, gin.H{"message": "Refresh token endpoint"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Implementation will be added in later tasks
	c.JSON(200, gin.H{"message": "Logout endpoint"})
}

func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"service": "auth-service",
	})
}

func (h *AuthHandler) ReadyCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ready",
		"service": "auth-service",
	})
}