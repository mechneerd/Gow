package logging

import (
	"gow/config"
	"gow/foundation"
	"log/slog"
)

// ServiceProvider registers the logging system.
type ServiceProvider struct {
	foundation.BaseServiceProvider
}

// Register registers the logger into the container.
func (p *ServiceProvider) Register(app *foundation.Application) {
	// We use Make to resolve the config repository
	cfg, err := app.Resolve(config.Repository{})
	
	level := "info"
	if err == nil {
		repo := cfg.(*config.Repository)
		level = repo.Get("LOG_LEVEL", "info")
	}

	logger := Setup(level)
	
	// Bind as instance
	app.Instance((*slog.Logger)(nil), logger)
}
