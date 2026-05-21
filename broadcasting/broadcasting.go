package broadcasting

// Driver interface for broadcasting events to channels.
type Driver interface {
	Broadcast(channels []string, event string, payload any) error
}

// Manager orchestrates broadcasting.
type Manager struct {
	drivers       map[string]Driver
	defaultDriver string
}

// NewManager creates a new Broadcasting manager.
func NewManager(defaultDriver string) *Manager {
	return &Manager{
		drivers:       make(map[string]Driver),
		defaultDriver: defaultDriver,
	}
}

// Extend registers a custom driver implementation.
func (m *Manager) Extend(name string, driver Driver) {
	m.drivers[name] = driver
}

// Driver gets the default or named driver.
func (m *Manager) Driver(name ...string) Driver {
	driverName := m.defaultDriver
	if len(name) > 0 {
		driverName = name[0]
	}
	return m.drivers[driverName]
}

// Broadcast dispatches the event using the default driver.
func (m *Manager) Broadcast(channels []string, event string, payload any) error {
	return m.Driver().Broadcast(channels, event, payload)
}

// PusherDriver implements the Driver interface for Pusher API.
type PusherDriver struct {
	// In a real implementation this would hold a *pusher.Client
}

func (d *PusherDriver) Broadcast(channels []string, event string, payload any) error {
	// e.g. return d.client.TriggerMulti(channels, event, payload)
	return nil
}

// RedisDriver implements the Driver interface using Redis Pub/Sub.
type RedisDriver struct {
	// In a real implementation this would hold a *redis.Client
}

func (d *RedisDriver) Broadcast(channels []string, event string, payload any) error {
	// We serialize payload to JSON and publish to a Redis channel
	// which is then consumed by Reverb or Soketi.
	return nil
}
