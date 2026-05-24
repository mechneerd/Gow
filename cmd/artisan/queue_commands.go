package artisan

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/mechneerd/gow/queue"
	"time"
)

// QueueManager holds the globally configured queue manager instance.
// This should be injected or set during framework bootstrap.
var QueueManager *queue.Manager

// QueueWorkCmd starts processing jobs on the queue as a daemon.
var QueueWorkCmd = &cobra.Command{
	Use:   "queue:work [connection]",
	Short: "Start processing jobs on the queue as a daemon",
	Run: func(cmd *cobra.Command, args []string) {
		connection := ""
		if len(args) > 0 {
			connection = args[0]
		}
		
		fmt.Printf("Starting queue worker for connection: %s\n", connection)
		
		if QueueManager == nil {
			fmt.Println("Error: QueueManager is not initialized.")
			return
		}

		driver := QueueManager.Connection(connection)
		if driver == nil {
			fmt.Printf("Error: queue driver for connection [%s] not found\n", connection)
			return
		}

		for {
			job, err := driver.Pop()
			if err != nil {
				fmt.Printf("Error popping job: %s\n", err)
				time.Sleep(3 * time.Second)
				continue
			}
			
			if job == nil {
				time.Sleep(1 * time.Second)
				continue
			}
			
			err = job.Handle()
			if err != nil {
				fmt.Printf("Job Failed: %s\n", err)
				job.Failed(err)
				// Here we would typically push to a failed_jobs table
			} else {
				fmt.Println("Job processed successfully.")
			}
		}
	},
}

// QueueRetryCmd retries a failed job.
var QueueRetryCmd = &cobra.Command{
	Use:   "queue:retry [id]",
	Short: "Retry a failed queue job",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		fmt.Printf("Retrying failed job %s...\n", id)
		// Full retry from failed_jobs table requires database queue driver wiring.
		fmt.Println("Job retry requested (implementation depends on queue driver).")
	},
}

