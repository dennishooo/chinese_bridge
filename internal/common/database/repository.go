package database

import (
	"context"
)

// Repository interface defines all database operations
type Repository interface {
	UserRepository
	RoomRepository
	GameRepository
	SessionRepository
	StatsRepository
}

// UserRepository interface for user operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByGoogleID(ctx context.Context, googleID string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
}

// RoomRepository interface for room operations
type RoomRepository interface {
	CreateRoom(ctx context.Context, room *Room) error
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	GetRoomsByStatus(ctx context.Context, status string, limit, offset int) ([]Room, error)
	UpdateRoom(ctx context.Context, room *Room) error
	DeleteRoom(ctx context.Context, id string) error
	AddRoomParticipant(ctx context.Context, participant *RoomParticipant) error
	RemoveRoomParticipant(ctx context.Context, roomID, userID string) error
	GetRoomParticipants(ctx context.Context, roomID string) ([]RoomParticipant, error)
}

// GameRepository interface for game operations
type GameRepository interface {
	CreateGame(ctx context.Context, game *Game) error
	GetGameByID(ctx context.Context, id string) (*Game, error)
	GetGameByRoomID(ctx context.Context, roomID string) (*Game, error)
	UpdateGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, id string) error
	GetUserGameHistory(ctx context.Context, userID string, limit, offset int) ([]Game, error)
	AddGameParticipant(ctx context.Context, participant *GameParticipant) error
	GetGameParticipants(ctx context.Context, gameID string) ([]GameParticipant, error)
}

// SessionRepository interface for session operations
type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id string) error
	DeleteExpiredSessions(ctx context.Context) error
}

// StatsRepository interface for statistics operations
type StatsRepository interface {
	CreateUserStats(ctx context.Context, stats *UserStats) error
	GetUserStats(ctx context.Context, userID string) (*UserStats, error)
	UpdateUserStats(ctx context.Context, stats *UserStats) error
	GetLeaderboard(ctx context.Context, limit int) ([]UserStats, error)
	GetTopPlayersByWins(ctx context.Context, limit int) ([]UserStats, error)
	GetTopPlayersByDeclarerWins(ctx context.Context, limit int) ([]UserStats, error)
}