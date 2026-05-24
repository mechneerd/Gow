package session

// Store represents a session storage driver.
type Store interface {
	// Read retrieves the session data by ID.
	Read(id string) (map[string]any, error)
	// Write saves the session data by ID.
	Write(id string, data map[string]any) error
	// Destroy removes the session data by ID.
	Destroy(id string) error
}

