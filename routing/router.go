package routing

import (
	"context"
	"errors"
	"github.com/mechneerd/gow/http/exception"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// HandlerFunc is the GoW handler function signature.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Router represents the HTTP router.
type Router struct {
	mu                sync.RWMutex
	trees             map[string]*node
	namedRoutes       map[string]*Route
	allRoutes         []*Route // flat list for route:list and introspection
	middlewares       []func(http.Handler) http.Handler
	groupPrefix       string
	middlewareAliases map[string]func(http.Handler) http.Handler
	middlewareGroups  map[string][]func(http.Handler) http.Handler

	// Route model binders: paramName => resolver func(value string) (any, error)
	binders map[string]func(string) (any, error)
}

// Route represents a registered route.
type Route struct {
	Method      string
	Path        string
	Handler     HandlerFunc
	Middlewares []func(http.Handler) http.Handler
	Name        string
}

// NewRouter creates a new router instance.
func NewRouter() *Router {
	return &Router{
		trees:             make(map[string]*node),
		namedRoutes:       make(map[string]*Route),
		allRoutes:         make([]*Route, 0),
		middlewareAliases: make(map[string]func(http.Handler) http.Handler),
		middlewareGroups:  make(map[string][]func(http.Handler) http.Handler),
		binders:           make(map[string]func(string) (any, error)),
	}
}

// node represents a node in the radix tree.
type node struct {
	path       string
	wildcard   bool
	paramName  string
	handler    HandlerFunc
	route      *Route
	children   map[string]*node
}

// AddRoute adds a route to the router.
func (r *Router) AddRoute(method, path string, handler HandlerFunc) *Route {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullPath := r.groupPrefix + path
	if fullPath != "/" {
		fullPath = strings.TrimRight(fullPath, "/")
	}

	route := &Route{
		Method:      method,
		Path:        fullPath,
		Handler:     handler,
		Middlewares: append([]func(http.Handler) http.Handler{}, r.middlewares...),
	}

	r.allRoutes = append(r.allRoutes, route)

	if r.trees[method] == nil {
		r.trees[method] = &node{children: make(map[string]*node)}
	}

	r.insert(method, fullPath, route)
	return route
}

// SetName sets the name of the route and stores it for reverse generation.
func (r *Router) SetName(route *Route, name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	route.Name = name
	r.namedRoutes[name] = route
}

// Name chainable method on Route.
func (route *Route) SetRouteName(router *Router, name string) *Route {
	router.SetName(route, name)
	return route
}

func (r *Router) insert(method, path string, route *Route) {
	root := r.trees[method]
	
	if path == "/" {
		root.handler = route.Handler
		root.route = route
		return
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	current := root

	for _, part := range parts {
		isWildcard := strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")
		
		var nodeKey string
		if isWildcard {
			nodeKey = "{}" // generic key for wildcard in children map
		} else {
			nodeKey = part
		}

		if current.children[nodeKey] == nil {
			newNode := &node{
				path:     part,
				children: make(map[string]*node),
			}
			if isWildcard {
				newNode.wildcard = true
				newNode.paramName = part[1 : len(part)-1]
			}
			current.children[nodeKey] = newNode
		}
		current = current.children[nodeKey]
	}

	current.handler = route.Handler
	current.route = route
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	root := r.trees[req.Method]
	r.mu.RUnlock()

	if root == nil {
		http.NotFound(w, req)
		return
	}

	path := req.URL.Path
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}

	if path == "/" {
		if root.handler != nil {
			r.executeRoute(root.route, w, req, nil)
			return
		}
		http.NotFound(w, req)
		return
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	current := root
	params := make(map[string]string)

	for _, part := range parts {
		if next, ok := current.children[part]; ok {
			current = next
		} else if next, ok := current.children["{}"]; ok {
			current = next
			params[current.paramName] = part
		} else {
			http.NotFound(w, req)
			return
		}
	}

	if current.handler != nil {
		r.executeRoute(current.route, w, req, params)
		return
	}

	http.NotFound(w, req)
}

// ContextKey used for storing params in request context
type ContextKey string

const (
	ParamsKey        ContextKey = "params"
	RouteBindingsKey ContextKey = "route_bindings"
)

func (router *Router) executeRoute(route *Route, w http.ResponseWriter, req *http.Request, params map[string]string) {
	ctx := req.Context()

	if len(params) > 0 {
		ctx = context.WithValue(ctx, ParamsKey, params)
	}

	// Resolve implicit/explicit route model bindings
	resolved := make(map[string]any)
	for key, rawValue := range params {
		if value, err := router.ResolveBinding(key, rawValue); err == nil {
			resolved[key] = value
		} else {
			// If binding fails, we can choose to 404 here for implicit binding
			http.NotFound(w, req)
			return
		}
	}
	if len(resolved) > 0 {
		ctx = context.WithValue(ctx, RouteBindingsKey, resolved)
	}

	req = req.WithContext(ctx)

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := route.Handler(w, r); err != nil {
			router.handleError(w, r, err)
		}
	})

	// Apply middlewares in reverse order so they execute in the order added
	for i := len(route.Middlewares) - 1; i >= 0; i-- {
		handler = route.Middlewares[i](handler)
	}

	handler.ServeHTTP(w, req)
}

