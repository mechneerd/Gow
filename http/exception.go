package http

import (
	"github.com/mechneerd/gow/http/exception"
)

// HttpException is an alias for exception.HttpException for backward compatibility.
// Use exception.HttpException directly for new code.
type HttpException = exception.HttpException

// Abort panics with an HttpException to stop execution and render an error response.
func Abort(status int, message string) {
	panic(exception.HttpException{
		Code:    status,
		Message: message,
	})
}

