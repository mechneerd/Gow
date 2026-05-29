package middleware

import (
	"context"
	"github.com/mechneerd/gow/auth"
	gowhttp "github.com/mechneerd/gow/http"
	"net/http"
)

// Authenticate verifies that the user is authenticated via the specified guard.
// Deprecated: Use auth.Middleware(guard) directly for new code.
func Authenticate(authManager *auth.Manager, guardName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			guard := authManager.Guard(guardName)
			if guard == nil || !guard.Check() {
				if r.Header.Get("Accept") == "application/json" {
					gowhttp.Abort(http.StatusUnauthorized, "Unauthenticated.")
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

// RedirectIfAuthenticated redirects the user if they are already logged in.
// Deprecated: Use auth.GuestMiddleware(guard, redirectTo) directly for new code.
func RedirectIfAuthenticated(authManager *auth.Manager, guardName string, redirectTo string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			guard := authManager.Guard(guardName)
			if guard != nil && guard.Check() {
				http.Redirect(w, r, redirectTo, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

