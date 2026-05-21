package queue

import (
	"log"
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
func (w *Worker) Work(connectionName string) {
	driver := w.manager.Connection(connectionName)
	if driver == nil {
		log.Fatalf("Queue connection [%s] not found.", connectionName)
	}

	log.Printf("Starting worker on connection [%s]...", connectionName)

	// Since SyncDriver uses channels, we can optimize by checking if it's the sync driver.
	if syncDriver, ok := driver.(*SyncDriver); ok {
		for job := range syncDriver.Channel() {
			w.processJob(job)
		}
	} else {
		// Generic polling for Redis/DB drivers
		for {
			job, err := driver.Pop()
			if err != nil {
				// Sleep/backoff logic would go here on error or empty queue
				continue
			}
			if job != nil {
				w.processJob(job)
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
