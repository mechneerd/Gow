package access

import (
	"fmt"
	"reflect"
)

// Gate represents the central authorization registry.
type Gate struct {
	policies map[string]any // Model Name -> Policy Struct
	gates    map[string]func(user any, args ...any) bool
	before   []func(user any, ability string, args ...any) *bool
	after    []func(user any, ability string, result bool, args ...any) *bool
}

// NewGate creates a new Gate instance.
func NewGate() *Gate {
	return &Gate{
		policies: make(map[string]any),
		gates:    make(map[string]func(user any, args ...any) bool),
	}
}

// Define registers an authorization closure.
func (g *Gate) Define(ability string, callback func(user any, args ...any) bool) {
	g.gates[ability] = callback
}

// Policy registers a policy struct for a model.
func (g *Gate) Policy(modelName string, policy any) {
	g.policies[modelName] = policy
}

// Before registers a callback to run before all other authorization checks.
func (g *Gate) Before(callback func(user any, ability string, args ...any) *bool) {
	g.before = append(g.before, callback)
}

// After registers a callback to run after all other authorization checks.
func (g *Gate) After(callback func(user any, ability string, result bool, args ...any) *bool) {
	g.after = append(g.after, callback)
}

// Allows determines if the given user is authorized for the given ability.
func (g *Gate) Allows(user any, ability string, args ...any) bool {
	// Execute Before hooks
	for _, hook := range g.before {
		if result := hook(user, ability, args...); result != nil {
			return *result
		}
	}

	result := g.check(user, ability, args...)

	// Execute After hooks
	for _, hook := range g.after {
		if hookResult := hook(user, ability, result, args...); hookResult != nil {
			return *hookResult
		}
	}

	return result
}

// check contains the core authorization logic
func (g *Gate) check(user any, ability string, args ...any) bool {
	// Check simple closures first
	if callback, exists := g.gates[ability]; exists {
		return callback(user, args...)
	}

	// Policy resolution
	if len(args) > 0 {
		model := args[0]
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}

		if policy, exists := g.policies[modelType.Name()]; exists {
			// Find method matching Ability on Policy
			methodName := ""
			if len(ability) > 0 {
				methodName = string(ability[0]-32) + ability[1:] // simple Titleize
			}

			policyVal := reflect.ValueOf(policy)
			method := policyVal.MethodByName(methodName)

			if method.IsValid() {
				callArgs := []reflect.Value{reflect.ValueOf(user)}
				
				// Pass all args to the method
				for _, arg := range args {
					callArgs = append(callArgs, reflect.ValueOf(arg))
				}

				// If method takes fewer arguments than provided, slice callArgs
				methodType := method.Type()
				if methodType.NumIn() < len(callArgs) {
					callArgs = callArgs[:methodType.NumIn()]
				}

				result := method.Call(callArgs)
				if len(result) > 0 {
					return result[0].Bool()
				}
			}
		}
	}

	return false
}

// Denies is the inverse of Allows.
func (g *Gate) Denies(user any, ability string, args ...any) bool {
	return !g.Allows(user, ability, args...)
}

// Authorize panics with an error if the user is not allowed.
func (g *Gate) Authorize(user any, ability string, args ...any) {
	if g.Denies(user, ability, args...) {
		panic(fmt.Sprintf("This action is unauthorized: %s", ability))
	}
}
