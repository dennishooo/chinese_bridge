package repository

import (
	"context"
	"errors"

	"chinese-bridge-game/internal/common/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *database.User) error
	GetUserByID(ctx context.Context, id string) (*database.User, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (*database.User, error)
	GetUserByEmail(ctx context.Context, email string) (*database.User, error)
	UpdateUser(ctx context.Context, user *database.User) error
	CreateSession(ctx context.Context, session *database.Session) error
	GetSessionByToken(ctx context.Context, token string) (*database.Session, error)
	DeleteSession(ctx context.Context, token string) error
	DeleteUserSessions(ctx context.Context, userID string) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *database.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) GetUserByID(ctx context.Context, id string) (*database.User, error) {
	var user database.User
	err := r.db.WithContext(ctx).
		Preload("Stats").
		First(&user, "id = ?", id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *authRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*database.User, error) {
	var user database.User
	err := r.db.WithContext(ctx).
		Preload("Stats").
		First(&user, "google_id = ?", googleID).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	var user database.User
	err := r.db.WithContext(ctx).
		Preload("Stats").
		First(&user, "email = ?", email).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *authRepository) UpdateUser(ctx context.Context, user *database.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) CreateSession(ctx context.Context, session *database.Session) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *authRepository) GetSessionByToken(ctx context.Context, token string) (*database.Session, error) {
	var session database.Session
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&session, "token = ?", token).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &session, nil
}

func (r *authRepository) DeleteSession(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&database.Session{}, "token = ?", token).Error
}

func (r *authRepository) DeleteUserSessions(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&database.Session{}, "user_id = ?", userID).Error
}