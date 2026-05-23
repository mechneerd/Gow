package inertia

import (
	"encoding/json"
	"net/http"
)

// Inertia renders Inertia.js responses.
type Inertia struct {
	RootView string
}

func NewInertia(rootView string) *Inertia {
	return &Inertia{RootView: rootView}
}

func (i *Inertia) Render(w http.ResponseWriter, component string, props map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Inertia", "true")
	json.NewEncoder(w).Encode(map[string]any{
		"component": component,
		"props":     props,
		"url":       "",
	})
}
