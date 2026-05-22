package auth

import (
	"net/http"
)

// RequireVerified is a middleware that ensures the authenticated user has verified their email.
func RequireVerified(guard Guard, redirectTo string) func(http.Handler) http.Handler {
	if redirectTo == "" {
		redirectTo = "/email/verify"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !guard.Check() {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			user := guard.User()
			if u, ok := user.(interface{ IsEmailVerified() bool }); ok {
				if !u.IsEmailVerified() {
					http.Redirect(w, r, redirectTo, http.StatusFound)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
