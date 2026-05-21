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

// Manager orchestrates application events.
type Manager struct {
	listeners map[string][]Listener
}

// NewManager creates a new event manager.
func NewManager() *Manager {
	return &Manager{
		listeners: make(map[string][]Listener),
	}
}

// Listen registers a listener for an event.
func (m *Manager) Listen(eventName string, listener Listener) {
	m.listeners[eventName] = append(m.listeners[eventName], listener)
}

// Subscribe registers an event subscriber.
func (m *Manager) Subscribe(subscriber Subscriber) {
	subscriber.Subscribe(m)
}

// Dispatch fires an event to all registered listeners.
func (m *Manager) Dispatch(event Event) {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	eventName := eventType.Name()

	if listeners, ok := m.listeners[eventName]; ok {
		for _, listener := range listeners {
			listener(event)
		}
	}
}

// QueueListen registers a listener that should be queued.
// In a full implementation, it would push to the Queue system instead of direct execution.
func (m *Manager) QueueListen(eventName string, listener Listener) {
	m.Listen(eventName, func(e Event) {
		// Simulate queuing the listener execution
		// queueManager.Push(...)
		listener(e)
	})
}
