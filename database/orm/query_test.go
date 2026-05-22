package orm

import (
	"context"
	"database/sql"
	"gow/database/dialect"
	"gow/database/query"
	"gow/database/schema"
	"testing"
	"time"
	"errors"

	_ "modernc.org/sqlite"
)

type User struct {
	ID        int       `db:"id" gow:"primaryKey,autoIncrement"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at" gow:"autoCreateTime"`
	UpdatedAt time.Time `db:"updated_at" gow:"autoUpdateTime"`
}

func TestORMAndSchema(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open sqlite memory db: %v", err)
	}
	defer conn.Close()

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

func TestORMTransactions(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open sqlite memory db: %v", err)
	}
	defer conn.Close()

	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(conn, d)
	_ = builder.Create("users", func(table *schema.Blueprint) {
		table.ID()
		table.String("name", 255)
		table.String("email", 255).Unique()
		table.Timestamps()
	})

	db := &DB{
		Conn:    conn,
		Builder: query.NewBuilder(conn, d),
	}

	// Test successful transaction
	err = db.Transaction(context.Background(), func(txDB *DB) error {
		q := NewQuery[User](txDB)
		return q.Insert(&User{Name: "TxUser", Email: "tx@example.com"})
	})
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}

	// Verify insertion
	fetched, _ := NewQuery[User](db).Where("email", "=", "tx@example.com").First()
	if fetched == nil {
		t.Error("Expected user to be created in transaction")
	}

	// Test rollback on error
	err = db.Transaction(context.Background(), func(txDB *DB) error {
		q := NewQuery[User](txDB)
		_ = q.Insert(&User{Name: "FailUser", Email: "fail@example.com"})
		return errors.New("rollback requested")
	})

	if err == nil || err.Error() != "rollback requested" {
		t.Errorf("Expected rollback error, got: %v", err)
	}

	// Verify rollback
	failedUser, _ := NewQuery[User](db).Where("email", "=", "fail@example.com").First()
	if failedUser != nil {
		t.Error("Expected user NOT to be created due to rollback")
	}
}
