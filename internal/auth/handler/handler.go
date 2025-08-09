package handler

import (
	"net/http"

	"chinese-bridge-game/internal/auth/dto"
	"chinese-bridge-game/internal/auth/service"
	"chinese-bridge-game/pkg/middleware"

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
	
	// Apply rate limiting to auth endpoints
	auth.Use(middleware.IPRateLimiter(5, 10)) // 5 requests per second, burst of 10
	
	{
		auth.GET("/google/url", h.GetGoogleOAuthURL)
		auth.POST("/google", h.GoogleOAuthCallback)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", middleware.JWTAuth(h.authService), h.Logout)
	}
}

// GetGoogleOAuthURL godoc
// @Summary Get Google OAuth URL
// @Description Get the Google OAuth authorization URL for login
// @Tags authentication
// @Accept json
// @Produce json
// @Param state query string false "OAuth state parameter"
// @Success 200 {object} map[string]string
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/google/url [get]
func (h *AuthHandler) GetGoogleOAuthURL(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		state = "default_state" // In production, generate a random state
	}

	url := h.authService.GetGoogleOAuthURL(state)
	
	c.JSON(http.StatusOK, gin.H{
		"url": url,
		"state": state,
	})
}

// GoogleOAuthCallback godoc
// @Summary Google OAuth callback
// @Description Handle Google OAuth callback and authenticate user
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body dto.GoogleOAuthRequest true "OAuth callback data"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/google [post]
func (h *AuthHandler) GoogleOAuthCallback(c *gin.Context) {
	var req dto.GoogleOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
			Details: err.Error(),
			TraceID: c.GetString("trace_id"),
		})
		return
	}

	authResponse, err := h.authService.GoogleOAuthLogin(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    "AUTHENTICATION_ERROR",
			Message: "Failed to authenticate with Google",
			Details: err.Error(),
			TraceID: c.GetString("trace_id"),
		})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh an expired JWT token using refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
			Details: err.Error(),
			TraceID: c.GetString("trace_id"),
		})
		return
	}

	tokenResponse, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    "AUTHENTICATION_ERROR",
			Message: "Failed to refresh token",
			Details: err.Error(),
			TraceID: c.GetString("trace_id"),
		})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate all sessions
// @Tags authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.MessageResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    "AUTHENTICATION_ERROR",
			Message: "User not authenticated",
			TraceID: c.GetString("trace_id"),
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to logout user",
			Details: err.Error(),
			TraceID: c.GetString("trace_id"),
		})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Successfully logged out",
	})
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the auth service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth-service",
	})
}

// ReadyCheck godoc
// @Summary Ready check
// @Description Check if the auth service is ready to serve requests
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ready [get]
func (h *AuthHandler) ReadyCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "auth-service",
	})
}