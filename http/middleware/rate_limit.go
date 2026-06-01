package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request is allowed
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	// Clean old requests
	requests := r.requests[key]
	validRequests := make([]time.Time, 0)
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	if len(validRequests) >= r.limit {
		return false
	}

	r.requests[key] = append(validRequests, now)
	return true
}

// Middleware returns HTTP middleware for rate limiting
func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := r.getKey(req)
		if !r.Allow(key) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, `{"error": "Too Many Requests"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func (r *RateLimiter) getKey(req *http.Request) string {
	// Use IP address as key
	ip := req.RemoteAddr
	if forwarded := req.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}
	return ip
}

// PerUser creates a rate limiter that limits per user
func (r *RateLimiter) PerUser(getUserID func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			key := getUserID(req)
			if !r.Allow(key) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, `{"error": "Too Many Requests"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
