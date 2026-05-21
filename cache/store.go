package cache

import "time"

// Store is the interface for cache drivers.
type Store interface {
	Get(key string) (any, error)
	Put(key string, value any, ttl time.Duration) error
	Increment(key string, value int) (int, error)
	Decrement(key string, value int) (int, error)
	Forever(key string, value any) error
	Forget(key string) error
	Flush() error
}
