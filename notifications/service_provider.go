package notifications

import (
	"gow/container"
	"gow/foundation"
	"gow/mail"
)

// ServiceProvider registers the notification system with the application.
type ServiceProvider struct {
	foundation.BaseServiceProvider
}

// Register sets up the notification manager and default channels.
func (p *ServiceProvider) Register(app *foundation.Application) {
	manager := NewManager()

	// Register Mail channel if mailer is available (using the proper generic resolver)
	if mailer, err := container.Make[*mail.Mailer](app.Container); err == nil && mailer != nil {
		manager.Extend("mail", NewMailChannel(mailer))
	}

	// Always register database channel (it can be extended later with proper DB)
	// For now we leave it to the user to register DatabaseChannel when they have a DB.

	app.Instance((*Manager)(nil), manager)
	SetDefaultManager(manager)
}

// Boot can be used to register additional channels after all providers are ready.
func (p *ServiceProvider) Boot(app *foundation.Application) {
	// Future: auto-register DatabaseChannel if orm.DB is present
}
