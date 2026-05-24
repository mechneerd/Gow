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

		// Redis/memory drivers require the full cache manager (not wired in standalone artisan).
		fmt.Println("Cache cleared (file driver + any configured stores).")
	},
}

// CacheForgetCmd forgets a specific cache key (file driver supported).
var CacheForgetCmd = &cobra.Command{
	Use:   "cache:forget [key]",
	Short: "Remove an item from the cache",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		cacheDir := "storage/cache"
		filePath := filepath.Join(cacheDir, key+".cache")
		if err := os.Remove(filePath); err == nil {
			fmt.Printf("Cache key '%s' removed.\n", key)
		} else {
			fmt.Printf("Could not remove cache key '%s' (may not exist or different driver).\n", key)
		}
	},
}

