package exception

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// HttpException represents a generic HTTP error.
type HttpException struct {
	Code    int
	Message string
	cause   error
}

func (e *HttpException) Error() string {
	return e.Message
}

func (e *HttpException) Unwrap() error {
	return e.cause
}

// Render outputs the exception as either JSON or HTML based on the request's Accept or Content-Type headers.
func (e *HttpException) Render(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(accept, "application/json") || strings.Contains(contentType, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Code)
		json.NewEncoder(w).Encode(map[string]any{
			"error":   true,
			"message": e.Message,
			"code":    e.Code,
		})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(e.Code)
	// Generic HTML fallback for frameworks (can be overridden by a custom error page later)
	fmt.Fprintf(w, "<html><head><title>%d %s</title></head><body><h1>%d %s</h1></body></html>", e.Code, e.Message, e.Code, e.Message)
}

// Constructors

// NotFound returns a 404 HttpException.
func NotFound(msg string) *HttpException {
	return &HttpException{Code: http.StatusNotFound, Message: msg}
}

// BadRequest returns a 400 HttpException.
func BadRequest(msg string) *HttpException {
	return &HttpException{Code: http.StatusBadRequest, Message: msg}
}

// Unauthorized returns a 401 HttpException.
func Unauthorized(msg string) *HttpException {
	return &HttpException{Code: http.StatusUnauthorized, Message: msg}
}

// Forbidden returns a 403 HttpException.
func Forbidden(msg string) *HttpException {
	return &HttpException{Code: http.StatusForbidden, Message: msg}
}

// InternalServerError returns a 500 HttpException.
func InternalServerError(msg string) *HttpException {
	return &HttpException{Code: http.StatusInternalServerError, Message: msg}
}

// Wrap creates an HttpException wrapping an underlying error.
func Wrap(code int, msg string, cause error) *HttpException {
	return &HttpException{
		Code:    code,
		Message: msg,
		cause:   cause,
	}
}
