package events

import (
	"fmt"
	"sync"
)

// Fake is a fake event dispatcher for testing.
type Fake struct {
	mu       sync.RWMutex
	events   []Event
	listeners map[string][]Listener
}

// NewFake creates a new fake event dispatcher.
func NewFake() *Fake {
	return &Fake{
		events:    make([]Event, 0),
		listeners: make(map[string][]Listener),
	}
}

// Dispatch captures the event instead of dispatching it.
func (f *Fake) Dispatch(event Event) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, event)
}

// Listen registers a listener on the fake for assertion purposes.
func (f *Fake) Listen(eventName string, listener Listener) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.listeners[eventName] = append(f.listeners[eventName], listener)
}

// GetEvents returns all captured events.
func (f *Fake) GetEvents() []Event {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make([]Event, len(f.events))
	copy(result, f.events)
	return result
}

// GetLastEvent returns the last captured event, or nil if none.
func (f *Fake) GetLastEvent() Event {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if len(f.events) == 0 {
		return nil
	}
	return f.events[len(f.events)-1]
}

// GetEventCount returns the number of events dispatched.
func (f *Fake) GetEventCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.events)
}

// HasDispatched checks if any event was dispatched.
func (f *Fake) HasDispatched() bool {
	return f.GetEventCount() > 0
}

// HasDispatchedEvent checks if an event of the given type was dispatched.
func (f *Fake) HasDispatchedEvent(event any) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, e := range f.events {
		if e == event {
			return true
		}
		// Check by type name
		if e != nil && event != nil {
			if tn := eventTypeName(e); tn == eventTypeName(event) {
				return true
			}
		}
	}
	return false
}

// AssertDispatched asserts that at least one event was dispatched.
func (f *Fake) AssertDispatched() bool {
	return f.HasDispatched()
}

// AssertNotDispatched asserts that no events were dispatched.
func (f *Fake) AssertNotDispatched() bool {
	return f.GetEventCount() == 0
}

// AssertDispatchedCount asserts the exact number of events dispatched.
func (f *Fake) AssertDispatchedCount(count int) bool {
	return f.GetEventCount() == count
}

// AssertDispatchedEvent asserts that a specific event type was dispatched.
func (f *Fake) AssertDispatchedEvent(event any) bool {
	return f.HasDispatchedEvent(event)
}

// AssertNotDispatchedEvent asserts that a specific event type was NOT dispatched.
func (f *Fake) AssertNotDispatchedEvent(event any) bool {
	return !f.HasDispatchedEvent(event)
}

// Clear resets all captured events.
func (f *Fake) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = make([]Event, 0)
}

func eventTypeName(v any) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%T", v)
}
