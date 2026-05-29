package queue

import (
	"context"
	"log"
	"math"
	"time"
)

const (
	maxRetries  = 3
	baseBackoff = 1 * time.Second
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
				w.processJobWithRetry(job)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (w *Worker) processJobWithRetry(job Job) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := job.Handle(); err != nil {
			lastErr = err
			backoff := time.Duration(math.Pow(2, float64(attempt))) * baseBackoff
			log.Printf("Job failed (attempt %d/%d): %v, retrying in %v", attempt+1, maxRetries+1, err, backoff)
			time.Sleep(backoff)
			continue
		}
		return // success
	}
	job.Failed(lastErr)
}

func (w *Worker) processJob(job Job) {
	w.processJobWithRetry(job)
}

