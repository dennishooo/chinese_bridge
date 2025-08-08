package repository

import (
	"gorm.io/gorm"
)

type AuthRepository interface {
	// Interface methods will be defined in later tasks
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}