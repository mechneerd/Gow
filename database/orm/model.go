package orm

import (
	"reflect"
)

// Scope represents a query scope that can be applied to a builder.
type Scope interface {
	Apply(builder *Builder) *Builder
}

// GlobalScopes registry.
var globalScopes = make(map[string][]Scope)

// AddGlobalScope registers a global scope for a given model type name.
func AddGlobalScope(model string, scope Scope) {
	globalScopes[model] = append(globalScopes[model], scope)
}

// Observer represents model event hooks.
type Observer interface {
	Creating(model any) bool
	Created(model any)
	Updating(model any) bool
	Updated(model any)
	Deleting(model any) bool
	Deleted(model any)
}

// Observers registry.
var observers = make(map[string][]Observer)

// Observe registers an observer for a model.
func Observe(model string, observer Observer) {
	observers[model] = append(observers[model], observer)
}

// Cast represents an attribute casting strategy.
type Cast interface {
	Get(value any) (any, error)
	Set(value any) (any, error)
}

// DispatchModelEvent fires registered observers for the specific event type.
// Returning false from a "ing" event (like Creating) halts the operation.
func DispatchModelEvent(model any, event string) bool {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	name := modelType.Name()

	if obsList, ok := observers[name]; ok {
		for _, obs := range obsList {
			switch event {
			case "creating":
				if !obs.Creating(model) {
					return false
				}
			case "created":
				obs.Created(model)
			case "updating":
				if !obs.Updating(model) {
					return false
				}
			case "updated":
				obs.Updated(model)
			case "deleting":
				if !obs.Deleting(model) {
					return false
				}
			case "deleted":
				obs.Deleted(model)
			}
		}
	}
	return true
}

// ModelConfig holds strictness rules.
var (
	PreventLazyLoading bool
	HasUuids           bool
)
