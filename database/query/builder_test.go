package query

import (
	"context"
	"database/sql"
	"gow/database/dialect"
	"strings"
	"testing"
)

// MockQueryExecer implements QueryExecer for testing SQL generation without a real DB.
type MockQueryExecer struct {
	LastQuery string
	LastArgs  []any
}

func (m *MockQueryExecer) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	m.LastQuery = query
	m.LastArgs = args
	return nil, nil // Return nil for rows since we just want to test SQL generation
}

func (m *MockQueryExecer) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	m.LastQuery = query
	m.LastArgs = args
	// Note: since it returns *sql.Row, we can't easily mock the return without a real DB connection.
	// But we can test ToSQL for generation.
	return nil
}

func (m *MockQueryExecer) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	m.LastQuery = query
	m.LastArgs = args
	return nil, nil
}

func TestBuilderSelectSQL(t *testing.T) {
	mockDB := &MockQueryExecer{}
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(mockDB, d)

	b.Table("users").
		Select("id", "name").
		Where("age", ">", 18).
		OrWhere("status", "=", "active").
		WhereIn("role", []any{"admin", "editor"}).
		WhereNull("deleted_at").
		WhereBetween("score", []any{50, 100}).
		OrderBy("created_at", "DESC").
		Limit(10).
		Offset(5)

	sqlQuery, args := b.ToSQL()
	
	expectedSQL := "SELECT \"id\", \"name\" FROM \"users\" WHERE \"age\" > ? OR \"status\" = ? AND \"role\" IN (?, ?) AND \"deleted_at\" IS NULL AND \"score\" BETWEEN ? AND ? ORDER BY \"created_at\" DESC LIMIT 10 OFFSET 5"
	if sqlQuery != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sqlQuery)
	}

	if len(args) != 6 {
		t.Errorf("Expected 6 args, got %d", len(args))
	}
}

func TestBuilderJoins(t *testing.T) {
	mockDB := &MockQueryExecer{}
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(mockDB, d)

	b.Table("users").
		Select("users.id", "posts.title").
		Join("posts", "users.id", "=", "posts.user_id").
		LeftJoin("profiles", "users.id", "=", "profiles.user_id")

	sqlQuery, _ := b.ToSQL()

	expectedSQL := "SELECT \"users.id\", \"posts.title\" FROM \"users\" INNER JOIN \"posts\" ON \"users.id\" = \"posts.user_id\" LEFT JOIN \"profiles\" ON \"users.id\" = \"profiles.user_id\""
	if sqlQuery != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sqlQuery)
	}
}

func TestBuilderDMLGeneration(t *testing.T) {
	mockDB := &MockQueryExecer{}
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(mockDB, d).Table("users")

	// Test Insert
	b.Insert(map[string]any{"name": "John", "age": 30})
	
	if !strings.HasPrefix(mockDB.LastQuery, "INSERT INTO \"users\"") {
		t.Errorf("Expected INSERT query, got: %s", mockDB.LastQuery)
	}
	if !strings.Contains(mockDB.LastQuery, "name") || !strings.Contains(mockDB.LastQuery, "age") {
		t.Errorf("Expected INSERT query to contain columns, got: %s", mockDB.LastQuery)
	}

	// Test Update
	b = NewBuilder(mockDB, d).Table("users").Where("id", "=", 1)
	b.Update(map[string]any{"status": "banned"})

	if !strings.HasPrefix(mockDB.LastQuery, "UPDATE \"users\" SET") {
		t.Errorf("Expected UPDATE query, got: %s", mockDB.LastQuery)
	}
	if !strings.Contains(mockDB.LastQuery, "\"status\" = ?") {
		t.Errorf("Expected UPDATE query to contain status, got: %s", mockDB.LastQuery)
	}
	if !strings.Contains(mockDB.LastQuery, "WHERE \"id\" = ?") {
		t.Errorf("Expected UPDATE query to contain WHERE, got: %s", mockDB.LastQuery)
	}

	// Test Delete
	b = NewBuilder(mockDB, d).Table("users").Where("id", "=", 1)
	b.Delete()

	if mockDB.LastQuery != "DELETE FROM \"users\" WHERE \"id\" = ?" {
		t.Errorf("Expected DELETE query, got: %s", mockDB.LastQuery)
	}
}
