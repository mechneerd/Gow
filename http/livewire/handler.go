package livewire

import (
	"encoding/json"
	"net/http"

	"github.com/mechneerd/gow/http/middleware"
)

// Handler returns an http.HandlerFunc that processes Livewire update requests.
// It validates the CSRF token from the X-CSRF-TOKEN header before processing.
func Handler(manager *Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate CSRF token
		session := middleware.GetSession(r)
		if session != nil {
			token := session.Get("_token")
			if token != nil {
				tokenStr, _ := token.(string)
				requestToken := r.Header.Get("X-CSRF-TOKEN")
				if tokenStr != "" && requestToken != tokenStr {
					http.Error(w, `{"error": "CSRF token mismatch"}`, 419)
					return
				}
			}
		}

		manager.Update(w, r)
	}
}

// MountAndRender is a helper to mount a component and return the initial payload
// (useful for server-side rendering the first time).
func MountAndRender(w http.ResponseWriter, component Component) {
	id, html, state := Mount(component)

	payload := map[string]any{
		"id":    id,
		"html":  html,
		"state": state,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

