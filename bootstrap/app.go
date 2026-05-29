package bootstrap

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	"github.com/mechneerd/gow/broadcasting"
	"github.com/mechneerd/gow/config"
	"github.com/mechneerd/gow/foundation"
	"github.com/mechneerd/gow/localization"
	"github.com/mechneerd/gow/logging"
	"github.com/mechneerd/gow/mail"
	"github.com/mechneerd/gow/notifications"
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

