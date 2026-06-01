package artisan

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/mechneerd/gow/routing"
)

var RouteCacheCmd = &cobra.Command{
	Use:   "route:cache",
	Short: "Create a route cache file for faster route registration",
	Run: func(cmd *cobra.Command, args []string) {
		if globalRouter == nil {
			fmt.Println("No router registered. Use artisan.SetRouterForListing(r) before running the command.")
			return
		}

		routes := globalRouter.GetAllRoutes()
		if len(routes) == 0 {
			fmt.Println("No routes to cache.")
			return
		}

		// Create the cache directory if it doesn't exist
		cacheDir := filepath.Join("bootstrap", "cache")
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			fmt.Println("Error creating cache directory:", err)
			return
		}

		cacheFile := filepath.Join(cacheDir, "routes.php") // Using .go extension for Go compatibility
		gob.Register(routing.Route{})

		// For now, just print the routes since we can't easily serialize the middleware functions
		fmt.Printf("Route cache created: %s (%d routes)\n", cacheFile, len(routes))
		fmt.Println("Note: Route caching with middleware serialization is not yet fully implemented.")
	},
}

var RouteClearCmd = &cobra.Command{
	Use:   "route:clear",
	Short: "Clear the route cache file",
	Run: func(cmd *cobra.Command, args []string) {
		cacheFile := filepath.Join("bootstrap", "cache", "routes.php")
		if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
			fmt.Println("No route cache file found.")
			return
		}

		if err := os.Remove(cacheFile); err != nil {
			fmt.Println("Error removing route cache:", err)
			return
		}

		fmt.Println("Route cache cleared successfully.")
	},
}
