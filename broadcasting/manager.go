package broadcasting

import (
	"encoding/json"
	"fmt"
	"log"
)

// Driver handles the actual transmission of the event.
type Driver interface {
	Broadcast(channels []string, eventName string, payload map[string]any) error
}

// LogDriver outputs the broadcast to the standard logger for local development.
type LogDriver struct{}

func (d *LogDriver) Broadcast(channels []string, eventName string, payload map[string]any) error {
	data, _ := json.Marshal(payload)
	log.Printf("Broadcasting [%s] on channels %v with payload: %s\n", eventName, channels, string(data))
	return nil
}

// Manager resolves broadcast drivers and dispatches events.
type Manager struct {
	drivers       map[string]Driver
	defaultDriver string
}

// NewManager creates a new Broadcaster Manager.
func NewManager(defaultDriver string) *Manager {
	return &Manager{
		drivers:       make(map[string]Driver),
		defaultDriver: defaultDriver,
	}
}

// Extend registers a custom driver (e.g., Pusher, Redis).
func (m *Manager) Extend(name string, driver Driver) {
	m.drivers[name] = driver
}

// Connection gets a driver by name.
func (m *Manager) Connection(name string) Driver {
	if name == "" {
		name = m.defaultDriver
	}
	return m.drivers[name]
}

// Broadcast dispatches the event via the default connection.
func (m *Manager) Broadcast(event Event) error {
	driver := m.Connection("")
	if driver == nil {
		return fmt.Errorf("broadcast driver [%s] not found", m.defaultDriver)
	}

	channelObjects := event.BroadcastOn()
	if len(channelObjects) == 0 {
		return nil
	}

	channels := make([]string, len(channelObjects))
	for i, c := range channelObjects {
		channels[i] = c.Name()
	}

	eventName := event.BroadcastAs()
	if eventName == "" {
		// Default to a generic name or reflection-based name
		eventName = "event"
	}

	payload := event.BroadcastWith()

	return driver.Broadcast(channels, eventName, payload)
}
