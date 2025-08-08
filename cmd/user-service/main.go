package main

import (
	"log"
	"os"

	"chinese-bridge-game/internal/common/config"
	"chinese-bridge-game/internal/common/database"
	"chinese-bridge-game/internal/user/handler"
	"chinese-bridge-game/internal/user/repository"
	"chinese-bridge-game/internal/user/service"
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
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, redisClient)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

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
			"service": "user-service",
		})
	})
	api.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ready",
			"service": "user-service",
		})
	})
	
	// Protected routes (auth required)
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	userHandler.RegisterRoutes(protected)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("User service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}