package database

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db *gorm.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// RunMigrations executes all database migrations
func (m *MigrationManager) RunMigrations(ctx context.Context) error {
	log.Println("Starting database migrations...")

	// Enable UUID extension for PostgreSQL
	if err := m.enableUUIDExtension(ctx); err != nil {
		return fmt.Errorf("failed to enable UUID extension: %w", err)
	}

	// Auto-migrate all models
	models := GetAllModels()
	for _, model := range models {
		if err := m.db.WithContext(ctx).AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
		log.Printf("Migrated model: %T", model)
	}

	// Create indexes for better performance
	if err := m.createIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// enableUUIDExtension enables the UUID extension in PostgreSQL
func (m *MigrationManager) enableUUIDExtension(ctx context.Context) error {
	// Check if we're using PostgreSQL
	if m.db.Dialector.Name() == "postgres" {
		return m.db.WithContext(ctx).Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	}
	// For other databases (like SQLite), skip UUID extension
	return nil
}

// createIndexes creates additional indexes for performance optimization
func (m *MigrationManager) createIndexes(ctx context.Context) error {
	indexes := []string{
		// User indexes
		"CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)",

		// Room indexes
		"CREATE INDEX IF NOT EXISTS idx_rooms_host_id ON rooms(host_id)",
		"CREATE INDEX IF NOT EXISTS idx_rooms_status ON rooms(status)",
		"CREATE INDEX IF NOT EXISTS idx_rooms_created_at ON rooms(created_at)",

		// Game indexes
		"CREATE INDEX IF NOT EXISTS idx_games_room_id ON games(room_id)",
		"CREATE INDEX IF NOT EXISTS idx_games_declarer_id ON games(declarer_id)",
		"CREATE INDEX IF NOT EXISTS idx_games_started_at ON games(started_at)",
		"CREATE INDEX IF NOT EXISTS idx_games_ended_at ON games(ended_at)",

		// Session indexes
		"CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)",

		// Statistics indexes
		"CREATE INDEX IF NOT EXISTS idx_user_stats_games_won ON user_stats(games_won)",
		"CREATE INDEX IF NOT EXISTS idx_user_stats_declarer_wins ON user_stats(declarer_wins)",
		"CREATE INDEX IF NOT EXISTS idx_user_stats_games_played ON user_stats(games_played)",

		// Junction table indexes
		"CREATE INDEX IF NOT EXISTS idx_room_participants_room_id ON room_participants(room_id)",
		"CREATE INDEX IF NOT EXISTS idx_room_participants_user_id ON room_participants(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_game_participants_game_id ON game_participants(game_id)",
		"CREATE INDEX IF NOT EXISTS idx_game_participants_user_id ON game_participants(user_id)",
	}

	for _, indexSQL := range indexes {
		if err := m.db.WithContext(ctx).Exec(indexSQL).Error; err != nil {
			log.Printf("Warning: Failed to create index: %s, Error: %v", indexSQL, err)
			// Continue with other indexes even if one fails
		}
	}

	return nil
}

// SeedData populates the database with initial test data
func (m *MigrationManager) SeedData(ctx context.Context) error {
	log.Println("Starting database seeding...")

	// Check if data already exists
	var userCount int64
	if err := m.db.WithContext(ctx).Model(&User{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if userCount > 0 {
		log.Println("Database already contains data, skipping seeding")
		return nil
	}

	// Create test users
	testUsers := []User{
		{
			GoogleID: "test_google_id_1",
			Email:    "player1@example.com",
			Name:     "Test Player 1",
			Avatar:   "https://example.com/avatar1.jpg",
		},
		{
			GoogleID: "test_google_id_2",
			Email:    "player2@example.com",
			Name:     "Test Player 2",
			Avatar:   "https://example.com/avatar2.jpg",
		},
		{
			GoogleID: "test_google_id_3",
			Email:    "player3@example.com",
			Name:     "Test Player 3",
			Avatar:   "https://example.com/avatar3.jpg",
		},
		{
			GoogleID: "test_google_id_4",
			Email:    "player4@example.com",
			Name:     "Test Player 4",
			Avatar:   "https://example.com/avatar4.jpg",
		},
	}

	// Create users and their stats
	for _, user := range testUsers {
		if err := m.db.WithContext(ctx).Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create test user %s: %w", user.Email, err)
		}

		// Create initial stats for each user
		stats := UserStats{
			UserID:          user.ID,
			GamesPlayed:     0,
			GamesWon:        0,
			GamesAsDeclarer: 0,
			DeclarerWins:    0,
			TotalPoints:     0,
			AverageBid:      0.0,
		}

		if err := m.db.WithContext(ctx).Create(&stats).Error; err != nil {
			return fmt.Errorf("failed to create stats for user %s: %w", user.Email, err)
		}

		log.Printf("Created test user: %s", user.Email)
	}

	// Create a test room
	if len(testUsers) >= 1 {
		testRoom := Room{
			Name:           "Test Room",
			HostID:         testUsers[0].ID,
			MaxPlayers:     4,
			CurrentPlayers: 1,
			Status:         "waiting",
		}

		if err := m.db.WithContext(ctx).Create(&testRoom).Error; err != nil {
			return fmt.Errorf("failed to create test room: %w", err)
		}

		// Add host as participant
		participant := RoomParticipant{
			RoomID:   testRoom.ID,
			UserID:   testUsers[0].ID,
			Position: 0,
		}

		if err := m.db.WithContext(ctx).Create(&participant).Error; err != nil {
			return fmt.Errorf("failed to add room participant: %w", err)
		}

		log.Printf("Created test room: %s", testRoom.Name)
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// DropAllTables drops all tables (useful for testing)
func (m *MigrationManager) DropAllTables(ctx context.Context) error {
	log.Println("Dropping all tables...")

	models := GetAllModels()
	// Reverse order to handle foreign key constraints
	for i := len(models) - 1; i >= 0; i-- {
		if err := m.db.WithContext(ctx).Migrator().DropTable(models[i]); err != nil {
			log.Printf("Warning: Failed to drop table for model %T: %v", models[i], err)
		}
	}

	log.Println("All tables dropped successfully")
	return nil
}

// RunMigrations is a convenience function to run migrations
func RunMigrations(db *gorm.DB) error {
	manager := NewMigrationManager(db)
	return manager.RunMigrations(context.Background())
}