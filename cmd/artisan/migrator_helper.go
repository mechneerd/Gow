package artisan

import (
	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgres driver
	_ "modernc.org/sqlite"

	"github.com/mechneerd/gow/database/migration"
)


// getMigrator creates a working Migrator instance from environment variables.
// This is a pragmatic implementation to make `gow migrate` and related commands functional.
// It now delegates to the new public migration API.
func getMigrator() (*migration.Migrator, error) {
	db, d, err := migration.ConnectFromEnv()
	if err != nil {
		return nil, err
	}

	// Use the default registry (populated by migration.Register calls in init())
	return migration.DefaultMigrator(db, d), nil
}

