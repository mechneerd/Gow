package http

import "fmt"

// HttpException represents an HTTP error that can be thrown via panic
// and caught by the recovery middleware.
type HttpException struct {
	StatusCode int
	Message    string
}

// Error implements the error interface.
func (e HttpException) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// Abort panics with an HttpException to stop execution and render an error response.
func Abort(status int, message string) {
	panic(HttpException{
		StatusCode: status,
		Message:    message,
	})
}

