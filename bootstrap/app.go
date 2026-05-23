package bootstrap

import (
	"gow/broadcasting"
	"gow/config"
	"gow/foundation"
	"gow/localization"
	"gow/logging"
	"gow/mail"
	"gow/notifications"
)

// NewApplication creates and boots a standard GoW Application with core providers.
func NewApplication(basePath string) *foundation.Application {
	app := foundation.NewApplication(basePath)

	// Register core service providers (order matters)
	app.RegisterProvider(&config.ServiceProvider{})
	app.RegisterProvider(&logging.ServiceProvider{})
	app.RegisterProvider(&broadcasting.ServiceProvider{})
	app.RegisterProvider(&mail.ServiceProvider{})
	app.RegisterProvider(&notifications.ServiceProvider{})
	app.RegisterProvider(&localization.ServiceProvider{})

	// Boot the application (freezes container)
	app.Boot()

	return app
}
