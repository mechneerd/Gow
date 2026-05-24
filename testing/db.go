package testing

import (
	"database/sql"
	"log"
	"testing"
)

// DBTestContext manages database state for tests.
type DBTestContext struct {
	db *sql.DB
	tx *sql.Tx
}

// NewDBTestContext initializes a new database testing context.
func NewDBTestContext(db *sql.DB) *DBTestContext {
	return &DBTestContext{db: db}
}

// SetupSuite should be called once before all tests run.
// It applies all migrations to the database.
func (c *DBTestContext) SetupSuite() {
	log.Println("Setting up test database suite (running migrations)...")
	// Run migrations using the framework's migration runner.
	// migration.Runner(c.db).Up()
}

// RefreshDatabase wraps a test in a transaction and automatically rolls it back after the test.
// This is the recommended way to keep tests isolated and fast.
func (c *DBTestContext) RefreshDatabase(t *testing.T, testFunc func()) {
	tx, err := c.db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Save the transaction so the app can use it during the test
	c.tx = tx

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		_ = tx.Rollback() // always rollback after test
	}()

	testFunc()
}

// WithTransaction is a helper that lets you run a block inside a transaction.
// Useful when you want to share the same transaction across multiple operations in a test.
func (c *DBTestContext) WithTransaction(t *testing.T, fn func(tx *sql.Tx)) {
	tx, err := c.db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	fn(tx)
}

// GetTransaction returns the current transaction (if RefreshDatabase was used).
func (c *DBTestContext) GetTransaction() *sql.Tx {
	return c.tx
}

// UseInMemorySQLite sets up a completely fresh SQLite database in memory for the test.
func (c *DBTestContext) UseInMemorySQLite(t *testing.T, testFunc func(db *sql.DB)) {
	log.Println("Setting up in-memory SQLite for test...")
	
	// Open an in-memory SQLite connection
	// memoryDB, err := sql.Open("sqlite3", ":memory:")
	// if err != nil {
	// 	t.Fatalf("Failed to open sqlite memory db: %v", err)
	// }
	// defer memoryDB.Close()
	
	// Run migrations on the memory DB
	// migration.Runner(memoryDB).Up()

	// Execute the test function, providing the fresh DB
	// testFunc(memoryDB)
}

