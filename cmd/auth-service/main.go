package main

import (
	"log"
	"os"

	"chinese-bridge-game/internal/auth/handler"
	"chinese-bridge-game/internal/auth/repository"
	"chinese-bridge-game/internal/auth/service"
	"chinese-bridge-game/internal/common/config"
	"chinese-bridge-game/internal/common/database"
	"chinese-bridge-game/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Redis
	redisClient := database.NewRedisClient(cfg.RedisURL)

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)

	// Initialize services
	authService := service.NewAuthService(authRepo, redisClient, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Setup routes
	api := router.Group("/api/v1")
	
	// Health check routes (no auth required)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "auth-service",
		})
	})
	api.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ready",
			"service": "auth-service",
		})
	})
	
	// Auth routes (no auth required for login)
	authHandler.RegisterRoutes(api)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Auth service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}