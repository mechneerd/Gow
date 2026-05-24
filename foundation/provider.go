package foundation

// ServiceProvider represents an abstraction for registering and booting services.
type ServiceProvider interface {
	Register(app *Application)
	Boot(app *Application)
}

// BaseServiceProvider provides empty default implementations for ServiceProvider.
type BaseServiceProvider struct{}

// Register is the default empty register method.
func (p *BaseServiceProvider) Register(app *Application) {}

// Boot is the default empty boot method.
func (p *BaseServiceProvider) Boot(app *Application) {}

// PublishableProvider extends ServiceProvider to support asset publishing (configs, views, etc).
// Implementers return a map of source (relative or absolute) => destination (relative to app base).
type PublishableProvider interface {
	ServiceProvider
	Publishes() map[string]string
}

