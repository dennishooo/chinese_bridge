package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

// CacheInvalidationStrategy defines cache invalidation policies
type CacheInvalidationStrategy interface {
	// InvalidateUserData invalidates all user-related cache entries
	InvalidateUserData(ctx context.Context, userID string) error

	// InvalidateRoomData invalidates all room-related cache entries
	InvalidateRoomData(ctx context.Context, roomID string) error

	// InvalidateGameData invalidates all game-related cache entries
	InvalidateGameData(ctx context.Context, gameID string) error

	// InvalidateLeaderboard invalidates leaderboard cache
	InvalidateLeaderboard(ctx context.Context) error

	// InvalidateExpiredEntries removes expired cache entries
	InvalidateExpiredEntries(ctx context.Context) error

	// SchedulePeriodicCleanup starts periodic cache cleanup
	SchedulePeriodicCleanup(ctx context.Context, interval time.Duration)
}

// cacheInvalidationManager implements cache invalidation strategies
type cacheInvalidationManager struct {
	cache Cache
}

// NewCacheInvalidationStrategy creates a new cache invalidation manager
func NewCacheInvalidationStrategy(cache Cache) CacheInvalidationStrategy {
	return &cacheInvalidationManager{
		cache: cache,
	}
}

// InvalidateUserData removes all user-related cache entries
func (c *cacheInvalidationManager) InvalidateUserData(ctx context.Context, userID string) error {
	var errors []error

	// Invalidate user session
	if err := c.cache.DeleteUserSession(ctx, userID); err != nil {
		errors = append(errors, fmt.Errorf("failed to invalidate user session: %w", err))
	}

	// Invalidate WebSocket connection
	if err := c.cache.DeleteWSConnection(ctx, userID); err != nil {
		errors = append(errors, fmt.Errorf("failed to invalidate WS connection: %w", err))
	}

	// Remove from matchmaking queue
	if err := c.cache.RemoveFromMatchmakingQueue(ctx, userID); err != nil {
		errors = append(errors, fmt.Errorf("failed to remove from matchmaking queue: %w", err))
	}

	// Invalidate leaderboard since user stats might have changed
	if err := c.cache.DeleteLeaderboard(ctx); err != nil {
		errors = append(errors, fmt.Errorf("failed to invalidate leaderboard: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple invalidation errors: %v", errors)
	}

	log.Printf("Invalidated cache data for user: %s", userID)
	return nil
}

// InvalidateRoomData removes all room-related cache entries
func (c *cacheInvalidationManager) InvalidateRoomData(ctx context.Context, roomID string) error {
	var errors []error

	// Invalidate room state
	if err := c.cache.DeleteRoomState(ctx, roomID); err != nil {
		errors = append(errors, fmt.Errorf("failed to invalidate room state: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple invalidation errors: %v", errors)
	}

	log.Printf("Invalidated cache data for room: %s", roomID)
	return nil
}

// InvalidateGameData removes all game-related cache entries
func (c *cacheInvalidationManager) InvalidateGameData(ctx context.Context, gameID string) error {
	var errors []error

	// Invalidate game state
	if err := c.cache.DeleteGameState(ctx, gameID); err != nil {
		errors = append(errors, fmt.Errorf("failed to invalidate game state: %w", err))
	}

	// Invalidate leaderboard since game completion affects stats
	if err := c.cache.DeleteLeaderboard(ctx); err != nil {
		errors = append(errors, fmt.Errorf("failed to invalidate leaderboard: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple invalidation errors: %v", errors)
	}

	log.Printf("Invalidated cache data for game: %s", gameID)
	return nil
}

// InvalidateLeaderboard removes leaderboard cache
func (c *cacheInvalidationManager) InvalidateLeaderboard(ctx context.Context) error {
	if err := c.cache.DeleteLeaderboard(ctx); err != nil {
		return fmt.Errorf("failed to invalidate leaderboard: %w", err)
	}

	log.Println("Invalidated leaderboard cache")
	return nil
}

// InvalidateExpiredEntries removes expired cache entries
// Note: Redis automatically handles TTL expiration, but this can be used for manual cleanup
func (c *cacheInvalidationManager) InvalidateExpiredEntries(ctx context.Context) error {
	log.Println("Expired entries cleanup completed (Redis handles TTL automatically)")
	return nil
}

// SchedulePeriodicCleanup starts periodic cache cleanup
func (c *cacheInvalidationManager) SchedulePeriodicCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Cache cleanup scheduler stopped")
				return
			case <-ticker.C:
				if err := c.InvalidateExpiredEntries(ctx); err != nil {
					log.Printf("Error during periodic cache cleanup: %v", err)
				}
			}
		}
	}()

	log.Printf("Started periodic cache cleanup with interval: %v", interval)
}

