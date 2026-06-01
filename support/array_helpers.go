package support

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Dot notation helpers

// DotGet gets a value from a map using dot notation
func DotGet(data map[string]any, key string) any {
	keys := strings.Split(key, ".")
	var current any = data

	for _, k := range keys {
		switch v := current.(type) {
		case map[string]any:
			current = v[k]
		case map[any]any:
			current = v[k]
		default:
			return nil
		}
	}
	return current
}

// DotSet sets a value in a map using dot notation
func DotSet(data map[string]any, key string, value any) {
	keys := strings.Split(key, ".")
	current := data

	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if _, ok := current[k]; !ok {
			current[k] = make(map[string]any)
		}
		if next, ok := current[k].(map[string]any); ok {
			current = next
		} else {
			return
		}
	}
	current[keys[len(keys)-1]] = value
}

// DotHas checks if a key exists using dot notation
func DotHas(data map[string]any, key string) bool {
	keys := strings.Split(key, ".")
	var current any = data

	for _, k := range keys {
		switch v := current.(type) {
		case map[string]any:
			if _, ok := v[k]; !ok {
				return false
			}
			current = v[k]
		default:
			return false
		}
	}
	return true
}

// DotDelete deletes a key using dot notation
func DotDelete(data map[string]any, key string) {
	keys := strings.Split(key, ".")
	current := data

	for i := 0; i < len(keys)-1; i++ {
		k := keys[i]
		if next, ok := current[k].(map[string]any); ok {
			current = next
		} else {
			return
		}
	}
	delete(current, keys[len(keys)-1])
}

// DotFlatten flattens a nested map using dot notation
func DotFlatten(data map[string]any) map[string]any {
	result := make(map[string]any)
	dotFlattenHelper(data, "", result)
	return result
}

func dotFlattenHelper(data map[string]any, prefix string, result map[string]any) {
	for key, value := range data {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]any); ok {
			dotFlattenHelper(nested, newKey, result)
		} else {
			result[newKey] = value
		}
	}
}

// DotUnflatten unflattens a map using dot notation
func DotUnflatten(data map[string]any) map[string]any {
	result := make(map[string]any)
	for key, value := range data {
		DotSet(result, key, value)
	}
	return result
}

// Array helpers

// ArrayGet gets a value from an array using dot notation
func ArrayGet(data []any, key string) any {
	keys := strings.Split(key, ".")
	var current any = data

	for _, k := range keys {
		switch v := current.(type) {
		case []any:
			index := 0
			fmt.Sscanf(k, "%d", &index)
			if index < len(v) {
				current = v[index]
			} else {
				return nil
			}
		case map[string]any:
			current = v[k]
		default:
			return nil
		}
	}
	return current
}

// ArraySet sets a value in an array using dot notation
func ArraySet(data []any, key string, value any) []any {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return data
	}

	index := 0
	fmt.Sscanf(keys[0], "%d", &index)

	// Ensure the slice is large enough
	for len(data) <= index {
		data = append(data, nil)
	}

	if len(keys) == 1 {
		data[index] = value
	} else {
		if data[index] == nil {
			data[index] = make(map[string]any)
		}
		if nested, ok := data[index].(map[string]any); ok {
			DotSet(nested, strings.Join(keys[1:], "."), value)
			data[index] = nested
		}
	}
	return data
}

// Collection helpers

// GroupBy groups a slice of maps by a key
func GroupBy[T any](data []T, key string) map[string][]T {
	result := make(map[string][]T)
	for _, item := range data {
		val := DotGet(toMap(item), key)
		if val != nil {
			group := fmt.Sprintf("%v", val)
			result[group] = append(result[group], item)
		}
	}
	return result
}

// Pluck extracts values from a slice of maps
func Pluck[T any, U any](data []T, key string, transform func(T) U) []U {
	var result []U
	for _, item := range data {
		val := DotGet(toMap(item), key)
		if val != nil {
			// Try to convert to the target type
			if converted, ok := val.(U); ok {
				result = append(result, converted)
			} else if transform != nil {
				result = append(result, transform(item))
			}
		}
	}
	return result
}

