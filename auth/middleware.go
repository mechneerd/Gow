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
