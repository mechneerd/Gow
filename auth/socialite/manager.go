package socialite

import (
	"context"
)

// User represents a user returned from an OAuth provider.
type User struct {
	ID        string
	Nickname  string
	Name      string
	Email     string
	Avatar    string
	Token     string
	RefreshToken string
	ExpiresIn  int
}

// Provider is the interface all OAuth providers must implement.
type Provider interface {
	Redirect(state string) string
	User(ctx context.Context, code string) (*User, error)
}

// Manager holds all registered providers.
type Manager struct {
	providers map[string]Provider
}

// NewManager creates a new Socialite manager.
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
	}
}

// Extend registers a new provider.
func (m *Manager) Extend(name string, provider Provider) {
	m.providers[name] = provider
}

// Driver returns a provider by name.
func (m *Manager) Driver(name string) Provider {
	return m.providers[name]
}

// Default providers can be registered via service provider later.

