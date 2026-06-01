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

// Start loads the session data from the store and ages flash data.
func (m *Manager) Start() error {
	data, err := m.store.Read(m.id)
	if err != nil {
		return err
	}
	m.data = data
	if m.data == nil {
		m.data = make(map[string]any)
	}
	
	m.ageFlashData()
	return nil
}

// ageFlashData cycles the flash data arrays (new -> old).
func (m *Manager) ageFlashData() {
	flashRaw, ok := m.data["_flash"]
	if !ok {
		m.data["_flash"] = map[string][]string{"old": {}, "new": {}}
		return
	}
	
	flash, ok := flashRaw.(map[string][]string)
	if !ok { // handle case where deserialization gave map[string]interface{}
		flash = map[string][]string{"old": {}, "new": {}}
		
		if typedMap, ok := flashRaw.(map[string]any); ok {
			if oldArr, ok := typedMap["old"].([]any); ok {
				var oldStr []string
				for _, v := range oldArr {
					if str, ok := v.(string); ok {
						oldStr = append(oldStr, str)
					}
				}
				flash["old"] = oldStr
			}
			if newArr, ok := typedMap["new"].([]any); ok {
				var newStr []string
				for _, v := range newArr {
					if str, ok := v.(string); ok {
						newStr = append(newStr, str)
					}
				}
				flash["new"] = newStr
			}
		}
	} else {
		// New becomes old, new becomes empty
		flash["old"] = flash["new"]
		flash["new"] = []string{}
	}
	
	m.data["_flash"] = flash
}

func (m *Manager) clearOldFlashData() {
	flashRaw, ok := m.data["_flash"]
	if !ok {
		return
	}
	
	if flash, ok := flashRaw.(map[string][]string); ok {
		newKeys := make(map[string]bool)
		for _, k := range flash["new"] {
			newKeys[k] = true
		}
		
		for _, key := range flash["old"] {
			if !newKeys[key] {
				delete(m.data, key)
			}
		}
		// old array is cleared before saving to avoid building up
		flash["old"] = []string{}
		m.data["_flash"] = flash
	}
}

// Save writes the session data to the store.
func (m *Manager) Save() error {
	m.clearOldFlashData()
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

// Flash stores a value in the session that will be available for the next request.
func (m *Manager) Flash(key string, value any) {
	m.Put(key, value)
	
	if flash, ok := m.data["_flash"].(map[string][]string); ok {
		flash["new"] = append(flash["new"], key)
		m.data["_flash"] = flash
	}
}

// Reflash keeps all current flash data for an additional request.
func (m *Manager) Reflash() {
	if flash, ok := m.data["_flash"].(map[string][]string); ok {
		flash["new"] = append(flash["new"], flash["old"]...)
		flash["old"] = []string{}
		m.data["_flash"] = flash
	}
}

// Keep keeps specific flash data for an additional request.
func (m *Manager) Keep(keys ...string) {
	if flash, ok := m.data["_flash"].(map[string][]string); ok {
		flash["new"] = append(flash["new"], keys...)
		m.data["_flash"] = flash
	}
}

// FlashInput flashes the current request input to the session.
func (m *Manager) FlashInput(input map[string]any) {
	m.Flash("_old_input", input)
}

// Old returns an old input value flashed from the previous request.
func (m *Manager) Old(key string, defaultValue any) any {
	oldInput, ok := m.Get("_old_input").(map[string]any)
	if !ok {
		return defaultValue
	}
	if val, exists := oldInput[key]; exists {
		return val
	}
	return defaultValue
}

// PreviousURL stores the URL the user came from (typically set by middleware).
func (m *Manager) PreviousURL() string {
	val := m.Get("_previous_url")
	if val == nil {
		return ""
	}
	if url, ok := val.(string); ok {
		return url
	}
	return ""
}

// SetPreviousURL stores the previous URL in the session.
func (m *Manager) SetPreviousURL(url string) {
	m.Put("_previous_url", url)
}

// Intended stores the intended URL the user was trying to access.
// After login, call Intended() to get the stored URL and redirect there.
func (m *Manager) Intended(defaultURL ...string) string {
	val := m.Get("_intended_url")
	if val == nil {
		if len(defaultURL) > 0 {
			return defaultURL[0]
		}
		return ""
	}
	if url, ok := val.(string); ok {
		return url
	}
	if len(defaultURL) > 0 {
		return defaultURL[0]
	}
	return ""
}

// SetIntendedURL stores the intended URL for post-login redirect.
func (m *Manager) SetIntendedURL(url string) {
	m.Put("_intended_url", url)
}

// Flush clears all session data.
func (m *Manager) Flush() {
	m.data = make(map[string]any)
}

// Pull removes a value from the session and returns it.
func (m *Manager) Pull(key string, defaultValue any) any {
	val := m.Get(key)
	if val == nil {
		return defaultValue
	}
	m.Forget(key)
	return val
}

// Forget removes a value from the session.
func (m *Manager) Forget(key string) {
	delete(m.data, key)
}

// Has checks if a key exists in the session.
func (m *Manager) Has(key string) bool {
	_, exists := m.data[key]
	return exists
}

// All returns all session data.
func (m *Manager) All() map[string]any {
	result := make(map[string]any, len(m.data))
	for k, v := range m.data {
		result[k] = v
	}
	return result
}

func generateID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate session ID: " + err.Error())
	}
	return hex.EncodeToString(b)
}

