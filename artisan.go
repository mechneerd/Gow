//go:build ignore

package main

import (
	"gow/bootstrap"
	"gow/console"

	"github.com/spf13/cobra"
)

func main() {
	// Use canonical bootstrap that registers all core providers
	app := bootstrap.NewApplication(".")

	// Create Console Kernel
	kernel := console.NewKernel(app)

	// Example basic command (for testing)
	kernel.RegisterCommand(&cobra.Command{
		Use:   "hello",
		Short: "Prints a greeting",
		Run: func(cmd *cobra.Command, args []string) {
			println("Hello from Artisan!")
		},
	})

	// Run the console kernel
	kernel.Run()
}
