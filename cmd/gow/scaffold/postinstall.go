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
	fmt.Println("→ Finalizing project setup...")

	if opts.RunGoModTidy {
		fmt.Println("   → Installing dependencies...")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = projectDir
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard

		if err := cmd.Run(); err != nil {
			fmt.Printf("   ⚠️  Dependency install failed: %v\n", err)
		} else {
			fmt.Println("   ✓ Dependencies ready")
		}
	}

	if opts.CopyEnv {
		envExample := filepath.Join(projectDir, ".env.example")
		envFile := filepath.Join(projectDir, ".env")

		if _, err := os.Stat(envExample); err == nil {
			if _, err := os.Stat(envFile); os.IsNotExist(err) {
				fmt.Println("   → Preparing environment file...")
				if err := copyFile(envExample, envFile); err != nil {
					fmt.Printf("   ⚠️  Failed to create .env: %v\n", err)
				} else {
					fmt.Println("   ✓ .env created")
				}
			}
		}
	}

	return nil
}

// PrintNextSteps shows a clean, professional message after successful project creation.
func PrintNextSteps(projectName string) {
	fmt.Printf(`
✅  Project "%s" created successfully!
📂  Location: ./%s

────────────────────────────────────────
Next steps:

  cd %s
  gow serve                 # Start the development server

  # Optional:
  gow migrate               # Run database migrations (if any in template)

────────────────────────────────────────
🎉  Happy coding with GoW!
`, projectName, projectName, projectName)
}
