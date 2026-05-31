package session

import (
	"database/sql"
	"encoding/json"
	"time"
)

// DatabaseDriver implements session storage using a database table.
type DatabaseDriver struct {
	db *sql.DB
}

// NewDatabaseDriver creates a new database session driver.
func NewDatabaseDriver(db *sql.DB) *DatabaseDriver {
	return &DatabaseDriver{db: db}
}

// Read reads session data by ID.
func (d *DatabaseDriver) Read(id string) (map[string]any, error) {
	var data string
	err := d.db.QueryRow("SELECT data FROM sessions WHERE id = ?", id).Scan(&data)
	if err != nil {
		return make(map[string]any), nil
	}

	// Deserialize JSON data
	var sessionData map[string]any
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return make(map[string]any), nil
	}
	return sessionData, nil
}

// Write writes session data to the database.
func (d *DatabaseDriver) Write(id string, data map[string]any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	expiration := time.Now().Add(24 * time.Hour)
	_, err = d.db.Exec(
		"INSERT OR REPLACE INTO sessions (id, data, expiration) VALUES (?, ?, ?)",
		id, string(jsonData), expiration,
	)
	return err
}

// Destroy removes a session from the database.
func (d *DatabaseDriver) Destroy(id string) error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

// GarbageCollect removes expired sessions.
func (d *DatabaseDriver) GarbageCollect(maxLifetime time.Duration) error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE expiration < ?", time.Now())
	return err
}
