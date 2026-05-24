package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mechneerd/gow/routing"
	"github.com/mechneerd/gow/support/collection"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type contextKey string

const jsonBodyKey contextKey = "gow_json_body"

// Param gets a URL path parameter from the router context.
func Param(r *http.Request, key string) string {
	if params, ok := r.Context().Value(routing.ParamsKey).(map[string]string); ok {
		return params[key]
	}
	return ""
}

// Query gets a query string value.
func Query(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// Input gets an input value from JSON body, form body, or query string (in that priority order).
func Input(r *http.Request, key string) string {
	all := All(r)
	if val, ok := all[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		default:
			// Convert non-string JSON values to a string representation for generic Input
			if bs, err := json.Marshal(v); err == nil {
				return string(bs)
			}
		}
	}
	return ""
}

// Has checks if an input value exists and is not empty.
func Has(r *http.Request, key string) bool {
	return Input(r, key) != ""
}

// All returns a merged map of all input data (JSON wins over form, form wins over query).
func All(r *http.Request) map[string]any {
	all := make(map[string]any)

	// 1. Query parameters
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			all[k] = v[0]
		}
	}

	// 2. Form data
	if err := r.ParseForm(); err == nil {
		for k, v := range r.PostForm {
			if len(v) > 0 {
				all[k] = v[0]
			}
		}
	}

	// 3. JSON data
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		jsonMap := parseJSONBody(r)
		for k, v := range jsonMap {
			all[k] = v
		}
	}

	return all
}

// Only returns only the specified keys from all inputs.
func Only(r *http.Request, keys ...string) map[string]any {
	all := All(r)
	filtered := make(map[string]any)
	for _, key := range keys {
		if val, ok := all[key]; ok {
			filtered[key] = val
		}
	}
	return filtered
}

// parseJSONBody reads the body, caches it in context, and restores r.Body.
func parseJSONBody(r *http.Request) map[string]any {
	// Check if already parsed
	if cached, ok := r.Context().Value(jsonBodyKey).(map[string]any); ok {
		return cached
	}

	jsonMap := make(map[string]any)
	if r.Body == nil {
		return jsonMap
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil || len(raw) == 0 {
		return jsonMap
	}

	// Unmarshal and cache
	if err := json.Unmarshal(raw, &jsonMap); err == nil {
		ctx := context.WithValue(r.Context(), jsonBodyKey, jsonMap)
		*r = *r.WithContext(ctx)
	}

	// Restore body
	r.Body = io.NopCloser(bytes.NewReader(raw))

	return jsonMap
}

// Except returns all input except the given keys.
func Except(r *http.Request, keys ...string) map[string]any {
	all := All(r)
	for _, key := range keys {
		delete(all, key)
	}
	return all
}

// Boolean returns the input value as bool (accepts "1", "true", "on", "yes").
func Boolean(r *http.Request, key string) bool {
	val := strings.ToLower(Input(r, key))
	return val == "1" || val == "true" || val == "on" || val == "yes"
}

// Integer returns the input value as int (0 on failure).
func Integer(r *http.Request, key string) int {
	val := Input(r, key)
	i, _ := strconv.Atoi(val)
	return i
}

// Float returns the input value as float64 (0 on failure).
func Float(r *http.Request, key string) float64 {
	val := Input(r, key)
	f, _ := strconv.ParseFloat(val, 64)
	return f
}

// Collect returns selected (or all) input as a generic Collection.
func Collect(r *http.Request, keys ...string) *collection.Collection[any] {
	var data []any
	if len(keys) == 0 {
		for _, v := range All(r) {
			data = append(data, v)
		}
	} else {
		only := Only(r, keys...)
		for _, v := range only {
			data = append(data, v)
		}
	}
	return collection.Collect(data)
}

// ExpectsJson returns true if the client expects a JSON response (Accept or X-Requested-With).
func ExpectsJson(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json") ||
		strings.Contains(r.Header.Get("X-Requested-With"), "XMLHttpRequest")
}

// WantsJson is an alias for ExpectsJson (Laravel naming).
func WantsJson(r *http.Request) bool {
	return ExpectsJson(r)
}

// Accepts checks if the request accepts a given content type.
func Accepts(r *http.Request, contentType string) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, contentType)
}

// Old retrieves old input flashed from the previous request.
// It first checks the request context (if middleware stored it), then falls back to empty string.
// Full integration with session flash is done via middleware/session.go + session.Manager.Old().
func Old(r *http.Request, key string, defaultValue ...string) string {
	// Check if old input was stored in context by session middleware
	if oldInput, ok := r.Context().Value("old_input").(map[string]any); ok {
		if v, exists := oldInput[key]; exists {
			if s, ok := v.(string); ok {
				return s
			}
			return fmt.Sprintf("%v", v)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

