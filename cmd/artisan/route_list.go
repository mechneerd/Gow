package artisan

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"gow/routing"
)

// RouteListCmd lists all registered routes.
var RouteListCmd = &cobra.Command{
	Use:   "route:list",
	Short: "List all registered routes",
	Run: func(cmd *cobra.Command, args []string) {
		if globalRouter == nil {
			fmt.Println("No router registered. Use artisan.SetRouterForListing(r) before running the command.")
			return
		}

		routes := globalRouter.GetAllRoutes()

		sort.Slice(routes, func(i, j int) bool {
			if routes[i].Method != routes[j].Method {
				return routes[i].Method < routes[j].Method
			}
			return routes[i].Path < routes[j].Path
		})

		fmt.Printf("%-8s %-45s %-20s\n", "Method", "URI", "Name")
		fmt.Println("--------------------------------------------------------------------------------")

		for _, route := range routes {
			name := route.Name
			if name == "" {
				name = "-"
			}
			fmt.Printf("%-8s %-45s %-20s\n", route.Method, route.Path, name)
		}
	},
}

// SetRouterForListing allows injecting the router for the route:list command (used by console kernel).
var globalRouter *routing.Router

func SetRouterForListing(r *routing.Router) {
	globalRouter = r
}

func init() {
	// This command will be properly wired once the console kernel has the router.
	// For now we keep the placeholder.
}
