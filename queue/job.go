package queue

import "time"

// Job defines the interface for an executable queue job.
type Job interface {
	// Handle executes the job logic.
	Handle() error
	// Failed handles job failure after all retries are exhausted.
	Failed(err error)
}

// ShouldBeUnique interface for jobs that must be unique in the queue.
// If a job with the same UniqueID is already pending, it will not be dispatched.
type ShouldBeUnique interface {
	UniqueID() string
}

// ShouldBeUniqueUntilProcessing is like ShouldBeUnique but releases the lock
// once the job begins processing (not after completion).
type ShouldBeUniqueUntilProcessing interface {
	ShouldBeUnique
	UniqueVia() string // optional: custom lock key prefix
}

// HasTimeout allows a job to specify its own timeout.
type HasTimeout interface {
	Timeout() time.Duration
}

// HasMaxRetries allows a job to specify its own max retry count.
type HasMaxRetries interface {
	MaxRetries() int
}

// HasBackoff allows a job to specify its own backoff strategy.
type HasBackoff interface {
	Backoff() time.Duration
}

// HasMaxAttempts allows a job to specify a hard limit on total attempts.
type HasMaxAttempts interface {
	MaxAttempts() int
}

// JobWithMiddleware allows a job to declare its own middleware.
type JobWithMiddleware interface {
	Middleware() []JobMiddleware
}

// ChainableJob allows a job to define middleware via a static method.
type ChainableJob interface {
	// Chain returns the middleware chain for this job type.
	Chain() []JobMiddleware
}

// HasPriority allows a job to specify its priority (0 = highest, 10 = lowest).
// Jobs with lower priority numbers are processed first.
type HasPriority interface {
	Priority() int
}

