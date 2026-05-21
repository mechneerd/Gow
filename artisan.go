//go:build ignore

package main

import (
	"gow/console"
	"gow/foundation"

	"github.com/spf13/cobra"
)

func main() {
	// 1. Create Application
	app := foundation.NewApplication(".")

	// 2. Create Console Kernel
	kernel := console.NewKernel(app)

	// Example: Register a basic command
	kernel.RegisterCommand(&cobra.Command{
		Use:   "hello",
		Short: "Prints a greeting",
		Run: func(cmd *cobra.Command, args []string) {
			println("Hello from Artisan!")
		},
	})

	// 3. Boot App
	app.Boot()

	// 4. Run Kernel
	kernel.Run()
}
