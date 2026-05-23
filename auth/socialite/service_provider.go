package socialite

import (
	"gow/foundation"
)

// ServiceProvider registers Socialite.
type ServiceProvider struct{}

func (p *ServiceProvider) Register(app *foundation.Application) {
	manager := NewManager()

	// Example: users can extend with their own providers
	// manager.Extend("google", NewGoogleProvider(...))

	app.Instance((*Manager)(nil), manager)
}

func (p *ServiceProvider) Boot(app *foundation.Application) {}
