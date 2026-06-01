package collection

import (
	"math/rand"
	"sync"
)

// GroupByFn groups items by a callback (generic version).
func (c *Collection[T]) GroupByFn(callback func(T) string) map[string]*Collection[T] {
	result := make(map[string]*Collection[T])
	for _, item := range c.items {
		key := callback(item)
		if _, ok := result[key]; !ok {
			result[key] = Collect([]T{})
		}
		result[key].items = append(result[key].items, item)
	}
	return result
}

// SortBy sorts the collection by a callback.
func (c *Collection[T]) SortBy(less func(a, b T) bool) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	// Using sort.SliceStable equivalent
	for i := 1; i < len(result); i++ {
		for j := i; j > 0 && less(result[j], result[j-1]); j-- {
			result[j], result[j-1] = result[j-1], result[j]
		}
	}
	return Collect(result)
}

// Sum returns the sum of all items.
func (c *Collection[T]) Sum(fn func(T) float64) float64 {
	var sum float64
	for _, item := range c.items {
		sum += fn(item)
	}
	return sum
}

// Average returns the average of all items.
func (c *Collection[T]) Average(fn func(T) float64) float64 {
	if len(c.items) == 0 {
		return 0
	}
	return c.Sum(fn) / float64(len(c.items))
}

// Where filters the collection by a key-value pair (for map items).
func (c *Collection[T]) Where(key string, value any) *Collection[T] {
	var result []T
	for _, item := range c.items {
		if m, ok := any(item).(map[string]any); ok {
			if m[key] == value {
				result = append(result, item)
			}
		}
	}
	return Collect(result)
}

// WhereIn filters by a key and list of values (for map items).
func (c *Collection[T]) WhereIn(key string, values []any) *Collection[T] {
	valueSet := make(map[any]bool)
	for _, v := range values {
		valueSet[v] = true
	}
	var result []T
	for _, item := range c.items {
		if m, ok := any(item).(map[string]any); ok {
			if valueSet[m[key]] {
				result = append(result, item)
			}
		}
	}
	return Collect(result)
}

// WhereNull filters where a key is nil (for map items).
func (c *Collection[T]) WhereNull(key string) *Collection[T] {
	var result []T
	for _, item := range c.items {
		if m, ok := any(item).(map[string]any); ok {
			if m[key] == nil {
				result = append(result, item)
			}
		}
	}
	return Collect(result)
}

// WhereNotNull filters where a key is not nil (for map items).
func (c *Collection[T]) WhereNotNull(key string) *Collection[T] {
	var result []T
	for _, item := range c.items {
		if m, ok := any(item).(map[string]any); ok {
			if m[key] != nil {
				result = append(result, item)
			}
		}
	}
	return Collect(result)
}

// Slice returns a subset of the collection.
func (c *Collection[T]) Slice(start, end int) *Collection[T] {
	if start < 0 {
		start = 0
	}
	if end > len(c.items) {
		end = len(c.items)
	}
	return Collect(c.items[start:end])
}

// SplitN splits the collection into n groups.
func (c *Collection[T]) SplitN(numGroups int) []*Collection[T] {
	if numGroups <= 0 {
		return nil
	}
	size := len(c.items) / numGroups
	var result []*Collection[T]
	for i := 0; i < numGroups; i++ {
		start := i * size
		end := start + size
		if i == numGroups-1 {
			end = len(c.items)
		}
		result = append(result, Collect(c.items[start:end]))
	}
	return result
}

// ZipN merges values from multiple collections.
func ZipN[T any](collections ...*Collection[T]) []any {
	if len(collections) == 0 {
		return nil
	}
	maxLen := 0
	for _, c := range collections {
		if len(c.items) > maxLen {
			maxLen = len(c.items)
		}
	}
	var result []any
	for i := 0; i < maxLen; i++ {
		var row []T
		for _, c := range collections {
			if i < len(c.items) {
				row = append(row, c.items[i])
			}
		}
		result = append(result, row)
	}
	return result
}

// CrossJoin joins with another collection.
func (c *Collection[T]) CrossJoin(other *Collection[T]) []any {
	var result []any
	for _, a := range c.items {
		for _, b := range other.items {
			result = append(result, []T{a, b})
		}
	}
	return result
}

