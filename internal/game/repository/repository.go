package repository

import (
	"gorm.io/gorm"
)

type GameRepository interface {
	// Interface methods will be defined in later tasks
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}