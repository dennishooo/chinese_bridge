package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) (*gorm.DB, Repository) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	migrationManager := NewMigrationManager(db)
	err = migrationManager.RunMigrations(context.Background())
	require.NoError(t, err)

	repo := NewGormRepository(db)
	return db, repo
}

func TestUserRepository(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	t.Run("CreateUser", func(t *testing.T) {
		user := &User{
			GoogleID: "test_google_id",
			Email:    "test@example.com",
			Name:     "Test User",
			Avatar:   "https://example.com/avatar.jpg",
		}

		err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		// Create a user first
		user := &User{
			GoogleID: "test_google_id_2",
			Email:    "test2@example.com",
			Name:     "Test User 2",
		}
		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Retrieve the user
		retrieved, err := repo.GetUserByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.ID)
		assert.Equal(t, user.Email, retrieved.Email)
		assert.Equal(t, user.Name, retrieved.Name)
	})

	t.Run("GetUserByGoogleID", func(t *testing.T) {
		user := &User{
			GoogleID: "test_google_id_3",
			Email:    "test3@example.com",
			Name:     "Test User 3",
		}
		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)

		retrieved, err := repo.GetUserByGoogleID(ctx, user.GoogleID)
		assert.NoError(t, err)
		assert.Equal(t, user.GoogleID, retrieved.GoogleID)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user := &User{
			GoogleID: "test_google_id_4",
			Email:    "test4@example.com",
			Name:     "Test User 4",
		}
		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Update the user
		user.Name = "Updated Name"
		err = repo.UpdateUser(ctx, user)
		assert.NoError(t, err)

		// Verify the update
		retrieved, err := repo.GetUserByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", retrieved.Name)
	})
}

