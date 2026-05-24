package events

import "reflect"

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
	m.broadcaster = b
}

// Listen registers a listener for a specific event type.
// You should pass an instance of the event type, e.g. Listen(UserCreated{}, func)
func (m *Manager) Listen(event Event, listener Listener) {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	m.listeners[eventType] = append(m.listeners[eventType], listener)
}

// ListenAny registers a wildcard listener that catches all events.
func (m *Manager) ListenAny(listener Listener) {
	m.wildcards = append(m.wildcards, listener)
}

// Subscribe registers an event subscriber.
func (m *Manager) Subscribe(subscriber Subscriber) {
	subscriber.Subscribe(m)
}

// Dispatch triggers all registered listeners for the given event.
func (m *Manager) Dispatch(event Event) {
	// Handle ShouldBroadcast
	if bEvent, ok := event.(ShouldBroadcast); ok && m.broadcaster != nil {
		m.broadcaster.Broadcast(bEvent.BroadcastOn(), bEvent.BroadcastAs(), bEvent.BroadcastWith())
	}

	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}

	// Execute wildcard listeners
	for _, listener := range m.wildcards {
		listener(event)
	}

	// Execute specific listeners
	if listeners, ok := m.listeners[eventType]; ok {
		for _, listener := range listeners {
			listener(event)
		}
	}
}

// QueueListen registers a listener that should be queued.
func (m *Manager) QueueListen(event Event, listener Listener) {
	m.Listen(event, func(e Event) {
		listener(e)
	})
}

