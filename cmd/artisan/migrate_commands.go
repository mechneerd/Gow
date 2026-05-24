package artisan

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mechneerd/gow/database/migration"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run the database migrations",
	Run: func(cmd *cobra.Command, args []string) {
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

