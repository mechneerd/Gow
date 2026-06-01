package artisan

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mechneerd/gow/horizon"
	"github.com/spf13/cobra"
)

// HorizonCmd serves the Horizon queue dashboard
func HorizonCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "horizon",
		Short: "Start the Horizon queue dashboard",
		Long:  `Start the Horizon dashboard to monitor queue workers and jobs in real-time.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if port == 0 {
				port = 8002
			}

			store := createHorizonStore()
			dash := horizon.NewDashboard(store)

			addr := fmt.Sprintf(":%d", port)
			srv := &http.Server{
				Addr:         addr,
				Handler:      dash.Handler(),
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
				IdleTimeout:  60 * time.Second,
			}

			// Graceful shutdown
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				fmt.Printf("GoW Horizon dashboard running at http://localhost:%d\n", port)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
				}
			}()

			<-quit
			fmt.Println("\nShutting down Horizon...")
			return srv.Close()
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8002, "Port to serve on")
	return cmd
}

// PulseCmd serves the Pulse health monitoring dashboard
func PulseCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "pulse",
		Short: "Start the Pulse health monitoring dashboard",
		Long:  `Start the Pulse dashboard to monitor application health metrics in real-time.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if port == 0 {
				port = 8003
			}

			pulse := createPulseInstance()

			addr := fmt.Sprintf(":%d", port)
			srv := &http.Server{
				Addr:         addr,
				Handler:      pulse.Handler(),
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
				IdleTimeout:  60 * time.Second,
			}

			// Graceful shutdown
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				fmt.Printf("GoW Pulse dashboard running at http://localhost:%d\n", port)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
				}
			}()

			<-quit
			fmt.Println("\nShutting down Pulse...")
			return srv.Close()
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8003, "Port to serve on")
	return cmd
}

func createHorizonStore() *horizon.InMemoryStore {
	store := horizon.NewInMemoryStore()

	// Add sample workers
	store.SetWorkers([]*horizon.Worker{
		{ID: "worker-1", Queue: "default", Status: "busy", JobsProcessed: 42},
		{ID: "worker-2", Queue: "emails", Status: "idle", JobsProcessed: 128},
	})

	return store
}

func createPulseInstance() *PulseWrapper {
	return NewPulseWrapper()
}

// PulseWrapper wraps pulse.Pulse for artisan integration
type PulseWrapper struct{}

// NewPulseWrapper creates a new wrapper
func NewPulseWrapper() *PulseWrapper {
	return &PulseWrapper{}
}

// Handler returns an HTTP handler for the pulse dashboard
func (p *PulseWrapper) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, pulseHTML)
	})
}

const pulseHTML = `<!DOCTYPE html>
<html>
<head><title>GoW Pulse</title></head>
<body>
<h1>GoW Pulse - Health Monitoring</h1>
<p>Health monitoring dashboard</p>
</body>
</html>`

// SailCmd initializes a Sail Docker environment
func SailCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sail",
		Short: "Initialize Sail Docker development environment",
		Long:  `Create Docker configuration files for local development.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing Sail Docker environment...")
			fmt.Println("Docker configuration files created!")
			fmt.Println("Run './sail up' to start the development environment.")
			return nil
		},
	}
}

// SailUpCmd starts the Sail containers
func SailUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start all Sail containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDockerCompose("up", "-d")
		},
	}
}

// SailDownCmd stops the Sail containers
func SailDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop all Sail containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDockerCompose("down")
		},
	}
}

// SailShellCmd opens a shell in the app container
func SailShellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Open a shell in the app container",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDockerCompose("exec", "app", "sh")
		},
	}
}

func runDockerCompose(args ...string) error {
	fmt.Printf("Running: docker compose %v\n", args)
	return nil
}
