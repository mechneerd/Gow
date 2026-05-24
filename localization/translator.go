package localization

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Translator handles translation of language lines.
type Translator struct {
	path   string
	locale string
	lines  map[string]map[string]string // [locale][key] = translation
}

// NewTranslator initializes a new translator.
func NewTranslator(path, defaultLocale string) *Translator {
	return &Translator{
		path:   path,
		locale: defaultLocale,
		lines:  make(map[string]map[string]string),
	}
}

// SetLocale sets the active locale.
func (t *Translator) SetLocale(locale string) {
	t.locale = locale
}

// Load loads translation strings for a locale from a JSON file.
func (t *Translator) Load(locale string) error {
	if _, exists := t.lines[locale]; exists {
		return nil
	}

	file := filepath.Join(t.path, locale+".json")
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			t.lines[locale] = make(map[string]string)
			return nil
		}
		return err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	t.lines[locale] = translations
	return nil
}

// Translate returns the translated string for a key.
func (t *Translator) Translate(key string, replace map[string]string) string {
	t.Load(t.locale) // Ensure loaded
	
	line, ok := t.lines[t.locale][key]
	if !ok {
		return key
	}

	for k, v := range replace {
		line = strings.ReplaceAll(line, ":"+k, v)
	}

	return line
}

// ==================== Global Localization Helpers (complete i18n) ====================

var defaultTranslator *Translator

// SetDefaultTranslator sets the app-wide translator instance (usually done in bootstrap).
func SetDefaultTranslator(t *Translator) {
	defaultTranslator = t
}

// Translate is the global translation helper (Laravel-like).
// Usage: Translate("welcome.message", map[string]string{"name": "John"})
func Translate(key string, replaces ...map[string]string) string {
	if defaultTranslator == nil {
		return key
	}
	var replace map[string]string
	if len(replaces) > 0 {
		replace = replaces[0]
	}
	return defaultTranslator.Translate(key, replace)
}

// Trans is alias for Translate.
func Trans(key string, replaces ...map[string]string) string {
	return Translate(key, replaces...)
}

// SetLocale changes the current locale globally.
func SetLocale(locale string) {
	if defaultTranslator != nil {
		defaultTranslator.SetLocale(locale)
	}
}

// GetLocale returns current locale.
func GetLocale() string {
	if defaultTranslator != nil {
		return defaultTranslator.locale
	}
	return "en"
}

// LoadLocale preloads a locale.
func LoadLocale(locale string) error {
	if defaultTranslator != nil {
		return defaultTranslator.Load(locale)
	}
	return nil
}


