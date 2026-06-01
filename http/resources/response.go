package resources

import (
	"encoding/json"
	"net/http"
)

// JsonResponse represents a standard JSON API response
type JsonResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status"`
}

// PaginatedResponse represents a paginated JSON response
type PaginatedResponse struct {
	Data  []map[string]any `json:"data"`
	Links *Links           `json:"links"`
	Meta  *Meta            `json:"meta"`
}

// Links contains pagination links
type Links struct {
	Self  string `json:"self"`
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// Meta contains pagination metadata
type Meta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

// CollectionResponse wraps a collection of resources
type CollectionResponse struct {
	Data []map[string]any `json:"data"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Success sends a success response
func Success(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, JsonResponse{
		Data:   data,
		Status: "success",
	})
}

// Created sends a created response
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, JsonResponse{
		Data:   data,
		Status: "created",
	})
}

// Error sends an error response
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, JsonResponse{
		Message: message,
		Status:  "error",
	})
}

// Paginated sends a paginated response
func Paginated(w http.ResponseWriter, data []map[string]any, links *Links, meta *Meta) {
	JSON(w, http.StatusOK, PaginatedResponse{
		Data:  data,
		Links: links,
		Meta:  meta,
	})
}

// NewCollectionResponse creates a collection response from maps
func NewCollectionResponse(data []map[string]any) *CollectionResponse {
	return &CollectionResponse{
		Data: data,
	}
}

// When adds a field conditionally
func When(condition bool, key string, value any) map[string]any {
	if condition {
		return map[string]any{key: value}
	}
	return nil
}

// Merge merges multiple maps
func Merge(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
