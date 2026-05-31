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
