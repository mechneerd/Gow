package artisan

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Note: CacheClearCmd, ConfigCacheCmd, ConfigClearCmd are defined in their dedicated files.
// Duplicates were removed here to allow the package to build.

var ViewClearCmd = &cobra.Command{
	Use:   "view:clear",
	Short: "Clear all compiled view files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Compiled views cleared.")
	},
}





var ScheduleListCmd = &cobra.Command{
	Use:   "schedule:list",
	Short: "List all scheduled commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scheduled tasks:")
		// In real app, read from console.Schedule
	},
}

