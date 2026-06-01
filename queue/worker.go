package queue

import (
	"context"
	"log"
	"math"
	"time"
)

const (
	defaultMaxRetries  = 3
	defaultBaseBackoff = 1 * time.Second
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

func (w *Worker) getMaxRetries(job Job) int {
	if hj, ok := job.(HasMaxRetries); ok {
		return hj.MaxRetries()
	}
	if hj, ok := job.(HasMaxAttempts); ok {
		return hj.MaxAttempts() - 1 // MaxAttempts includes the first attempt
	}
	return defaultMaxRetries
}

func (w *Worker) getBackoff(job Job, attempt int) time.Duration {
	if hj, ok := job.(HasBackoff); ok {
		return hj.Backoff()
	}
	return time.Duration(math.Pow(2, float64(attempt))) * defaultBaseBackoff
}

func (w *Worker) getTimeout(job Job) time.Duration {
	if hj, ok := job.(HasTimeout); ok {
		return hj.Timeout()
	}
	return 0 // no timeout by default
}

func (w *Worker) processJobWithRetry(job Job) {
	maxRetries := w.getMaxRetries(job)
	timeout := w.getTimeout(job)
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		var err error
		if timeout > 0 {
			err = w.executeWithTimeout(job, timeout)
		} else {
			err = job.Handle()
		}

		if err == nil {
			return // success
		}

		lastErr = err
		if attempt < maxRetries {
			backoff := w.getBackoff(job, attempt)
			log.Printf("Job %T failed (attempt %d/%d): %v, retrying in %v", job, attempt+1, maxRetries+1, err, backoff)
			time.Sleep(backoff)
		}
	}

	log.Printf("Job %T failed after %d attempts: %v", job, maxRetries+1, lastErr)
	job.Failed(lastErr)
}

func (w *Worker) executeWithTimeout(job Job, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		done <- job.Handle()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return context.DeadlineExceeded
	}
}

func (w *Worker) processJob(job Job) {
	w.processJobWithRetry(job)
}

