package artisan

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Display the current environment configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Environment Information")
		fmt.Println("-----------------------")

		// Try to read .env
		if content, err := os.ReadFile(".env"); err == nil {
			fmt.Println(string(content))
		} else {
			fmt.Println("No .env file found in current directory.")
		}

		fmt.Println("\nGo Environment:")
		fmt.Printf("  GOOS:   %s\n", os.Getenv("GOOS"))
		fmt.Printf("  GOARCH: %s\n", os.Getenv("GOARCH"))
	},
}
