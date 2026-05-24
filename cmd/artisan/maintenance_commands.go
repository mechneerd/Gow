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

