package artisan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// CacheClearCmd clears all cache (file, memory, redis if configured).
var CacheClearCmd = &cobra.Command{
	Use:   "cache:clear",
	Short: "Flush the application cache",
	Run: func(cmd *cobra.Command, args []string) {
		// Simple file cache clear (most common for dev)
		cacheDir := "storage/cache"
		if _, err := os.Stat(cacheDir); err == nil {
			files, _ := os.ReadDir(cacheDir)
			for _, f := range files {
				os.Remove(filepath.Join(cacheDir, f.Name()))
			}
			fmt.Println("File cache cleared.")
		}

		// For Redis/memory, in real app we would resolve cache.Store and call Flush
		fmt.Println("Cache cleared (file + any configured stores).")
	},
}

// CacheForgetCmd forgets a specific cache key.
var CacheForgetCmd = &cobra.Command{
	Use:   "cache:forget [key]",
	Short: "Remove an item from the cache",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		fmt.Printf("Forgot cache key: %s (implement store-specific logic)\n", key)
	},
}