// Every returns true if all items pass the test.
func (c *Collection[T]) Every(fn func(T) bool) bool {
	for _, item := range c.items {
		if !fn(item) {
			return false
		}
	}
	return true
}

// Some returns true if at least one item passes the test.
func (c *Collection[T]) Some(fn func(T) bool) bool {
	for _, item := range c.items {
		if fn(item) {
			return true
		}
	}
	return false
}

// EachWithBreak iterates with early termination.
func (c *Collection[T]) EachWithBreak(fn func(T) bool) {
	for _, item := range c.items {
		if !fn(item) {
			break
		}
	}
}

// MapWith returns a collection of a different type.
func MapWith[T, U any](c *Collection[T], fn func(T) U) *Collection[U] {
	result := make([]U, len(c.items))
	for i, item := range c.items {
		result[i] = fn(item)
	}
	return Collect(result)
}

// ShuffleRandom shuffles the collection using math/rand.
func (c *Collection[T]) ShuffleRandom() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return Collect(result)
}

// UniqueByFn returns unique items by a callback.
func (c *Collection[T]) UniqueByFn(callback func(T) string) *Collection[T] {
	seen := make(map[string]bool)
	var result []T
	for _, item := range c.items {
		key := callback(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return Collect(result)
}

// LazyCollection provides lazy evaluation for large datasets.
type LazyCollection[T any] struct {
	generator func(yield func(T) bool)
}

// Lazy creates a new lazy collection.
func Lazy[T any](generator func(yield func(T) bool)) *LazyCollection[T] {
	return &LazyCollection[T]{generator: generator}
}

// Collect materializes the lazy collection.
func (lc *LazyCollection[T]) Collect() *Collection[T] {
	var items []T
	lc.generator(func(item T) bool {
		items = append(items, item)
		return true
	})
	return Collect(items)
}

// Each iterates over the lazy collection.
func (lc *LazyCollection[T]) Each(callback func(T)) {
	lc.generator(func(item T) bool {
		callback(item)
		return true
	})
}

// EachWithBreak iterates with early termination.
func (lc *LazyCollection[T]) EachWithBreak(callback func(T) bool) {
	lc.generator(callback)
}

// Filter filters the lazy collection.
func (lc *LazyCollection[T]) Filter(callback func(T) bool) *LazyCollection[T] {
	return Lazy(func(yield func(T) bool) {
		lc.generator(func(item T) bool {
			if callback(item) {
				return yield(item)
			}
			return true
		})
	})
}

// Map transforms each element of the lazy collection.
func MapLazy[T, U any](lc *LazyCollection[T], fn func(T) U) *LazyCollection[U] {
	return Lazy(func(yield func(U) bool) {
		lc.generator(func(item T) bool {
			return yield(fn(item))
		})
	})
}

// Take takes the first n items from the lazy collection.
func (lc *LazyCollection[T]) Take(n int) *LazyCollection[T] {
	return Lazy(func(yield func(T) bool) {
		count := 0
		lc.generator(func(item T) bool {
			if count >= n {
				return false
			}
			count++
			return yield(item)
		})
	})
}

// Skip skips the first n items of the lazy collection.
func (lc *LazyCollection[T]) Skip(n int) *LazyCollection[T] {
	return Lazy(func(yield func(T) bool) {
		skipped := 0
		lc.generator(func(item T) bool {
			if skipped < n {
				skipped++
				return true
			}
			return yield(item)
		})
	})
}

// SyncCollection provides thread-safe access to a collection.
type SyncCollection[T any] struct {
	collection *Collection[T]
	mu         sync.RWMutex
}

// SyncCol wraps a collection for thread-safe access.
func SyncCol[T any](c *Collection[T]) *SyncCollection[T] {
	return &SyncCollection[T]{collection: c}
}

func (sc *SyncCollection[T]) All() []T {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.collection.All()
}

func (sc *SyncCollection[T]) Count() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.collection.Count()
}

func (sc *SyncCollection[T]) Push(item T) *SyncCollection[T] {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.collection.Push(item)
	return sc
}

func (sc *SyncCollection[T]) Each(callback func(T)) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	sc.collection.Each(callback)
}
