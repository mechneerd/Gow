package middleware

import (
	"net/http"
	"os"
	"strings"
)

// TrustHostsMiddleware validates that the request Host header matches trusted hosts.
// If no trusted hosts are configured, all hosts are allowed.
func TrustHostsMiddleware(trustedHosts ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(trustedHosts) == 0 {
				// Check for trusted_hosts config
				if configHosts := os.Getenv("TRUSTED_HOSTS"); configHosts != "" {
					trustedHosts = strings.Split(configHosts, ",")
				} else {
					next.ServeHTTP(w, r)
					return
				}
			}

			host := r.Host
			// Remove port if present
			if idx := strings.LastIndex(host, ":"); idx != -1 {
				host = host[:idx]
			}

			for _, trusted := range trustedHosts {
				trusted = strings.TrimSpace(trusted)
				if trusted == host {
					next.ServeHTTP(w, r)
					return
				}
				// Check wildcard patterns like *.example.com
				if strings.HasPrefix(trusted, "*.") {
					suffix := trusted[1:] // .example.com
					if strings.HasSuffix(host, suffix) {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			http.Error(w, "Host not trusted", http.StatusBadGateway)
		})
	}
}