func (router *Router) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var httpErr *exception.HttpException
	if errors.As(err, &httpErr) {
		httpErr.Render(w, r)
		return
	}

	// For non-HttpExceptions, log the error and return 500
	log.Printf("[Router Error] %s: %s - %v", r.Method, r.URL.Path, err)

	w.WriteHeader(http.StatusInternalServerError)
	if os.Getenv("APP_DEBUG") == "true" {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("Internal Server Error"))
	}
}

// Get helper
func (r *Router) Get(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodGet, path, handler)
}

// Post helper
func (r *Router) Post(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodPost, path, handler)
}

// Put helper
func (r *Router) Put(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodPut, path, handler)
}

// Patch helper
func (r *Router) Patch(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodPatch, path, handler)
}

// Delete helper
func (r *Router) Delete(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodDelete, path, handler)
}

// Options helper
func (r *Router) Options(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodOptions, path, handler)
}

// Head helper
func (r *Router) Head(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodHead, path, handler)
}

// Group creates a new route group
func (r *Router) Group(prefix string, callback func(*Router)) {
	r.mu.Lock()
	oldPrefix := r.groupPrefix
	r.groupPrefix = oldPrefix + prefix
	r.mu.Unlock()
	
	callback(r)
	
	r.mu.Lock()
	r.groupPrefix = oldPrefix
	r.mu.Unlock()
}

// Use adds a middleware to the router for subsequent routes
func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, mw)
}

// Alias registers a short name for a middleware (e.g. "auth" => auth.Middleware).
func (r *Router) Alias(name string, mw func(http.Handler) http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewareAliases[name] = mw
}

// GroupMiddleware registers a named group of middlewares (e.g. "web", "api").
func (r *Router) GroupMiddleware(name string, mws ...func(http.Handler) http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewareGroups[name] = append(r.middlewareGroups[name], mws...)
}

// Middleware applies one or more named middlewares (from aliases or groups) to the current group.
func (r *Router) Middleware(names ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, name := range names {
		// Check alias first
		if mw, ok := r.middlewareAliases[name]; ok {
			r.middlewares = append(r.middlewares, mw)
			continue
		}
		// Check group
		if group, ok := r.middlewareGroups[name]; ok {
			r.middlewares = append(r.middlewares, group...)
		}
	}
}

// Bind registers a model binder for a route parameter.
// Example: router.Bind("user", func(id string) (any, error) { return user.Find(id) })
func (r *Router) Bind(param string, resolver func(string) (any, error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.binders[param] = resolver
}

// ResolveBinding attempts to resolve a route parameter using a registered binder.
func (r *Router) ResolveBinding(param string, value string) (any, error) {
	r.mu.RLock()
	resolver, ok := r.binders[param]
	r.mu.RUnlock()

	if !ok {
		return value, nil // return raw value if no binder
	}
	return resolver(value)
}

// GetAllRoutes returns a copy of all registered routes (for route:list and introspection).
func (r *Router) GetAllRoutes() []*Route {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routes := make([]*Route, len(r.allRoutes))
	copy(routes, r.allRoutes)
	return routes
}

// Binding retrieves a resolved route model binding by parameter name.
func Binding(r *http.Request, key string) (any, bool) {
	if bindings, ok := r.Context().Value(RouteBindingsKey).(map[string]any); ok {
		if val, exists := bindings[key]; exists {
			return val, true
		}
	}
	return nil, false
}

// Model is a generic helper to retrieve a strongly-typed bound model from the route.
// Usage:
//
//	user, ok := routing.Model[models.User](req, "user")
func Model[T any](r *http.Request, key string) (T, bool) {
	var zero T
	val, ok := Binding(r, key)
	if !ok {
		return zero, false
	}
	if typed, ok := val.(T); ok {
		return typed, true
	}
	return zero, false
}

// ==================== PHASE 4: Redirect Routes ====================

// Redirect registers a route that redirects from one URL to another.
func (r *Router) Redirect(from, to string, statusCode ...int) {
	code := http.StatusMovedPermanently
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	r.Get(from, func(w http.ResponseWriter, req *http.Request) error {
		http.Redirect(w, req, to, code)
		return nil
	})
}

// PermanentRedirect registers a route that sends a 301 permanent redirect.
func (r *Router) PermanentRedirect(from, to string) {
	r.Redirect(from, to, http.StatusMovedPermanently)
}

// TemporaryRedirect registers a route that sends a 307 temporary redirect.
func (r *Router) TemporaryRedirect(from, to string) {
	r.Redirect(from, to, http.StatusTemporaryRedirect)
}

// ==================== PHASE 4: Fallback Routes ====================

// fallbackHandler stores the catch-all handler.
var fallbackHandler HandlerFunc

// SetFallback registers a catch-all route that handles any request not matched.
func (r *Router) SetFallback(handler HandlerFunc) {
	fallbackHandler = handler
}

// GetFallback returns the registered fallback handler.
func (r *Router) GetFallback() HandlerFunc {
	return fallbackHandler
}

