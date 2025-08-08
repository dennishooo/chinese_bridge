package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// gormRepository implements the Repository interface using GORM
type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM repository instance
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// User operations
func (r *gormRepository) CreateUser(ctx context.Context, user *User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Preload("Stats").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Preload("Stats").
		First(&user, "google_id = ?", googleID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Preload("Stats").
		First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormRepository) UpdateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *gormRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error
}

// Room operations
func (r *gormRepository) CreateRoom(ctx context.Context, room *Room) error {
	if room.ID == "" {
		room.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *gormRepository) GetRoomByID(ctx context.Context, id string) (*Room, error) {
	var room Room
	err := r.db.WithContext(ctx).
		Preload("Host").
		Preload("Participants").
		Preload("Participants.User").
		First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *gormRepository) GetRoomsByStatus(ctx context.Context, status string, limit, offset int) ([]Room, error) {
	var rooms []Room
	err := r.db.WithContext(ctx).
		Preload("Host").
		Preload("Participants").
		Preload("Participants.User").
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Find(&rooms).Error
	return rooms, err
}

func (r *gormRepository) UpdateRoom(ctx context.Context, room *Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *gormRepository) DeleteRoom(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Room{}, "id = ?", id).Error
}

func (r *gormRepository) AddRoomParticipant(ctx context.Context, participant *RoomParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

func (r *gormRepository) RemoveRoomParticipant(ctx context.Context, roomID, userID string) error {
	return r.db.WithContext(ctx).
		Delete(&RoomParticipant{}, "room_id = ? AND user_id = ?", roomID, userID).Error
}

func (r *gormRepository) GetRoomParticipants(ctx context.Context, roomID string) ([]RoomParticipant, error) {
	var participants []RoomParticipant
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("room_id = ?", roomID).
		Find(&participants).Error
	return participants, err
}

// Game operations
func (r *gormRepository) CreateGame(ctx context.Context, game *Game) error {
	if game.ID == "" {
		game.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(game).Error
}

func (r *gormRepository) GetGameByID(ctx context.Context, id string) (*Game, error) {
	var game Game
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("Declarer").
		Preload("Participants").
		Preload("Participants.User").
		First(&game, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gormRepository) GetGameByRoomID(ctx context.Context, roomID string) (*Game, error) {
	var game Game
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("Declarer").
		Preload("Participants").
		Preload("Participants.User").
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gormRepository) UpdateGame(ctx context.Context, game *Game) error {
	return r.db.WithContext(ctx).Save(game).Error
}

func (r *gormRepository) DeleteGame(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Game{}, "id = ?", id).Error
}

func (r *gormRepository) GetUserGameHistory(ctx context.Context, userID string, limit, offset int) ([]Game, error) {
	var games []Game
	err := r.db.WithContext(ctx).
		Preload("Room").
		Preload("Declarer").
		Preload("Participants").
		Preload("Participants.User").
		Joins("JOIN game_participants ON games.id = game_participants.game_id").
		Where("game_participants.user_id = ?", userID).
		Order("games.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&games).Error
	return games, err
}

func (r *gormRepository) AddGameParticipant(ctx context.Context, participant *GameParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

func (r *gormRepository) GetGameParticipants(ctx context.Context, gameID string) ([]GameParticipant, error) {
	var participants []GameParticipant
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("game_id = ?", gameID).
		Find(&participants).Error
	return participants, err
}

// Session operations
func (r *gormRepository) CreateSession(ctx context.Context, session *Session) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *gormRepository) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	var session Session
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *gormRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]Session, error) {
	var sessions []Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Find(&sessions).Error
	return sessions, err
}

func (r *gormRepository) UpdateSession(ctx context.Context, session *Session) error {
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *gormRepository) DeleteSession(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Session{}, "id = ?", id).Error
}

func (r *gormRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Delete(&Session{}, "expires_at <= ?", time.Now()).Error
}

// Statistics operations
func (r *gormRepository) CreateUserStats(ctx context.Context, stats *UserStats) error {
	return r.db.WithContext(ctx).Create(stats).Error
}

func (r *gormRepository) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	var stats UserStats
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&stats, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *gormRepository) UpdateUserStats(ctx context.Context, stats *UserStats) error {
	return r.db.WithContext(ctx).Save(stats).Error
}

func (r *gormRepository) GetLeaderboard(ctx context.Context, limit int) ([]UserStats, error) {
	var stats []UserStats
	err := r.db.WithContext(ctx).
		Preload("User").
		Order("games_won DESC, games_played ASC").
		Limit(limit).
		Find(&stats).Error
	return stats, err
}

func (r *gormRepository) GetTopPlayersByWins(ctx context.Context, limit int) ([]UserStats, error) {
	var stats []UserStats
	err := r.db.WithContext(ctx).
		Preload("User").
		Order("games_won DESC").
		Limit(limit).
		Find(&stats).Error
	return stats, err
}

func (r *gormRepository) GetTopPlayersByDeclarerWins(ctx context.Context, limit int) ([]UserStats, error) {
	var stats []UserStats
	err := r.db.WithContext(ctx).
		Preload("User").
		Order("declarer_wins DESC").
		Limit(limit).
		Find(&stats).Error
	return stats, err
}