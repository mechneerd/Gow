package artisan

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var AboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Display basic information about the GoW application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoW Framework")
		fmt.Println("-------------")
		fmt.Printf("Go Version:     %s\n", runtime.Version())
		fmt.Printf("OS / Arch:      %s / %s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Println("Environment:    (read from .env when available)")
		fmt.Println("Cache Driver:   file (default)")
		fmt.Println("Queue Driver:   sync (default)")
		fmt.Println("")
		fmt.Println("Tip: Use 'gow list' or 'artisan' to see all available commands.")
	},
}
