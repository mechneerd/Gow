package scaffold

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// PostInstallOptions controls what happens after scaffolding.
type PostInstallOptions struct {
	RunGoModTidy bool
	CopyEnv      bool
}

// RunPostInstall performs common post-scaffolding tasks.
func RunPostInstall(projectDir string, opts PostInstallOptions) error {
	if opts.RunGoModTidy {
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = projectDir
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard

		if err := cmd.Run(); err != nil {
			fmt.Printf("   ⚠️  Dependency install failed: %v\n", err)
		}
	}

	if opts.CopyEnv {
		envExample := filepath.Join(projectDir, ".env.example")
		envFile := filepath.Join(projectDir, ".env")

		if _, err := os.Stat(envExample); err == nil {
			if _, err := os.Stat(envFile); os.IsNotExist(err) {
				if err := copyFile(envExample, envFile); err != nil {
					fmt.Printf("   ⚠️  Failed to create .env: %v\n", err)
				}
			}
		}
	}

	return nil
}

// PrintNextSteps shows a clean, professional message after successful project creation.
func PrintNextSteps(projectName string) {
	green := "\033[32m"
	cyan := "\033[36m"
	reset := "\033[0m"

	fmt.Printf(`
%s✅  Project "%s" created successfully!%s
📂  Location: ./%s

────────────────────────────────────────
Next steps:

  cd %s
  gow serve                 # Start the development server

  # Optional:
  gow migrate               # Run all pending migrations
  gow migrate:rollback      # Rollback last migration (or --step=N)
  gow migrate:run <name>    # Run one specific migration file
  gow migrate:status        # Show migration status
  gow db:seed               # Seed roles + create superadmin (auth kits)
  gow key:generate          # Generate secure APP_KEY
  gow about                 # Framework & environment info
  gow make:seeder           # Create database seeders
  gow make:request          # Create FormRequest validation classes
  gow make:controller       # Create controllers (--resource or --api)
  gow make:command          # Create your own custom gow commands
  gow env                   # Show current environment
  gow list                  # List all available commands

────────────────────────────────────────
%s🎉  Happy coding with GoW!%s
`, green, projectName, reset, projectName, projectName, cyan, reset)
}
