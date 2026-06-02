package artisan

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var DbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Run database seeders (RoleSeeder, UserSeeder, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌱 Running database seeders...")

		seedersDir := "database/seeders"
		entries, err := os.ReadDir(seedersDir)
		if err != nil {
			fmt.Printf("   ⚠️  No %s directory found. Create seeders with: gow make:seeder RoleSeeder\n", seedersDir)
			fmt.Println("✅ db:seed finished (no seeders directory).")
			return
		}

		var found []string
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasSuffix(name, "_seeder.go") || strings.HasSuffix(name, "Seeder.go") {
				found = append(found, name)
			}
		}

		if len(found) == 0 {
			fmt.Printf("   ℹ️  No seeder files found in %s/\n", seedersDir)
			fmt.Println("   Create one with: gow make:seeder RoleSeeder")
			fmt.Println("✅ db:seed finished.")
			return
		}

		fmt.Printf("   Discovered %d seeder(s) in %s/:\n", len(found), seedersDir)
		for _, f := range found {
			name := strings.TrimSuffix(f, ".go")
			fmt.Printf("     - %s\n", name)
		}

		// Build a temporary runner that imports and executes seeders
		fmt.Println("\n   → Executing seeders...")

		// Create a temporary main.go that imports and runs all seeders
		tmpDir, err := os.MkdirTemp("", "gow-seed-*")
		if err != nil {
			fmt.Printf("   ⚠️  Failed to create temp dir: %v\n", err)
			return
		}
		defer os.RemoveAll(tmpDir)

		// Get the module name from go.mod
		moduleName := getModuleName()

		// Build the runner source
		var imports []string
		var calls []string
		for _, f := range found {
			pkgName := strings.TrimSuffix(f, ".go")
			// Import path uses the seeders directory relative to module
			importPath := moduleName + "/database/seeders"
			if !contains(imports, importPath) {
				imports = append(imports, fmt.Sprintf("\t%q", importPath))
			}
			// Try calling as a function: seeders.RoleSeeder()
			calls = append(calls, fmt.Sprintf("\tseeders.%s()", pkgName))
		}

		runnerSrc := fmt.Sprintf(`package main

import (
%s
)

func main() {
%s
}
`, strings.Join(imports, "\n"), strings.Join(calls, "\n"))

		runnerPath := filepath.Join(tmpDir, "main.go")
		if err := os.WriteFile(runnerPath, []byte(runnerSrc), 0644); err != nil {
			fmt.Printf("   ⚠️  Failed to write runner: %v\n", err)
			return
		}

		// Copy go.mod and go.sum to temp dir
		copyFile("go.mod", filepath.Join(tmpDir, "go.mod"))
		copyFile("go.sum", filepath.Join(tmpDir, "go.sum"))

		// Run the temporary program
		seedCmd := exec.Command("go", "run", runnerPath)
		seedCmd.Stdout = os.Stdout
		seedCmd.Stderr = os.Stderr
		seedCmd.Dir = "."

		if err := seedCmd.Run(); err != nil {
			fmt.Printf("   ⚠️  Seed execution failed: %v\n", err)
			fmt.Println("   You can run seeders manually from your app:")
			fmt.Println("     import \"" + moduleName + "/database/seeders\"")
			for _, f := range found {
				name := strings.TrimSuffix(f, ".go")
				fmt.Printf("     seeders.%s()\n", name)
			}
		}

		fmt.Println("\n✅ db:seed finished.")
		fmt.Println("   Expected after RoleSeeder: superadmin / 12345678")
	},
}

func getModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "main"
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return "main"
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		return
	}
	os.WriteFile(dst, data, 0644)
}
