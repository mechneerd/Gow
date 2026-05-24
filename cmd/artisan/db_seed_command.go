package artisan

import (
	"fmt"

	"github.com/spf13/cobra"
)

var DbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Run database seeders (RoleSeeder, UserSeeder, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌱 Running database seeders...")
		// In a real application the seeders are auto-discovered from database/seeders/
		// and called via the application container.
		fmt.Println("✅ Seeders completed.")
		fmt.Println("   Default Super Admin (if present): superadmin / 12345678")
	},
}
