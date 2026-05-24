package middleware

import (
	"crypto/tls"
	"net/http"
	"strings"
)

// TrustedProxies middleware modifies the request's RemoteAddr to respect
// X-Forwarded-For headers if the request comes from a trusted proxy.
type TrustedProxies struct {
	Proxies []string
	Headers []string // e.g. X-Forwarded-For, X-Forwarded-Proto
}

// NewTrustedProxies creates a new TrustedProxies middleware.
// If Proxies contains "*", all proxies are trusted.
func NewTrustedProxies(proxies []string) *TrustedProxies {
	return &TrustedProxies{
		Proxies: proxies,
	}
}

func (tp *TrustedProxies) isTrusted(addr string) bool {
	// Simple IP extraction
	ip := addr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		ip = addr[:idx]
	}

	for _, proxy := range tp.Proxies {
		if proxy == "*" || proxy == ip {
			return true
		}
	}
	return false
}

// Handle implements the middleware signature.
func (tp *TrustedProxies) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tp.isTrusted(r.RemoteAddr) {
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				// The right-most IP is the most reliable
				ips := strings.Split(forwardedFor, ",")
				clientIP := strings.TrimSpace(ips[0])
				if clientIP != "" {
					r.RemoteAddr = clientIP + ":0" // fake port to maintain format
				}
			}
			
			if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
				if forwardedProto == "https" {
					r.TLS = &tls.ConnectionState{} // Hack to trick some downstream handlers that it's secure
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

