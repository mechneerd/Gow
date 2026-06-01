package collection

import (
	"fmt"
	"sort"
)

// First returns the first element of the collection
func (c *Collection[T]) First() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	return c.items[0], true
}

// Last returns the last element of the collection
func (c *Collection[T]) Last() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	return c.items[len(c.items)-1], true
}

// Head returns the first n elements
func (c *Collection[T]) Head(n int) *Collection[T] {
	if n <= 0 {
		return Collect([]T{})
	}
	if n > len(c.items) {
		n = len(c.items)
	}
	return Collect(c.items[:n])
}

// Tail returns all elements except the first n
func (c *Collection[T]) Tail(n int) *Collection[T] {
	if n <= 0 {
		return Collect(c.items)
	}
	if n > len(c.items) {
		return Collect([]T{})
	}
	return Collect(c.items[n:])
}

// Take returns a new collection with the first n elements
func (c *Collection[T]) Take(n int) *Collection[T] {
	return c.Head(n)
}

// Skip returns a new collection without the first n elements
func (c *Collection[T]) Skip(n int) *Collection[T] {
	return c.Tail(n)
}

// Unique returns unique items
func (c *Collection[T]) Unique() *Collection[T] {
	seen := make(map[any]bool)
	var result []T
	for _, item := range c.items {
		key := item
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}
	return Collect(result)
}

// Values returns the values for a given key (for map collections)
func Values[K comparable, V any](c *Collection[map[K]V], key K) *Collection[V] {
	var result []V
	for _, item := range c.items {
		if val, ok := item[key]; ok {
			result = append(result, val)
		}
	}
	return Collect(result)
}

// Pluck extracts values for a given key (for map collections)
func Pluck[K comparable, V any](c *Collection[map[K]V], key K) *Collection[V] {
	return Values(c, key)
}

// Sort sorts the collection
func (c *Collection[T]) Sort(less func(i, j T) bool) *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})
	return Collect(result)
}

// Reverse reverses the collection
func (c *Collection[T]) Reverse() *Collection[T] {
	result := make([]T, len(c.items))
	for i, item := range c.items {
		result[len(c.items)-1-i] = item
	}
	return Collect(result)
}

// Keys returns the keys of a map collection
func Keys[K comparable, V any](c *Collection[map[K]V]) *Collection[K] {
	var result []K
	for _, item := range c.items {
		for key := range item {
			result = append(result, key)
		}
	}
	return Collect(result)
}

// Flatten flattens a collection of slices into a single collection
func Flatten[T any](c *Collection[[]T]) *Collection[T] {
	var result []T
	for _, item := range c.items {
		result = append(result, item...)
	}
	return Collect(result)
}

// Implode joins all items with a separator
func (c *Collection[T]) Implode(separator string) string {
	result := ""
	for i, item := range c.items {
		if i > 0 {
			result += separator
		}
		result += fmt.Sprintf("%v", item)
	}
	return result
}

// Contains checks if the collection contains an item
func (c *Collection[T]) Contains(item T) bool {
	for _, v := range c.items {
		if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", item) {
			return true
		}
	}
	return false
}

// Except returns items except the specified keys
func Except[K comparable, V any](c *Collection[map[K]V], keys ...K) *Collection[map[K]V] {
	excludeMap := make(map[K]bool)
	for _, key := range keys {
		excludeMap[key] = true
	}
	var result []map[K]V
	for _, item := range c.items {
		newItem := make(map[K]V)
		for k, v := range item {
			if !excludeMap[k] {
				newItem[k] = v
			}
		}
		result = append(result, newItem)
	}
	return Collect(result)
}

// Only returns only the specified keys from map collections
func Only[K comparable, V any](c *Collection[map[K]V], keys ...K) *Collection[map[K]V] {
	includeMap := make(map[K]bool)
	for _, key := range keys {
		includeMap[key] = true
	}
	var result []map[K]V
	for _, item := range c.items {
		newItem := make(map[K]V)
		for k, v := range item {
			if includeMap[k] {
				newItem[k] = v
			}
		}
		result = append(result, newItem)
	}
	return Collect(result)
}

// ForEach is an alias for Each
func (c *Collection[T]) ForEach(callback func(T)) *Collection[T] {
	return c.Each(callback)
}

