package orm

import (
	"gow/database/query"
	"reflect"
)

// MassAssignable allows models to control mass assignment protection,
// similar to Laravel's $fillable and $guarded properties.
type MassAssignable interface {
	// Fillable returns the list of fields that are allowed to be mass assigned.
	// If empty, all fields are allowed unless Guarded is specified.
	Fillable() []string

	// Guarded returns the list of fields that are not allowed to be mass assigned.
	// Use ["*"] to guard all fields by default.
	Guarded() []string
}

// Scope represents a query scope that can be applied to a builder.
type Scope interface {
	Apply(builder *query.Builder) *query.Builder
}

// ScopeFunc is a convenience wrapper to use a function as a Scope.
type ScopeFunc func(builder *query.Builder) *query.Builder

func (f ScopeFunc) Apply(builder *query.Builder) *query.Builder {
	return f(builder)
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

// EventDispatcher interface allows ORM to dispatch global events without strong coupling to the events package.
type EventDispatcher interface {
	Dispatch(event any)
}

var globalEventManager EventDispatcher

// SetEventManager injects the global event manager.
func SetEventManager(m EventDispatcher) {
	globalEventManager = m
}

// DispatchModelEvent fires registered observers, interface hooks, and global events.
// Returning false from a "ing" event (like Creating) halts the operation.
func DispatchModelEvent(model any, event string) bool {
	// 1. Interface Hooks
	switch event {
	case "creating":
		if hook, ok := model.(BeforeCreateHook); ok {
			if err := hook.BeforeCreate(); err != nil {
				return false
			}
		}
		if hook, ok := model.(BeforeSaveHook); ok {
			if err := hook.BeforeSave(); err != nil {
				return false
			}
		}
	case "created":
		if hook, ok := model.(AfterCreateHook); ok {
			hook.AfterCreate()
		}
		if hook, ok := model.(AfterSaveHook); ok {
			hook.AfterSave()
		}
	case "updating":
		if hook, ok := model.(BeforeUpdateHook); ok {
			if err := hook.BeforeUpdate(); err != nil {
				return false
			}
		}
		if hook, ok := model.(BeforeSaveHook); ok {
			if err := hook.BeforeSave(); err != nil {
				return false
			}
		}
	case "updated":
		if hook, ok := model.(AfterUpdateHook); ok {
			hook.AfterUpdate()
		}
		if hook, ok := model.(AfterSaveHook); ok {
			hook.AfterSave()
		}
	case "deleting":
		if hook, ok := model.(BeforeDeleteHook); ok {
			if err := hook.BeforeDelete(); err != nil {
				return false
			}
		}
	case "deleted":
		if hook, ok := model.(AfterDeleteHook); ok {
			hook.AfterDelete()
		}
	case "restoring":
		if hook, ok := model.(BeforeRestoreHook); ok {
			if err := hook.BeforeRestore(); err != nil {
				return false
			}
		}
	case "restored":
		if hook, ok := model.(AfterRestoreHook); ok {
			hook.AfterRestore()
		}
	}

	// 2. Global Events (events.Manager)
	if globalEventManager != nil {
		switch event {
		case "creating":
			globalEventManager.Dispatch(ModelCreating{Model: model})
			globalEventManager.Dispatch(ModelSaving{Model: model})
		case "created":
			globalEventManager.Dispatch(ModelCreated{Model: model})
			globalEventManager.Dispatch(ModelSaved{Model: model})
		case "updating":
			globalEventManager.Dispatch(ModelUpdating{Model: model})
			globalEventManager.Dispatch(ModelSaving{Model: model})
		case "updated":
			globalEventManager.Dispatch(ModelUpdated{Model: model})
			globalEventManager.Dispatch(ModelSaved{Model: model})
		case "deleting":
			globalEventManager.Dispatch(ModelDeleting{Model: model})
		case "deleted":
			globalEventManager.Dispatch(ModelDeleted{Model: model})
		case "restoring":
			globalEventManager.Dispatch(ModelRestoring{Model: model})
		case "restored":
			globalEventManager.Dispatch(ModelRestored{Model: model})
		}
	}

	// 3. Legacy String Observers (optional, can be deprecated)
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
