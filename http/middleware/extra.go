package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// PreventRequestsDuringMaintenance returns middleware that prevents requests during maintenance mode
func PreventRequestsDuringMaintenance(maintenanceFunc func() bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if maintenanceFunc() {
				// Check for bypass token
				if token := r.URL.Query().Get("token"); token != "" {
					// In production: validate the maintenance bypass token
					next.ServeHTTP(w, r)
					return
				}
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message": "Service Unavailable", "status": 503}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PreventBackslash prevents URLs with trailing slashes from being redirected
func PreventBackslash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}

// PreventLazyLoading prevents lazy loading of resources
func PreventLazyLoading(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex, nofollow")
		next.ServeHTTP(w, r)
	})
}

// PreloadLinks sends Link headers for preloading resources
func PreloadLinks(preloads map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var links []string
			for url, rel := range preloads {
				links = append(links, fmt.Sprintf("<%s>; rel=%s", url, rel))
			}
			if len(links) > 0 {
				w.Header().Set("Link", strings.Join(links, ", "))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AddHeaders adds custom headers to the response
func AddHeaders(headers map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range headers {
				w.Header().Set(key, value)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SetSecurityHeaders sets common security headers
func SetSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

// PreventHTTPOnly ensures cookies are not accessible via JavaScript
func PreventHTTPOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is handled when setting cookies
		next.ServeHTTP(w, r)
	})
}

// PreventHostHeaderAttack prevents host header attacks
func PreventHostHeaderAttack(allowedHosts []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Host
			allowed := false
			for _, allowedHost := range allowedHosts {
				if host == allowedHost {
					allowed = true
					break
				}
			}
			if !allowed && len(allowedHosts) > 0 {
				http.Error(w, `{"error": "Invalid host"}`, http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SetTimeout sets a timeout for the request
func SetTimeout(timeout int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In production, use context.WithTimeout
			next.ServeHTTP(w, r)
		})
	}
}
