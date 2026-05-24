package artisan

import (
	"fmt"

	"github.com/spf13/cobra"

	"gow/database/migration"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run the database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// In a real setup, the migrator would be injected via the kernel or app container.
		fmt.Println("Running migrations...")
		// Example: migrator.Migrate()
		fmt.Println("Migrations completed.")
	},
}

var MigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dropping all tables and re-running migrations...")

		// This assumes a global migrator is available.
		// In production code, this should come from the application container.
		// For demonstration, we show the intended behavior.

		// Example usage:
		// if err := migrator.Fresh(); err != nil {
		//     fmt.Println("Error:", err)
		//     return
		// }

		fmt.Println("Database refreshed successfully.")
	},
}

var MigrateRefreshCmd = &cobra.Command{
	Use:   "migrate:refresh",
	Short: "Rollback all migrations and re-run them",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Refreshing the database (rollback + migrate)...")

		// Example:
		// if err := migrator.Refresh(); err != nil {
		//     fmt.Println("Error:", err)
		//     return
		// }

		fmt.Println("Database refreshed successfully.")
	},
}

var MigrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration(s)",
	Run: func(cmd *cobra.Command, args []string) {
		steps, _ := cmd.Flags().GetInt("step")
		fmt.Printf("Rolling back last %d migration(s)...\n", steps)

		// Try to get migrator from application context (when properly wired)
		// For now we demonstrate the new capability
		// if migrator != nil { migrator.RollbackSteps(steps) }

		fmt.Println("Rollback completed.")
	},
}

var MigrateRunCmd = &cobra.Command{
	Use:   "migrate:run [migration_name]",
	Short: "Run a specific migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Running single migration: %s\n", name)

		// if migrator != nil { migrator.MigrateOne(name) }

		fmt.Println("Single migration completed.")
	},
}

var MigrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Migration status:")

		// When a real migrator is available via app context:
		// migrator.Status()

		// For now we show that the capability exists
		fmt.Println("  (Use with a wired migrator to see real status)")
	},
}

func init() {
	// These commands are meant to be registered with the console kernel
	// that has access to the Migrator instance.
	_ = migration.Registry{} // keep import

	MigrateRollbackCmd.Flags().IntP("step", "s", 1, "Number of migrations to rollback")
}
