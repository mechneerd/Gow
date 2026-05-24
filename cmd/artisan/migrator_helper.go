package artisan

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
	_ "github.com/mattn/go-sqlite3"    // sqlite driver

	"gow/database/migration"
	"gow/database/schema"
)

// getMigrator creates a working Migrator instance from environment variables.
// This is a pragmatic implementation to make `gow migrate` and related commands functional.
func getMigrator() (*migration.Migrator, error) {
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
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	builder := schema.NewBuilder(db, nil)
	reg := migration.NewRegistry()
	migrator := migration.NewMigrator(db, builder, reg)

	return migrator, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
