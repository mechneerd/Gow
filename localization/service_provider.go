package localization

import (
	"github.com/mechneerd/gow/foundation"
	"path/filepath"
)

// ServiceProvider registers the localization system.
type ServiceProvider struct {
	app *foundation.Application
}

func (p *ServiceProvider) Register(app *foundation.Application) {
	p.app = app
}

func (p *ServiceProvider) Boot(app *foundation.Application) {
	locale := "en"
	langPath := "lang"

	// Resolve absolute lang path from base path if relative
	basePath := "."
	if bp := app.BasePath(); bp != "" {
		basePath = bp
	}
	fullPath := filepath.Join(basePath, langPath)

	t := NewTranslator(fullPath, locale)
	t.Load(locale)
	t.Load("en")

	SetDefaultTranslator(t)

	// Bind to container (use pointer type)
	app.Singleton((*Translator)(nil), t)
}

