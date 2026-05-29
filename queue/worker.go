package queue

import (
	"context"
	"log"
	"time"
)

// Worker processes jobs from a queue connection.
type Worker struct {
	manager *Manager
}

// NewWorker initializes a new Worker.
func NewWorker(manager *Manager) *Worker {
	return &Worker{manager: manager}
}

// Work starts a worker daemon for the given connection.
// It runs until the context is cancelled.
func (w *Worker) Work(ctx context.Context, connectionName string) {
	driver := w.manager.Connection(connectionName)
	if driver == nil {
		log.Fatalf("Queue connection [%s] not found.", connectionName)
	}

	log.Printf("Starting worker on connection [%s]...", connectionName)

	// Since SyncDriver uses channels, we can optimize by checking if it's the sync driver.
	if syncDriver, ok := driver.(*SyncDriver); ok {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Worker on [%s] shutting down.", connectionName)
				return
			case job, ok := <-syncDriver.Channel():
				if !ok {
					return
				}
				w.processJob(job)
			}
		}
	} else {
		// Generic polling for Redis/DB drivers
		for {
			select {
			case <-ctx.Done():
				log.Printf("Worker on [%s] shutting down.", connectionName)
				return
			default:
			}

			job, err := driver.Pop()
			if err != nil {
				log.Printf("Queue pop error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			if job != nil {
				w.processJob(job)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (w *Worker) processJob(job Job) {
	err := job.Handle()
	if err != nil {
		job.Failed(err)
	}
}

