package artisan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "Put the application into maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		path := filepath.Join("storage", "framework", "down")
		os.MkdirAll(filepath.Dir(path), 0755)
		
		err := os.WriteFile(path, []byte(`{"message": "Service Unavailable"}`), 0644)
		if err != nil {
			fmt.Println("Error putting application into maintenance mode:", err)
			return
		}
		fmt.Println("Application is now in maintenance mode.")
	},
}

var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring the application out of maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		path := filepath.Join("storage", "framework", "down")
		if _, err := os.Stat(path); err == nil {
			err := os.Remove(path)
			if err != nil {
				fmt.Println("Error bringing application out of maintenance mode:", err)
				return
			}
		}
		fmt.Println("Application is now live.")
	},
}

var StorageLinkCmd = &cobra.Command{
	Use:   "storage:link",
	Short: "Create the symbolic link for the storage directory",
	Run: func(cmd *cobra.Command, args []string) {
		linkPath := filepath.Join("public", "storage")
		targetPath := filepath.Join("..", "storage", "app", "public")

		// Check if link already exists
		if _, err := os.Lstat(linkPath); err == nil {
			fmt.Println("storage link already exists.")
			return
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
			fmt.Println("Error creating public directory:", err)
			return
		}

		// Create the symbolic link
		if err := os.Symlink(targetPath, linkPath); err != nil {
			fmt.Println("Error creating storage link:", err)
			fmt.Println("Note: You may need to run this command as administrator on Windows.")
			return
		}

		fmt.Println("The [public/storage] link has been connected to [storage/app/public].")
	},
}

var EventListCmd = &cobra.Command{
	Use:   "event:list",
	Short: "List all registered events and their listeners",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Event listing requires application context.")
		fmt.Println("Use event.ListEvents() in your application code.")
	},
}

