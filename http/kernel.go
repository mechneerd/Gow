package http

import (
	"gow/foundation"
	"gow/routing"
	"net/http"
)

// Kernel handles incoming HTTP requests.
type Kernel struct {
	app         *foundation.Application
	router      *routing.Router
	middlewares []func(http.Handler) http.Handler
}

// NewKernel creates a new HTTP Kernel.
func NewKernel(app *foundation.Application, router *routing.Router) *Kernel {
	return &Kernel{
		app:    app,
		router: router,
	}
}

// Use adds global middleware to the kernel.
func (k *Kernel) Use(mw func(http.Handler) http.Handler) {
	k.middlewares = append(k.middlewares, mw)
}

// Handle handles an incoming request and returns a response.
// In Go, this typically means implementing the http.Handler interface.
func (k *Kernel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Build the middleware pipeline
	var handler http.Handler = k.router

	for i := len(k.middlewares) - 1; i >= 0; i-- {
		handler = k.middlewares[i](handler)
	}

	handler.ServeHTTP(w, r)
}
