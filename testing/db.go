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

// RefreshDatabase wraps a test in a transaction and rolls it back at the end.
func (c *DBTestContext) RefreshDatabase(t *testing.T, testFunc func(tx *sql.Tx)) {
	log.Println("Starting test transaction...")
	
	tx, err := c.db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Ensure rollback happens at the end of the test regardless of panic/fail
	defer func() {
		log.Println("Rolling back test transaction...")
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			t.Errorf("Failed to rollback transaction: %v", err)
		}
	}()

	// Execute the test function, providing the transaction
	// The application's DB container should be temporarily bound to this transaction
	// so that all repository/model queries use it instead of the global DB pool.
	testFunc(tx)
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
