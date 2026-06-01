package cache

import "time"

// Store is the interface for cache drivers.
type Store interface {
	Get(key string) (any, error)
	Put(key string, value any, ttl time.Duration) error
	Increment(key string, value int) (int, error)
	Decrement(key string, value int) (int, error)
	Forever(key string, value any) error
	Forget(key string) error
	Has(key string) bool
	Flush() error
}

// EventType represents the type of cache event.
type EventType string

const (
	EventHit     EventType = "hit"
	EventMiss    EventType = "miss"
	EventWrite   EventType = "write"
	EventDelete  EventType = "delete"
	EventFlush   EventType = "flush"
)

// CacheEvent represents a cache event.
type CacheEvent struct {
	Type EventType
	Key  string
}

// CacheListener is a function that handles cache events.
type CacheListener func(CacheEvent)

// EventDispatcher dispatches cache events.
type EventDispatcher struct {
	listeners []CacheListener
}

// NewEventDispatcher creates a new cache event dispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		listeners: make([]CacheListener, 0),
	}
}

// Listen registers a listener for cache events.
func (d *EventDispatcher) Listen(listener CacheListener) {
	d.listeners = append(d.listeners, listener)
}

// Dispatch sends a cache event to all listeners.
func (d *EventDispatcher) Dispatch(event CacheEvent) {
	for _, listener := range d.listeners {
		listener(event)
	}
}

