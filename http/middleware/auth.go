package middleware

import (
	"context"
	"gow/auth"
	gowhttp "gow/http"
	"net/http"
)

// Authenticate verifies that the user is authenticated via the specified guard.
func Authenticate(authManager *auth.Manager, guardName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			guard := authManager.Guard(guardName)
			if guard == nil || !guard.Check() {
				// For simplicity, we check if it's an API request vs Web.
				if r.Header.Get("Accept") == "application/json" {
					gowhttp.Abort(http.StatusUnauthorized, "Unauthenticated.")
				} else {
					http.Redirect(w, r, "/login", http.StatusFound)
				}
				return
			}

			// Add the authenticated user to the request context
			// so downstream handlers can easily fetch the user.
			ctx := context.WithValue(r.Context(), auth.UserContextKey, guard.User())
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// RedirectIfAuthenticated redirects the user if they are already logged in.
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
