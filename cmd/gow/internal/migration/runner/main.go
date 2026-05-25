package main

import (
	"fmt"
	"os"

	"github.com/mechneerd/gow/cmd/artisan"
)

// This is a minimal runner used by `gow migrate` via `go run -C`.
// It imports the project's database/migrations package (triggering all init() registrations)
// and then runs the migration commands.
func main() {
	fmt.Println("→ Executing migrations via runner...")

	// We reuse the existing migrator commands.
	// In a real implementation we would call the migrator directly,
	// but for the smallest working slice we invoke the command logic.
	if len(os.Args) > 1 && os.Args[1] == "fresh" {
		artisan.MigrateFreshCmd.Run(artisan.MigrateFreshCmd, []string{})
	} else {
		artisan.MigrateCmd.Run(artisan.MigrateCmd, []string{})
	}
}
