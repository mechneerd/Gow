package request

import (
	"bytes"
	"context"
	"encoding/json"
	"gow/routing"
	"io"
	"net/http"
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
