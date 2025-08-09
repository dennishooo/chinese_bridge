package service

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient interface for Redis operations
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
}

// Ensure redis.Client implements RedisClient interface
var _ RedisClient = (*redis.Client)(nil)