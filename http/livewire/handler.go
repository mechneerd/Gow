package livewire

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.HandlerFunc that processes Livewire update requests.
// Usage in routes:
//
//   router.Post("/livewire/update", livewire.Handler(livewireManager))
func Handler(manager *Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

