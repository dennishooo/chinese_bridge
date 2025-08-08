package service

import (
	"chinese-bridge-game/internal/user/repository"

	"github.com/go-redis/redis/v8"
)

type UserService interface {
	// Interface methods will be defined in later tasks
}

type userService struct {
	repo        repository.UserRepository
	redisClient *redis.Client
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		repo:        repo,
		redisClient: redisClient,
	}
}