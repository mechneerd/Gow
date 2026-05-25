package artisan

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	gowmigration "github.com/mechneerd/gow/cmd/gow/migration"
	"github.com/mechneerd/gow/database/migration"
)


var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run the database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		if err := generateMigrationRegister(); err != nil {
			fmt.Println("Warning: could not generate migration register:", err)
		}

		// Generate a small runner inside the project (this avoids module resolution issues)
		runnerFile := filepath.Join(cwd, ".gow", "migrate_runner.go")
		if err := os.MkdirAll(filepath.Dir(runnerFile), 0755); err != nil {
			fmt.Println("Warning: could not create runner dir:", err)
		}

		modulePath := readModulePath(filepath.Join(cwd, "go.mod"))
		if modulePath == "" {
			modulePath = "unknown"
		}

		runnerCode := fmt.Sprintf(`package main

import (
	_ "%s/database/migrations"
	"log"

	"github.com/mechneerd/gow/migration"
)

func main() {
	db, dialect, err := migration.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := migration.RunPending(db, dialect); err != nil {
		log.Fatal(err)
	}
}
`, modulePath)

		if err := os.WriteFile(runnerFile, []byte(runnerCode), 0644); err != nil {
			fmt.Println("Warning: could not write migration runner:", err)
		}

		execCmd := exec.Command("go", "run", "-C", cwd, runnerFile)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		fmt.Println("→ Running migrations with auto-discovered files...")
		if err := execCmd.Run(); err != nil {
			fmt.Println("Migration execution failed:", err)
			return
		}

		_ = os.Remove(runnerFile)
		_ = os.Remove(filepath.Dir(runnerFile))
	},
}

var MigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		if err := generateMigrationRegister(); err != nil {
			fmt.Println("Warning:", err)
		}

		runnerFile := filepath.Join(cwd, ".gow", "migrate_runner.go")
		if err := os.MkdirAll(filepath.Dir(runnerFile), 0755); err != nil {
			fmt.Println("Warning:", err)
		}

		modulePath := readModulePath(filepath.Join(cwd, "go.mod"))
		if modulePath == "" {
			modulePath = "unknown"
		}

		runnerCode := fmt.Sprintf(`package main

import (
	_ "%s/database/migrations"
	"log"

	"github.com/mechneerd/gow/migration"
)

func main() {
	db, dialect, err := migration.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := migration.Fresh(db, dialect); err != nil {
		log.Fatal(err)
	}
}
`, modulePath)

		_ = os.WriteFile(runnerFile, []byte(runnerCode), 0644)

		execCmd := exec.Command("go", "run", "-C", cwd, runnerFile)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		fmt.Println("→ Running migrate:fresh with auto-discovered files...")
		if err := execCmd.Run(); err != nil {
			fmt.Println("migrate:fresh failed:", err)
			return
		}

		_ = os.Remove(runnerFile)
		_ = os.Remove(filepath.Dir(runnerFile))
	},
}


var MigrateRefreshCmd = &cobra.Command{
	Use:   "migrate:refresh",
	Short: "Rollback all migrations and re-run them",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		if err := generateMigrationRegister(); err != nil {
			fmt.Println("Warning:", err)
		}

		runnerFile := filepath.Join(cwd, ".gow", "migrate_runner.go")
		if err := os.MkdirAll(filepath.Dir(runnerFile), 0755); err != nil {
			fmt.Println("Warning:", err)
		}

		modulePath := readModulePath(filepath.Join(cwd, "go.mod"))
		if modulePath == "" {
			modulePath = "unknown"
		}

		runnerCode := fmt.Sprintf(`package main

import (
	_ "%s/database/migrations"
	"log"

	"github.com/mechneerd/gow/migration"
)

func main() {
	db, dialect, err := migration.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	// Rollback everything, then run all migrations
	if err := migration.Rollback(db, dialect, 0); err != nil {
		log.Fatal(err)
	}
	if err := migration.RunPending(db, dialect); err != nil {
		log.Fatal(err)
	}
}
`, modulePath)

		_ = os.WriteFile(runnerFile, []byte(runnerCode), 0644)

		execCmd := exec.Command("go", "run", "-C", cwd, runnerFile)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		fmt.Println("→ Running migrate:refresh with auto-discovered files...")
		if err := execCmd.Run(); err != nil {
			fmt.Println("migrate:refresh failed:", err)
			return
		}

		_ = os.Remove(runnerFile)
		_ = os.Remove(filepath.Dir(runnerFile))
	},
}

var MigrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration(s)",
	Run: func(cmd *cobra.Command, args []string) {
		steps, _ := cmd.Flags().GetInt("step")
		cwd, _ := os.Getwd()

		if err := generateMigrationRegister(); err != nil {
			fmt.Println("Warning:", err)
		}

		runnerFile := filepath.Join(cwd, ".gow", "migrate_runner.go")
		if err := os.MkdirAll(filepath.Dir(runnerFile), 0755); err != nil {
			fmt.Println("Warning:", err)
		}

		modulePath := readModulePath(filepath.Join(cwd, "go.mod"))
		if modulePath == "" {
			modulePath = "unknown"
		}

		runnerCode := fmt.Sprintf(`package main

import (
	_ "%s/database/migrations"
	"log"

	"github.com/mechneerd/gow/migration"
)

func main() {
	db, dialect, err := migration.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := migration.Rollback(db, dialect, %d); err != nil {
		log.Fatal(err)
	}
}
`, modulePath, steps)

		_ = os.WriteFile(runnerFile, []byte(runnerCode), 0644)

		execCmd := exec.Command("go", "run", "-C", cwd, runnerFile)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		fmt.Printf("→ Running migrate:rollback (--step=%d) with auto-discovered files...\n", steps)
		if err := execCmd.Run(); err != nil {
			fmt.Println("migrate:rollback failed:", err)
			return
		}

		_ = os.Remove(runnerFile)
		_ = os.Remove(filepath.Dir(runnerFile))
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
		cwd, _ := os.Getwd()

		if err := generateMigrationRegister(); err != nil {
			fmt.Println("Warning:", err)
		}

		runnerFile := filepath.Join(cwd, ".gow", "migrate_runner.go")
		if err := os.MkdirAll(filepath.Dir(runnerFile), 0755); err != nil {
			fmt.Println("Warning:", err)
		}

		modulePath := readModulePath(filepath.Join(cwd, "go.mod"))
		if modulePath == "" {
			modulePath = "unknown"
		}

		runnerCode := fmt.Sprintf(`package main

import (
	_ "%s/database/migrations"
	"log"

	"github.com/mechneerd/gow/migration"
)

func main() {
	db, dialect, err := migration.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := migration.Status(db, dialect); err != nil {
		log.Fatal(err)
	}
}
`, modulePath)

		_ = os.WriteFile(runnerFile, []byte(runnerCode), 0644)

		execCmd := exec.Command("go", "run", "-C", cwd, runnerFile)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		fmt.Println("→ Running migrate:status with auto-discovered files...")
		if err := execCmd.Run(); err != nil {
			fmt.Println("migrate:status failed:", err)
			return
		}

		_ = os.Remove(runnerFile)
		_ = os.Remove(filepath.Dir(runnerFile))
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
		return nil
	}

	found, err := gowmigration.FindMigrations(migrationsDir)
	if err != nil {
		return err
	}

	if err := gowmigration.GenerateRegisterFile(migrationsDir, found); err != nil {
		return err
	}

	_ = gowmigration.EnsureRegisterGoIsGitignored(cwd)

	if len(found) > 0 {
		fmt.Printf("→ Discovered %d migration(s). Generated database/migrations/register.go\n", len(found))
	}

	return nil
}

// readModulePath extracts the module name from go.mod (e.g. "myapp")
func readModulePath(goModPath string) string {
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}


