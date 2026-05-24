package health

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"
)

// Checker interface represents a system component that can be checked for health.
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

// Manager manages health checks.
type Manager struct {
	checkers []Checker
}

// NewManager creates a new health check manager.
func NewManager() *Manager {
	return &Manager{}
}

// Add registers a new checker.
func (m *Manager) Add(checker Checker) {
	m.checkers = append(m.checkers, checker)
}

// Handler returns an http.HandlerFunc that performs the health checks.
func (m *Manager) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		results := make(map[string]string)
		status := http.StatusOK

		for _, checker := range m.checkers {
			err := checker.Check(ctx)
			if err != nil {
				results[checker.Name()] = "down: " + err.Error()
				status = http.StatusServiceUnavailable
			} else {
				results[checker.Name()] = "up"
			}
		}

		response := map[string]any{
			"status":    "ok",
			"time":      time.Now().UTC().Format(time.RFC3339),
			"version":   getAppVersion(),
			"components": results,
		}
		
		if status != http.StatusOK {
			response["status"] = "error"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
	}
}

// getAppVersion returns the module version from build info or "dev".
func getAppVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return "dev"
}

