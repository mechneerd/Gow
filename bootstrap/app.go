package bootstrap

import (
	"gow/broadcasting"
	"gow/config"
	"gow/foundation"
	"gow/logging"
	"gow/mail"
)

// NewApplication creates and boots a standard GoW Application with core providers.
func NewApplication(basePath string) *foundation.Application {
	app := foundation.NewApplication(basePath)

	// Register core service providers (order matters)
	app.RegisterProvider(&config.ServiceProvider{})
	app.RegisterProvider(&logging.ServiceProvider{})
	app.RegisterProvider(&broadcasting.ServiceProvider{})
	app.RegisterProvider(&mail.ServiceProvider{})

	// Boot the application (freezes container)
	app.Boot()

	return app
}
