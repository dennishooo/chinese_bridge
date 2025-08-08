package handler

import (
	"chinese-bridge-game/internal/game/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService service.GameService
}

func NewGameHandler(gameService service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Room-related routes
	rooms := router.Group("/rooms")
	{
		rooms.POST("/:roomId/start", h.StartGame)
	}
	
	// Game-related routes
	games := router.Group("/games")
	{
		games.GET("/:gameId", h.GetGameState)
		games.POST("/:gameId/bid", h.PlaceBid)
		games.POST("/:gameId/trump", h.DeclareTrump)
		games.POST("/:gameId/kitty", h.ExchangeKitty)
		games.POST("/:gameId/play", h.PlayCards)
	}
}

func (h *GameHandler) StartGame(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Start game endpoint"})
}

func (h *GameHandler) PlaceBid(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Place bid endpoint"})
}

func (h *GameHandler) DeclareTrump(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Declare trump endpoint"})
}

func (h *GameHandler) ExchangeKitty(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Exchange kitty endpoint"})
}

func (h *GameHandler) PlayCards(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Play cards endpoint"})
}

func (h *GameHandler) GetGameState(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get game state endpoint"})
}

func (h *GameHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"service": "game-service",
	})
}

func (h *GameHandler) ReadyCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ready",
		"service": "game-service",
	})
}