package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// CorsConfig holds the configuration for CORS middleware.
type CorsConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
	MaxAge         int
	AllowCredentials bool
}

// DefaultCorsConfig provides a permissive default configuration.
func DefaultCorsConfig() CorsConfig {
	return CorsConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-TOKEN"},
		MaxAge:         86400,
	}
}

// Cors handles Cross-Origin Resource Sharing.
func Cors(config CorsConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				// Handle Allowed Origins
				allowOrigin := ""
				for _, o := range config.AllowedOrigins {
					if o == "*" || o == origin {
						allowOrigin = o
						break
					}
				}
				if allowOrigin != "" {
					w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				}

				// Handle Credentials
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				// Handle Preflight OPTIONS request
				if r.Method == "OPTIONS" {
					if len(config.AllowedMethods) > 0 {
						w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
					}
					if len(config.AllowedHeaders) > 0 {
						w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
					}
					w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
					
					w.WriteHeader(http.StatusNoContent)
					return
				}

				// Expose Headers
				if len(config.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
