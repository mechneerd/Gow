package resources

import (
	"encoding/json"
	"net/http"
)

// ResourceResponse represents the standard JSON envelope.
type ResourceResponse struct {
	Data  any            `json:"data"`
	Meta  map[string]any `json:"meta,omitempty"`
	Links map[string]any `json:"links,omitempty"`
}

// Resource defines how a model should be transformed into JSON.
type Resource[T any] interface {
	ToMap(req *http.Request, model T) map[string]any
}

// Transform applies a resource transformer to a single model.
func Transform[T any](req *http.Request, model T, transformer Resource[T]) map[string]any {
	return transformer.ToMap(req, model)
}

// Collection applies a resource transformer to a slice of models.
func Collection[T any](req *http.Request, models []T, transformer Resource[T]) []map[string]any {
	var result []map[string]any
	for _, model := range models {
		result = append(result, transformer.ToMap(req, model))
	}
	return result
}

// Respond writes a single resource out as JSON to the response writer inside a standard envelope.
func Respond[T any](w http.ResponseWriter, req *http.Request, model T, transformer Resource[T], status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	payload := ResourceResponse{
		Data: Transform(req, model, transformer),
	}

	json.NewEncoder(w).Encode(payload)
}

// RespondCollection writes a collection out as JSON inside a standard envelope.
func RespondCollection[T any](w http.ResponseWriter, req *http.Request, models []T, transformer Resource[T], status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	payload := ResourceResponse{
		Data: Collection(req, models, transformer),
	}

	json.NewEncoder(w).Encode(payload)
}

// RespondWithMeta writes a resource or collection with custom meta and links.
func RespondWithMeta(w http.ResponseWriter, data any, meta, links map[string]any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	payload := ResourceResponse{
		Data:  data,
		Meta:  meta,
		Links: links,
	}

	json.NewEncoder(w).Encode(payload)
}
