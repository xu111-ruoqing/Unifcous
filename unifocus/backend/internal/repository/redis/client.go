package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/unifocus/backend/internal/config"
	"github.com/unifocus/backend/pkg/logger"
)

// Client wraps redis.Client with namespace management and helper methods
type Client struct {
	*redis.Client
	namespace string // Key namespace to avoid key conflicts
}

// NewClient creates a new Redis client with connection pool and namespace
// The namespace is prepended to all keys to avoid conflicts in multi-tenant scenarios
func NewClient(cfg *config.RedisConfig, namespace string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Infof("Redis connection established: %s (namespace: %s)", cfg.GetAddr(), namespace)

	return &Client{
		Client:    rdb,
		namespace: namespace,
	}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
}

// Key returns a namespaced key
// Example: Key("user:123") -> "unifocus:user:123"
func (c *Client) Key(key string) string {
	if c.namespace == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", c.namespace, key)
}

// Cache operations

// Set stores a value in cache with expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, c.Key(key), value, expiration).Err()
}

// Get retrieves a value from cache
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, c.Key(key)).Result()
}

// GetBytes retrieves a value as bytes from cache
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.Client.Get(ctx, c.Key(key)).Bytes()
}

// Delete removes a key from cache
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	namespacedKeys := make([]string, len(keys))
	for i, key := range keys {
		namespacedKeys[i] = c.Key(key)
	}
	return c.Client.Del(ctx, namespacedKeys...).Err()
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.Client.Exists(ctx, c.Key(key)).Result()
	return count > 0, err
}

// Queue operations

// Enqueue adds an item to a queue (list)
func (c *Client) Enqueue(ctx context.Context, queueName string, value interface{}) error {
	return c.Client.LPush(ctx, c.Key(queueName), value).Err()
}

// Dequeue removes and returns an item from a queue (blocking)
func (c *Client) Dequeue(ctx context.Context, queueName string, timeout time.Duration) (string, error) {
	result, err := c.Client.BRPop(ctx, timeout, c.Key(queueName)).Result()
	if err != nil {
		return "", err
	}
	if len(result) < 2 {
		return "", fmt.Errorf("invalid queue result")
	}
	return result[1], nil
}

// QueueLength returns the length of a queue
func (c *Client) QueueLength(ctx context.Context, queueName string) (int64, error) {
	return c.Client.LLen(ctx, c.Key(queueName)).Result()
}

// Distributed lock operations

// Lock acquires a distributed lock with expiration
// Returns true if lock was acquired, false if already locked
func (c *Client) Lock(ctx context.Context, lockKey string, expiration time.Duration) (bool, error) {
	key := c.Key(fmt.Sprintf("lock:%s", lockKey))
	result, err := c.Client.SetNX(ctx, key, "1", expiration).Result()
	return result, err
}

// Unlock releases a distributed lock
func (c *Client) Unlock(ctx context.Context, lockKey string) error {
	key := c.Key(fmt.Sprintf("lock:%s", lockKey))
	return c.Client.Del(ctx, key).Err()
}

// HealthCheck performs a health check on the Redis connection
func (c *Client) HealthCheck(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}
