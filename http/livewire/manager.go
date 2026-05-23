package livewire

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

// Manager handles all Livewire components.
type Manager struct {
	components map[string]Component
	mu         sync.RWMutex
}

// NewManager creates a new Livewire manager.
func NewManager() *Manager {
	return &Manager{
		components: make(map[string]Component),
	}
}

// Register creates a new instance of a component and returns its ID.
func (m *Manager) Register(component Component) string {
	id := uuid.New().String()
	component.SetID(id)

	// Call Mount lifecycle hook if it exists
	if mounter, ok := component.(interface{ Mount() }); ok {
		mounter.Mount()
	}

	m.mu.Lock()
	m.components[id] = component
	m.mu.Unlock()

	return id
}

// Get retrieves a component by ID.
func (m *Manager) Get(id string) (Component, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	comp, ok := m.components[id]
	return comp, ok
}

// Remove deletes a component.
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	delete(m.components, id)
	m.mu.Unlock()
}

// Update handles an incoming Livewire update request from the frontend.
func (m *Manager) Update(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID     string                 `json:"id"`
		Method string                 `json:"method"`
		Params []any                  `json:"params"`
		State  map[string]any         `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	comp, ok := m.Get(payload.ID)
	if !ok {
		http.Error(w, "Component not found", http.StatusNotFound)
		return
	}

	// Hydrate state from frontend
	comp.Hydrate(payload.State)

	// Call the method if provided
	if payload.Method != "" {
		if err := callMethod(comp, payload.Method, payload.Params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Render the updated component
	html := comp.Render()

	response := map[string]any{
		"id":    payload.ID,
		"html":  html,
		"state": comp.GetState(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// callMethod uses reflection to call a method on the component.
func callMethod(component Component, methodName string, params []any) error {
	v := reflect.ValueOf(component)
	method := v.MethodByName(methodName)

	if !method.IsValid() {
		return fmt.Errorf("method %s not found on component", methodName)
	}

	// Convert params to reflect.Value
	args := make([]reflect.Value, len(params))
	for i, p := range params {
		args[i] = reflect.ValueOf(p)
	}

	// Call the method
	method.Call(args)
	return nil
}

// RenderComponent is a helper to render a component for the first time.
func (m *Manager) RenderComponent(w http.ResponseWriter, component Component) {
	id := m.Register(component)
	html := component.Render()

	response := map[string]any{
		"id":    id,
		"html":  html,
		"state": component.GetState(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Mount is a convenience helper to create and render a component in one call.
func Mount(component Component) (id string, html string, state map[string]any) {
	m := NewManager()
	id = m.Register(component)
	html = component.Render()
	state = component.GetState()
	return
}
