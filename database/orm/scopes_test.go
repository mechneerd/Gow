package orm

import (
	"database/sql"
	"gow/database/dialect"
	"gow/database/query"
	"testing"
)

type ScopeUser struct {
	ID     int    `db:"id" gow:"primaryKey"`
	Name   string `db:"name"`
	Active bool   `db:"active"`
}

func (ScopeUser) TableName() string { return "scope_users" }

func TestGlobalScopes(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Exec(`
		CREATE TABLE scope_users (id INTEGER PRIMARY KEY, name TEXT, active BOOLEAN);
		INSERT INTO scope_users (id, name, active) VALUES (1, 'Alice', 1), (2, 'Bob', 0), (3, 'Charlie', 1);
	`)
	if err != nil {
		t.Fatal(err)
	}

	db := &DB{
		Conn:    conn,
		Builder: query.NewBuilder(conn, &dialect.SQLiteDialect{}),
	}

	// Add global scope to only get active users
	AddGlobalScope("ScopeUser", ScopeFunc(func(b *query.Builder) *query.Builder {
		return b.Where("active", "=", 1)
	}))

	users, err := NewQuery[ScopeUser](db).Get()
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("Expected 2 active users, got %d", len(users))
	}

	for _, u := range users {
		if !u.Active {
			t.Errorf("Expected active user, got inactive: %s", u.Name)
		}
	}
}
