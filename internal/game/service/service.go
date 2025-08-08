package service

import (
	"chinese-bridge-game/internal/game/repository"

	"github.com/go-redis/redis/v8"
)

type GameService interface {
	// Interface methods will be defined in later tasks
}

type gameService struct {
	repo        repository.GameRepository
	redisClient *redis.Client
}

func NewGameService(repo repository.GameRepository, redisClient *redis.Client) GameService {
	return &gameService{
		repo:        repo,
		redisClient: redisClient,
	}
}