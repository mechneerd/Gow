package artisan

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available Artisan commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoW Artisan Commands")
		fmt.Println("=====================")

		commands := []*cobra.Command{
			MigrateCmd,
			MigrateFreshCmd,
			MigrateRefreshCmd,
			MigrateRollbackCmd,
			MigrateRunCmd,
			MigrateStatusCmd,
			DbSeedCmd,
			KeyGenerateCmd,
			AboutCmd,
			EnvCmd,
			MakeModelCmd,
			MakeMigrationCmd,
			MakeControllerCmd,
			MakeRequestCmd,
			MakeSeederCmd,
			MakeCommandCmd,
			RouteListCmd,
			MakeMiddlewareCmd,
			CacheClearCmd,
			CacheForgetCmd,
			ViewClearCmd,
			ScheduleListCmd,
		}

		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Use < commands[j].Use
		})

		for _, c := range commands {
			if c != nil {
				fmt.Printf("  %-25s %s\n", c.Use, c.Short)
			}
		}

		fmt.Println("\nUse 'gow <command> --help' for more information.")
	},
}

