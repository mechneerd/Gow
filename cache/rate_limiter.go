package cache

import (
	"fmt"
	"time"
)

// RateLimiter provides rate limiting functionality using the Cache Store.
type RateLimiter struct {
	cache Store
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(cache Store) *RateLimiter {
	return &RateLimiter{cache: cache}
}

// TooManyAttempts determines if the given key has exceeded the max attempts.
func (rl *RateLimiter) TooManyAttempts(key string, maxAttempts int) bool {
	attempts, err := rl.Attempts(key)
	if err != nil {
		return false
	}
	return attempts >= maxAttempts
}

// Hit increments the number of attempts for a given key.
func (rl *RateLimiter) Hit(key string, decaySeconds int) (int, error) {
	val, err := rl.cache.Get(key)
	if err != nil {
		return 0, err
	}

	if val == nil {
		err = rl.cache.Put(key, 1, time.Duration(decaySeconds)*time.Second)
		if err != nil {
			return 0, err
		}
		// Also store a timer key
		rl.cache.Put(key+":timer", time.Now().Unix()+int64(decaySeconds), time.Duration(decaySeconds)*time.Second)
		return 1, nil
	}

	return rl.cache.Increment(key, 1)
}

// Attempts gets the number of attempts for the given key.
func (rl *RateLimiter) Attempts(key string) (int, error) {
	val, err := rl.cache.Get(key)
	if err != nil || val == nil {
		return 0, err
	}
	
	switch v := val.(type) {
	case int:
		return v, nil
	case float64: // JSON unmarshals numbers as float64
		return int(v), nil
	default:
		return 0, fmt.Errorf("invalid type for rate limiter attempts")
	}
}

// Reset resets the attempts for the given key.
func (rl *RateLimiter) Reset(key string) error {
	rl.cache.Forget(key)
	rl.cache.Forget(key + ":timer")
	return nil
}
