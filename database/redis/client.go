package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Config represents the Redis connection configuration.
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Client wraps the go-redis/v9 UniversalClient.
type Client struct {
	client redis.UniversalClient
}

// NewClient creates a new Redis client instance.
func NewClient(cfg Config) *Client {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Client{
		client: c,
	}
}

// Get retrieves a value by key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set stores a value by key with an optional expiration.
func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Del deletes one or more keys.
func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Ping pings the Redis server.
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// RawClient returns the underlying universal client for advanced usage.
func (c *Client) RawClient() redis.UniversalClient {
	return c.client
}
