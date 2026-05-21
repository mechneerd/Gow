package foundation

import (
	"gow/container"
)

// Application represents the GoW framework application.
type Application struct {
	*container.Container
	basePath  string
	booted    bool
	providers []ServiceProvider
}

// NewApplication creates a new Application instance.
func NewApplication(basePath string) *Application {
	app := &Application{
		Container: container.New(),
		basePath:  basePath,
	}
	app.registerBaseBindings()
	return app
}

func (app *Application) registerBaseBindings() {
	app.Instance((*Application)(nil), app)
	app.Instance((*container.Container)(nil), app.Container)
}

// BasePath returns the application's base path.
func (app *Application) BasePath() string {
	return app.basePath
}

// Bootstrap runs the given bootstrapper functions.
func (app *Application) Bootstrap(bootstrappers ...func(*Application)) {
	for _, bootstrapper := range bootstrappers {
		bootstrapper(app)
	}
}

// RegisterProvider registers a service provider with the application.
func (app *Application) RegisterProvider(p ServiceProvider) {
	p.Register(app)
	app.providers = append(app.providers, p)
}

// Boot boots the application, preventing further container bindings.
func (app *Application) Boot() {
	if app.booted {
		return
	}

	for _, p := range app.providers {
		p.Boot(app)
	}

	app.booted = true
	app.Freeze()
}
