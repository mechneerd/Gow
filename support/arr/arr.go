package arr

import (
	"reflect"
)

// IsAssoc checks if a map is essentially associative (all maps in Go are associative, 
// but this can be used to distinguish between slices and maps or specific map structures).
// Since Go is strongly typed, this checks if the reflect kind is a Map.
func IsAssoc(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

// Wrap wraps a value in a slice if it's not already a slice or array.
func Wrap[T any](v any) []T {
	if v == nil {
		return []T{}
	}
	
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		var result []T
		for i := 0; i < val.Len(); i++ {
			if item, ok := val.Index(i).Interface().(T); ok {
				result = append(result, item)
			}
		}
		return result
	}
	
	if item, ok := v.(T); ok {
		return []T{item}
	}
	
	return []T{}
}

// Only returns a subset of a map with only the specified keys.
func Only(m map[string]any, keys ...string) map[string]any {
	result := make(map[string]any)
	for _, k := range keys {
		if val, ok := m[k]; ok {
			result[k] = val
		}
	}
	return result
}

// Except returns a subset of a map excluding the specified keys.
func Except(m map[string]any, keys ...string) map[string]any {
	result := make(map[string]any)
	
	excludes := make(map[string]bool)
	for _, k := range keys {
		excludes[k] = true
	}
	
	for k, v := range m {
		if !excludes[k] {
			result[k] = v
		}
	}
	return result
}

// Flatten flattens a multi-dimensional slice into a single-dimensional slice.
// Since Go uses typed slices, this function accepts any and returns []any.
func Flatten(v any) []any {
	var result []any
	val := reflect.ValueOf(v)
	
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return []any{v}
	}
	
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i).Interface()
		itemVal := reflect.ValueOf(item)
		
		if itemVal.Kind() == reflect.Slice || itemVal.Kind() == reflect.Array {
			result = append(result, Flatten(item)...)
		} else {
			result = append(result, item)
		}
	}
	
	return result
}

