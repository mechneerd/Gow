package artisan

import (
	"fmt"
	"gow/console"

	"github.com/spf13/cobra"
)

// Assuming the user configures their schedule in a central place (e.g., app/console/kernel.go)
// we would inject the Schedule object here. For demonstration, we create it directly.
var scheduleInstance *console.Schedule

func init() {
	scheduleInstance = console.NewSchedule()
	// Example registration (this would typically be done by the user in app/console/kernel.go)
	// scheduleInstance.Command("emails:send").EveryMinute().WithoutOverlapping()
}

// ScheduleRunCmd defines the artisan schedule:run command.
var ScheduleRunCmd = &cobra.Command{
	Use:   "schedule:run",
	Short: "Run the scheduled commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running schedule...")
		scheduleInstance.Run()
	},
}
