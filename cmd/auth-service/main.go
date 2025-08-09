package main

import (
	"log"
	"os"

	"chinese-bridge-game/docs"
	"chinese-bridge-game/internal/auth/handler"
	"chinese-bridge-game/internal/auth/repository"
	"chinese-bridge-game/internal/auth/service"
	"chinese-bridge-game/internal/common/config"
	"chinese-bridge-game/internal/common/database"
	"chinese-bridge-game/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Chinese Bridge Auth Service API
// @version 1.0
// @description Authentication service for Chinese Bridge card game platform
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	// Database migrations are handled manually for now
	log.Println("Skipping automatic migrations - using manual schema")

	// Initialize Redis
	redisClient := database.NewRedisClient(cfg.RedisURL)

	// Test Redis connection
	if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)

	// Initialize services
	authService := service.NewAuthService(authRepo, redisClient, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.Default()
	
	// Apply global middleware
	router.Use(middleware.TraceID())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.RateLimiter(100, 200)) // Global rate limit: 100 req/sec, burst 200

	// Swagger documentation
	docs.SwaggerInfo.Host = "localhost:" + getPort()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	api := router.Group("/api/v1")
	
	// Health check routes (no auth required)
	api.GET("/health", authHandler.HealthCheck)
	api.GET("/ready", authHandler.ReadyCheck)
	
	// Auth routes
	authHandler.RegisterRoutes(api)

	// Start server
	port := getPort()
	log.Printf("Auth service starting on port %s", port)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}