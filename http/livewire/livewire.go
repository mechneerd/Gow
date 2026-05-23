package livewire

import "net/http"

// Livewire is a minimal Livewire-like component system.
type Component interface {
	Render() string
}

type Manager struct{}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) Handle(w http.ResponseWriter, r *http.Request, component Component) {
	w.Write([]byte(component.Render()))
}