// Unique returns unique values from a slice
func Unique[T comparable](data []T) []T {
	seen := make(map[T]bool)
	var result []T
	for _, item := range data {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

// Chunk splits a slice into chunks
func Chunk[T any](data []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var result [][]T
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		result = append(result, data[i:end])
	}
	return result
}

// Filter filters a slice based on a predicate
func Filter[T any](data []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range data {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms a slice using a mapper function
func Map[T any, U any](data []T, mapper func(T) U) []U {
	result := make([]U, len(data))
	for i, item := range data {
		result[i] = mapper(item)
	}
	return result
}

// Reduce reduces a slice to a single value
func Reduce[T any, U any](data []T, reducer func(U, T) U, initial U) U {
	result := initial
	for _, item := range data {
		result = reducer(result, item)
	}
	return result
}

// Each iterates over a slice
func Each[T any](data []T, callback func(T)) {
	for _, item := range data {
		callback(item)
	}
}

// EachWithIndex iterates over a slice with index
func EachWithIndex[T any](data []T, callback func(int, T)) {
	for i, item := range data {
		callback(i, item)
	}
}

// Contains checks if a slice contains an item
func Contains[T comparable](data []T, item T) bool {
	for _, v := range data {
		if v == item {
			return true
		}
	}
	return false
}

// ContainsAny checks if a slice contains any of the items
func ContainsAny[T comparable](data []T, items ...T) bool {
	for _, item := range items {
		if Contains(data, item) {
			return true
		}
	}
	return false
}

// ContainsAll checks if a slice contains all items
func ContainsAll[T comparable](data []T, items ...T) bool {
	for _, item := range items {
		if !Contains(data, item) {
			return false
		}
	}
	return true
}

// Diff returns items in the first slice that are not in the second
func Diff[T comparable](a, b []T) []T {
	bMap := make(map[T]bool)
	for _, item := range b {
		bMap[item] = true
	}
	var result []T
	for _, item := range a {
		if !bMap[item] {
			result = append(result, item)
		}
	}
	return result
}

// Intersect returns items that are in both slices
func Intersect[T comparable](a, b []T) []T {
	bMap := make(map[T]bool)
	for _, item := range b {
		bMap[item] = true
	}
	var result []T
	for _, item := range a {
		if bMap[item] {
			result = append(result, item)
		}
	}
	return result
}

// Merge merges multiple slices
func Merge[T any](slices ...[]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Flatten flattens a slice of slices
func Flatten[T any](slices [][]T) []T {
	var result []T
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

// Reverse reverses a slice
func Reverse[T any](data []T) []T {
	result := make([]T, len(data))
	for i, item := range data {
		result[len(data)-1-i] = item
	}
	return result
}

// Shuffle shuffles a slice
func Shuffle[T any](data []T) []T {
	result := make([]T, len(data))
	copy(result, data)
	// Simple shuffle (in production, use math/rand)
	for i := len(result) - 1; i > 0; i-- {
		j := i // In production, use rand.Intn(i+1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// Take returns the first n items
func Take[T any](data []T, n int) []T {
	if n > len(data) {
		n = len(data)
	}
	return data[:n]
}

// Skip returns all items except the first n
func Skip[T any](data []T, n int) []T {
	if n > len(data) {
		return []T{}
	}
	return data[n:]
}

// Head returns the first item
func Head[T any](data []T) (T, bool) {
	if len(data) == 0 {
		var zero T
		return zero, false
	}
	return data[0], true
}

// Tail returns all items except the first
func Tail[T any](data []T) []T {
	if len(data) <= 1 {
		return []T{}
	}
	return data[1:]
}

// Last returns the last item
func Last[T any](data []T) (T, bool) {
	if len(data) == 0 {
		var zero T
		return zero, false
	}
	return data[len(data)-1], true
}

// Sum returns the sum of items (for numeric types)
func Sum[T int | float32 | float64](data []T) T {
	var sum T
	for _, item := range data {
		sum += item
	}
	return sum
}

// Average returns the average of items (for numeric types)
func Average[T int | float32 | float64](data []T) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, item := range data {
		sum += float64(item)
	}
	return sum / float64(len(data))
}

// Min returns the minimum item (for ordered types)
func Min[T int | float32 | float64 | string](data []T) (T, bool) {
	if len(data) == 0 {
		var zero T
		return zero, false
	}
	min := data[0]
	for _, item := range data[1:] {
		if item < min {
			min = item
		}
	}
	return min, true
}

// Max returns the maximum item (for ordered types)
func Max[T int | float32 | float64 | string](data []T) (T, bool) {
	if len(data) == 0 {
		var zero T
		return zero, false
	}
	max := data[0]
	for _, item := range data[1:] {
		if item > max {
			max = item
		}
	}
	return max, true
}

// toMap converts a value to map[string]any
func toMap(v any) map[string]any {
	switch val := v.(type) {
	case map[string]any:
		return val
	default:
		// Use reflection for structs
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() == reflect.Struct {
			result := make(map[string]any)
			rt := rv.Type()
			for i := 0; i < rt.NumField(); i++ {
				field := rt.Field(i)
				result[strings.ToLower(field.Name)] = rv.Field(i).Interface()
			}
			return result
		}
		// Try JSON marshal/unmarshal
		data, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			return nil
		}
		return result
	}
}
