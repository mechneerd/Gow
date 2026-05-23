package middleware

import (
	"gow/cache"
	gowhttp "gow/http"
	"net/http"
	"strconv"
	"strings"
)

// ThrottleRequests limits requests. Supports named limiters via key prefix.
func ThrottleRequests(limiter *cache.RateLimiter, maxAttempts int, decayMinutes int, name ...string) func(http.Handler) http.Handler {
	prefix := "throttle"
	if len(name) > 0 && name[0] != "" {
		prefix = "throttle:" + name[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Key: named limiter + IP (or user ID if authenticated in future)
			ip := r.RemoteAddr
			key := prefix + ":" + ip
			// Support user-based if "user:" in context or header (simple)
			if userID := r.Header.Get("X-User-ID"); userID != "" {
				key = prefix + ":user:" + userID
			}

			if limiter.TooManyAttempts(key, maxAttempts) {
				w.Header().Set("Retry-After", strconv.Itoa(decayMinutes*60))
				gowhttp.Abort(http.StatusTooManyRequests, "Too Many Requests")
				return
			}

			limiter.Hit(key, decayMinutes*60)

			attempts, _ := limiter.Attempts(key)
			remaining := maxAttempts - attempts
			if remaining < 0 {
				remaining = 0
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(maxAttempts))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			if len(name) > 0 {
				w.Header().Set("X-RateLimit-Name", name[0])
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Throttle is a convenience wrapper for simple "throttle:60,1" style (max, decayMinutes).
func Throttle(maxAttempts, decayMinutes int) func(http.Handler) http.Handler {
	return ThrottleRequests(cache.NewRateLimiter(), maxAttempts, decayMinutes)
}

// ThrottleNamed allows "throttle:login:5,1" style named limiters.
func ThrottleNamed(name string, maxAttempts, decayMinutes int) func(http.Handler) http.Handler {
	return ThrottleRequests(cache.NewRateLimiter(), maxAttempts, decayMinutes, name)
}
