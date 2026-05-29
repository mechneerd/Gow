package middleware

import (
	"crypto/rand"
	"encoding/hex"
	gowhttp "github.com/mechneerd/gow/http"
	"net/http"
)

// VerifyCsrfToken validates the CSRF token on mutating requests.
func VerifyCsrfToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			manager := GetSession(r)
			if manager == nil {
				gowhttp.Abort(http.StatusInternalServerError, "Session not initialized")
				return
			}

			// Generate token if not exists
			token := manager.Get("_token")
		if token == nil {
				token = generateToken()
				manager.Put("_token", token)
			}

			// Convert token to string safely
			tokenStr, ok := token.(string)
			if !ok {
				// Token is not a string (e.g., deserialized as different type), regenerate
				token = generateToken()
				manager.Put("_token", token)
				tokenStr = token.(string)
			}

			// For read operations, skip validation
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// Check token
			requestToken := r.Header.Get("X-CSRF-TOKEN")
			if requestToken == "" {
				requestToken = r.FormValue("_token")
			}

			if requestToken != tokenStr {
				gowhttp.Abort(419, "CSRF token mismatch")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

