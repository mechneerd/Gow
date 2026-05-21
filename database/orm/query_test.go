package orm

import (
	"database/sql"
	"gow/database/dialect"
	"gow/database/query"
	"gow/database/schema"
	"testing"
	// _ "github.com/mattn/go-sqlite3" // Removed to avoid CGO/Network issues
)

type User struct {
	ID        int       `db:"id" gow:"primaryKey,autoIncrement"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at" gow:"autoCreateTime"`
}

func TestORMAndSchema(t *testing.T) {
	t.Skip("Skipping execution test to avoid sqlite3 dependency.")

	d := &dialect.SQLiteDialect{}
	
	// 1. Test Schema Builder
	builder := schema.NewBuilder(conn, d)
	err = builder.Create("users", func(table *schema.Blueprint) {
		table.ID()
		table.String("name", 255)
		table.String("email", 255).Unique()
		table.Timestamps()
	})
	if err != nil {
		t.Fatalf("Schema Create failed: %v", err)
	}

	// 2. Test ORM Insert
	db := &DB{
		Conn:    conn,
		Builder: query.NewBuilder(conn, d),
	}

	q := NewQuery[User](db)
	user := &User{
		Name:  "Alice",
		Email: "alice@example.com",
	}

	err = q.Insert(user)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("Expected ID to be populated")
	}

	// 3. Test ORM Get/First
	q2 := NewQuery[User](db)
	fetched, err := q2.Where("email", "=", "alice@example.com").First()
	if err != nil {
		t.Fatalf("First() failed: %v", err)
	}

	if fetched.Name != "Alice" {
		t.Errorf("Expected Name 'Alice', got '%s'", fetched.Name)
	}
	if fetched.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, fetched.ID)
	}
}