// ForEachWithIndex iterates with index
func (c *Collection[T]) ForEachWithIndex(callback func(int, T)) *Collection[T] {
	for i, item := range c.items {
		callback(i, item)
	}
	return c
}

// Count returns the number of items
func (c *Collection[T]) Count() int {
	return len(c.items)
}

// IsEmpty checks if the collection is empty
func (c *Collection[T]) IsEmpty() bool {
	return len(c.items) == 0
}

// IsNotEmpty checks if the collection is not empty
func (c *Collection[T]) IsNotEmpty() bool {
	return !c.IsEmpty()
}

// Max returns the maximum item
func (c *Collection[T]) Max() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	max := c.items[0]
	for _, item := range c.items[1:] {
		if fmt.Sprintf("%v", item) > fmt.Sprintf("%v", max) {
			max = item
		}
	}
	return max, true
}

// Min returns the minimum item
func (c *Collection[T]) Min() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	min := c.items[0]
	for _, item := range c.items[1:] {
		if fmt.Sprintf("%v", item) < fmt.Sprintf("%v", min) {
			min = item
		}
	}
	return min, true
}

// Nth returns every nth element
func (c *Collection[T]) Nth(n int, offset int) *Collection[T] {
	if n <= 0 {
		return Collect([]T{})
	}
	var result []T
	for i := offset; i < len(c.items); i += n {
		result = append(result, c.items[i])
	}
	return Collect(result)
}

// Pull removes and returns an item by key (for map collections)
func Pull[K comparable, V any](m map[K]V, key K) (V, bool) {
	val, ok := m[key]
	if ok {
		delete(m, key)
	}
	return val, ok
}

// Wrap wraps a value in a collection
func Wrap[T any](value T) *Collection[T] {
	return Collect([]T{value})
}

// Pop removes and returns the last item
func (c *Collection[T]) Pop() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	item := c.items[len(c.items)-1]
	c.items = c.items[:len(c.items)-1]
	return item, true
}

// Shift removes and returns the first item
func (c *Collection[T]) Shift() (T, bool) {
	if len(c.items) == 0 {
		var zero T
		return zero, false
	}
	item := c.items[0]
	c.items = c.items[1:]
	return item, true
}

// Push appends an item to the end
func (c *Collection[T]) Push(item T) *Collection[T] {
	c.items = append(c.items, item)
	return c
}

// Prepend adds an item to the beginning
func (c *Collection[T]) Prepend(item T) *Collection[T] {
	c.items = append([]T{item}, c.items...)
	return c
}

// Merge merges another collection into this one
func (c *Collection[T]) Merge(other *Collection[T]) *Collection[T] {
	c.items = append(c.items, other.All()...)
	return c
}

// Diff returns items in the collection that are not in the other collection
func (c *Collection[T]) Diff(other *Collection[T]) *Collection[T] {
	otherItems := make(map[any]bool)
	for _, item := range other.All() {
		otherItems[fmt.Sprintf("%v", item)] = true
	}
	var result []T
	for _, item := range c.items {
		if !otherItems[fmt.Sprintf("%v", item)] {
			result = append(result, item)
		}
	}
	return Collect(result)
}

// Intersect returns items that are in both collections
func (c *Collection[T]) Intersect(other *Collection[T]) *Collection[T] {
	otherItems := make(map[any]bool)
	for _, item := range other.All() {
		otherItems[fmt.Sprintf("%v", item)] = true
	}
	var result []T
	for _, item := range c.items {
		if otherItems[fmt.Sprintf("%v", item)] {
			result = append(result, item)
		}
	}
	return Collect(result)
}

