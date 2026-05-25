package artisan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	gowmigration "github.com/mechneerd/gow/cmd/gow/migration"
	"github.com/mechneerd/gow/database/migration"
)


var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run the database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// Auto-discover and generate register.go before running
		if err := generateMigrationRegister(); err != nil {
			fmt.Println("Warning: could not generate migration register:", err)
		}

		migrator, err := getMigrator()
		if err != nil {
			fmt.Println("Error initializing migrator:", err)
			return
		}

		if err := migrator.Migrate(); err != nil {
			fmt.Println("Migration failed:", err)
			return
		}
		fmt.Println("Migrations completed successfully.")
	},
}


var MigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrator, err := getMigrator()
		if err != nil {
			fmt.Println("Error initializing migrator:", err)
			return
		}

		if err := migrator.Fresh(); err != nil {
			fmt.Println("migrate:fresh failed:", err)
			return
		}

		fmt.Println("Database refreshed successfully.")
	},
}

var MigrateRefreshCmd = &cobra.Command{
	Use:   "migrate:refresh",
	Short: "Rollback all migrations and re-run them",
	Run: func(cmd *cobra.Command, args []string) {
		migrator, err := getMigrator()
		if err != nil {
			fmt.Println("Error initializing migrator:", err)
			return
		}

		if err := migrator.Refresh(); err != nil {
			fmt.Println("migrate:refresh failed:", err)
			return
		}

		fmt.Println("Database refreshed successfully.")
	},
}

var MigrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration(s)",
	Run: func(cmd *cobra.Command, args []string) {
		steps, _ := cmd.Flags().GetInt("step")
		fmt.Printf("Rolling back last %d migration(s)...\n", steps)

		migrator, err := getMigrator()
		if err != nil {
			fmt.Println("Error initializing migrator:", err)
			return
		}

		if err := migrator.RollbackSteps(steps); err != nil {
			fmt.Println("Rollback failed:", err)
			return
		}
		fmt.Println("Rollback completed successfully.")
	},
}

var MigrateRunCmd = &cobra.Command{
	Use:   "migrate:run [migration_name]",
	Short: "Run a specific migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Running single migration: %s\n", name)

		migrator, err := getMigrator()
		if err != nil {
			fmt.Println("Error initializing migrator:", err)
			return
		}

		if err := migrator.MigrateOne(name); err != nil {
			fmt.Println("Migration failed:", err)
			return
		}
		fmt.Println("Single migration completed successfully.")
	},
}

var MigrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrator, err := getMigrator()
		if err != nil {
			fmt.Println("Error initializing migrator:", err)
			return
		}

		if err := migrator.Status(); err != nil {
			fmt.Println("Failed to get migration status:", err)
			return
		}
	},
}

func init() {
	// These commands are meant to be registered with the console kernel
	// that has access to the Migrator instance.
	_ = migration.Registry{} // keep import

	MigrateRollbackCmd.Flags().IntP("step", "s", 1, "Number of migrations to rollback")
}

// generateMigrationRegister auto-discovers migrations and writes register.go
func generateMigrationRegister() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationsDir := filepath.Join(cwd, "database", "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return nil // no migrations dir
	}

	found, err := gowmigration.FindMigrations(migrationsDir)
	if err != nil {
		return err
	}

	if err := gowmigration.GenerateRegisterFile(migrationsDir, found); err != nil {
		return err
	}

	// Make sure it's gitignored
	_ = gowmigration.EnsureRegisterGoIsGitignored(cwd)

	if len(found) > 0 {
		fmt.Printf("→ Discovered %d migration(s). Generated database/migrations/register.go\n", len(found))
	}

	return nil
}


