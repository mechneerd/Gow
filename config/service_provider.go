package config

import (
	"gow/foundation"
)

// ServiceProvider registers the configuration repository with the application.
type ServiceProvider struct {
	foundation.BaseServiceProvider
}

// Register registers the config repository into the container.
func (p *ServiceProvider) Register(app *foundation.Application) {
	repo := NewRepository(app.BasePath())
	
	// Bind as instance
	app.Instance((*Repository)(nil), repo)
}
