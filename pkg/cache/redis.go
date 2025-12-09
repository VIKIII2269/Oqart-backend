package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/oqart/backend/config"
	"github.com/oqart/backend/pkg/logger"
)

type Cache struct {
	client *redis.Client
	config *config.RedisConfig
	logger *logger.Logger
}

func NewRedis(cfg *config.RedisConfig, log *logger.Logger) (*Cache, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Successfully connected to Redis",
		"host", cfg.Host,
		"port", cfg.Port,
		"db", cfg.DB)

	return &Cache{
		client: client,
		config: cfg,
		logger: log,
	}, nil
}

func (c *Cache) Close() error {
	c.logger.Info("Closing Redis connection")
	return c.client.Close()
}

// Set stores a value in cache with expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SetString stores a string value
func (c *Cache) SetString(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// GetString retrieves a string value
func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Increment increments a counter
func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// IncrementBy increments a counter by a specific value
func (c *Cache) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

// Decrement decrements a counter
func (c *Cache) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// Expire sets expiration on a key
func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL gets the remaining time to live
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// FlushAll clears all keys
func (c *Cache) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

// FlushDB clears current database
func (c *Cache) FlushDB(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// SetNX sets a value only if the key doesn't exist
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.client.SetNX(ctx, key, data, expiration).Result()
}

// HSet sets a field in a hash
func (c *Cache) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return c.client.HSet(ctx, key, field, value).Err()
}

// HGet gets a field from a hash
func (c *Cache) HGet(ctx context.Context, key string, field string) (string, error) {
	return c.client.HGet(ctx, key, field).Result()
}

// HGetAll gets all fields from a hash
func (c *Cache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel deletes fields from a hash
func (c *Cache) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

// SAdd adds members to a set
func (c *Cache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SAdd(ctx, key, members...).Err()
}

// SMembers gets all members of a set
func (c *Cache) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

// SRem removes members from a set
func (c *Cache) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SRem(ctx, key, members...).Err()
}

// ZAdd adds members to a sorted set
func (c *Cache) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return c.client.ZAdd(ctx, key, members...).Err()
}

// ZRange gets a range from a sorted set
func (c *Cache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// Health check
func (c *Cache) Health(ctx context.Context) error {
	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}
	return nil
}

// Custom errors
var (
	ErrCacheMiss = fmt.Errorf("cache miss")
)
