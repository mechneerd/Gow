package migration

import (
	"database/sql"
	"testing"

	"github.com/mechneerd/gow/database/dialect"
	"github.com/mechneerd/gow/database/schema"

	_ "modernc.org/sqlite"
)

type createUsersTable struct{}

func (m *createUsersTable) Up(b *schema.Builder) error {
	return b.Create("users", func(table *schema.Blueprint) {
		table.ID()
		table.String("name", 255)
		table.String("email", 255)
	})
}

func (m *createUsersTable) Down(b *schema.Builder) error {
	return b.DropIfExists("users")
}

type createPostsTable struct{}

func (m *createPostsTable) Up(b *schema.Builder) error {
	return b.Create("posts", func(table *schema.Blueprint) {
		table.ID()
		table.String("title", 255)
		table.Integer("user_id")
	})
}

func (m *createPostsTable) Down(b *schema.Builder) error {
	return b.DropIfExists("posts")
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	return db
}

func newTestRegistry() *Registry {
	reg := NewRegistry()
	reg.Register("001_create_users", &createUsersTable{})
	reg.Register("002_create_posts", &createPostsTable{})
	return reg
}

func TestMigrateAndRollback(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	reg := newTestRegistry()
	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(db, d)
	m := NewMigrator(db, builder, reg)

	if err := m.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 'users' table to exist")
	}

	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='posts'").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 'posts' table to exist")
	}

	var ranCount int
	db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&ranCount)
	if ranCount != 2 {
		t.Errorf("Expected 2 migrations logged, got %d", ranCount)
	}

	if err := m.Rollback(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='posts'").Scan(&count)
	if count != 0 {
		t.Errorf("Expected 'posts' table to be dropped after rollback")
	}

	db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&ranCount)
	if ranCount != 0 {
		t.Errorf("Expected 0 migrations after rollback, got %d", ranCount)
	}
}

func TestMigrateIdempotent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	reg := newTestRegistry()
	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(db, d)
	m := NewMigrator(db, builder, reg)

	if err := m.Migrate(); err != nil {
		t.Fatalf("First Migrate failed: %v", err)
	}
	if err := m.Migrate(); err != nil {
		t.Fatalf("Second Migrate failed: %v", err)
	}

	var ranCount int
	db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&ranCount)
	if ranCount != 2 {
		t.Errorf("Expected 2 migrations after running twice, got %d", ranCount)
	}
}

func TestFresh(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	reg := newTestRegistry()
	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(db, d)
	m := NewMigrator(db, builder, reg)

	// Fresh from empty state should work
	if err := m.Fresh(); err != nil {
		t.Fatalf("Fresh failed: %v", err)
	}

	var ranCount int
	db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&ranCount)
	if ranCount != 2 {
		t.Errorf("Expected 2 migrations after Fresh, got %d", ranCount)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 'users' table to exist after Fresh")
	}

	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='posts'").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 'posts' table to exist after Fresh")
	}
}

func TestRollbackSteps(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	reg := newTestRegistry()
	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(db, d)
	m := NewMigrator(db, builder, reg)

	m.Migrate()

	if err := m.RollbackSteps(1); err != nil {
		t.Fatalf("RollbackSteps failed: %v", err)
	}

	var ranCount int
	db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&ranCount)
	if ranCount != 1 {
		t.Errorf("Expected 1 migration after RollbackSteps(1), got %d", ranCount)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='posts'").Scan(&count)
	if count != 0 {
		t.Errorf("Expected 'posts' table to be dropped")
	}
}

func TestRollbackMigration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	reg := newTestRegistry()
	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(db, d)
	m := NewMigrator(db, builder, reg)

	m.Migrate()

	if err := m.RollbackMigration("002_create_posts"); err != nil {
		t.Fatalf("RollbackMigration failed: %v", err)
	}

	var ranCount int
	db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&ranCount)
	if ranCount != 1 {
		t.Errorf("Expected 1 migration, got %d", ranCount)
	}

	var name string
	db.QueryRow("SELECT migration FROM migrations LIMIT 1").Scan(&name)
	if name != "001_create_users" {
		t.Errorf("Expected remaining migration to be '001_create_users', got '%s'", name)
	}
}

func TestMigrateOne(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	reg := newTestRegistry()
	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(db, d)
	m := NewMigrator(db, builder, reg)

	if err := m.MigrateOne("001_create_users"); err != nil {
		t.Fatalf("MigrateOne failed: %v", err)
	}

	var ranCount int
	db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&ranCount)
	if ranCount != 1 {
		t.Errorf("Expected 1 migration, got %d", ranCount)
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if count != 1 {
		t.Errorf("Expected 'users' table to exist")
	}

	db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='posts'").Scan(&count)
	if count != 0 {
		t.Errorf("Expected 'posts' table to NOT exist")
	}
}

func TestStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	reg := newTestRegistry()
	d := &dialect.SQLiteDialect{}
	builder := schema.NewBuilder(db, d)
	m := NewMigrator(db, builder, reg)

	if err := m.Status(); err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	m.MigrateOne("001_create_users")
	if err := m.Status(); err != nil {
		t.Fatalf("Status after partial migration failed: %v", err)
	}
}

func TestRegistryCustom(t *testing.T) {
	reg := NewRegistry()

	if len(reg.migrations) != 0 {
		t.Errorf("Expected empty registry")
	}

	reg.Register("test", &createUsersTable{})
	if len(reg.migrations) != 1 {
		t.Errorf("Expected 1 migration in registry")
	}
}
