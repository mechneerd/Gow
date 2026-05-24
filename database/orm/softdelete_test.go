package orm

import (
	"database/sql"
	"github.com/mechneerd/gow/database/dialect"
	"github.com/mechneerd/gow/database/query"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

type SoftDeleteUser struct {
	ID        int           `db:"id" gow:"primaryKey"`
	Name      string        `db:"name"`
	DeletedAt *time.Time    `db:"deleted_at" gow:"softDelete"`
}

func (SoftDeleteUser) TableName() string { return "sd_users" }

func setupSDDB(t *testing.T) *DB {
	os.Remove("./test_sd.db")
	conn, err := sql.Open("sqlite", "./test_sd.db")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE sd_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			deleted_at DATETIME
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

func TestSoftDeletes(t *testing.T) {
	db := setupSDDB(t)
	defer os.Remove("./test_sd.db")

	user := &SoftDeleteUser{Name: "Alice"}
	NewQuery[SoftDeleteUser](db).Insert(user)

	user2 := &SoftDeleteUser{Name: "Bob"}
	NewQuery[SoftDeleteUser](db).Insert(user2)

	// 1. Soft delete Alice
	err := NewQuery[SoftDeleteUser](db).Delete(user)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Ensure DeletedAt was updated on the model struct
	if user.DeletedAt == nil {
		t.Error("Expected DeletedAt to be set on struct")
	}

	// 2. Querying should only return Bob
	users, _ := NewQuery[SoftDeleteUser](db).Get()
	if len(users) != 1 || users[0].Name != "Bob" {
		t.Errorf("Expected only Bob, got %d users", len(users))
	}

	// 3. WithTrashed should return both
	trashedUsers, _ := NewQuery[SoftDeleteUser](db).WithTrashed().Get()
	if len(trashedUsers) != 2 {
		t.Errorf("Expected 2 users with trashed, got %d", len(trashedUsers))
	}

	// 4. OnlyTrashed should return only Alice
	onlyTrashed, _ := NewQuery[SoftDeleteUser](db).OnlyTrashed().Get()
	if len(onlyTrashed) != 1 || onlyTrashed[0].Name != "Alice" {
		t.Errorf("Expected only Alice as trashed, got %d users", len(onlyTrashed))
	}

	// 5. Restore Alice
	err = NewQuery[SoftDeleteUser](db).Restore(user)
	if err != nil {
		t.Fatalf("Failed to restore user: %v", err)
	}
	
	if user.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil on struct after restore")
	}

	// Now both should be active
	all, _ := NewQuery[SoftDeleteUser](db).Get()
	if len(all) != 2 {
		t.Errorf("Expected 2 active users after restore, got %d", len(all))
	}

	// 6. ForceDelete Bob
	err = NewQuery[SoftDeleteUser](db).ForceDelete(user2)
	if err != nil {
		t.Fatalf("Failed to force delete Bob: %v", err)
	}

	// Bob should be completely gone
	allWithTrashed, _ := NewQuery[SoftDeleteUser](db).WithTrashed().Get()
	if len(allWithTrashed) != 1 || allWithTrashed[0].Name != "Alice" {
		t.Errorf("Expected only Alice after force delete, got %d users", len(allWithTrashed))
	}
}

