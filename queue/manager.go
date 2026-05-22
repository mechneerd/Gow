package queue

import (
	"fmt"
)

// Driver interface abstracts queue interactions.
type Driver interface {
	Push(job Job) error
	Pop() (Job, error)
}

// Manager resolves queue drivers.
type Manager struct {
	drivers map[string]Driver
	defaultQueue string
}

// NewManager initializes a new Queue Manager.
func NewManager(defaultQueue string) *Manager {
	m := &Manager{
		drivers:      make(map[string]Driver),
		defaultQueue: defaultQueue,
	}
	// Register default internal drivers
	m.AddDriver("memory", NewMemoryDriver(10000))
	return m
}

// AddDriver registers a new queue driver.
func (m *Manager) AddDriver(name string, driver Driver) {
	m.drivers[name] = driver
}

// Connection gets a driver by name.
func (m *Manager) Connection(name string) Driver {
	if name == "" {
		name = m.defaultQueue
	}
	return m.drivers[name]
}

// Push pushes a job to the default connection.
func (m *Manager) Push(job Job) error {
	driver := m.Connection("")
	if driver == nil {
		return fmt.Errorf("queue driver [%s] not found", m.defaultQueue)
	}
	return driver.Push(job)
}
