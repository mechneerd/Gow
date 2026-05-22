package database

import (
	"errors"
	"gow/config"
	"sync"
)

// Manager handles multiple database connections.
type Manager struct {
	mu          sync.RWMutex
	connections map[string]*Connection
	defaultConn string
	config      *config.Repository
}

// NewManager creates a new database connection manager.
func NewManager(cfg *config.Repository) *Manager {
	return &Manager{
		connections: make(map[string]*Connection),
		config:      cfg,
		defaultConn: cfg.Get("DB_CONNECTION", "sqlite"),
	}
}

// AddConnection adds a pre-configured connection to the manager.
// It automatically applies connection pool settings from the config repository
// using PoolConfigFromEnv and ConfigurePool.
func (m *Manager) AddConnection(name string, conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[name] = conn

	// Apply pool configuration automatically
	pc := PoolConfigFromEnv(m.config)
	ConfigurePool(conn.DB, pc)
}

// Connection gets a database connection by name.
func (m *Manager) Connection(name string) (*Connection, error) {
	if name == "" {
		name = m.defaultConn
	}

	m.mu.RLock()
	conn, exists := m.connections[name]
	m.mu.RUnlock()

	if !exists {
		// Attempt to resolve configuration and build the connection
		// For Phase 2, we will lazy-load connections if possible
		return nil, errors.New("database connection not found: " + name)
	}

	return conn, nil
}

// DB gets the default database connection.
func (m *Manager) DB() (*Connection, error) {
	return m.Connection("")
}

// Close closes all managed database connections.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for _, conn := range m.connections {
		if err := conn.DB.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.New("failed to close one or more database connections")
	}
	return nil
}
