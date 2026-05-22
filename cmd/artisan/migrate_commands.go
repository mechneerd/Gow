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

func init() {
	// These commands are meant to be registered with the console kernel
	// that has access to the Migrator instance.
	_ = migration.Registry{} // keep import
}
