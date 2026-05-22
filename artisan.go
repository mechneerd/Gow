//go:build ignore

package main

import (
	"gow/bootstrap"
	"gow/cmd/artisan"
	"gow/console"

	"github.com/spf13/cobra"
)

func main() {
	// Use canonical bootstrap that registers all core providers
	app := bootstrap.NewApplication(".")

	// Create Console Kernel
	kernel := console.NewKernel(app)

	// Register built-in commands
	kernel.RegisterCommand(artisan.TinkerCmd)
	kernel.RegisterCommand(artisan.RouteListCmd)
	kernel.RegisterCommand(artisan.MigrateCmd)
	kernel.RegisterCommand(artisan.MigrateFreshCmd)
	kernel.RegisterCommand(artisan.MigrateRefreshCmd)

	// Run the console kernel
	kernel.Run()
}