// CacheWarmupStrategy defines cache warming policies
type CacheWarmupStrategy interface {
	// WarmupUserData preloads frequently accessed user data
	WarmupUserData(ctx context.Context, userID string) error

	// WarmupLeaderboard preloads leaderboard data
	WarmupLeaderboard(ctx context.Context) error

	// WarmupActiveRooms preloads active room data
	WarmupActiveRooms(ctx context.Context) error
}

// cacheWarmupManager implements cache warming strategies
type cacheWarmupManager struct {
	cache      Cache
	repository Repository
}

// NewCacheWarmupStrategy creates a new cache warmup manager
func NewCacheWarmupStrategy(cache Cache, repository Repository) CacheWarmupStrategy {
	return &cacheWarmupManager{
		cache:      cache,
		repository: repository,
	}
}

// WarmupUserData preloads user session and stats data
func (c *cacheWarmupManager) WarmupUserData(ctx context.Context, userID string) error {
	// Get user from database
	user, err := c.repository.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user from database: %w", err)
	}

	// Cache user session data
	sessionData := CachedUserSession{
		UserID:    user.ID,
		UpdatedAt: time.Now(),
	}

	if err := c.cache.SetUserSession(ctx, userID, sessionData, DefaultUserSessionTTL); err != nil {
		return fmt.Errorf("failed to cache user session: %w", err)
	}

	log.Printf("Warmed up cache for user: %s", userID)
	return nil
}

// WarmupLeaderboard preloads leaderboard data from database
func (c *cacheWarmupManager) WarmupLeaderboard(ctx context.Context) error {
	// Get leaderboard from database
	stats, err := c.repository.GetLeaderboard(ctx, 100) // Top 100 players
	if err != nil {
		return fmt.Errorf("failed to get leaderboard from database: %w", err)
	}

	// Convert to cached format
	var entries []LeaderboardEntry
	for _, stat := range stats {
		winRate := 0.0
		if stat.GamesPlayed > 0 {
			winRate = float64(stat.GamesWon) / float64(stat.GamesPlayed)
		}

		entries = append(entries, LeaderboardEntry{
			UserID:      stat.UserID,
			Name:        stat.User.Name,
			Avatar:      stat.User.Avatar,
			GamesWon:    stat.GamesWon,
			GamesPlayed: stat.GamesPlayed,
			WinRate:     winRate,
		})
	}

	leaderboard := CachedLeaderboard{
		Players:   entries,
		UpdatedAt: time.Now(),
	}

	if err := c.cache.SetLeaderboard(ctx, leaderboard, DefaultLeaderboardTTL); err != nil {
		return fmt.Errorf("failed to cache leaderboard: %w", err)
	}

	log.Printf("Warmed up leaderboard cache with %d entries", len(entries))
	return nil
}

// WarmupActiveRooms preloads active room data
func (c *cacheWarmupManager) WarmupActiveRooms(ctx context.Context) error {
	// Get active rooms from database
	rooms, err := c.repository.GetRoomsByStatus(ctx, "waiting", 50, 0)
	if err != nil {
		return fmt.Errorf("failed to get active rooms from database: %w", err)
	}

	// Cache each room
	for _, room := range rooms {
		// Get participants
		participants, err := c.repository.GetRoomParticipants(ctx, room.ID)
		if err != nil {
			log.Printf("Warning: failed to get participants for room %s: %v", room.ID, err)
			continue
		}

		var playerIDs []string
		for _, participant := range participants {
			playerIDs = append(playerIDs, participant.UserID)
		}

		roomState := CachedRoomState{
			ID:             room.ID,
			Name:           room.Name,
			HostID:         room.HostID,
			Players:        playerIDs,
			Status:         room.Status,
			CurrentPlayers: room.CurrentPlayers,
			MaxPlayers:     room.MaxPlayers,
			UpdatedAt:      time.Now(),
		}

		if err := c.cache.SetRoomState(ctx, room.ID, roomState, DefaultRoomStateTTL); err != nil {
			log.Printf("Warning: failed to cache room %s: %v", room.ID, err)
			continue
		}
	}

	log.Printf("Warmed up cache for %d active rooms", len(rooms))
	return nil
}