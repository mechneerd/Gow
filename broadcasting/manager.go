package broadcasting

import (
	"encoding/json"
	"fmt"
	"log"
)

// Broadcaster handles the actual transmission of the event.
type Broadcaster interface {
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
	drivers       map[string]Broadcaster
	defaultDriver string
}

// NewManager creates a new Broadcaster Manager.
func NewManager(defaultDriver string) *Manager {
	return &Manager{
		drivers:       make(map[string]Broadcaster),
		defaultDriver: defaultDriver,
	}
}

// Extend registers a custom driver (e.g., Pusher, Redis, WebSocket).
func (m *Manager) Extend(name string, driver Broadcaster) {
	m.drivers[name] = driver
}

// Connection gets a driver by name.
func (m *Manager) Connection(name string) Broadcaster {
	if name == "" {
		name = m.defaultDriver
	}
	return m.drivers[name]
}

// Broadcast dispatches the event via the default connection.
func (m *Manager) Broadcast(channels []string, eventName string, payload map[string]any) error {
	driver := m.Connection("")
	if driver == nil {
		return fmt.Errorf("broadcast driver [%s] not found", m.defaultDriver)
	}

	if len(channels) == 0 {
		return nil
	}

	if eventName == "" {
		eventName = "event"
	}

	return driver.Broadcast(channels, eventName, payload)
}
