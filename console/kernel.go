package console

import (
	"context"
	"gow/foundation"
	"os"

	"github.com/spf13/cobra"
)

// Kernel handles incoming console commands.
type Kernel struct {
	app     *foundation.Application
	rootCmd *cobra.Command
}

// NewKernel creates a new Console Kernel.
func NewKernel(app *foundation.Application) *Kernel {
	rootCmd := &cobra.Command{
		Use:   "artisan",
		Short: "GoW Framework CLI",
		Long:  `Artisan is the command-line interface included with GoW.`,
	}

	return &Kernel{
		app:     app,
		rootCmd: rootCmd,
	}
}

// RegisterCommand registers a cobra command.
func (k *Kernel) RegisterCommand(cmd *cobra.Command) {
	k.rootCmd.AddCommand(cmd)
}

// Run executes the console kernel.
func (k *Kernel) Run() {
	// Make the Application available to all commands via context
	ctx := context.WithValue(k.rootCmd.Context(), "app", k.app)
	k.rootCmd.SetContext(ctx)

	if err := k.rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
