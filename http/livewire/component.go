package livewire

import (
	"reflect"
	"sync"
)

// BaseComponent should be embedded in your Livewire components.
type BaseComponent struct {
	ID    string
	Dirty map[string]bool
	mu    sync.RWMutex
}

// Lifecycle methods (no-op defaults so embedding structs satisfy the interface)
func (b *BaseComponent) Mount()                  {}
func (b *BaseComponent) Updated(property string) {}
func (b *BaseComponent) Rendering()              {}

// GetID returns the component's unique ID.
func (b *BaseComponent) GetID() string {
	return b.ID
}

// SetID sets the component ID (used internally).
func (b *BaseComponent) SetID(id string) {
	b.ID = id
}

// MarkDirty marks a property as changed.
func (b *BaseComponent) MarkDirty(field string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Dirty == nil {
		b.Dirty = make(map[string]bool)
	}
	b.Dirty[field] = true
}

// GetState returns all exported fields as a map (for sending to frontend).
func (b *BaseComponent) GetState() map[string]any {
	state := make(map[string]any)
	v := reflect.ValueOf(b).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue // unexported
		}
		state[field.Name] = v.Field(i).Interface()
	}
	return state
}

// Hydrate sets the component's state from a map (used when receiving updates from frontend).
func (b *BaseComponent) Hydrate(data map[string]any) {
	v := reflect.ValueOf(b).Elem()

	for key, val := range data {
		field := v.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			rv := reflect.ValueOf(val)
			if rv.Type().ConvertibleTo(field.Type()) {
				field.Set(rv.Convert(field.Type()))
			}
		}
	}
}

// Component is the interface all Livewire components must implement.
type Component interface {
	Render() string
	GetID() string
	SetID(id string)
	GetState() map[string]any
	Hydrate(data map[string]any)
	MarkDirty(field string)

	// Lifecycle hooks (optional)
	Mount()                  // Called when component is first created
	Updated(property string) // Called after a property is updated
	Rendering()              // Called before Render()
}

// Action represents a method call from the frontend.
type Action struct {
	Name   string
	Params []any
}
