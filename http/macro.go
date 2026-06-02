package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// ResponseMacro provides a way to add custom response methods.
type ResponseMacro struct {
	macros map[string]ResponseMacroFunc
	mu     sync.RWMutex
}

// ResponseMacroFunc is a function that creates an HTTP response.
type ResponseMacroFunc func(args ...any) *MacroResponse

// MacroResponse represents an HTTP response for macros.
type MacroResponse struct {
	StatusCode int
	Body       any
	Headers    map[string]string
}

// NewResponseMacro creates a new ResponseMacro.
func NewResponseMacro() *ResponseMacro {
	return &ResponseMacro{
		macros: make(map[string]ResponseMacroFunc),
	}
}

// Macro registers a new response macro.
func (rm *ResponseMacro) Macro(name string, fn ResponseMacroFunc) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.macros[name] = fn
}

// Call calls a registered macro.
func (rm *ResponseMacro) Call(name string, args ...any) (*MacroResponse, error) {
	rm.mu.RLock()
	fn, ok := rm.macros[name]
	rm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("macro '%s' not found", name)
	}
	return fn(args...), nil
}

// Has checks if a macro is registered.
func (rm *ResponseMacro) Has(name string) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	_, ok := rm.macros[name]
	return ok
}

// List returns all registered macro names.
func (rm *ResponseMacro) List() []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	names := make([]string, 0, len(rm.macros))
	for name := range rm.macros {
		names = append(names, name)
	}
	return names
}

// Global response macro instance
var globalResponseMacros = NewResponseMacro()

// Macro registers a global response macro.
func Macro(name string, fn ResponseMacroFunc) {
	globalResponseMacros.Macro(name, fn)
}

// CallMacro calls a global response macro.
func CallMacro(name string, args ...any) (*MacroResponse, error) {
	return globalResponseMacros.Call(name, args...)
}

// WriteMacroResponse writes a MacroResponse to an http.ResponseWriter.
func WriteMacroResponse(w http.ResponseWriter, resp *MacroResponse) {
	for key, value := range resp.Headers {
		w.Header().Set(key, value)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp.Body)
}

// Response helpers for common macros

// ApiResponse creates a JSON API response wrapper.
func ApiResponse(data any, meta ...map[string]any) *MacroResponse {
	resp := map[string]any{
		"data": data,
	}
	if len(meta) > 0 {
		resp["meta"] = meta[0]
	}
	return &MacroResponse{StatusCode: 200, Body: resp}
}

// ErrorResponse creates an error response.
func ErrorResponse(message string, code int, errors ...map[string]any) *MacroResponse {
	resp := map[string]any{
		"error": map[string]any{
			"message": message,
			"code":    code,
		},
	}
	if len(errors) > 0 {
		resp["error"].(map[string]any)["errors"] = errors[0]
	}
	return &MacroResponse{StatusCode: code, Body: resp}
}

// PaginatedResponse creates a paginated response.
func PaginatedResponse(data any, total int, perPage int, currentPage int) *MacroResponse {
	totalPages := (total + perPage - 1) / perPage
	return &MacroResponse{
		StatusCode: 200,
		Body: map[string]any{
			"data": data,
			"meta": map[string]any{
				"total":        total,
				"per_page":     perPage,
				"current_page": currentPage,
				"last_page":    totalPages,
			},
			"links": map[string]any{
				"first": fmt.Sprintf("?page=1&per_page=%d", perPage),
				"last":  fmt.Sprintf("?page=%d&per_page=%d", totalPages, perPage),
				"prev":  fmt.Sprintf("?page=%d&per_page=%d", currentPage-1, perPage),
				"next":  fmt.Sprintf("?page=%d&per_page=%d", currentPage+1, perPage),
			},
		},
	}
}

// CreatedResponse creates a 201 Created response.
func CreatedResponse(data any) *MacroResponse {
	return &MacroResponse{
		StatusCode: 201,
		Body:       map[string]any{"data": data},
	}
}

// NoContentResponse creates a 204 No Content response.
func NoContentResponse() *MacroResponse {
	return &MacroResponse{StatusCode: 204}
}

// AcceptedResponse creates a 202 Accepted response.
func AcceptedResponse(data any, message ...string) *MacroResponse {
	resp := map[string]any{"data": data}
	if len(message) > 0 {
		resp["message"] = message[0]
	}
	return &MacroResponse{StatusCode: 202, Body: resp}
}

// BadRequestResponse creates a 400 Bad Request response.
func BadRequestResponse(message string) *MacroResponse {
	return ErrorResponse(message, 400)
}

// UnauthorizedResponse creates a 401 Unauthorized response.
func UnauthorizedResponse(message ...string) *MacroResponse {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}
	return ErrorResponse(msg, 401)
}

// ForbiddenResponse creates a 403 Forbidden response.
func ForbiddenResponse(message ...string) *MacroResponse {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	return ErrorResponse(msg, 403)
}

// NotFoundResponse creates a 404 Not Found response.
func NotFoundResponse(message ...string) *MacroResponse {
	msg := "Not Found"
	if len(message) > 0 {
		msg = message[0]
	}
	return ErrorResponse(msg, 404)
}

// ValidationFailedResponse creates a 422 Validation Failed response.
func ValidationFailedResponse(errors map[string][]string) *MacroResponse {
	// Convert map[string][]string to map[string]any
	errorMap := make(map[string]any)
	for k, v := range errors {
		errorMap[k] = v
	}
	return ErrorResponse("Validation failed", 422, errorMap)
}

// TooManyRequestsResponse creates a 429 Too Many Requests response.
func TooManyRequestsResponse(retryAfter ...int) *MacroResponse {
	resp := ErrorResponse("Too Many Requests", 429)
	if len(retryAfter) > 0 {
		resp.Headers = map[string]string{
			"Retry-After": fmt.Sprintf("%d", retryAfter[0]),
		}
	}
	return resp
}
