package repository

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	// Interface methods will be defined in later tasks
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}