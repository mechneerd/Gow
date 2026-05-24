package artisan

import (
	"fmt"
	"sort"
	"strings"

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

		pathFilter, _ := cmd.Flags().GetString("path")
		methodFilter, _ := cmd.Flags().GetString("method")

		routes := globalRouter.GetAllRoutes()

		sort.Slice(routes, func(i, j int) bool {
			if routes[i].Method != routes[j].Method {
				return routes[i].Method < routes[j].Method
			}
			return routes[i].Path < routes[j].Path
		})

		fmt.Printf("%-8s %-50s %-25s %-30s\n", "Method", "URI", "Name", "Middleware")
		fmt.Println("-----------------------------------------------------------------------------------------------------------")

		for _, route := range routes {
			if pathFilter != "" && !strings.Contains(route.Path, pathFilter) {
				continue
			}
			if methodFilter != "" && !strings.EqualFold(route.Method, methodFilter) {
				continue
			}

			name := route.Name
			if name == "" {
				name = "-"
			}
			middleware := "-"
			if len(route.Middlewares) > 0 {
				mwNames := make([]string, len(route.Middlewares))
				for i, mw := range route.Middlewares {
					mwNames[i] = fmt.Sprintf("%T", mw)
				}
				middleware = strings.Join(mwNames, ",")
			}
			fmt.Printf("%-8s %-50s %-25s %-30s\n", route.Method, route.Path, name, middleware)
		}
	},
}

// SetRouterForListing allows injecting the router for the route:list command (used by console kernel).
var globalRouter *routing.Router

func SetRouterForListing(r *routing.Router) {
	globalRouter = r
}

func init() {
	RouteListCmd.Flags().StringP("path", "p", "", "Filter routes by path")
	RouteListCmd.Flags().StringP("method", "m", "", "Filter routes by HTTP method")
}
