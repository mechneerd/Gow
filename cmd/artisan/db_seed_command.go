package artisan

import (
	"fmt"
	"os"
	"os/exec"
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
				found = append(found, strings.TrimSuffix(name, ".go"))
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
			fmt.Printf("     - %s\n", f)
		}

		// Try to run seeders via the project's main.go with "seed" argument
		fmt.Println("\n   → Executing seeders...")
		seedCmd := exec.Command("go", "run", "main.go", "seed")
		seedCmd.Stdout = os.Stdout
		seedCmd.Stderr = os.Stderr
		if err := seedCmd.Run(); err != nil {
			// Fallback: print manual instructions
			fmt.Printf("   ⚠️  Auto-execution failed: %v\n", err)
			fmt.Println("   You can run seeders manually:")
			for _, f := range found {
				fmt.Printf("     seeders.%s(nil)\n", f)
			}
		}

		fmt.Println("\n✅ db:seed finished.")
		fmt.Println("   Expected after RoleSeeder: superadmin / 12345678")
	},
}
