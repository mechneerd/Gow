package events

import (
	"reflect"
	"sync"
)

// Event represents an application event.
type Event any

// Listener is a function that handles an event.
type Listener func(Event)

// Subscriber represents a class that subscribes to multiple events.
type Subscriber interface {
	Subscribe(dispatcher *Manager)
}

// ShouldBroadcast indicates an event should be broadcasted over WebSockets.
type ShouldBroadcast interface {
	BroadcastOn() []string // Channel names
	BroadcastAs() string   // Event name
	BroadcastWith() map[string]any // JSON payload
}

// Broadcaster interface decouples events from the broadcasting package.
type Broadcaster interface {
	Broadcast(channels []string, eventName string, payload map[string]any) error
}

// Manager orchestrates application events.
type Manager struct {
	mu          sync.RWMutex
	listeners   map[reflect.Type][]Listener
	wildcards   []Listener
	broadcaster Broadcaster
}

// NewManager creates a new event manager.
func NewManager() *Manager {
	return &Manager{
		listeners: make(map[reflect.Type][]Listener),
		wildcards: make([]Listener, 0),
	}
}

// SetBroadcaster assigns the broadcasting manager for ShouldBroadcast events.
func (m *Manager) SetBroadcaster(b Broadcaster) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.broadcaster = b
}

// Listen registers a listener for a specific event type.
func (m *Manager) Listen(event Event, listener Listener) {
	m.mu.Lock()
	defer m.mu.Unlock()
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	m.listeners[eventType] = append(m.listeners[eventType], listener)
}

// ListenAny registers a wildcard listener that catches all events.
func (m *Manager) ListenAny(listener Listener) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wildcards = append(m.wildcards, listener)
}

// Subscribe registers an event subscriber.
func (m *Manager) Subscribe(subscriber Subscriber) {
	subscriber.Subscribe(m)
}

// Dispatch triggers all registered listeners for the given event.
func (m *Manager) Dispatch(event Event) {
	m.mu.RLock()
	broadcaster := m.broadcaster
	wildcards := make([]Listener, len(m.wildcards))
	copy(wildcards, m.wildcards)
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	specificListeners := make([]Listener, len(m.listeners[eventType]))
	copy(specificListeners, m.listeners[eventType])
	m.mu.RUnlock()

	// Handle ShouldBroadcast
	if bEvent, ok := event.(ShouldBroadcast); ok && broadcaster != nil {
		broadcaster.Broadcast(bEvent.BroadcastOn(), bEvent.BroadcastAs(), bEvent.BroadcastWith())
	}

	// Execute wildcard listeners
	for _, listener := range wildcards {
		listener(event)
	}

	// Execute specific listeners
	for _, listener := range specificListeners {
		listener(event)
	}
}

// QueueListen registers a listener that should be queued.
func (m *Manager) QueueListen(event Event, listener Listener) {
	m.Listen(event, func(e Event) {
		listener(e)
	})
}

