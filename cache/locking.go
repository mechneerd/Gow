package cache

import (
	"sync"
	"time"
)

// DistributedLock provides distributed locking using the cache store.
type DistributedLock struct {
	store Store
	key   string
	value string
	ttl   time.Duration
}

// NewDistributedLock creates a new distributed lock.
func NewDistributedLock(store Store, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		store: store,
		key:   "lock:" + key,
		value: generateLockValue(),
		ttl:   ttl,
	}
}

// Acquire attempts to acquire the lock.
func (dl *DistributedLock) Acquire() bool {
	existing, _ := dl.store.Get(dl.key)
	if existing != nil {
		return false
	}
	err := dl.store.Put(dl.key, dl.value, dl.ttl)
	return err == nil
}

// AcquireBlocking attempts to acquire the lock with retry.
func (dl *DistributedLock) AcquireBlocking(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if dl.Acquire() {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// Release releases the lock if owned by this instance.
func (dl *DistributedLock) Release() bool {
	val, _ := dl.store.Get(dl.key)
	if val == dl.value {
		dl.store.Forget(dl.key)
		return true
	}
	return false
}

// ForceRelease releases the lock regardless of ownership.
func (dl *DistributedLock) ForceRelease() {
	dl.store.Forget(dl.key)
}

// IsHeld checks if the lock is currently held.
func (dl *DistributedLock) IsHeld() bool {
	val, _ := dl.store.Get(dl.key)
	return val != nil
}

// GetRemainingTTL returns the remaining TTL of the lock (approximate).
func (dl *DistributedLock) GetRemainingTTL() time.Duration {
	// Simplified - in production track acquisition time
	return dl.ttl
}

func generateLockValue() string {
	return time.Now().Format("20060102150405.000000000")
}

// TagStore extends TaggedCache with additional operations.
type TagStore interface {
	GetTagsForKey(key string) []string
	SetTagsForKey(key string, tags []string)
	GetKeysForTag(tag string) []string
	AddKeyToTag(tag, key string)
	RemoveKeyFromTag(tag, key string)
	FlushTag(tag string)
}

// InMemoryTagStore is an in-memory implementation of TagStore.
type InMemoryTagStore struct {
	tagKeys map[string]map[string]bool
	keyTags map[string][]string
	mu      sync.RWMutex
}

// NewInMemoryTagStore creates a new in-memory tag store.
func NewInMemoryTagStore() *InMemoryTagStore {
	return &InMemoryTagStore{
		tagKeys: make(map[string]map[string]bool),
		keyTags: make(map[string][]string),
	}
}

func (s *InMemoryTagStore) GetTagsForKey(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.keyTags[key]
}

func (s *InMemoryTagStore) SetTagsForKey(key string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keyTags[key] = tags
	for _, tag := range tags {
		if s.tagKeys[tag] == nil {
			s.tagKeys[tag] = make(map[string]bool)
		}
		s.tagKeys[tag][key] = true
	}
}

func (s *InMemoryTagStore) GetKeysForTag(tag string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var keys []string
	for key := range s.tagKeys[tag] {
		keys = append(keys, key)
	}
	return keys
}

func (s *InMemoryTagStore) AddKeyToTag(tag, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tagKeys[tag] == nil {
		s.tagKeys[tag] = make(map[string]bool)
	}
	s.tagKeys[tag][key] = true
	s.keyTags[key] = append(s.keyTags[key], tag)
}

func (s *InMemoryTagStore) RemoveKeyFromTag(tag, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tagKeys[tag] != nil {
		delete(s.tagKeys[tag], key)
	}
}

func (s *InMemoryTagStore) FlushTag(tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tagKeys, tag)
}

// AtomicCounter provides atomic increment/decrement operations.
type AtomicCounter struct {
	store Store
}

// NewAtomicCounter creates a new atomic counter.
func NewAtomicCounter(store Store) *AtomicCounter {
	return &AtomicCounter{store: store}
}

// Increment atomically increments a key by the given amount.
func (ac *AtomicCounter) Increment(key string, amount int) (int, error) {
	return ac.store.Increment(key, amount)
}

// Decrement atomically decrements a key by the given amount.
func (ac *AtomicCounter) Decrement(key string, amount int) (int, error) {
	return ac.store.Decrement(key, amount)
}

// Get returns the current value of the counter.
func (ac *AtomicCounter) Get(key string) (int, error) {
	val, err := ac.store.Get(key)
	if err != nil || val == nil {
		return 0, err
	}
	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, nil
	}
}

// Set sets the counter to a specific value.
func (ac *AtomicCounter) Set(key string, value int) error {
	return ac.store.Put(key, value, 0)
}

// Reset resets the counter to zero.
func (ac *AtomicCounter) Reset(key string) error {
	return ac.store.Put(key, 0, 0)
}

// CacheStats provides statistics about cache usage.
type CacheStats struct {
	Hits       int64
	Misses     int64
	Sets       int64
	Deletes    int64
	Flushes    int64
	Evictions  int64
	mu         sync.RWMutex
}

// NewCacheStats creates a new CacheStats instance.
func NewCacheStats() *CacheStats {
	return &CacheStats{}
}

// RecordHit records a cache hit.
func (s *CacheStats) RecordHit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Hits++
}

// RecordMiss records a cache miss.
func (s *CacheStats) RecordMiss() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Misses++
}

// RecordSet records a cache set operation.
func (s *CacheStats) RecordSet() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sets++
}

// RecordDelete records a cache delete operation.
func (s *CacheStats) RecordDelete() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Deletes++
}

// RecordFlush records a cache flush operation.
func (s *CacheStats) RecordFlush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Flushes++
}

// HitRate returns the cache hit rate as a percentage.
func (s *CacheStats) HitRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total) * 100
}

// Reset resets all statistics.
func (s *CacheStats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Hits = 0
	s.Misses = 0
	s.Sets = 0
	s.Deletes = 0
	s.Flushes = 0
	s.Evictions = 0
}

// Snapshot returns a copy of the current stats.
func (s *CacheStats) Snapshot() CacheStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return CacheStats{
		Hits:      s.Hits,
		Misses:    s.Misses,
		Sets:      s.Sets,
		Deletes:   s.Deletes,
		Flushes:   s.Flushes,
		Evictions: s.Evictions,
	}
}
