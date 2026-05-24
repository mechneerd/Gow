package view

// Composer represents a callback that is executed before a view is rendered.
type Composer func(data map[string]any)

// ComposerRegistry stores view composers for specific views.
type ComposerRegistry struct {
	composers map[string][]Composer
}

// NewComposerRegistry creates a new composer registry.
func NewComposerRegistry() *ComposerRegistry {
	return &ComposerRegistry{
		composers: make(map[string][]Composer),
	}
}

// Register adds a composer for a specific view or wildcard (*).
func (r *ComposerRegistry) Register(view string, composer Composer) {
	r.composers[view] = append(r.composers[view], composer)
}

// Compose executes all registered composers for a given view, modifying the data map in-place.
func (r *ComposerRegistry) Compose(view string, data map[string]any) {
	// Execute wildcard composers
	if wildcards, ok := r.composers["*"]; ok {
		for _, c := range wildcards {
			c(data)
		}
	}

	// Execute specific view composers
	if specific, ok := r.composers[view]; ok {
		for _, c := range specific {
			c(data)
		}
	}
}

