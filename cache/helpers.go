package cache

import (
	"sync"
	"time"
)

// Remember retrieves an item from the cache and executes the closure if the
// key doesn't exist, storing the closure's result. This is Laravel's
// cache()->remember() pattern.
func Remember(store Store, key string, ttl time.Duration, callback func() (any, error)) (any, error) {
	val, err := store.Get(key)
	if err == nil && val != nil {
		return val, nil
	}

	val, err = callback()
	if err != nil {
		return nil, err
	}

	if err := store.Put(key, val, ttl); err != nil {
		return val, err
	}

	return val, nil
}

// RememberForever retrieves an item from the cache and executes the closure
// if the key doesn't exist, storing the result permanently.
func RememberForever(store Store, key string, callback func() (any, error)) (any, error) {
	val, err := store.Get(key)
	if err == nil && val != nil {
		return val, nil
	}

	val, err = callback()
	if err != nil {
		return nil, err
	}

	if err := store.Forever(key, val); err != nil {
		return val, err
	}

	return val, nil
}

// Many retrieves multiple cache values by their keys.
func Many(store Store, keys []string) map[string]any {
	results := make(map[string]any)
	for _, key := range keys {
		val, err := store.Get(key)
		if err == nil && val != nil {
			results[key] = val
		}
	}
	return results
}

// PutMany stores multiple key-value pairs in the cache with the same TTL.
func PutMany(store Store, values map[string]any, ttl time.Duration) error {
	for key, value := range values {
		if err := store.Put(key, value, ttl); err != nil {
			return err
		}
	}
	return nil
}

// PutManyForever stores multiple key-value pairs permanently.
func PutManyForever(store Store, values map[string]any) error {
	for key, value := range values {
		if err := store.Forever(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Lock represents a distributed lock backed by the cache.
type Lock struct {
	store    Store
	name     string
	owner    string
	ttl      time.Duration
	acquired bool
	mu       sync.Mutex
}

// NewLock creates a new lock instance.
func NewLock(store Store, name string, ttl time.Duration) *Lock {
	return &Lock{
		store: store,
		name:  "lock:" + name,
		owner: randomString(16),
		ttl:   ttl,
	}
}

// Get attempts to acquire the lock. Returns true if successful.
func (l *Lock) Get() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.acquired {
		return true
	}

	// Try to set the lock key - only if it doesn't exist
	existing, _ := l.store.Get(l.name)
	if existing != nil {
		return false
	}

	if err := l.store.Put(l.name, l.owner, l.ttl); err != nil {
		return false
	}

	l.acquired = true
	return true
}

// GetBlocking attempts to acquire the lock, retrying until timeout.
func (l *Lock) GetBlocking(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if l.Get() {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// ForceGet acquires the lock without checking existing locks (force).
func (l *Lock) ForceGet() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.store.Put(l.name, l.owner, l.ttl); err != nil {
		return false
	}
	l.acquired = true
	return true
}

// Release releases the lock only if the current owner holds it.
func (l *Lock) Release() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.acquired {
		return false
	}

	val, _ := l.store.Get(l.name)
	if val != l.owner {
		l.acquired = false
		return false
	}

	l.store.Forget(l.name)
	l.acquired = false
	return true
}

// GetRemaining returns the remaining time-to-live for the lock.
func (l *Lock) GetRemaining() time.Duration {
	val, _ := l.store.Get(l.name)
	if val == nil {
		return 0
	}
	// Lock exists but we don't track exact TTL in memory driver
	// This is an approximation
	return l.ttl
}

// IsHeld checks if the lock is currently held.
func (l *Lock) IsHeld() bool {
	val, _ := l.store.Get(l.name)
	return val != nil
}

// Block is a convenience method that tries to acquire the lock, blocking
// until successful or timeout. Alias for GetBlocking.
func (l *Lock) Block(timeout time.Duration) bool {
	return l.GetBlocking(timeout)
}

// randomString generates a random string of the given length.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond) // Ensure different values
	}
	return string(b)
}

// AcquireLock attempts to acquire a named lock with the given TTL.
// Returns the Lock object. Call lock.Get() to try acquiring, lock.Release() to free.
func AcquireLock(store Store, name string, ttl time.Duration) *Lock {
	return NewLock(store, name, ttl)
}

// GetOrLock tries to get a cached value, or acquires a lock and executes
// the callback if the value doesn't exist. This prevents cache stampede.
func GetOrLock(store Store, key string, ttl time.Duration, callback func() (any, error)) (any, error) {
	// Try cache first
	val, err := store.Get(key)
	if err == nil && val != nil {
		return val, nil
	}

	// Try to acquire lock
	lock := NewLock(store, key+":lock", ttl)
	if !lock.Get() {
		// Another process is building the cache, wait and retry
		time.Sleep(100 * time.Millisecond)
		return store.Get(key)
	}
	defer lock.Release()

	// Double-check after acquiring lock (another process may have populated it)
	val, err = store.Get(key)
	if err == nil && val != nil {
		return val, nil
	}

	// Execute callback
	val, err = callback()
	if err != nil {
		return nil, err
	}

	// Store result
	if putErr := store.Put(key, val, ttl); putErr != nil {
		return val, putErr
	}

	return val, nil
}
