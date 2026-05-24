package cache

import (
	"context"
	"time"

	"github.com/mechneerd/gow/database/redis"
)

// RedisStore is a cache driver backed by Redis.
type RedisStore struct {
	client *redis.Client
	prefix string
	ctx    context.Context
}

// NewRedisStore creates a new Redis-backed cache store.
func NewRedisStore(client *redis.Client, prefix string) *RedisStore {
	return &RedisStore{
		client: client,
		prefix: prefix,
		ctx:    context.Background(),
	}
}

func (s *RedisStore) itemKey(key string) string {
	return s.prefix + key
}

// Get retrieves a value by key.
func (s *RedisStore) Get(key string) (any, error) {
	val, err := s.client.Get(s.ctx, s.itemKey(key))
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Put stores a value by key with an expiration time.
func (s *RedisStore) Put(key string, value any, ttl time.Duration) error {
	return s.client.Set(s.ctx, s.itemKey(key), value, ttl)
}

// Increment increments an integer value in the cache.
func (s *RedisStore) Increment(key string, value int) (int, error) {
	res, err := s.client.RawClient().IncrBy(s.ctx, s.itemKey(key), int64(value)).Result()
	return int(res), err
}

// Decrement decrements an integer value in the cache.
func (s *RedisStore) Decrement(key string, value int) (int, error) {
	res, err := s.client.RawClient().DecrBy(s.ctx, s.itemKey(key), int64(value)).Result()
	return int(res), err
}

// Forever stores a value by key indefinitely.
func (s *RedisStore) Forever(key string, value any) error {
	return s.client.Set(s.ctx, s.itemKey(key), value, 0)
}

// Forget removes a specific key from the cache.
func (s *RedisStore) Forget(key string) error {
	return s.client.Del(s.ctx, s.itemKey(key))
}

// Flush clears the entire cache (current database).
func (s *RedisStore) Flush() error {
	return s.client.RawClient().FlushDB(s.ctx).Err()
}

