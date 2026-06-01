package events

import (
	"reflect"
	"sync"
)

// Event represents an application event.
type Event any

// Listener is a function that handles an event.
type Listener func(Event)

// StoppableEvent is an event that can stop propagation to subsequent listeners.
type StoppableEvent interface {
	StopPropagation()
	IsPropagationStopped() bool
}

// eventStoppable is a base struct that implements StoppableEvent.
type eventStoppable struct {
	propagationStopped bool
}

func (e *eventStoppable) StopPropagation() {
	e.propagationStopped = true
}

func (e *eventStoppable) IsPropagationStopped() bool {
	return e.propagationStopped
}

// ShouldQueue indicates that an event listener should be queued for async processing.
type ShouldQueue interface {
	HandleQueued(event Event) error
}

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

// QueueDispatcher is an interface for dispatching events to a queue.
type QueueDispatcher interface {
	DispatchToQueue(event Event, handler ShouldQueue) error
}

// Manager orchestrates application events.
type Manager struct {
	mu             sync.RWMutex
	listeners      map[reflect.Type][]Listener
	queuedListeners map[reflect.Type][]ShouldQueue
	wildcards      []Listener
	broadcaster    Broadcaster
	queueDispatcher QueueDispatcher
	fakes          []*Fake
}

// NewManager creates a new event manager.
func NewManager() *Manager {
	return &Manager{
		listeners:       make(map[reflect.Type][]Listener),
		queuedListeners: make(map[reflect.Type][]ShouldQueue),
		wildcards:       make([]Listener, 0),
	}
}

// SetBroadcaster assigns the broadcasting manager for ShouldBroadcast events.
func (m *Manager) SetBroadcaster(b Broadcaster) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.broadcaster = b
}

// SetQueueDispatcher sets the queue dispatcher for queued listeners.
func (m *Manager) SetQueueDispatcher(qd QueueDispatcher) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queueDispatcher = qd
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

// QueueListener registers a queued listener for a specific event type.
func (m *Manager) QueueListener(event Event, handler ShouldQueue) {
	m.mu.Lock()
	defer m.mu.Unlock()
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	m.queuedListeners[eventType] = append(m.queuedListeners[eventType], handler)
}

// Subscribe registers an event subscriber.
func (m *Manager) Subscribe(subscriber Subscriber) {
	subscriber.Subscribe(m)
}

// Dispatch triggers all registered listeners for the given event.
// If the event implements StoppableEvent, propagation stops when StopPropagation() is called.
func (m *Manager) Dispatch(event Event) {
	// Check for fakes
	m.mu.RLock()
	if len(m.fakes) > 0 {
		for _, fake := range m.fakes {
			fake.Dispatch(event)
		}
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()

	m.mu.RLock()
	broadcaster := m.broadcaster
	queueDispatcher := m.queueDispatcher
	wildcards := make([]Listener, len(m.wildcards))
	copy(wildcards, m.wildcards)
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	specificListeners := make([]Listener, len(m.listeners[eventType]))
	copy(specificListeners, m.listeners[eventType])
	queuedListeners := make([]ShouldQueue, len(m.queuedListeners[eventType]))
	copy(queuedListeners, m.queuedListeners[eventType])
	m.mu.RUnlock()

	// Handle ShouldBroadcast
	if bEvent, ok := event.(ShouldBroadcast); ok && broadcaster != nil {
		broadcaster.Broadcast(bEvent.BroadcastOn(), bEvent.BroadcastAs(), bEvent.BroadcastWith())
	}

	// Execute wildcard listeners
	for _, listener := range wildcards {
		listener(event)
		if stopEvent, ok := event.(StoppableEvent); ok && stopEvent.IsPropagationStopped() {
			return
		}
	}

	// Execute specific listeners
	for _, listener := range specificListeners {
		listener(event)
		if stopEvent, ok := event.(StoppableEvent); ok && stopEvent.IsPropagationStopped() {
			return
		}
	}

	// Execute queued listeners
	for _, handler := range queuedListeners {
		if queueDispatcher != nil {
			queueDispatcher.DispatchToQueue(event, handler)
		} else {
			// Fallback to synchronous execution
			handler.HandleQueued(event)
		}
	}
}

// Until dispatches an event and stops as soon as a listener returns a non-nil value.
func (m *Manager) Until(event Event) any {
	m.mu.RLock()
	wildcards := make([]Listener, len(m.wildcards))
	copy(wildcards, m.wildcards)
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	specificListeners := make([]Listener, len(m.listeners[eventType]))
	copy(specificListeners, m.listeners[eventType])
	m.mu.RUnlock()

	for _, listener := range wildcards {
		listener(event)
		if stopEvent, ok := event.(StoppableEvent); ok && stopEvent.IsPropagationStopped() {
			return nil
		}
	}

	for _, listener := range specificListeners {
		listener(event)
		if stopEvent, ok := event.(StoppableEvent); ok && stopEvent.IsPropagationStopped() {
			return nil
		}
	}

	return nil
}

// Forget removes all listeners for a specific event type.
func (m *Manager) Forget(event Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	delete(m.listeners, eventType)
	delete(m.queuedListeners, eventType)
}

// ForgetListener removes a specific listener from an event type.
func (m *Manager) ForgetListener(event Event, listener Listener) {
	m.mu.Lock()
	defer m.mu.Unlock()
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}

	listeners := m.listeners[eventType]
	for i, l := range listeners {
		if reflect.ValueOf(l).Pointer() == reflect.ValueOf(listener).Pointer() {
			m.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
			return
		}
	}
}

// HasListeners checks if there are any listeners registered for the given event type.
func (m *Manager) HasListeners(event Event) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}

	return len(m.listeners[eventType]) > 0 || len(m.wildcards) > 0 || len(m.queuedListeners[eventType]) > 0
}

// ListenerCount returns the number of listeners registered for the given event type.
func (m *Manager) ListenerCount(event Event) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}

	return len(m.listeners[eventType]) + len(m.wildcards) + len(m.queuedListeners[eventType])
}

// Fake sets a fake event dispatcher for testing.
func (m *Manager) Fake(fake *Fake) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fakes = append(m.fakes, fake)
}

// QueueListen registers a listener that should be queued.
// Deprecated: Use QueueListener instead.
func (m *Manager) QueueListen(event Event, listener Listener) {
	m.Listen(event, func(e Event) {
		listener(e)
	})
}
