package middleware

import (
	"github.com/mechneerd/gow/cache"
	gowhttp "github.com/mechneerd/gow/http"
	"net"
	"net/http"
	"strconv"
)

// sharedDefaultLimiter is a package-level rate limiter shared across all Throttle() calls.
var sharedDefaultLimiter = cache.NewRateLimiter(cache.NewMemoryDriver())

// ThrottleRequests limits requests. Supports named limiters via key prefix.
func ThrottleRequests(limiter *cache.RateLimiter, maxAttempts int, decayMinutes int, name ...string) func(http.Handler) http.Handler {
	prefix := "throttle"
	if len(name) > 0 && name[0] != "" {
		prefix = "throttle:" + name[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Parse IP from RemoteAddr (strips port)
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr // fallback if no port
			}
			key := prefix + ":" + ip
			// Support user-based rate limiting if session user ID is available
			// (do NOT trust client-supplied headers for rate limiting)

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
// Uses a shared package-level rate limiter so all routes share the same rate limit state.
func Throttle(maxAttempts, decayMinutes int) func(http.Handler) http.Handler {
	return ThrottleRequests(sharedDefaultLimiter, maxAttempts, decayMinutes)
}

// ThrottleNamed allows "throttle:login:5,1" style named limiters.
func ThrottleNamed(name string, maxAttempts, decayMinutes int) func(http.Handler) http.Handler {
	return ThrottleRequests(sharedDefaultLimiter, maxAttempts, decayMinutes, name)
}