func TestRoomRepository(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create a test user first
	user := &User{
		GoogleID: "host_google_id",
		Email:    "host@example.com",
		Name:     "Host User",
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("CreateRoom", func(t *testing.T) {
		room := &Room{
			Name:           "Test Room",
			HostID:         user.ID,
			MaxPlayers:     4,
			CurrentPlayers: 1,
			Status:         "waiting",
		}

		err := repo.CreateRoom(ctx, room)
		assert.NoError(t, err)
		assert.NotEmpty(t, room.ID)
		assert.NotZero(t, room.CreatedAt)
	})

	t.Run("GetRoomByID", func(t *testing.T) {
		room := &Room{
			Name:           "Test Room 2",
			HostID:         user.ID,
			MaxPlayers:     4,
			CurrentPlayers: 1,
			Status:         "waiting",
		}
		err := repo.CreateRoom(ctx, room)
		require.NoError(t, err)

		retrieved, err := repo.GetRoomByID(ctx, room.ID)
		assert.NoError(t, err)
		assert.Equal(t, room.Name, retrieved.Name)
		assert.Equal(t, room.HostID, retrieved.HostID)
		assert.NotNil(t, retrieved.Host)
		assert.Equal(t, user.Name, retrieved.Host.Name)
	})

	t.Run("AddRoomParticipant", func(t *testing.T) {
		room := &Room{
			Name:           "Test Room 3",
			HostID:         user.ID,
			MaxPlayers:     4,
			CurrentPlayers: 1,
			Status:         "waiting",
		}
		err := repo.CreateRoom(ctx, room)
		require.NoError(t, err)

		participant := &RoomParticipant{
			RoomID:   room.ID,
			UserID:   user.ID,
			Position: 0,
		}

		err = repo.AddRoomParticipant(ctx, participant)
		assert.NoError(t, err)

		participants, err := repo.GetRoomParticipants(ctx, room.ID)
		assert.NoError(t, err)
		assert.Len(t, participants, 1)
		assert.Equal(t, user.ID, participants[0].UserID)
	})
}

func TestGameRepository(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test user and room
	user := &User{
		GoogleID: "game_user_google_id",
		Email:    "gameuser@example.com",
		Name:     "Game User",
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	room := &Room{
		Name:           "Game Room",
		HostID:         user.ID,
		MaxPlayers:     4,
		CurrentPlayers: 1,
		Status:         "waiting",
	}
	err = repo.CreateRoom(ctx, room)
	require.NoError(t, err)

	t.Run("CreateGame", func(t *testing.T) {
		now := time.Now()
		game := &Game{
			RoomID:     room.ID,
			DeclarerID: &user.ID,
			TrumpSuit:  stringPtr("spades"),
			Contract:   120,
			FinalScore: 0,
			StartedAt:  &now,
		}

		err := repo.CreateGame(ctx, game)
		assert.NoError(t, err)
		assert.NotEmpty(t, game.ID)
		assert.NotZero(t, game.CreatedAt)
	})

	t.Run("GetGameByID", func(t *testing.T) {
		now := time.Now()
		game := &Game{
			RoomID:     room.ID,
			DeclarerID: &user.ID,
			TrumpSuit:  stringPtr("hearts"),
			Contract:   115,
			FinalScore: 0,
			StartedAt:  &now,
		}
		err := repo.CreateGame(ctx, game)
		require.NoError(t, err)

		retrieved, err := repo.GetGameByID(ctx, game.ID)
		assert.NoError(t, err)
		assert.Equal(t, game.Contract, retrieved.Contract)
		assert.Equal(t, *game.TrumpSuit, *retrieved.TrumpSuit)
		assert.NotNil(t, retrieved.Room)
		assert.NotNil(t, retrieved.Declarer)
	})

	t.Run("AddGameParticipant", func(t *testing.T) {
		game := &Game{
			RoomID:   room.ID,
			Contract: 110,
		}
		err := repo.CreateGame(ctx, game)
		require.NoError(t, err)

		participant := &GameParticipant{
			GameID:         game.ID,
			UserID:         user.ID,
			Position:       0,
			Role:           "declarer",
			PointsCaptured: 0,
		}

		err = repo.AddGameParticipant(ctx, participant)
		assert.NoError(t, err)

		participants, err := repo.GetGameParticipants(ctx, game.ID)
		assert.NoError(t, err)
		assert.Len(t, participants, 1)
		assert.Equal(t, "declarer", participants[0].Role)
	})
}

func TestSessionRepository(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test user
	user := &User{
		GoogleID: "session_user_google_id",
		Email:    "sessionuser@example.com",
		Name:     "Session User",
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("CreateSession", func(t *testing.T) {
		session := &Session{
			UserID:    user.ID,
			Token:     "test_token_123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err := repo.CreateSession(ctx, session)
		assert.NoError(t, err)
		assert.NotEmpty(t, session.ID)
	})

	t.Run("GetSessionByToken", func(t *testing.T) {
		session := &Session{
			UserID:    user.ID,
			Token:     "test_token_456",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		err := repo.CreateSession(ctx, session)
		require.NoError(t, err)

		retrieved, err := repo.GetSessionByToken(ctx, session.Token)
		assert.NoError(t, err)
		assert.Equal(t, session.Token, retrieved.Token)
		assert.Equal(t, session.UserID, retrieved.UserID)
		assert.NotNil(t, retrieved.User)
	})

	t.Run("DeleteExpiredSessions", func(t *testing.T) {
		// Create an expired session
		expiredSession := &Session{
			UserID:    user.ID,
			Token:     "expired_token",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		}
		err := repo.CreateSession(ctx, expiredSession)
		require.NoError(t, err)

		// Create a valid session
		validSession := &Session{
			UserID:    user.ID,
			Token:     "valid_token",
			ExpiresAt: time.Now().Add(1 * time.Hour), // Expires in 1 hour
		}
		err = repo.CreateSession(ctx, validSession)
		require.NoError(t, err)

		// Delete expired sessions
		err = repo.DeleteExpiredSessions(ctx)
		assert.NoError(t, err)

		// Verify expired session is gone
		_, err = repo.GetSessionByToken(ctx, expiredSession.Token)
		assert.Error(t, err)

		// Verify valid session still exists
		retrieved, err := repo.GetSessionByToken(ctx, validSession.Token)
		assert.NoError(t, err)
		assert.Equal(t, validSession.Token, retrieved.Token)
	})
}

func TestStatsRepository(t *testing.T) {
	_, repo := setupTestDB(t)
	ctx := context.Background()

	// Create test user
	user := &User{
		GoogleID: "stats_user_google_id",
		Email:    "statsuser@example.com",
		Name:     "Stats User",
	}
	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)

	t.Run("CreateUserStats", func(t *testing.T) {
		stats := &UserStats{
			UserID:          user.ID,
			GamesPlayed:     10,
			GamesWon:        6,
			GamesAsDeclarer: 4,
			DeclarerWins:    2,
			TotalPoints:     1200,
			AverageBid:      115.5,
		}

		err := repo.CreateUserStats(ctx, stats)
		assert.NoError(t, err)
		assert.NotZero(t, stats.CreatedAt)
	})

	t.Run("GetUserStats", func(t *testing.T) {
		// Create a new user for this test
		testUser := &User{
			GoogleID: "stats_user_google_id_2",
			Email:    "statsuser2@example.com",
			Name:     "Stats User 2",
		}
		err := repo.CreateUser(ctx, testUser)
		require.NoError(t, err)

		stats := &UserStats{
			UserID:          testUser.ID,
			GamesPlayed:     15,
			GamesWon:        9,
			GamesAsDeclarer: 6,
			DeclarerWins:    4,
			TotalPoints:     1800,
			AverageBid:      118.0,
		}
		err = repo.CreateUserStats(ctx, stats)
		require.NoError(t, err)

		retrieved, err := repo.GetUserStats(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, stats.GamesPlayed, retrieved.GamesPlayed)
		assert.Equal(t, stats.GamesWon, retrieved.GamesWon)
		assert.Equal(t, stats.AverageBid, retrieved.AverageBid)
		assert.NotNil(t, retrieved.User)
	})

	t.Run("UpdateUserStats", func(t *testing.T) {
		// Create a new user for this test
		testUser := &User{
			GoogleID: "stats_user_google_id_3",
			Email:    "statsuser3@example.com",
			Name:     "Stats User 3",
		}
		err := repo.CreateUser(ctx, testUser)
		require.NoError(t, err)

		stats := &UserStats{
			UserID:      testUser.ID,
			GamesPlayed: 20,
			GamesWon:    12,
		}
		err = repo.CreateUserStats(ctx, stats)
		require.NoError(t, err)

		// Update stats
		stats.GamesPlayed = 25
		stats.GamesWon = 15
		err = repo.UpdateUserStats(ctx, stats)
		assert.NoError(t, err)

		// Verify update
		retrieved, err := repo.GetUserStats(ctx, testUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, 25, retrieved.GamesPlayed)
		assert.Equal(t, 15, retrieved.GamesWon)
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}