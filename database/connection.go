package database

import (
	"database/sql"
	"gow/database/dialect"
)

// Connection represents a database connection along with its dialect.
type Connection struct {
	DB      *sql.DB
	Dialect dialect.Dialect
}

// NewConnection creates a new database connection wrapper.
func NewConnection(db *sql.DB, d dialect.Dialect) *Connection {
	return &Connection{
		DB:      db,
		Dialect: d,
	}
}
