package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache interface defines caching operations
type Cache interface {
	// User session caching
	SetUserSession(ctx context.Context, userID string, sessionData interface{}, ttl time.Duration) error
	GetUserSession(ctx context.Context, userID string) (string, error)
	DeleteUserSession(ctx context.Context, userID string) error

	// Room state caching
	SetRoomState(ctx context.Context, roomID string, roomState interface{}, ttl time.Duration) error
	GetRoomState(ctx context.Context, roomID string) (string, error)
	DeleteRoomState(ctx context.Context, roomID string) error

	// Game state caching
	SetGameState(ctx context.Context, gameID string, gameState interface{}, ttl time.Duration) error
	GetGameState(ctx context.Context, gameID string) (string, error)
	DeleteGameState(ctx context.Context, gameID string) error

	// Leaderboard caching
	SetLeaderboard(ctx context.Context, leaderboardData interface{}, ttl time.Duration) error
	GetLeaderboard(ctx context.Context) (string, error)
	DeleteLeaderboard(ctx context.Context) error

	// WebSocket connection mapping
	SetWSConnection(ctx context.Context, userID string, connectionID string, ttl time.Duration) error
	GetWSConnection(ctx context.Context, userID string) (string, error)
	DeleteWSConnection(ctx context.Context, userID string) error

	// Matchmaking queue
	AddToMatchmakingQueue(ctx context.Context, userID string, userData interface{}) error
	RemoveFromMatchmakingQueue(ctx context.Context, userID string) error
	GetMatchmakingQueue(ctx context.Context, limit int) ([]string, error)

	// Generic operations
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetTTL(ctx context.Context, key string, ttl time.Duration) error
}

// redisCache implements the Cache interface using Redis
type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{client: client}
}

// Cache key constants
const (
	UserSessionKeyPrefix    = "session:user:"
	RoomStateKeyPrefix      = "room:state:"
	GameStateKeyPrefix      = "game:state:"
	LeaderboardKey          = "leaderboard:global"
	WSConnectionKeyPrefix   = "ws:user:"
	MatchmakingQueueKey     = "queue:matchmaking"
)

// Default TTL values
const (
	DefaultUserSessionTTL = 24 * time.Hour
	DefaultRoomStateTTL   = 30 * time.Minute
	DefaultGameStateTTL   = 2 * time.Hour
	DefaultLeaderboardTTL = 5 * time.Minute
	DefaultWSConnectionTTL = 1 * time.Hour
)

// User session operations
func (c *redisCache) SetUserSession(ctx context.Context, userID string, sessionData interface{}, ttl time.Duration) error {
	key := UserSessionKeyPrefix + userID
	return c.Set(ctx, key, sessionData, ttl)
}

func (c *redisCache) GetUserSession(ctx context.Context, userID string) (string, error) {
	key := UserSessionKeyPrefix + userID
	return c.Get(ctx, key)
}

func (c *redisCache) DeleteUserSession(ctx context.Context, userID string) error {
	key := UserSessionKeyPrefix + userID
	return c.Delete(ctx, key)
}

// Room state operations
func (c *redisCache) SetRoomState(ctx context.Context, roomID string, roomState interface{}, ttl time.Duration) error {
	key := RoomStateKeyPrefix + roomID
	return c.Set(ctx, key, roomState, ttl)
}

func (c *redisCache) GetRoomState(ctx context.Context, roomID string) (string, error) {
	key := RoomStateKeyPrefix + roomID
	return c.Get(ctx, key)
}

func (c *redisCache) DeleteRoomState(ctx context.Context, roomID string) error {
	key := RoomStateKeyPrefix + roomID
	return c.Delete(ctx, key)
}

// Game state operations
func (c *redisCache) SetGameState(ctx context.Context, gameID string, gameState interface{}, ttl time.Duration) error {
	key := GameStateKeyPrefix + gameID
	return c.Set(ctx, key, gameState, ttl)
}

func (c *redisCache) GetGameState(ctx context.Context, gameID string) (string, error) {
	key := GameStateKeyPrefix + gameID
	return c.Get(ctx, key)
}

func (c *redisCache) DeleteGameState(ctx context.Context, gameID string) error {
	key := GameStateKeyPrefix + gameID
	return c.Delete(ctx, key)
}

// Leaderboard operations
func (c *redisCache) SetLeaderboard(ctx context.Context, leaderboardData interface{}, ttl time.Duration) error {
	return c.Set(ctx, LeaderboardKey, leaderboardData, ttl)
}

func (c *redisCache) GetLeaderboard(ctx context.Context) (string, error) {
	return c.Get(ctx, LeaderboardKey)
}

func (c *redisCache) DeleteLeaderboard(ctx context.Context) error {
	return c.Delete(ctx, LeaderboardKey)
}

// WebSocket connection operations
func (c *redisCache) SetWSConnection(ctx context.Context, userID string, connectionID string, ttl time.Duration) error {
	key := WSConnectionKeyPrefix + userID
	return c.Set(ctx, key, connectionID, ttl)
}

func (c *redisCache) GetWSConnection(ctx context.Context, userID string) (string, error) {
	key := WSConnectionKeyPrefix + userID
	return c.Get(ctx, key)
}

func (c *redisCache) DeleteWSConnection(ctx context.Context, userID string) error {
	key := WSConnectionKeyPrefix + userID
	return c.Delete(ctx, key)
}

// Matchmaking queue operations
func (c *redisCache) AddToMatchmakingQueue(ctx context.Context, userID string, userData interface{}) error {
	data, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	// Add to sorted set with current timestamp as score
	score := float64(time.Now().Unix())
	return c.client.ZAdd(ctx, MatchmakingQueueKey, &redis.Z{
		Score:  score,
		Member: userID + ":" + string(data),
	}).Err()
}

func (c *redisCache) RemoveFromMatchmakingQueue(ctx context.Context, userID string) error {
	// Remove all entries that start with userID
	members, err := c.client.ZRange(ctx, MatchmakingQueueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, member := range members {
		if len(member) > len(userID) && member[:len(userID)+1] == userID+":" {
			if err := c.client.ZRem(ctx, MatchmakingQueueKey, member).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *redisCache) GetMatchmakingQueue(ctx context.Context, limit int) ([]string, error) {
	return c.client.ZRange(ctx, MatchmakingQueueKey, 0, int64(limit-1)).Result()
}

// Generic operations
func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return result, err
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	return count > 0, err
}

func (c *redisCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

// CachedData structures for type-safe caching
type CachedUserSession struct {
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CachedRoomState struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	HostID         string    `json:"host_id"`
	Players        []string  `json:"players"`
	Status         string    `json:"status"`
	CurrentPlayers int       `json:"current_players"`
	MaxPlayers     int       `json:"max_players"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CachedGameState struct {
	ID           string                 `json:"id"`
	RoomID       string                 `json:"room_id"`
	Phase        string                 `json:"phase"`
	Players      []string               `json:"players"`
	Declarer     *string                `json:"declarer,omitempty"`
	TrumpSuit    *string                `json:"trump_suit,omitempty"`
	Contract     int                    `json:"contract"`
	GameData     map[string]interface{} `json:"game_data"`
	LastActivity time.Time              `json:"last_activity"`
}

type CachedLeaderboard struct {
	Players   []LeaderboardEntry `json:"players"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type LeaderboardEntry struct {
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	GamesWon    int     `json:"games_won"`
	GamesPlayed int     `json:"games_played"`
	WinRate     float64 `json:"win_rate"`
}

type CachedMatchmakingUser struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	SkillLevel int      `json:"skill_level"`
	JoinedAt  time.Time `json:"joined_at"`
}