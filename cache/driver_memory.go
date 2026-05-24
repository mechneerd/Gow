package cache

import (
	"sync"
	"time"
)

type memoryItem struct {
	value   any
	expires int64
}

// MemoryDriver implements the cache Store in memory.
type MemoryDriver struct {
	mu    sync.RWMutex
	items map[string]memoryItem
}

// NewMemoryDriver creates a new in-memory cache driver.
func NewMemoryDriver() *MemoryDriver {
	return &MemoryDriver{
		items: make(map[string]memoryItem),
	}
}

func (d *MemoryDriver) Get(key string) (any, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	item, exists := d.items[key]
	if !exists {
		return nil, nil // Cache miss
	}

	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		return nil, nil // Expired
	}

	return item.value, nil
}

func (d *MemoryDriver) Put(key string, value any, ttl time.Duration) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	var expires int64
	if ttl > 0 {
		expires = time.Now().Add(ttl).UnixNano()
	}

	d.items[key] = memoryItem{
		value:   value,
		expires: expires,
	}
	return nil
}

func (d *MemoryDriver) Increment(key string, value int) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item, exists := d.items[key]
	var current int
	if exists {
		if val, ok := item.value.(int); ok {
			current = val
		}
	}
	current += value
	d.items[key] = memoryItem{value: current, expires: item.expires}
	return current, nil
}

func (d *MemoryDriver) Decrement(key string, value int) (int, error) {
	return d.Increment(key, -value)
}

func (d *MemoryDriver) Forever(key string, value any) error {
	return d.Put(key, value, 0)
}

func (d *MemoryDriver) Forget(key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.items, key)
	return nil
}

func (d *MemoryDriver) Flush() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = make(map[string]memoryItem)
	return nil
}

