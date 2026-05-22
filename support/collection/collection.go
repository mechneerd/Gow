package collection

// Collection provides a fluent, generic wrapper for working with arrays of data.
type Collection[T any] struct {
	items []T
}

// Collect creates a new Collection instance from a slice.
func Collect[T any](items []T) *Collection[T] {
	return &Collection[T]{
		items: items,
	}
}

// All returns the underlying slice of items.
func (c *Collection[T]) All() []T {
	return c.items
}

// Map applies a callback to each item and returns a new collection of the same type.
func (c *Collection[T]) Map(callback func(T) T) *Collection[T] {
	result := make([]T, len(c.items))
	for i, item := range c.items {
		result[i] = callback(item)
	}
	return Collect(result)
}

// Filter filters the collection by a given callback.
func (c *Collection[T]) Filter(callback func(T) bool) *Collection[T] {
	var result []T
	for _, item := range c.items {
		if callback(item) {
			result = append(result, item)
		}
	}
	return Collect(result)
}

// Reduce reduces the collection to a single value.
func (c *Collection[T]) Reduce(callback func(accumulator any, item T) any, initial any) any {
	result := initial
	for _, item := range c.items {
		result = callback(result, item)
	}
	return result
}

// Each iterates over the items in the collection and passes each to a callback.
func (c *Collection[T]) Each(callback func(T)) *Collection[T] {
	for _, item := range c.items {
		callback(item)
	}
	return c
}

// Chunk breaks the collection into multiple, smaller collections of a given size.
func (c *Collection[T]) Chunk(size int) []*Collection[T] {
	var result []*Collection[T]
	if size <= 0 {
		return result
	}

	for i := 0; i < len(c.items); i += size {
		end := i + size
		if end > len(c.items) {
			end = len(c.items)
		}
		result = append(result, Collect(c.items[i:end]))
	}
	return result
}
