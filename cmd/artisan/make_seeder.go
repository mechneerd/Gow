package artisan

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var MakeSeederCmd = &cobra.Command{
	Use:   "make:seeder [name]",
	Short: "Create a new database seeder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !strings.HasSuffix(name, "Seeder") {
			name += "Seeder"
		}

		timestamp := time.Now().Format("20060102150405")
		filename := fmt.Sprintf("%s_%s.go", timestamp, strings.ToLower(name))
		path := fmt.Sprintf("database/seeders/%s", filename) // Standardized to database/seeders

		stub := `package seeders

import "fmt"

func ` + name + `() {
	fmt.Println("🌱 Running ` + name + `...")

	// TODO: Add your seeding logic here
	// Example:
	// user := Models.User{Name: "Example", Email: "example@test.com"}
	// user.Save()

	fmt.Println("✅ ` + name + ` completed.")
}
`
		generateFile(path, stub)
		fmt.Printf("Seeder created: %s\n", path)
	},
}

