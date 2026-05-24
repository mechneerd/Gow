package foundation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ProviderRegistry manages lists of service providers including publishables.
type ProviderRegistry struct {
	providers    []ServiceProvider
	publishables []PublishableProvider
}

// NewProviderRegistry creates a new empty registry.
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{}
}

// Register adds a provider, tracking if it is publishable.
func (r *ProviderRegistry) Register(p ServiceProvider) {
	r.providers = append(r.providers, p)
	if pp, ok := p.(PublishableProvider); ok {
		r.publishables = append(r.publishables, pp)
	}
}

// Providers returns all registered service providers.
func (r *ProviderRegistry) Providers() []ServiceProvider {
	return r.providers
}

// Publishables returns only the publishable providers.
func (r *ProviderRegistry) Publishables() []PublishableProvider {
	return r.publishables
}

// AutoDiscover performs auto-discovery of providers for the given application.
// This is intentionally a no-op in Go (no reliable runtime package scanning).
// All providers must be explicitly registered in the application's bootstrap.
func AutoDiscover(app *Application) {
	// Intentionally empty. Extend only if a config-driven provider list is added.
}

// PublishAssets copies the files declared by the PublishableProvider into the
// application's base path using the Publishes() map (src => destRel).
func PublishAssets(provider PublishableProvider, basePath string) error {
	for src, dstRel := range provider.Publishes() {
		dst := filepath.Join(basePath, dstRel)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		srcFile, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("publish source not found: %s: %w", src, err)
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	}
	return nil
}
