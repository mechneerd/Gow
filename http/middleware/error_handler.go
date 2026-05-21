package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	gowhttp "gow/http"
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
					
					// Type assert for our custom HttpException
					if httpErr, ok := err.(gowhttp.HttpException); ok {
						status = httpErr.StatusCode
						message = httpErr.Message
					} else {
						// For other panics, preserve the original error message if in debug mode
						if debugMode {
							message = fmt.Sprintf("%v\n\n%s", err, string(debug.Stack()))
						}
					}

					// We could try to render a view here (e.g., view.Make("errors.500"))
					// For this middleware, we'll write the status and message.
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.WriteHeader(status)
					w.Write([]byte(fmt.Sprintf("%d | %s", status, message)))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
