package events

import (
	"reflect"
)

// Listener represents a callback that reacts to an event.
type Listener func(event any)

// Dispatcher manages the registration and dispatching of events.
type Dispatcher struct {
	listeners map[reflect.Type][]Listener
	wildcards []Listener
}

// NewDispatcher creates a new Event Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[reflect.Type][]Listener),
		wildcards: make([]Listener, 0),
	}
}

// Listen registers a listener for a specific event type.
func (d *Dispatcher) Listen(event any, listener Listener) {
	eventType := reflect.TypeOf(event)
	d.listeners[eventType] = append(d.listeners[eventType], listener)
}

// ListenAny registers a wildcard listener that catches all events.
func (d *Dispatcher) ListenAny(listener Listener) {
	d.wildcards = append(d.wildcards, listener)
}

// Dispatch triggers all registered listeners for the given event.
func (d *Dispatcher) Dispatch(event any) {
	eventType := reflect.TypeOf(event)

	// Execute wildcard listeners
	for _, listener := range d.wildcards {
		listener(event)
	}

	// Execute specific listeners
	if listeners, exists := d.listeners[eventType]; exists {
		for _, listener := range listeners {
			listener(event)
		}
	}
}

// Note: A "QueuedListener" integration would wrap the callback into a Queue Job
// and push it via the queue.Manager instead of executing it synchronously.
