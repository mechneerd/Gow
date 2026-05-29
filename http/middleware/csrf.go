package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	gowhttp "github.com/mechneerd/gow/http"
)

// CsrfOptions configures the CSRF middleware.
type CsrfOptions struct {
	// URIs is a list of URI prefixes to exclude from CSRF checks (e.g., "/api/").
	URIs []string
}

// VerifyCsrfToken validates the CSRF token on mutating requests.
func VerifyCsrfToken(opts ...CsrfOptions) func(http.Handler) http.Handler {
	var options CsrfOptions
	if len(opts) > 0 {
		options = opts[0]
	}

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

			tokenStr, ok := token.(string)
			if !ok {
				token = generateToken()
				manager.Put("_token", token)
				tokenStr = token.(string)
			}

			// For read operations, skip validation
			if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
				// Inject token into request context for template access
				ctx := r.Context()
				ctx = context.WithValue(ctx, csrfTokenKey, tokenStr)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Check URI exceptions
			for _, prefix := range options.URIs {
				if strings.HasPrefix(r.URL.Path, prefix) {
					next.ServeHTTP(w, r)
					return
				}
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

// CsrfTokenKey is the context key for the CSRF token.
type csrfTokenKeyType string

const csrfTokenKey csrfTokenKeyType = "csrf_token"

// CsrfTokenFromContext retrieves the CSRF token from the request context.
func CsrfTokenFromContext(r *http.Request) string {
	if token, ok := r.Context().Value(csrfTokenKey).(string); ok {
		return token
	}
	return ""
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

