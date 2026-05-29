package routing

import (
	"fmt"
	"net/http"
	"sync"
	"reflect"
	"strings"
)

// Resource registers standard CRUD routes for a controller using reflection.
// It maps methods like Index, Create, Store, Show, Edit, Update, Destroy.
// Optional middleware can be passed as additional arguments.
func (r *Router) Resource(name string, controller any, middleware ...func(http.Handler) http.Handler) {
	r.registerResource(name, controller, false, middleware...)
}

// ApiResource registers standard API CRUD routes (excluding Create and Edit).
func (r *Router) ApiResource(name string, controller any, middleware ...func(http.Handler) http.Handler) {
	r.registerResource(name, controller, true, middleware...)
}

func (r *Router) registerResource(name string, controller any, apiOnly bool, middleware ...func(http.Handler) http.Handler) {
	val := reflect.ValueOf(controller)
	if val.Kind() != reflect.Ptr && val.Kind() != reflect.Struct {
		panic("controller must be a struct or pointer to struct")
	}

	basePath := "/" + strings.Trim(name, "/")
	paramName := "{" + singularize(name) + "}" // e.g. "photos" -> "{photo}"

	methods := map[string]struct {
		method string
		path   string
		isApi  bool
	}{
		"Index":   {http.MethodGet, basePath, true},
		"Create":  {http.MethodGet, basePath + "/create", false},
		"Store":   {http.MethodPost, basePath, true},
		"Show":    {http.MethodGet, basePath + "/" + paramName, true},
		"Edit":    {http.MethodGet, basePath + "/" + paramName + "/edit", false},
		"Update":  {http.MethodPut, basePath + "/" + paramName, true},
		"Destroy": {http.MethodDelete, basePath + "/" + paramName, true},
	}

	for methodName, routeInfo := range methods {
		if apiOnly && !routeInfo.isApi {
			continue
		}

		methodVal := val.MethodByName(methodName)
		if methodVal.IsValid() {
			handler, ok := methodVal.Interface().(func(http.ResponseWriter, *http.Request) error)
			if ok {
				route := r.AddRoute(routeInfo.method, routeInfo.path, handler)
				routeName := fmt.Sprintf("%s.%s", name, strings.ToLower(methodName))
				r.SetName(route, routeName)
				// Apply resource-specific middleware
				for i := len(middleware) - 1; i >= 0; i-- {
					currentHandler := route.Handler
					route.Handler = func(w http.ResponseWriter, r *http.Request) error {
						middleware[i](http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							_ = currentHandler(w, r)
						})).ServeHTTP(w, r)
						return nil
					}
				}
			}
		}
	}
}

// Macro support for Router (adds dynamic methods basically)
var (
	routerMacros = make(map[string]func(*Router, ...any) any)
	macroMu      sync.RWMutex
)

// Macro registers a custom macro on the router
func (r *Router) Macro(name string, macro func(*Router, ...any) any) {
	macroMu.Lock()
	defer macroMu.Unlock()
	routerMacros[name] = macro
}

// CallMacro executes a registered macro
func (r *Router) CallMacro(name string, args ...any) (any, error) {
	macroMu.RLock()
	macro, exists := routerMacros[name]
	macroMu.RUnlock()
	if exists {
		return macro(r, args...), nil
	}
	return nil, fmt.Errorf("macro %s not found", name)
}

// singularize converts a plural name to singular using basic English rules.
func singularize(name string) string {
	lower := strings.ToLower(name)
	// Handle common irregular forms
	if lower == "children" {
		return "child"
	}
	if lower == "people" {
		return "person"
	}
	// HandleIES: categories -> category, entries -> entry
	if strings.HasSuffix(lower, "ies") && len(lower) > 3 {
		return lower[:len(lower)-3] + "y"
	}
	// HandleSES/SHES/CHES/XES/ZES: boxes -> box, churches -> church
	if strings.HasSuffix(lower, "ses") || strings.HasSuffix(lower, "shes") ||
		strings.HasSuffix(lower, "ches") || strings.HasSuffix(lower, "xes") || strings.HasSuffix(lower, "zes") {
		return lower[:len(lower)-2]
	}
	// Handle simple S: photos -> photo, users -> user
	if strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") {
		return lower[:len(lower)-1]
	}
	return lower
}

