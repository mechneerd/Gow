package mail

import (
	"github.com/mechneerd/gow/config"
	"github.com/mechneerd/gow/container"
	"github.com/mechneerd/gow/foundation"
)

// ServiceProvider registers the mail system with the application.
type ServiceProvider struct {
	foundation.BaseServiceProvider
}

// Register sets up the mail configuration and binds the Mailer.
func (p *ServiceProvider) Register(app *foundation.Application) {
	cfg, _ := container.Make[config.Repository](app.Container)

	mailer := "log"
	host := cfg.Get("MAIL_HOST", "localhost")
	port := cfg.Get("MAIL_PORT", "1025")
	username := cfg.Get("MAIL_USERNAME", "")
	password := cfg.Get("MAIL_PASSWORD", "")
	from := cfg.Get("MAIL_FROM_ADDRESS", "hello@example.com")

	var driver Driver

	switch mailer {
	case "smtp":
		driver = &SmtpDriver{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
		}
	default:
		driver = &LogDriver{}
	}

	mailerInstance := NewMailer(driver)
	mailerInstance.SetFrom(from)

	app.Instance((*Mailer)(nil), mailerInstance)
}

