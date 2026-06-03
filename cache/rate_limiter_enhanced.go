package cache

import (
	"sync"
	"time"
)

// SlidingWindowRateLimiter uses a sliding window algorithm for rate limiting.
type SlidingWindowRateLimiter struct {
	store    Store
	mu       sync.RWMutex
	windows  map[string]*slidingWindow
	window   time.Duration
	limit    int
}

type slidingWindow struct {
	events   []time.Time
	window   time.Duration
	limit    int
}

// NewSlidingWindowRateLimiter creates a new sliding window rate limiter.
func NewSlidingWindowRateLimiter(store Store, window time.Duration, limit int) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		store:   store,
		windows: make(map[string]*slidingWindow),
		window:  window,
		limit:   limit,
	}
}

// Allow checks if a request is allowed under the rate limit.
func (sw *SlidingWindowRateLimiter) Allow(key string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	window, exists := sw.windows[key]

	if !exists {
		sw.windows[key] = &slidingWindow{
			events:  []time.Time{now},
			window:  sw.window,
			limit:   sw.limit,
		}
		return true
	}

	// Remove expired events
	cutoff := now.Add(-window.window)
	validEvents := make([]time.Time, 0)
	for _, event := range window.events {
		if event.After(cutoff) {
			validEvents = append(validEvents, event)
		}
	}
	window.events = validEvents

	// Check limit
	if len(window.events) >= window.limit {
		return false
	}

	window.events = append(window.events, now)
	return true
}

// Remaining returns the number of remaining requests in the window.
func (sw *SlidingWindowRateLimiter) Remaining(key string) int {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	window, exists := sw.windows[key]
	if !exists {
		return 10
	}

	now := time.Now()
	cutoff := now.Add(-window.window)
	count := 0
	for _, event := range window.events {
		if event.After(cutoff) {
			count++
		}
	}

	remaining := window.limit - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset resets the rate limiter for a key.
func (sw *SlidingWindowRateLimiter) Reset(key string) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	delete(sw.windows, key)
}

// TokenBucketRateLimiter uses a token bucket algorithm for rate limiting.
type TokenBucketRateLimiter struct {
	tokens    map[string]*tokenBucket
	mu        sync.RWMutex
	rate      int           // tokens per second
	capacity  int           // max tokens
}

type tokenBucket struct {
	tokens    int
	lastTime  time.Time
	rate      int
	capacity  int
}

// NewTokenBucketRateLimiter creates a new token bucket rate limiter.
func NewTokenBucketRateLimiter(rate, capacity int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		tokens:   make(map[string]*tokenBucket),
		rate:     rate,
		capacity: capacity,
	}
}

// Allow checks if a request is allowed.
func (tb *TokenBucketRateLimiter) Allow(key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	bucket, exists := tb.tokens[key]

	if !exists {
		tb.tokens[key] = &tokenBucket{
			tokens:   tb.capacity - 1,
			lastTime: now,
			rate:     tb.rate,
			capacity: tb.capacity,
		}
		return true
	}

	// Calculate tokens to add based on elapsed time
	elapsed := now.Sub(bucket.lastTime)
	tokensToAdd := int(elapsed.Seconds() * float64(bucket.rate))

	bucket.tokens += tokensToAdd
	if bucket.tokens > bucket.capacity {
		bucket.tokens = bucket.capacity
	}
	bucket.lastTime = now

	// Check if token available
	if bucket.tokens <= 0 {
		return false
	}

	bucket.tokens--
	return true
}

// Remaining returns the number of remaining tokens.
func (tb *TokenBucketRateLimiter) Remaining(key string) int {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	bucket, exists := tb.tokens[key]
	if !exists {
		return tb.capacity
	}
	return bucket.tokens
}

// Reset resets the token bucket for a key.
func (tb *TokenBucketRateLimiter) Reset(key string) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	delete(tb.tokens, key)
}

// Wait blocks until a token is available or timeout.
func (tb *TokenBucketRateLimiter) Wait(key string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for {
		if tb.Allow(key) {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(10 * time.Millisecond)
	}
}
