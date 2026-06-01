package session

import "sync"

// ArrayDriver is an in-memory session driver for testing.
type ArrayDriver struct {
	mu       sync.RWMutex
	sessions map[string]map[string]any
}

// NewArrayDriver creates a new in-memory session driver.
func NewArrayDriver() *ArrayDriver {
	return &ArrayDriver{
		sessions: make(map[string]map[string]any),
	}
}

func (d *ArrayDriver) Read(id string) (map[string]any, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	data, exists := d.sessions[id]
	if !exists {
		return nil, nil
	}

	// Return a copy
	result := make(map[string]any, len(data))
	for k, v := range data {
		result[k] = v
	}
	return result, nil
}

func (d *ArrayDriver) Write(id string, data map[string]any) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Store a copy
	d.sessions[id] = make(map[string]any, len(data))
	for k, v := range data {
		d.sessions[id][k] = v
	}
	return nil
}

func (d *ArrayDriver) Destroy(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.sessions, id)
	return nil
}

// Flush removes all sessions (useful in testing).
func (d *ArrayDriver) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sessions = make(map[string]map[string]any)
}

// Count returns the number of active sessions.
func (d *ArrayDriver) Count() int {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return len(d.sessions)
}
