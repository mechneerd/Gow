package routing

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

// HandlerFunc is the GoW handler function signature.
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Router represents the HTTP router.
type Router struct {
	mu           sync.RWMutex
	trees        map[string]*node
	namedRoutes  map[string]*Route
	middlewares  []func(http.Handler) http.Handler
	groupPrefix  string
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
		trees:       make(map[string]*node),
		namedRoutes: make(map[string]*Route),
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
const ParamsKey ContextKey = "params"

func (r *Router) executeRoute(route *Route, w http.ResponseWriter, req *http.Request, params map[string]string) {
	if len(params) > 0 {
		ctx := context.WithValue(req.Context(), ParamsKey, params)
		req = req.WithContext(ctx)
	}

	var handler http.Handler = http.HandlerFunc(route.Handler)
	
	// Apply middlewares in reverse order so they execute in the order added
	for i := len(route.Middlewares) - 1; i >= 0; i-- {
		handler = route.Middlewares[i](handler)
	}

	handler.ServeHTTP(w, req)
}

// Get helper
func (r *Router) Get(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodGet, path, handler)
}

// Post helper
func (r *Router) Post(path string, handler HandlerFunc) *Route {
	return r.AddRoute(http.MethodPost, path, handler)
}

// Group creates a new route group
func (r *Router) Group(prefix string, callback func(*Router)) {
	oldPrefix := r.groupPrefix
	
	r.mu.Lock()
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
