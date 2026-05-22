package notifications

import (
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

	// Register Database channel if a DB connection is available
	if dbIface, err := app.Resolve((*interface{})(nil)); err == nil {
		// This is a loose check — in real usage we would resolve *orm.DB properly
		_ = dbIface
	}

	// Register Mail channel if mailer is available
	if mailerIface, err := app.Resolve((*mail.Mailer)(nil)); err == nil {
		if mailer, ok := mailerIface.(*mail.Mailer); ok {
			manager.Extend("mail", NewMailChannel(mailer))
		}
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
