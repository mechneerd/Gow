package cache

import (
	"sync"
	"time"
)

// TaggedCache provides cache tag-based invalidation.
type TaggedCache struct {
	store Store
	tags  []string
	mu    sync.RWMutex
}

// Tag creates a tagged cache instance.
func Tag(store Store, tags ...string) *TaggedCache {
	return &TaggedCache{store: store, tags: tags}
}

// tagKey generates a unique key for a tag.
func (tc *TaggedCache) tagKey(key string) string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	result := key
	for _, tag := range tc.tags {
		result = tag + ":" + result
	}
	return result
}

// Get retrieves a value by key with tag prefix.
func (tc *TaggedCache) Get(key string) (any, error) {
	return tc.store.Get(tc.tagKey(key))
}

// Put stores a value with tag prefix.
func (tc *TaggedCache) Put(key string, value any, duration time.Duration) error {
	return tc.store.Put(tc.tagKey(key), value, duration)
}

// Increment increments a counter with tag prefix.
func (tc *TaggedCache) Increment(key string, amount int) (int, error) {
	return tc.store.Increment(tc.tagKey(key), amount)
}

// Decrement decrements a counter with tag prefix.
func (tc *TaggedCache) Decrement(key string, amount int) (int, error) {
	return tc.store.Decrement(tc.tagKey(key), amount)
}

// Forever stores a value permanently with tag prefix.
func (tc *TaggedCache) Forever(key string, value any) error {
	return tc.store.Forever(tc.tagKey(key), value)
}

// Forget removes a specific key with tag prefix.
func (tc *TaggedCache) Forget(key string) error {
	return tc.store.Forget(tc.tagKey(key))
}

// Flush removes all keys with any of the tags.
func (tc *TaggedCache) Flush() error {
	return tc.store.Flush()
}

// Has checks if a key exists with tag prefix.
func (tc *TaggedCache) Has(key string) bool {
	val, _ := tc.store.Get(tc.tagKey(key))
	return val != nil
}
