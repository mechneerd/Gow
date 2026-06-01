package queue

import (
	"fmt"
	"sync"
	"time"
)

// JobMiddleware wraps job execution with custom logic.
type JobMiddleware interface {
	Handle(job Job, next func(Job) error) error
}

// MiddlewareFunc is a function adapter for JobMiddleware.
type MiddlewareFunc func(job Job, next func(Job) error) error

func (m MiddlewareFunc) Handle(job Job, next func(Job) error) error {
	return m(job, next)
}

// RetryMiddleware adds retry logic with exponential backoff.
func RetryMiddleware(maxRetries int) JobMiddleware {
	return MiddlewareFunc(func(job Job, next func(Job) error) error {
		var lastErr error
		for i := 0; i <= maxRetries; i++ {
			lastErr = next(job)
			if lastErr == nil {
				return nil
			}
			if i < maxRetries {
				backoff := time.Duration(1<<uint(i)) * time.Second
				time.Sleep(backoff)
			}
		}
		return fmt.Errorf("job failed after %d retries: %w", maxRetries, lastErr)
	})
}

// TimeoutMiddleware adds a timeout to job execution.
func TimeoutMiddleware(timeout time.Duration) JobMiddleware {
	return MiddlewareFunc(func(job Job, next func(Job) error) error {
		done := make(chan error, 1)
		go func() {
			done <- next(job)
		}()

		select {
		case err := <-done:
			return err
		case <-time.After(timeout):
			return fmt.Errorf("job timed out after %v", timeout)
		}
	})
}

// LoggingMiddleware logs job execution.
func LoggingMiddleware() JobMiddleware {
	return MiddlewareFunc(func(job Job, next func(Job) error) error {
		start := time.Now()
		fmt.Printf("[Queue] Processing %T\n", job)
		err := next(job)
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("[Queue] %T failed after %v: %v\n", job, duration, err)
		} else {
			fmt.Printf("[Queue] %T completed in %v\n", job, duration)
		}
		return err
	})
}

// RateLimitMiddleware limits job execution to N jobs per time period.
// Uses a token bucket algorithm.
func RateLimitMiddleware(maxJobs int, per time.Duration) JobMiddleware {
	var (
		mu       sync.Mutex
		tokens   = maxJobs
		lastTime = time.Now()
	)

	return MiddlewareFunc(func(job Job, next func(Job) error) error {
		mu.Lock()
		now := time.Now()
		elapsed := now.Sub(lastTime)

		// Refill tokens based on elapsed time
		if elapsed > per {
			tokens = maxJobs
			lastTime = now
		} else if tokens <= 0 {
			// Wait until we have a token
			waitTime := per - elapsed
			mu.Unlock()
			time.Sleep(waitTime)
			mu.Lock()
			tokens = maxJobs
			lastTime = time.Now()
		}

		tokens--
		mu.Unlock()

		return next(job)
	})
}

// RateLimitByKeyMiddleware limits jobs by a key extracted from the job.
// Useful for per-user or per-tenant rate limiting.
func RateLimitByKeyMiddleware(keyFunc func(Job) string, maxJobs int, per time.Duration) JobMiddleware {
	type keyTokenBucket struct {
		tokens   int
		lastTime time.Time
	}

	var (
		mu      sync.Mutex
		buckets = make(map[string]*keyTokenBucket)
	)

	return MiddlewareFunc(func(job Job, next func(Job) error) error {
		key := keyFunc(job)

		mu.Lock()
		bucket, exists := buckets[key]
		if !exists {
			bucket = &keyTokenBucket{
				tokens:   maxJobs,
				lastTime: time.Now(),
			}
			buckets[key] = bucket
		}

		now := time.Now()
		elapsed := now.Sub(bucket.lastTime)

		if elapsed > per {
			bucket.tokens = maxJobs
			bucket.lastTime = now
		} else if bucket.tokens <= 0 {
			waitTime := per - elapsed
			mu.Unlock()
			time.Sleep(waitTime)
			mu.Lock()
			bucket.tokens = maxJobs
			bucket.lastTime = time.Now()
		}

		bucket.tokens--
		mu.Unlock()

		return next(job)
	})
}

// WithoutRetry disables retry logic for a job.
func WithoutRetry() JobMiddleware {
	return MiddlewareFunc(func(job Job, next func(Job) error) error {
		return next(job)
	})
}

// EnsureUnique ensures only one instance of a job with the same UniqueID runs at a time.
func EnsureUnique(store UniqueJobStore) JobMiddleware {
	return MiddlewareFunc(func(job Job, next func(Job) error) error {
		if uniqueJob, ok := job.(ShouldBeUnique); ok {
			uid := uniqueJob.UniqueID()
			if !store.Acquire(uid) {
				return fmt.Errorf("job %s is already running", uid)
			}
			defer store.Release(uid)
		}
		return next(job)
	})
}

// UniqueJobStore is an interface for storing job uniqueness locks.
type UniqueJobStore interface {
	Acquire(key string) bool
	Release(key string)
}

// InMemoryUniqueStore is an in-memory implementation of UniqueJobStore.
type InMemoryUniqueStore struct {
	mu    sync.Mutex
	locks map[string]bool
}

func NewInMemoryUniqueStore() *InMemoryUniqueStore {
	return &InMemoryUniqueStore{
		locks: make(map[string]bool),
	}
}

func (s *InMemoryUniqueStore) Acquire(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.locks[key] {
		return false
	}
	s.locks[key] = true
	return true
}

func (s *InMemoryUniqueStore) Release(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.locks, key)
}

// middlewareChain manages job middleware.
type middlewareChain struct {
	middleware []JobMiddleware
	mu        sync.RWMutex
}

// Use adds middleware to the chain.
func (mc *middlewareChain) Use(mw JobMiddleware) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.middleware = append(mc.middleware, mw)
}

// Run executes the middleware chain for a job.
func (mc *middlewareChain) Run(job Job, handler func(Job) error) error {
	mc.mu.RLock()
	mws := make([]JobMiddleware, len(mc.middleware))
	copy(mws, mc.middleware)
	mc.mu.RUnlock()

	if len(mws) == 0 {
		return handler(job)
	}

	// Build chain from inside out
	chain := handler
	for i := len(mws) - 1; i >= 0; i-- {
		mw := mws[i]
		next := chain
		chain = func(j Job) error {
			return mw.Handle(j, next)
		}
	}

	return chain(job)
}
