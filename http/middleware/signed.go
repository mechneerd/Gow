package middleware

import (
	"gow/http"
	"gow/routing"
	"net/http"
)

// ValidateSignature is a middleware that requires the request to have a valid signed URL signature.
// Use with routes that were generated via SignedRoute or TemporarySignedRoute.
func ValidateSignature(generator *routing.URLGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !generator.HasValidSignature(r) {
				http.Abort(http.StatusForbidden, "Invalid signature")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
