package service

import (
	"chinese-bridge-game/internal/auth/repository"
	"chinese-bridge-game/internal/common/config"

	"github.com/go-redis/redis/v8"
)

type AuthService interface {
	// Interface methods will be defined in later tasks
}

type authService struct {
	repo        repository.AuthRepository
	redisClient *redis.Client
	config      *config.Config
}

func NewAuthService(repo repository.AuthRepository, redisClient *redis.Client, config *config.Config) AuthService {
	return &authService{
		repo:        repo,
		redisClient: redisClient,
		config:      config,
	}
}