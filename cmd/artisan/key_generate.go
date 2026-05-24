package artisan

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var KeyGenerateCmd = &cobra.Command{
	Use:   "key:generate",
	Short: "Generate a new APP_KEY and write it to .env",
	Run: func(cmd *cobra.Command, args []string) {
		key := generateAppKey(32)
		fmt.Println("Application key generated successfully.")

		// Try to update .env file
		envPath := ".env"
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			fmt.Println("No .env file found. Please create one manually with:")
			fmt.Printf("APP_KEY=%s\n", key)
			return
		}

		content, err := os.ReadFile(envPath)
		if err != nil {
			fmt.Println("Could not read .env:", err)
			return
		}

		lines := strings.Split(string(content), "\n")
		updated := false
		for i, line := range lines {
			if strings.HasPrefix(line, "APP_KEY=") {
				lines[i] = "APP_KEY=" + key
				updated = true
				break
			}
		}

		if !updated {
			lines = append(lines, "APP_KEY="+key)
		}

		if err := os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			fmt.Println("Failed to write .env:", err)
			return
		}

		fmt.Println("APP_KEY written to .env file.")
	},
}

func generateAppKey(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return "base64:" + base64.StdEncoding.EncodeToString(b)
}

