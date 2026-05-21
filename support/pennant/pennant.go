package pennant

import (
	"context"
	"fmt"
)

// Store defines how feature flags are persisted.
type Store interface {
	Get(ctx context.Context, feature string, scope any) (bool, error)
	Set(ctx context.Context, feature string, scope any, value bool) error
	Flush(ctx context.Context, feature string) error
}

// Manager orchestrates feature flags.
type Manager struct {
	stores       map[string]Store
	defaultStore string
	features     map[string]func(scope any) bool
}

// NewManager creates a new Pennant manager.
func NewManager(defaultStore string) *Manager {
	return &Manager{
		stores:       make(map[string]Store),
		defaultStore: defaultStore,
		features:     make(map[string]func(scope any) bool),
	}
}

// Extend adds a custom store.
func (m *Manager) Extend(name string, store Store) {
	m.stores[name] = store
}

// Define registers a default resolution closure for a feature.
func (m *Manager) Define(feature string, resolver func(scope any) bool) {
	m.features[feature] = resolver
}

// Active determines if a feature is active for the given scope.
func (m *Manager) Active(ctx context.Context, feature string, scope any) bool {
	store := m.stores[m.defaultStore]
	if store == nil {
		return false
	}

	val, err := store.Get(ctx, feature, scope)
	if err == nil {
		return val
	}

	// Resolve using closure if not in store
	if resolver, exists := m.features[feature]; exists {
		result := resolver(scope)
		// Persist result
		store.Set(ctx, feature, scope, result)
		return result
	}

	return false
}

// InMemoryStore for testing.
type InMemoryStore struct {
	flags map[string]map[string]bool
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		flags: make(map[string]map[string]bool),
	}
}

func (s *InMemoryStore) scopeKey(scope any) string {
	return fmt.Sprintf("%v", scope)
}

func (s *InMemoryStore) Get(ctx context.Context, feature string, scope any) (bool, error) {
	if fMap, ok := s.flags[feature]; ok {
		if val, ok := fMap[s.scopeKey(scope)]; ok {
			return val, nil
		}
	}
	return false, fmt.Errorf("not found")
}

func (s *InMemoryStore) Set(ctx context.Context, feature string, scope any, value bool) error {
	if s.flags[feature] == nil {
		s.flags[feature] = make(map[string]bool)
	}
	s.flags[feature][s.scopeKey(scope)] = value
	return nil
}

func (s *InMemoryStore) Flush(ctx context.Context, feature string) error {
	delete(s.flags, feature)
	return nil
}

// In a real implementation, we would define RedisStore (using HSET/HGET) 
// and DatabaseStore (using a feature_flags table) here or in dedicated driver files.
