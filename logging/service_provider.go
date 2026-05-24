package logging

import (
	"gow/config"
	"gow/container"
	"gow/foundation"
	"log/slog"
)

// ServiceProvider registers the logging system.
type ServiceProvider struct {
	foundation.BaseServiceProvider
}

// Register registers the logger into the container.
func (p *ServiceProvider) Register(app *foundation.Application) {
	// Use the generic Make helper (the correct way to resolve from container)
	repo, err := container.Make[config.Repository](app.Container)
	
	level := "info"
	if err == nil {
		level = repo.Get("LOG_LEVEL", "info")
	}

	logger := Setup(Config{
		Level:   level,
		Channel: repo.Get("LOG_CHANNEL", "single"),
		Path:    repo.Get("LOG_PATH", "storage/logs"),
	})
	
	// Bind as instance
	app.Instance((*slog.Logger)(nil), logger)
}
