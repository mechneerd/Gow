package migration

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mechneerd/gow/database/dialect"
	_ "modernc.org/sqlite"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// ConnectFromEnv reads standard environment variables (DB_CONNECTION, DB_DATABASE, etc.)
// and returns a ready *sql.DB connection along with the matching Dialect.
//
// It does **not** automatically load .env files — callers (especially the CLI generated runner)
// should call godotenv.Load() first if they want .env support.
func ConnectFromEnv() (*sql.DB, dialect.Dialect, error) {
	// Best-effort .env loading (graceful fallback if godotenv is not used or .env is missing)
	_ = godotenv.Load()
	_ = godotenv.Load(".env.local")

	driver := os.Getenv("DB_CONNECTION")
	if driver == "" {
		driver = "sqlite"
	}

	var dsn string
	switch driver {
	case "sqlite":
		dsn = os.Getenv("DB_DATABASE")
		if dsn == "" {
			dsn = "database.sqlite"
		}
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			getEnv("DB_USERNAME", "root"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_HOST", "127.0.0.1"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_DATABASE", ""),
		)
	case "postgres", "pgsql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "127.0.0.1"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USERNAME", "postgres"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_DATABASE", ""),
		)
		driver = "postgres"
	default:
		return nil, nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	var d dialect.Dialect
	switch driver {
	case "sqlite":
		d = &dialect.SQLiteDialect{}
	case "mysql":
		d = &dialect.MySQLDialect{}
	case "postgres":
		d = &dialect.PostgresDialect{}
	default:
		d = &dialect.SQLiteDialect{}
	}

	return db, d, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// RunPending runs all pending migrations using the provided DB and Dialect.
func RunPending(db *sql.DB, d dialect.Dialect) error {
	m := DefaultMigrator(db, d)
	return m.Migrate()
}

// Fresh drops the migrations table (if it exists) and re-runs all migrations from scratch.
func Fresh(db *sql.DB, d dialect.Dialect) error {
	m := DefaultMigrator(db, d)
	return m.Fresh()
}

// Rollback rolls back the last N batches of migrations.
func Rollback(db *sql.DB, d dialect.Dialect, steps int) error {
	m := DefaultMigrator(db, d)
	return m.RollbackSteps(steps)
}

// Status prints the current migration status to stdout.
func Status(db *sql.DB, d dialect.Dialect) error {
	m := DefaultMigrator(db, d)
	return m.Status()
}
