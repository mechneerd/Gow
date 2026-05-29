package livewire

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/mechneerd/gow/routing"
)

var (
	defaultManager *Manager
	managerOnce    sync.Once
)

// getManager returns a singleton Manager instance.
func getManager() *Manager {
	managerOnce.Do(func() {
		defaultManager = NewManager()
	})
	return defaultManager
}

// RegisterRoutes registers the standard Livewire endpoints.
// Usage in your routes file:
//
//	router := routing.NewRouter()
//	livewire.RegisterRoutes(router)
func RegisterRoutes(router *routing.Router) {
	manager := getManager()

	// Main Livewire update endpoint (used by livewire.js)
	router.Post("/livewire/update", func(w http.ResponseWriter, r *http.Request) error {
		Handler(manager)(w, r)
		return nil
	})

	// Demo endpoint for the landing page counter
	router.Get("/livewire/counter", func(w http.ResponseWriter, r *http.Request) error {
		comp := &demoCounter{}
		comp.Mount()
		id, html, state := Mount(comp)

		payload := map[string]any{
			"id":    id,
			"html":  html,
			"state": state,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
		return nil
	})
}

// demoCounter is a minimal self-contained component for the landing page demo only.
type demoCounter struct {
	BaseComponent
	Count int
}

func (c *demoCounter) Mount() {
	c.Count = 0
}

func (c *demoCounter) Render() string {
	return `<div wire:id="` + c.GetID() + `" class="p-6 bg-zinc-900 border border-zinc-700 rounded-3xl max-w-sm mx-auto">
		<div class="text-center">
			<div class="text-sm text-zinc-400 mb-1">Livewire Counter (Demo)</div>
			<div class="text-6xl font-semibold tabular-nums tracking-tighter mb-6">` + "0" + `</div>
			<div class="text-xs text-zinc-500">Replace this with your real component</div>
		</div>
	</div>`
}

