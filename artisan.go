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

	// Artisan core
	kernel.RegisterCommand(artisan.CacheClearCmd)
	kernel.RegisterCommand(artisan.CacheForgetCmd)
	kernel.RegisterCommand(artisan.ConfigCacheCmd)
	kernel.RegisterCommand(artisan.ConfigClearCmd)
	kernel.RegisterCommand(artisan.ScheduleRunCmd)
	kernel.RegisterCommand(artisan.QueueWorkCmd)
	kernel.RegisterCommand(artisan.QueueRetryCmd)
	kernel.RegisterCommand(artisan.VendorPublishCmd)
	kernel.RegisterCommand(artisan.DownCmd)
	kernel.RegisterCommand(artisan.UpCmd)

	// Generators (make:*)
	kernel.RegisterCommand(artisan.MakeControllerCmd)
	kernel.RegisterCommand(artisan.MakeModelCmd)
	kernel.RegisterCommand(artisan.MakeMigrationCmd)
	kernel.RegisterCommand(artisan.MakeMiddlewareCmd)
	kernel.RegisterCommand(artisan.MakeJobCmd)
	kernel.RegisterCommand(artisan.MakeCommandCmd)
	kernel.RegisterCommand(artisan.MakeAuthCmd)

	// Additional generators (Wave 4 completion)
	kernel.RegisterCommand(artisan.MakeMailCmd)
	kernel.RegisterCommand(artisan.MakeEventCmd)
	kernel.RegisterCommand(artisan.MakeListenerCmd)
	kernel.RegisterCommand(artisan.MakePolicyCmd)
	kernel.RegisterCommand(artisan.MakeResourceCmd)
	kernel.RegisterCommand(artisan.MakeJobCmd)
	kernel.RegisterCommand(artisan.MakeNotificationCmd)
	kernel.RegisterCommand(artisan.MakeTestCmd)

	// Cache / Config / View / Queue / Schedule commands
	kernel.RegisterCommand(artisan.CacheClearCmd)
	kernel.RegisterCommand(artisan.ConfigCacheCmd)
	kernel.RegisterCommand(artisan.ConfigClearCmd)
	kernel.RegisterCommand(artisan.ViewClearCmd)
	kernel.RegisterCommand(artisan.QueueRetryCmd)
	kernel.RegisterCommand(artisan.QueueFailedCmd)
	kernel.RegisterCommand(artisan.QueueFlushCmd)
	kernel.RegisterCommand(artisan.ScheduleListCmd)

	// Run the console kernel
	kernel.Run()
}
