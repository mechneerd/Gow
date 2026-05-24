package exception

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
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

	// Professional debug error page (Whoops-style) when in development
	stack := string(debug.Stack())
	debugHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>%d %s</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #1e1e1e; color: #ddd; padding: 20px; }
        h1 { color: #f55; }
        pre { background: #2d2d2d; padding: 15px; border-radius: 6px; overflow: auto; font-size: 13px; }
        .stack { color: #8f8; }
        .info { color: #8af; }
    </style>
</head>
<body>
    <h1>%d %s</h1>
    <p class="info">GoW Framework - Debug Error Page</p>
    <h3>Message</h3>
    <pre>%s</pre>
    <h3>Stack Trace</h3>
    <pre class="stack">%s</pre>
</body>
</html>`, e.Code, e.Message, e.Code, e.Message, e.Message, stack)

	// Simple production fallback (no stack)
	prodHTML := fmt.Sprintf("<html><head><title>%d %s</title></head><body><h1>%d %s</h1><p>Something went wrong.</p></body></html>", e.Code, e.Message, e.Code, e.Message)

	// For now always show debug (in real app check config app.debug)
	fmt.Fprint(w, debugHTML)
	_ = prodHTML // keep for future conditional
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

