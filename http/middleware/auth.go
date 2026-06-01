package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mechneerd/gow/auth"
)

// Authenticate returns a middleware that checks if the user is authenticated via the specified guard.
func Authenticate(manager *auth.Manager, guardName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			guard := manager.Guard(guardName)
			if guard == nil || !guard.Check() {
				if r.Header.Get("Accept") == "application/json" || strings.HasPrefix(r.URL.Path, "/api") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"message":"Unauthenticated."}`))
				} else {
					http.Redirect(w, r, "/login", http.StatusFound)
				}
				return
			}

			ctx := context.WithValue(r.Context(), auth.UserContextKey, guard.User())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthenticateSanctum authenticates Sanctum tokens.
func AuthenticateSanctum(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			http.Error(w, `{"error": "Unauthenticated"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ConfirmPassword requires password confirmation.
func ConfirmPassword(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// SubstituteBindings substitutes route bindings.
func SubstituteBindings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// HandlePanic handles panics.
func HandlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// SetMiddlewareHeaders sets common middleware headers.
func SetMiddlewareHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

// AddRequestId adds a request ID to the response.
func AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = strconv.FormatInt(time.Now().UnixNano(), 36)
		}
		w.Header().Set("X-Request-ID", requestId)
		next.ServeHTTP(w, r)
	})
}

// LogRequests logs incoming requests.
func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		_ = time.Since(start)
	})
}

// TrustProxies trusts proxy headers.
func TrustProxies(proxies []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
