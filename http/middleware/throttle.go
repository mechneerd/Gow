package middleware

import (
	"fmt"
	"gow/cache"
	gowhttp "gow/http"
	"net/http"
	"strconv"
)

// ThrottleRequests limits the number of requests a client can make within a timeframe.
func ThrottleRequests(limiter *cache.RateLimiter, maxAttempts int, decayMinutes int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In a real app, use the IP address or authenticated user ID.
			// For simplicity we use IP here.
			ip := r.RemoteAddr
			key := "throttle:" + ip

			if limiter.TooManyAttempts(key, maxAttempts) {
				w.Header().Set("Retry-After", strconv.Itoa(decayMinutes*60))
				gowhttp.Abort(http.StatusTooManyRequests, "Too Many Requests")
				return
			}

			limiter.Hit(key, decayMinutes*60)

			// Add headers for remaining attempts
			attempts, _ := limiter.Attempts(key)
			remaining := maxAttempts - attempts
			if remaining < 0 {
				remaining = 0
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(maxAttempts))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			next.ServeHTTP(w, r)
		})
	}
}
