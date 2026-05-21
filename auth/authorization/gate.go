package authorization

import (
	"fmt"
	"reflect"
)

// Gate represents the central authorization registry.
type Gate struct {
	policies map[string]any // Model Name -> Policy Struct
	gates    map[string]func(user any, args ...any) bool
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
// Model name can be derived from the model struct or explicitly passed.
func (g *Gate) Policy(modelName string, policy any) {
	g.policies[modelName] = policy
}

// Allows determines if the given user is authorized for the given ability.
func (g *Gate) Allows(user any, ability string, args ...any) bool {
	// Check simple closures first
	if callback, exists := g.gates[ability]; exists {
		return callback(user, args...)
	}

	// Policy resolution
	// Ability is usually like "update", args[0] is the model instance
	if len(args) > 0 {
		model := args[0]
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		
		policy, exists := g.policies[modelType.Name()]
		if exists {
			// Find method matching Ability on Policy (e.g. "Update")
			// Make ability title case
			methodName := ""
			if len(ability) > 0 {
				methodName = string(ability[0]-32) + ability[1:] // simple Titleize
			}
			
			policyVal := reflect.ValueOf(policy)
			method := policyVal.MethodByName(methodName)
			
			if method.IsValid() {
				// Build args
				callArgs := []reflect.Value{reflect.ValueOf(user), reflect.ValueOf(model)}
				// Add extra args if necessary
				for i := 1; i < len(args); i++ {
					callArgs = append(callArgs, reflect.ValueOf(args[i]))
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

// Authorize panics with an AuthorizationException if the user is not allowed.
func (g *Gate) Authorize(user any, ability string, args ...any) {
	if g.Denies(user, ability, args...) {
		panic(fmt.Sprintf("This action is unauthorized: %s", ability))
	}
}
