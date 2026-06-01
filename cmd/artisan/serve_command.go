package artisan

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application on the PHP development server",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		public, _ := cmd.Flags().GetBool("public")

		if host == "" {
			host = os.Getenv("APP_HOST")
			if host == "" {
				host = "127.0.0.1"
			}
		}

		if port == 0 {
			portStr := os.Getenv("APP_PORT")
			if portStr != "" {
				p, err := strconv.Atoi(portStr)
				if err == nil {
					port = p
				}
			}
			if port == 0 {
				port = 8000
			}
		}

		addr := fmt.Sprintf("%s:%d", host, port)

		// Check if we should use the built-in Go server or external
		if public {
			servePublic(host, port)
			return
		}

		fmt.Printf("GoW development server started on http://%s\n", addr)
		fmt.Println("Press Ctrl+C to stop the server")

		// Try to start with 'go run' if main.go exists
		if _, err := os.Stat("main.go"); err == nil {
			startWithGoRun(host, port)
			return
		}

		// Try to start with 'go run cmd/gow/main.go'
		if _, err := os.Stat("cmd/gow/main.go"); err == nil {
			startWithModule(host, port)
			return
		}

		// Fallback: simple HTTP server serving public directory
		serveFallback(addr)
	},
}

func servePublic(host string, port int) {
	publicDir := "public"
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		fmt.Println("public directory not found, serving from current directory")
		publicDir = "."
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Server running at http://%s\n", addr)

	fs := http.FileServer(http.Dir(publicDir))
	if err := http.ListenAndServe(addr, fs); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}

func serveFallback(addr string) {
	publicDir := "public"
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		publicDir = "."
	}

	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir(publicDir))
	mux.Handle("/", fs)

	// Add a basic handler for SPA routing
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","framework":"gow"}`))
	})

	fmt.Printf("Listening on http://%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func startWithGoRun(host string, port int) {
	env := os.Environ()
	env = append(env, fmt.Sprintf("APP_HOST=%s", host))
	env = append(env, fmt.Sprintf("APP_PORT=%d", port))

	cmd := exec.Command("go", "run", "main.go")
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func startWithModule(host string, port int) {
	env := os.Environ()
	env = append(env, fmt.Sprintf("APP_HOST=%s", host))
	env = append(env, fmt.Sprintf("APP_PORT=%d", port))

	cmd := exec.Command("go", "run", "cmd/gow/main.go", "serve")
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// WatchFiles watches for file changes and restarts the server.
// This is a simplified version - in production you'd use fsnotify.
func WatchFiles(callback func()) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	fileModTimes := make(map[string]time.Time)

	for range ticker.C {
		changed := false
		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			// Skip vendor, .git, etc.
			if strings.HasPrefix(path, "vendor") || strings.HasPrefix(path, ".git") {
				return nil
			}
			if strings.HasSuffix(path, ".go") {
				lastMod, exists := fileModTimes[path]
				if !exists || info.ModTime().After(lastMod) {
					fileModTimes[path] = info.ModTime()
					if exists {
						changed = true
					}
				}
			}
			return nil
		})
		if err != nil {
			continue
		}
		if changed {
			callback()
		}
	}
}

func init() {
	ServeCmd.Flags().StringP("host", "H", "", "The host address to serve on (default: 127.0.0.1)")
	ServeCmd.Flags().IntP("port", "p", 0, "The port to serve on (default: 8000)")
	ServeCmd.Flags().Bool("public", false, "Serve the public directory directly")
}
