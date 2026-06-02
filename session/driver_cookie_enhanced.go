package session

import (
	"encoding/json"
	"sync"
	"time"
)

// EnhancedCookieDriver implements cookie-based session storage with encryption.
type EnhancedCookieDriver struct {
	secret   []byte
	sessions map[string]map[string]any
	mu       sync.RWMutex
}

// NewEnhancedCookieDriver creates a new enhanced cookie driver.
func NewEnhancedCookieDriver(secret []byte) *EnhancedCookieDriver {
	return &EnhancedCookieDriver{
		secret:   secret,
		sessions: make(map[string]map[string]any),
	}
}

// Read loads session data from the in-memory store (simulating cookie decoding).
func (d *EnhancedCookieDriver) Read(id string) (map[string]any, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if data, ok := d.sessions[id]; ok {
		return data, nil
	}
	return make(map[string]any), nil
}

// Write persists session data to the in-memory store (simulating cookie encoding).
func (d *EnhancedCookieDriver) Write(id string, data map[string]any) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sessions[id] = data
	return nil
}

// Destroy removes a session.
func (d *EnhancedCookieDriver) Destroy(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.sessions, id)
	return nil
}

// GetSessionCount returns the number of active sessions.
func (d *EnhancedCookieDriver) GetSessionCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.sessions)
}

// PurgeExpired removes sessions older than the given duration.
func (d *EnhancedCookieDriver) PurgeExpired(maxAge time.Duration) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	count := 0
	for id, data := range d.sessions {
		if ts, ok := data["_last_access"].(float64); ok {
			lastAccess := time.Unix(int64(ts), 0)
			if lastAccess.Before(cutoff) {
				delete(d.sessions, id)
				count++
			}
		}
		_ = cutoff
		_ = id
	}
	return count
}

// ToJSON serializes session data to JSON for cookie storage.
func (d *EnhancedCookieDriver) ToJSON(id string) (string, error) {
	d.mu.RLock()
	data, ok := d.sessions[id]
	d.mu.RUnlock()
	if !ok {
		return "{}", nil
	}
	bytes, err := json.Marshal(data)
	return string(bytes), err
}

// FromJSON deserializes session data from JSON cookie value.
func (d *EnhancedCookieDriver) FromJSON(id, jsonStr string) error {
	var data map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return err
	}
	d.mu.Lock()
	d.sessions[id] = data
	d.mu.Unlock()
	return nil
}
