package database

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// setupTestRedis creates a Redis client for testing
// Note: This requires a running Redis instance for integration tests
func setupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1, // Use DB 1 for testing
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		t.Skip("Redis not available for testing, skipping cache tests")
	}

	// Clean up test database
	client.FlushDB(ctx)

	return client
}

func TestRedisCache_UserSession(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	userID := "test-user-123"
	sessionData := CachedUserSession{
		UserID:    userID,
		Token:     "test-token-456",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	t.Run("SetUserSession", func(t *testing.T) {
		err := cache.SetUserSession(ctx, userID, sessionData, DefaultUserSessionTTL)
		assert.NoError(t, err)
	})

	t.Run("GetUserSession", func(t *testing.T) {
		result, err := cache.GetUserSession(ctx, userID)
		assert.NoError(t, err)
		assert.Contains(t, result, sessionData.Token)
		assert.Contains(t, result, sessionData.UserID)
	})

	t.Run("DeleteUserSession", func(t *testing.T) {
		err := cache.DeleteUserSession(ctx, userID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = cache.GetUserSession(ctx, userID)
		assert.Error(t, err)
	})
}

func TestRedisCache_RoomState(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	roomID := "test-room-123"
	roomState := CachedRoomState{
		ID:             roomID,
		Name:           "Test Room",
		HostID:         "host-user-123",
		Players:        []string{"player1", "player2"},
		Status:         "waiting",
		CurrentPlayers: 2,
		MaxPlayers:     4,
		UpdatedAt:      time.Now(),
	}

	t.Run("SetRoomState", func(t *testing.T) {
		err := cache.SetRoomState(ctx, roomID, roomState, DefaultRoomStateTTL)
		assert.NoError(t, err)
	})

	t.Run("GetRoomState", func(t *testing.T) {
		result, err := cache.GetRoomState(ctx, roomID)
		assert.NoError(t, err)
		assert.Contains(t, result, roomState.Name)
		assert.Contains(t, result, roomState.HostID)
		assert.Contains(t, result, "waiting")
	})

	t.Run("DeleteRoomState", func(t *testing.T) {
		err := cache.DeleteRoomState(ctx, roomID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = cache.GetRoomState(ctx, roomID)
		assert.Error(t, err)
	})
}

func TestRedisCache_GameState(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	gameID := "test-game-123"
	gameState := CachedGameState{
		ID:           gameID,
		RoomID:       "test-room-123",
		Phase:        "bidding",
		Players:      []string{"player1", "player2", "player3", "player4"},
		Contract:     120,
		GameData:     map[string]interface{}{"current_bid": 120, "bidder": "player1"},
		LastActivity: time.Now(),
	}

	t.Run("SetGameState", func(t *testing.T) {
		err := cache.SetGameState(ctx, gameID, gameState, DefaultGameStateTTL)
		assert.NoError(t, err)
	})

	t.Run("GetGameState", func(t *testing.T) {
		result, err := cache.GetGameState(ctx, gameID)
		assert.NoError(t, err)
		assert.Contains(t, result, gameState.Phase)
		assert.Contains(t, result, gameState.RoomID)
		assert.Contains(t, result, "120")
	})

	t.Run("DeleteGameState", func(t *testing.T) {
		err := cache.DeleteGameState(ctx, gameID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = cache.GetGameState(ctx, gameID)
		assert.Error(t, err)
	})
}

func TestRedisCache_Leaderboard(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	leaderboard := CachedLeaderboard{
		Players: []LeaderboardEntry{
			{
				UserID:      "player1",
				Name:        "Player One",
				Avatar:      "avatar1.jpg",
				GamesWon:    10,
				GamesPlayed: 15,
				WinRate:     0.67,
			},
			{
				UserID:      "player2",
				Name:        "Player Two",
				Avatar:      "avatar2.jpg",
				GamesWon:    8,
				GamesPlayed: 12,
				WinRate:     0.67,
			},
		},
		UpdatedAt: time.Now(),
	}

	t.Run("SetLeaderboard", func(t *testing.T) {
		err := cache.SetLeaderboard(ctx, leaderboard, DefaultLeaderboardTTL)
		assert.NoError(t, err)
	})

	t.Run("GetLeaderboard", func(t *testing.T) {
		result, err := cache.GetLeaderboard(ctx)
		assert.NoError(t, err)
		assert.Contains(t, result, "Player One")
		assert.Contains(t, result, "Player Two")
		assert.Contains(t, result, "0.67")
	})

	t.Run("DeleteLeaderboard", func(t *testing.T) {
		err := cache.DeleteLeaderboard(ctx)
		assert.NoError(t, err)

		// Verify deletion
		_, err = cache.GetLeaderboard(ctx)
		assert.Error(t, err)
	})
}

func TestRedisCache_WSConnection(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	userID := "test-user-123"
	connectionID := "ws-connection-456"

	t.Run("SetWSConnection", func(t *testing.T) {
		err := cache.SetWSConnection(ctx, userID, connectionID, DefaultWSConnectionTTL)
		assert.NoError(t, err)
	})

	t.Run("GetWSConnection", func(t *testing.T) {
		result, err := cache.GetWSConnection(ctx, userID)
		assert.NoError(t, err)
		assert.Contains(t, result, connectionID)
	})

	t.Run("DeleteWSConnection", func(t *testing.T) {
		err := cache.DeleteWSConnection(ctx, userID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = cache.GetWSConnection(ctx, userID)
		assert.Error(t, err)
	})
}

func TestRedisCache_MatchmakingQueue(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	user1ID := "user1"
	user1Data := CachedMatchmakingUser{
		UserID:     user1ID,
		Name:       "User One",
		SkillLevel: 1200,
		JoinedAt:   time.Now(),
	}

	user2ID := "user2"
	user2Data := CachedMatchmakingUser{
		UserID:     user2ID,
		Name:       "User Two",
		SkillLevel: 1150,
		JoinedAt:   time.Now(),
	}

	t.Run("AddToMatchmakingQueue", func(t *testing.T) {
		err := cache.AddToMatchmakingQueue(ctx, user1ID, user1Data)
		assert.NoError(t, err)

		err = cache.AddToMatchmakingQueue(ctx, user2ID, user2Data)
		assert.NoError(t, err)
	})

	t.Run("GetMatchmakingQueue", func(t *testing.T) {
		queue, err := cache.GetMatchmakingQueue(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, queue, 2)

		// Check that both users are in the queue
		found1, found2 := false, false
		for _, entry := range queue {
			if len(entry) > len(user1ID) && entry[:len(user1ID)] == user1ID {
				found1 = true
			}
			if len(entry) > len(user2ID) && entry[:len(user2ID)] == user2ID {
				found2 = true
			}
		}
		assert.True(t, found1, "User1 should be in queue")
		assert.True(t, found2, "User2 should be in queue")
	})

	t.Run("RemoveFromMatchmakingQueue", func(t *testing.T) {
		err := cache.RemoveFromMatchmakingQueue(ctx, user1ID)
		assert.NoError(t, err)

		queue, err := cache.GetMatchmakingQueue(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, queue, 1)

		// Check that only user2 remains
		assert.True(t, len(queue[0]) > len(user2ID) && queue[0][:len(user2ID)] == user2ID)
	})
}

func TestRedisCache_GenericOperations(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	key := "test-key"
	value := map[string]interface{}{
		"name":  "test",
		"value": 123,
		"flag":  true,
	}

	t.Run("Set", func(t *testing.T) {
		err := cache.Set(ctx, key, value, 1*time.Hour)
		assert.NoError(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		result, err := cache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Contains(t, result, "test")
		assert.Contains(t, result, "123")
		assert.Contains(t, result, "true")
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = cache.Exists(ctx, "non-existent-key")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("SetTTL", func(t *testing.T) {
		err := cache.SetTTL(ctx, key, 30*time.Second)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := cache.Delete(ctx, key)
		assert.NoError(t, err)

		// Verify deletion
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestRedisCache_TTLExpiration(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	key := "ttl-test-key"
	value := "test-value"

	t.Run("ShortTTL", func(t *testing.T) {
		// Set with very short TTL
		err := cache.Set(ctx, key, value, 100*time.Millisecond)
		assert.NoError(t, err)

		// Should exist immediately
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)

		// Wait for expiration
		time.Sleep(150 * time.Millisecond)

		// Should not exist after expiration
		exists, err = cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}