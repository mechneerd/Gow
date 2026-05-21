package queue

// Job defines the interface for an executable queue job.
type Job interface {
	// Handle executes the job logic.
	Handle() error
	// Failed handles job failure after all retries are exhausted.
	Failed(err error)
}
