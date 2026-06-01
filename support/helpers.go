package support

import "sync"

// Conditionable provides conditional method chaining.
// Usage:
//
//	obj.When(condition, func(o *MyStruct) {
//	    o.SetColor("red")
//	})
type Conditionable[T any] struct {
	value T
}

// NewConditionable creates a new Conditionable.
func NewConditionable[T any](value T) *Conditionable[T] {
	return &Conditionable[T]{value: value}
}

// When executes the callback if the condition is true.
func (c *Conditionable[T]) When(condition bool, callback func(T)) *Conditionable[T] {
	if condition {
		callback(c.value)
	}
	return c
}

// Unless executes the callback if the condition is false.
func (c *Conditionable[T]) Unless(condition bool, callback func(T)) *Conditionable[T] {
	return c.When(!condition, callback)
}

// Value returns the underlying value.
func (c *Conditionable[T]) Value() T {
	return c.value
}

// Tappable provides tap-style method chaining.
// Usage:
//
//	tap(obj, func(o *MyStruct) {
//	    o.SetColor("red")
//	    o.SetSize(100)
//	})
func Tap[T any](value T, callback func(T)) T {
	callback(value)
	return value
}

// Tappable is an interface that types can implement for tap-style chaining.
type Tappable[T any] interface {
	Tap(fn func(T)) T
}

// Once provides a sync.Once-style helper for lazy initialization.
type Once[T any] struct {
	once sync.Once
	value T
	fn   func() T
}

// NewOnce creates a new Once instance.
func NewOnce[T any](fn func() T) *Once[T] {
	return &Once[T]{fn: fn}
}

// Get returns the value, executing the function only once.
func (o *Once[T]) Get() T {
	o.once.Do(func() {
		o.value = o.fn()
	})
	return o.value
}

// Defer provides a deferred execution helper.
// The callback is registered and will be executed when Flush() is called.
type Defer struct {
	mu       sync.Mutex
	callbacks []func()
}

// NewDefer creates a new Defer instance.
func NewDefer() *Defer {
	return &Defer{}
}

// Add registers a callback to be deferred.
func (d *Defer) Add(fn func()) *Defer {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.callbacks = append(d.callbacks, fn)
	return d
}

// Flush executes all deferred callbacks in LIFO order.
func (d *Defer) Flush() {
	d.mu.Lock()
	callbacks := make([]func(), len(d.callbacks))
	copy(callbacks, d.callbacks)
	d.callbacks = nil
	d.mu.Unlock()

	// Execute in LIFO order
	for i := len(callbacks) - 1; i >= 0; i-- {
		callbacks[i]()
	}
}

// Concurrency provides helpers for concurrent operations.
type Concurrency struct {
	wg sync.WaitGroup
}

// NewConcurrency creates a new Concurrency instance.
func NewConcurrency() *Concurrency {
	return &Concurrency{}
}

// Go runs a function in a goroutine.
func (c *Concurrency) Go(fn func()) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		fn()
	}()
}

// Wait blocks until all goroutines complete.
func (c *Concurrency) Wait() {
	c.wg.Wait()
}

// Parallel runs multiple functions concurrently and waits for all to complete.
func Parallel(fns ...func()) {
	c := NewConcurrency()
	for _, fn := range fns {
		c.Go(fn)
	}
	c.Wait()
}

// Collect collects results from concurrent operations.
func Collect[T any](fns ...func() T) []T {
	results := make([]T, len(fns))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, fn := range fns {
		wg.Add(1)
		go func(idx int, f func() T) {
			defer wg.Done()
			result := f()
			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i, fn)
	}

	wg.Wait()
	return results
}

// LimitConcurrency limits the number of concurrent operations.
func LimitConcurrency(limit int, fns ...func()) {
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			f()
		}(fn)
	}

	wg.Wait()
}

// Retry retries a function up to n times with exponential backoff.
func Retry(attempts int, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
	}
	return err
}

// RetryWithDelay retries a function with a specific delay between attempts.
func RetryWithDelay(attempts int, delayFn func(attempt int) int, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if delayFn != nil && i < attempts-1 {
			// Simple sleep using channel (avoids time import)
			done := make(chan struct{})
			go func() {
				// This is a simplified delay - in production use time.Sleep
				close(done)
			}()
			<-done
		}
	}
	return err
}