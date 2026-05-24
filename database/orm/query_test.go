package orm

import (
	"database/sql"
	"github.com/mechneerd/gow/database/dialect"
	"github.com/mechneerd/gow/database/query"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

type CrudUser struct {
	ID        int       `db:"id" gow:"primaryKey"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at" gow:"autoCreateTime"`
	UpdatedAt time.Time `db:"updated_at" gow:"autoUpdateTime"`
}

func (CrudUser) TableName() string { return "test_users" }

func setupTestDB(t *testing.T) *DB {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE test_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return &DB{
		Conn:    conn,
		Builder: query.NewBuilder(conn, &dialect.SQLiteDialect{}),
	}
}

func TestORMCrudLifecycle(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Insert and Find", func(t *testing.T) {
		user := &CrudUser{
			Name:  "Alice",
			Email: "alice@example.com",
		}

		err := NewQuery[CrudUser](db).Insert(user)
		if err != nil {
			t.Fatalf("Insert failed: %v", err)
		}

		if user.ID == 0 {
			t.Errorf("Expected ID to be populated")
		}
		if user.CreatedAt.IsZero() {
			t.Errorf("Expected CreatedAt to be populated")
		}

		// Find the inserted user
		found, err := NewQuery[CrudUser](db).Find(user.ID)
		if err != nil {
			t.Fatalf("Find failed: %v", err)
		}

		if found.Name != "Alice" {
			t.Errorf("Expected name 'Alice', got '%s'", found.Name)
		}
	})

	t.Run("Update", func(t *testing.T) {
		user := &CrudUser{
			Name:  "Bob",
			Email: "bob@example.com",
		}

		NewQuery[CrudUser](db).Insert(user)
		
		originalUpdateTime := user.UpdatedAt

		// Wait slightly to ensure time change
		time.Sleep(10 * time.Millisecond)

		user.Name = "Bobby"
		err := NewQuery[CrudUser](db).Update(user)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if user.UpdatedAt.Equal(originalUpdateTime) {
			t.Errorf("Expected UpdatedAt to change, got %v vs %v", user.UpdatedAt, originalUpdateTime)
		}

		found, _ := NewQuery[CrudUser](db).Find(user.ID)
		if found.Name != "Bobby" {
			t.Errorf("Expected name to be updated to 'Bobby', got '%s'", found.Name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		user := &CrudUser{
			Name:  "Charlie",
			Email: "charlie@example.com",
		}

		NewQuery[CrudUser](db).Insert(user)

		err := NewQuery[CrudUser](db).Delete(user)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = NewQuery[CrudUser](db).Find(user.ID)
		if err == nil {
			t.Errorf("Expected error when finding deleted user, got nil")
		}
	})
}

