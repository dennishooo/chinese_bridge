package database

import (
	"time"

	"gorm.io/datatypes"
)

// User model with GORM tags
type User struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	GoogleID  string    `json:"google_id" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Stats             *UserStats          `json:"stats,omitempty" gorm:"foreignKey:UserID"`
	HostedRooms       []Room             `json:"hosted_rooms,omitempty" gorm:"foreignKey:HostID"`
	RoomParticipants  []RoomParticipant  `json:"room_participants,omitempty" gorm:"foreignKey:UserID"`
	GameParticipants  []GameParticipant  `json:"game_participants,omitempty" gorm:"foreignKey:UserID"`
}

// UserStats model for tracking player statistics
type UserStats struct {
	UserID          string  `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	GamesPlayed     int     `json:"games_played" gorm:"default:0"`
	GamesWon        int     `json:"games_won" gorm:"default:0"`
	GamesAsDeclarer int     `json:"games_as_declarer" gorm:"default:0"`
	DeclarerWins    int     `json:"declarer_wins" gorm:"default:0"`
	TotalPoints     int     `json:"total_points" gorm:"default:0"`
	AverageBid      float64 `json:"average_bid" gorm:"type:decimal(5,2);default:0"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Association
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Room model for game rooms
type Room struct {
	ID             string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	HostID         string    `json:"host_id" gorm:"type:varchar(36);not null"`
	MaxPlayers     int       `json:"max_players" gorm:"default:4"`
	CurrentPlayers int       `json:"current_players" gorm:"default:0"`
	Status         string    `json:"status" gorm:"default:'waiting'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Host         User               `json:"host" gorm:"foreignKey:HostID"`
	Participants []RoomParticipant  `json:"participants" gorm:"foreignKey:RoomID"`
	Games        []Game             `json:"games,omitempty" gorm:"foreignKey:RoomID"`
}

// RoomParticipant junction table for room membership
type RoomParticipant struct {
	RoomID   string    `json:"room_id" gorm:"type:varchar(36);primaryKey"`
	UserID   string    `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	Position int       `json:"position"` // 0-3 for seating position
	JoinedAt time.Time `json:"joined_at" gorm:"autoCreateTime"`

	// Associations
	Room Room `json:"room" gorm:"foreignKey:RoomID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Game model for individual game instances
type Game struct {
	ID          string     `json:"id" gorm:"type:varchar(36);primaryKey"`
	RoomID      string     `json:"room_id" gorm:"type:varchar(36);not null"`
	DeclarerID  *string    `json:"declarer_id" gorm:"type:varchar(36)"`
	TrumpSuit   *string    `json:"trump_suit"`
	Contract    int        `json:"contract"`
	FinalScore  int        `json:"final_score"`
	WinnerTeam  *string    `json:"winner_team"` // 'declarer' or 'defenders'
	GameData    datatypes.JSON `json:"game_data" gorm:"type:jsonb"` // Complete game state
	StartedAt   *time.Time `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Room         Room               `json:"room" gorm:"foreignKey:RoomID"`
	Declarer     *User              `json:"declarer,omitempty" gorm:"foreignKey:DeclarerID"`
	Participants []GameParticipant  `json:"participants" gorm:"foreignKey:GameID"`
}

// GameParticipant junction table for game participation
type GameParticipant struct {
	GameID         string `json:"game_id" gorm:"type:varchar(36);primaryKey"`
	UserID         string `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	Position       int    `json:"position"` // 0-3 for game position
	Role           string `json:"role"`     // 'declarer' or 'defender'
	PointsCaptured int    `json:"points_captured" gorm:"default:0"`

	// Associations
	Game Game `json:"game" gorm:"foreignKey:GameID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Session model for user authentication sessions
type Session struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Token     string    `json:"token" gorm:"not null;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// Association
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// GetAllModels returns all models for migration
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&UserStats{},
		&Room{},
		&RoomParticipant{},
		&Game{},
		&GameParticipant{},
		&Session{},
	}
}