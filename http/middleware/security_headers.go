package middleware

import (
	"net/http"
)

// SecurityHeaders injects standard secure HTTP headers into the response.
func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Protects against clickjacking by preventing the page from being rendered in a frame.
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")

			// Protects against MIME sniffing.
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Enables the Cross-Site Scripting (XSS) filter built into most recent web browsers.
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Enforces secure (HTTP over SSL/TLS) connections — only over HTTPS
			if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			// Prevents browsers from referring to this page when following a link to another site.
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Optional Content-Security-Policy (can be overridden by the user config)
			w.Header().Set("Content-Security-Policy", "default-src 'self'")

			next.ServeHTTP(w, r)
		})
	}
}

