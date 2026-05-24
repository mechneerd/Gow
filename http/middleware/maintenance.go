package middleware

import (
	"net/http"
	"os"
	"path/filepath"
)

// CheckForMaintenanceMode middleware returns a 503 Service Unavailable if the app is down.
func CheckForMaintenanceMode(basePath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			downPath := filepath.Join(basePath, "storage", "framework", "down")
			if _, err := os.Stat(downPath); err == nil {
				// We could also read the JSON and return a custom message/view
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 Service Unavailable"))
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

