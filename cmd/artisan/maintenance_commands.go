package artisan

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type maintenancePayload struct {
	Message  string   `json:"message"`
	Allowed  []string `json:"allowed,omitempty"`
	Retry    int      `json:"retry,omitempty"`
}

var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "Put the application into maintenance mode",
	Run: func(cmd *cobra.Command, args []string) {
		msg, _ := cmd.Flags().GetString("message")
		retry, _ := cmd.Flags().GetInt("retry")
		allowIPs, _ := cmd.Flags().GetStringSlice("allow")

		if msg == "" {
			msg = "Service Unavailable"
		}

		path := filepath.Join("storage", "framework", "down")
		os.MkdirAll(filepath.Dir(path), 0755)

		payload := maintenancePayload{
			Message: msg,
			Allowed: allowIPs,
			Retry:   retry,
		}

		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Println("Error creating maintenance payload:", err)
			return
		}

		err = os.WriteFile(path, data, 0644)
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

// IsDown checks if the application is in maintenance mode.
func IsDown() bool {
	path := filepath.Join("storage", "framework", "down")
	_, err := os.Stat(path)
	return err == nil
}

// IsAllowedIP checks if an IP is allowed to bypass maintenance mode.
func IsAllowedIP(ip string) bool {
	path := filepath.Join("storage", "framework", "down")
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var payload maintenancePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return false
	}

	if len(payload.Allowed) == 0 {
		return false
	}

	for _, allowed := range payload.Allowed {
		if allowed == ip || allowed == "*" {
			return true
		}
		// Check CIDR ranges
		if strings.Contains(allowed, "/") {
			_, cidr, err := net.ParseCIDR(allowed)
			if err != nil {
				continue
			}
			parsedIP := net.ParseIP(ip)
			if parsedIP != nil && cidr.Contains(parsedIP) {
				return true
			}
		}
	}
	return false
}

// GetRetryAfter returns the retry-after header value from the maintenance file.
func GetRetryAfter() int {
	path := filepath.Join("storage", "framework", "down")
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}

	var payload maintenancePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return 0
	}

	return payload.Retry
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

func init() {
	DownCmd.Flags().StringP("message", "m", "", "The message to display during maintenance")
	DownCmd.Flags().IntP("retry", "r", 0, "The number of seconds to wait before retrying (Retry-After header)")
	DownCmd.Flags().StringSlice("allow", nil, "IP addresses or CIDR ranges allowed to bypass maintenance")
}

