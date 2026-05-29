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
	chain           http.Handler // built middleware chain, cached
	chainDirty      bool        // true when middlewares changed since last build
	shutdownHooks   []func()
	shutdownTimeout time.Duration
	mu              sync.Mutex
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
	k.middlewares = append(k.middlewares, mw)
	k.chainDirty = true
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

// ServeHTTP implements http.Handler. It uses the cached middleware pipeline.
func (k *Kernel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if k.chainDirty || k.chain == nil {
		k.buildChain()
	}
	k.chain.ServeHTTP(w, r)
}

