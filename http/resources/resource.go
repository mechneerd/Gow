package resources

import (
	"encoding/json"
	"net/http"
)

// Resource represents an API transformer that shapes data.
type Resource interface {
	// ToArray transforms the underlying model into a map for JSON serialization.
	ToArray() map[string]any
}

// JsonResource base struct to embed for resources.
type JsonResource struct {
	Resource any
}

// Collection transforms a slice of models into a slice of resource maps.
func Collection(items []any, transformer func(any) Resource) []map[string]any {
	var result []map[string]any
	for _, item := range items {
		result = append(result, transformer(item).ToArray())
	}
	return result
}

// Respond writes the resource out as JSON to the response writer.
func Respond(w http.ResponseWriter, r Resource, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	// Default envelope
	payload := map[string]any{
		"data": r.ToArray(),
	}
	
	json.NewEncoder(w).Encode(payload)
}

// RespondCollection writes a collection out as JSON.
func RespondCollection(w http.ResponseWriter, items []any, transformer func(any) Resource, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	payload := map[string]any{
		"data": Collection(items, transformer),
	}
	
	json.NewEncoder(w).Encode(payload)
}
