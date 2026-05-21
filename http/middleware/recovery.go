package middleware

import (
	gowhttp "gow/http"
	"gow/http/response"
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery returns a middleware that recovers from panics.
// If the panic is an HttpException, it returns the appropriate HTTP status.
// Otherwise, it returns a 500 Internal Server Error.
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Check if it's our custom HttpException
					if httpErr, ok := err.(gowhttp.HttpException); ok {
						response.Error(w, httpErr.StatusCode, httpErr.Message)
						return
					}

					// Otherwise, log the stack trace and return 500
					log.Printf("panic recovered: %v\n%s", err, debug.Stack())
					response.Error(w, http.StatusInternalServerError, "Internal Server Error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