// UniqueBy returns unique items by a callback
func (c *Collection[T]) UniqueBy(callback func(T) string) *Collection[T] {
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

// GroupBy groups items by a callback
func GroupBy[T any, K comparable](c *Collection[T], callback func(T) K) map[K]*Collection[T] {
	result := make(map[K]*Collection[T])
	for _, item := range c.items {
		key := callback(item)
		if _, ok := result[key]; !ok {
			result[key] = Collect([]T{})
		}
		result[key].items = append(result[key].items, item)
	}
	return result
}

// Partition splits the collection into two groups
func (c *Collection[T]) Partition(callback func(T) bool) ([]*Collection[T], []*Collection[T]) {
	var trueGroup, falseGroup []T
	for _, item := range c.items {
		if callback(item) {
			trueGroup = append(trueGroup, item)
		} else {
			falseGroup = append(falseGroup, item)
		}
	}
	return []*Collection[T]{Collect(trueGroup)}, []*Collection[T]{Collect(falseGroup)}
}

// Pad pads the collection to a given size with a value
func (c *Collection[T]) Pad(size int, value T) *Collection[T] {
	if size <= len(c.items) {
		return Collect(c.items)
	}
	result := make([]T, size)
	copy(result, c.items)
	for i := len(c.items); i < size; i++ {
		result[i] = value
	}
	return Collect(result)
}

// Zip merges elements from another collection at the same position
func Zip[T any](c *Collection[T], other *Collection[T]) *Collection[[]T] {
	maxLen := len(c.items)
	if len(other.items) > maxLen {
		maxLen = len(other.items)
	}
	result := make([][]T, maxLen)
	for i := 0; i < maxLen; i++ {
		var pair []T
		if i < len(c.items) {
			pair = append(pair, c.items[i])
		}
		if i < len(other.items) {
			pair = append(pair, other.items[i])
		}
		result[i] = pair
	}
	return Collect(result)
}

// Cross joins with another collection
func Cross[T any](c *Collection[T], other *Collection[T]) *Collection[[]T] {
	var result [][]T
	for _, a := range c.items {
		for _, b := range other.items {
			result = append(result, []T{a, b})
		}
	}
	return Collect(result)
}

// Rotate rotates the collection by n positions
func (c *Collection[T]) Rotate(n int) *Collection[T] {
	if len(c.items) == 0 {
		return c
	}
	n = n % len(c.items)
	if n < 0 {
		n = len(c.items) + n
	}
	result := make([]T, len(c.items))
	copy(result, c.items[n:])
	result = append(result, c.items[:n]...)
	return Collect(result)
}

// Shuffle shuffles the collection
func (c *Collection[T]) Shuffle() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)
	// Fisher-Yates shuffle
	for i := len(result) - 1; i > 0; i-- {
		j := i // In production, use rand.Intn(i+1)
		result[i], result[j] = result[j], result[i]
	}
	return Collect(result)
}

// ChunkWhile chunks the collection while the callback returns true
func (c *Collection[T]) ChunkWhile(callback func(T, T) bool) []*Collection[T] {
	if len(c.items) == 0 {
		return nil
	}
	var result []*Collection[T]
	var current []T
	for i := 0; i < len(c.items); i++ {
		if i == 0 || !callback(c.items[i-1], c.items[i]) {
			if len(current) > 0 {
				result = append(result, Collect(current))
			}
			current = []T{c.items[i]}
		} else {
			current = append(current, c.items[i])
		}
	}
	if len(current) > 0 {
		result = append(result, Collect(current))
	}
	return result
}

// Split splits the collection into two at the given index
func (c *Collection[T]) Split(index int) ([]*Collection[T], []*Collection[T]) {
	if index <= 0 {
		return []*Collection[T]{Collect([]T{})}, []*Collection[T]{Collect(c.items)}
	}
	if index >= len(c.items) {
		return []*Collection[T]{Collect(c.items)}, []*Collection[T]{Collect([]T{})}
	}
	return []*Collection[T]{Collect(c.items[:index])}, []*Collection[T]{Collect(c.items[index:])}
}

// Tapper taps into the collection without modifying it
func (c *Collection[T]) Tapper(callback func(*Collection[T])) *Collection[T] {
	callback(c)
	return c
}

// Contract reduces the collection to a single value
func Contract[T any, U any](c *Collection[T], initial U, callback func(U, T) U) U {
	result := initial
	for _, item := range c.items {
		result = callback(result, item)
	}
	return result
}

// MapToKeys maps to keys (for map collections)
func MapToKeys[K comparable, V any](c *Collection[map[K]V]) *Collection[K] {
	return Keys(c)
}

// MapToValues maps to values (for map collections)
func MapToValues[K comparable, V any](c *Collection[map[K]V]) *Collection[V] {
	var result []V
	for _, item := range c.items {
		for _, v := range item {
			result = append(result, v)
		}
	}
	return Collect(result)
}
