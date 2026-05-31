package http

import (
	"context"
	"fmt"
	"github.com/mechneerd/gow/foundation"
	"github.com/mechneerd/gow/routing"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Kernel handles incoming HTTP requests.
type Kernel struct {
	app             *foundation.Application
	router          *routing.Router
	middlewares     []func(http.Handler) http.Handler
	terminateHooks  []func(http.ResponseWriter, *http.Request) // per-request terminate callbacks
	chain           http.Handler
	chainDirty      bool
	shutdownHooks   []func()
	shutdownTimeout time.Duration
	mu              sync.RWMutex
}

// TerminableMiddleware is an optional interface that middleware can implement
// to run code after the response has been sent to the client.
type TerminableMiddleware interface {
	Terminate(w http.ResponseWriter, r *http.Request)
}

// NewKernel creates a new HTTP Kernel.
func NewKernel(app *foundation.Application, router *routing.Router) *Kernel {
	return &Kernel{
		app:             app,
		router:          router,
		shutdownTimeout: 30 * time.Second,
	}
}

// SetShutdownTimeout configures the graceful shutdown timeout.
// Default is 30 seconds.
func (k *Kernel) SetShutdownTimeout(d time.Duration) {
	k.shutdownTimeout = d
}

// OnShutdown registers a function to be called during graceful shutdown.
// Hooks are called in the order they were registered, after the HTTP server
// stops accepting new connections but before the process exits.
// Use this to close database connections, flush logs, drain queue workers, etc.
func (k *Kernel) OnShutdown(fn func()) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.shutdownHooks = append(k.shutdownHooks, fn)
}

// Serve starts the HTTP server and blocks until a shutdown signal is received.
// It listens for SIGINT and SIGTERM, then gracefully shuts down the server,
// allowing in-flight requests to complete within the configured timeout.
func (k *Kernel) Serve(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: k,
	}

	// Channel for server errors
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("[GoW] Server started on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Printf("[GoW] Received signal %v, shutting down gracefully...", sig)
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), k.shutdownTimeout)
	defer cancel()

	// Shutdown the HTTP server (stops accepting new connections, waits for in-flight)
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[GoW] HTTP server shutdown error: %v", err)
	}

	// Run shutdown hooks
	k.mu.Lock()
	hooks := make([]func(), len(k.shutdownHooks))
	copy(hooks, k.shutdownHooks)
	k.mu.Unlock()

	for i, hook := range hooks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[GoW] Shutdown hook %d panicked: %v", i, r)
				}
			}()
			hook()
		}()
	}

	log.Println("[GoW] Server stopped.")
	return nil
}

// Use adds global middleware to the kernel.
func (k *Kernel) Use(mw func(http.Handler) http.Handler) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.middlewares = append(k.middlewares, mw)
	k.chainDirty = true
}

// Terminate registers a callback to run after every response is sent.
// Use this for post-response cleanup: flushing logs, closing connections, etc.
func (k *Kernel) Terminate(fn func(http.ResponseWriter, *http.Request)) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.terminateHooks = append(k.terminateHooks, fn)
}

// buildChain builds the middleware pipeline once and caches it.
func (k *Kernel) buildChain() {
	var handler http.Handler = k.router
	for i := len(k.middlewares) - 1; i >= 0; i-- {
		handler = k.middlewares[i](handler)
	}
	k.chain = handler
	k.chainDirty = false
}

// ServeHTTP implements http.Handler. It uses the cached middleware pipeline
// and runs terminate hooks after the response is sent.
func (k *Kernel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k.mu.RLock()
	dirty := k.chainDirty
	k.mu.RUnlock()

	if dirty || k.chain == nil {
		k.mu.Lock()
		if k.chainDirty || k.chain == nil {
			k.buildChain()
		}
		k.mu.Unlock()
	}

	// Use a response wrapper to detect when the response is complete
	rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	k.chain.ServeHTTP(rw, r)

	// Run terminate hooks
	k.mu.RLock()
	hooks := make([]func(http.ResponseWriter, *http.Request), len(k.terminateHooks))
	copy(hooks, k.terminateHooks)
	k.mu.RUnlock()

	for _, hook := range hooks {
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					log.Printf("[GoW] Terminate hook panicked: %v", rec)
				}
			}()
			hook(w, r)
		}()
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

