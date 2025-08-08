package handler

import (
	"chinese-bridge-game/internal/user/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.GET("/stats", h.GetStats)
		users.GET("/history", h.GetHistory)
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get profile endpoint"})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update profile endpoint"})
}

func (h *UserHandler) GetStats(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get stats endpoint"})
}

func (h *UserHandler) GetHistory(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get history endpoint"})
}

func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"service": "user-service",
	})
}

func (h *UserHandler) ReadyCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ready",
		"service": "user-service",
	})
}