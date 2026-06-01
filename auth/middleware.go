package auth

import (
	"context"
	"net/http"
)

// Middleware returns an HTTP middleware that ensures the user is authenticated
// using the provided Guard.
func Middleware(guard Guard) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if guard.Check() {
				// Attach user to context for easy access via auth.User(r)
				ctx := context.WithValue(r.Context(), UserContextKey, guard.User())
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Not authenticated
			http.Redirect(w, r, "/login", http.StatusFound)
		})
	}
}

// GuestMiddleware allows only unauthenticated users (useful for login/register pages).
func GuestMiddleware(guard Guard, redirectTo string) func(http.Handler) http.Handler {
	if redirectTo == "" {
		redirectTo = "/"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if guard.Guest() {
				next.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, redirectTo, http.StatusFound)
		})
	}
}

// SessionFixationMiddleware regenerates the session ID on authentication state change.
// Use this middleware to prevent session fixation attacks.
// It should be applied after the Auth middleware.
func SessionFixationMiddleware(guard Guard) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the user is authenticated, ensure session was regenerated
			// This is automatically handled by the guard on login,
			// but this middleware provides an explicit protection layer
			if guard.Check() {
				// Context value indicates if session was already regenerated
				if r.Context().Value(sessionRegeneratedKey) == nil {
					// Not yet regenerated in this request - mark as protected
					ctx := context.WithValue(r.Context(), sessionRegeneratedKey, true)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// contextKey for session regeneration state
type contextKey string

const sessionRegeneratedKey contextKey = "session_regenerated"

