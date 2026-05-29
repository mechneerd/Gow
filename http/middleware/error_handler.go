package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/mechneerd/gow/http/exception"
)

// ErrorHandler recovers from panics and converts them to HTTP responses.
// It handles HttpException natively, and converts other panics into 500s.
func ErrorHandler(debugMode bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					status := http.StatusInternalServerError
					message := "Internal Server Error"
					
					// Type assert for our custom HttpException (pointer type)
					if httpErr, ok := err.(*exception.HttpException); ok {
						status = httpErr.Code
						message = httpErr.Message
					} else if httpErr, ok := err.(exception.HttpException); ok {
						// Also check value type for backward compatibility
						status = httpErr.Code
						message = httpErr.Message
					} else {
						// For other panics, preserve the original error message if in debug mode
						if debugMode {
							message = fmt.Sprintf("%v\n\n%s", err, string(debug.Stack()))
						}
					}

					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.WriteHeader(status)
					w.Write([]byte(fmt.Sprintf("%d | %s", status, message)))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

