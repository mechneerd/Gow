package session

import (
	"context"
	"encoding/json"
	"time"

	"gow/database/redis"
)

// RedisStore is a session driver backed by Redis.
type RedisStore struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
	ctx    context.Context
}

// NewRedisStore creates a new Redis session store.
func NewRedisStore(client *redis.Client, prefix string, ttl time.Duration) *RedisStore {
	return &RedisStore{
		client: client,
		prefix: prefix,
		ttl:    ttl,
		ctx:    context.Background(),
	}
}

func (s *RedisStore) sessionKey(id string) string {
	return s.prefix + id
}

// Read retrieves the session data by ID.
func (s *RedisStore) Read(id string) (map[string]any, error) {
	val, err := s.client.Get(s.ctx, s.sessionKey(id))
	if err != nil {
		// return empty session rather than error if key doesn't exist
		return make(map[string]any), nil
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return make(map[string]any), nil
	}
	return data, nil
}

// Write saves the session data by ID.
func (s *RedisStore) Write(id string, data map[string]any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.client.Set(s.ctx, s.sessionKey(id), string(bytes), s.ttl)
}

// Destroy removes the session data by ID.
func (s *RedisStore) Destroy(id string) error {
	return s.client.Del(s.ctx, s.sessionKey(id))
}
