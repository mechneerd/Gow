package access

// Gate handles closure-based authorization checks.
type Gate struct {
	abilities map[string]func(user any, models ...any) bool
	policies  map[string]any // maps model name to a policy instance
}

// NewGate creates a new Gate instance.
func NewGate() *Gate {
	return &Gate{
		abilities: make(map[string]func(user any, models ...any) bool),
		policies:  make(map[string]any),
	}
}

// Define registers a new ability using a closure.
func (g *Gate) Define(ability string, callback func(user any, models ...any) bool) {
	g.abilities[ability] = callback
}

// Policy registers a policy for a specific model type.
func (g *Gate) Policy(modelName string, policy any) {
	g.policies[modelName] = policy
}

// Allows determines if the user has the given ability.
func (g *Gate) Allows(ability string, user any, models ...any) bool {
	// First check defined closures
	if callback, exists := g.abilities[ability]; exists {
		return callback(user, models...)
	}

	// In a complete implementation, this would use reflection to check
	// if a registered policy has a method matching the ability name.
	// For example: if policy for Post exists, check if policy.Update(user, post) returns true.

	return false
}

// Denies determines if the user does NOT have the given ability.
func (g *Gate) Denies(ability string, user any, models ...any) bool {
	return !g.Allows(ability, user, models...)
}
