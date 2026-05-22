package middleware

import (
	"gow/auth"
	"gow/auth/access"
	gowhttp "gow/http"
	"net/http"
)

// Authorize returns a middleware that checks if the authenticated user
// has the specified ability via the Gate. If the ability check requires models,
// they must be passed or resolved (this is a simplified implementation that only checks abilities without models,
// or abilities where the model can be resolved later, e.g. via Route Model Binding).
func Authorize(gate *access.Gate, ability string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := auth.User(r)

			if user == nil {
				if r.Header.Get("Accept") == "application/json" {
					gowhttp.Abort(http.StatusUnauthorized, "Unauthenticated.")
				} else {
					http.Redirect(w, r, "/login", http.StatusFound)
				}
				return
			}

			if gate.Denies(user, ability) {
				if r.Header.Get("Accept") == "application/json" {
					gowhttp.Abort(http.StatusForbidden, "This action is unauthorized.")
				} else {
					// Normally redirect or show 403 error page
					gowhttp.Abort(http.StatusForbidden, "Forbidden.")
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
