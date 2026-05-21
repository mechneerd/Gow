package session

import (
	"crypto/rand"
	"encoding/hex"
)

// Manager manages the current session state.
type Manager struct {
	store Store
	id    string
	data  map[string]any
}

// NewManager creates a new Session Manager.
func NewManager(store Store, id string) *Manager {
	if id == "" {
		id = generateID()
	}
	return &Manager{
		store: store,
		id:    id,
		data:  make(map[string]any),
	}
}

// Start loads the session data from the store.
func (m *Manager) Start() error {
	data, err := m.store.Read(m.id)
	if err != nil {
		return err
	}
	m.data = data
	return nil
}

// Save writes the session data to the store.
func (m *Manager) Save() error {
	return m.store.Write(m.id, m.data)
}

// Get retrieves a value from the session.
func (m *Manager) Get(key string) any {
	return m.data[key]
}

// Put stores a value in the session.
func (m *Manager) Put(key string, value any) {
	m.data[key] = value
}

// ID returns the current session ID.
func (m *Manager) ID() string {
	return m.id
}

// Regenerate generates a new session ID and migrates the data.
func (m *Manager) Regenerate() error {
	err := m.store.Destroy(m.id)
	if err != nil {
		return err
	}
	m.id = generateID()
	return m.Save()
}

func generateID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
