package exception

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

// Renderable is an interface for exceptions that can render themselves.
type Renderable interface {
	error
	Render(w http.ResponseWriter, r *http.Request)
}

// Reportable is an interface for exceptions that can report themselves.
type Reportable interface {
	error
	Report() []string // Returns additional context for logging
}

// HttpException represents a generic HTTP error.
type HttpException struct {
	Code     int
	Message  string
	cause    error
	headers  map[string]string
	reporter func() []string
}

func (e *HttpException) Error() string {
	return e.Message
}

func (e *HttpException) Unwrap() error {
	return e.cause
}

// WithHeader adds a response header to the exception.
func (e *HttpException) WithHeader(key, value string) *HttpException {
	if e.headers == nil {
		e.headers = make(map[string]string)
	}
	e.headers[key] = value
	return e
}

// WithReport adds a reporter function for custom log context.
func (e *HttpException) WithReport(fn func() []string) *HttpException {
	e.reporter = fn
	return e
}

// Report returns additional context for logging.
func (e *HttpException) Report() []string {
	if e.reporter != nil {
		return e.reporter()
	}
	return []string{fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)}
}

// Render outputs the exception as either JSON or HTML based on the request's Accept or Content-Type headers.
func (e *HttpException) Render(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	contentType := r.Header.Get("Content-Type")

	// Set custom headers
	for k, v := range e.headers {
		w.Header().Set(k, v)
	}

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

	stack := html.EscapeString(string(debug.Stack()))
	safeMessage := html.EscapeString(e.Message)
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
</html>`, e.Code, safeMessage, e.Code, safeMessage, safeMessage, stack)

	prodHTML := fmt.Sprintf("<html><head><title>%d %s</title></head><body><h1>%d %s</h1><p>Something went wrong.</p></body></html>", e.Code, safeMessage, e.Code, safeMessage)

	if os.Getenv("APP_DEBUG") == "true" {
		fmt.Fprint(w, debugHTML)
	} else {
		fmt.Fprint(w, prodHTML)
	}
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

// Conflict returns a 409 HttpException.
func Conflict(msg string) *HttpException {
	return &HttpException{Code: http.StatusConflict, Message: msg}
}

// Gone returns a 410 HttpException.
func Gone(msg string) *HttpException {
	return &HttpException{Code: http.StatusGone, Message: msg}
}

// UnprocessableEntity returns a 422 HttpException.
func UnprocessableEntity(msg string) *HttpException {
	return &HttpException{Code: http.StatusUnprocessableEntity, Message: msg}
}

// TooManyRequests returns a 429 HttpException.
func TooManyRequests(msg string) *HttpException {
	return &HttpException{Code: http.StatusTooManyRequests, Message: msg}
}

// InternalServerError returns a 500 HttpException.
func InternalServerError(msg string) *HttpException {
	return &HttpException{Code: http.StatusInternalServerError, Message: msg}
}

// NotImplemented returns a 501 HttpException.
func NotImplemented(msg string) *HttpException {
	return &HttpException{Code: http.StatusNotImplemented, Message: msg}
}

// BadGateway returns a 502 HttpException.
func BadGateway(msg string) *HttpException {
	return &HttpException{Code: http.StatusBadGateway, Message: msg}
}

// ServiceUnavailable returns a 503 HttpException.
func ServiceUnavailable(msg string) *HttpException {
	return &HttpException{Code: http.StatusServiceUnavailable, Message: msg}
}

// GatewayTimeout returns a 504 HttpException.
func GatewayTimeout(msg string) *HttpException {
	return &HttpException{Code: http.StatusGatewayTimeout, Message: msg}
}

// Wrap creates an HttpException wrapping an underlying error.
func Wrap(code int, msg string, cause error) *HttpException {
	return &HttpException{
		Code:    code,
		Message: msg,
		cause:   cause,
	}
}

// WrapWithReport creates an HttpException wrapping an error with a reporter.
func WrapWithReport(code int, msg string, cause error, reporter func() []string) *HttpException {
	return &HttpException{
		Code:     code,
		Message:  msg,
		cause:    cause,
		reporter: reporter,
	}
}

// Global exception handler middleware
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch t := rec.(type) {
				case error:
					err = t
				case string:
					err = fmt.Errorf("%s", t)
				default:
					err = fmt.Errorf("%v", t)
				}

				// Check if the error is renderable
				if renderable, ok := err.(Renderable); ok {
					renderable.Render(w, r)
					return
				}

				// Default to 500
				e := InternalServerError("Internal Server Error")
				e.Render(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RenderError is a convenience function to render any error as an HTTP response.
func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	if renderable, ok := err.(Renderable); ok {
		renderable.Render(w, r)
		return
	}
	e := Wrap(http.StatusInternalServerError, err.Error(), err)
	e.Render(w, r)
}

