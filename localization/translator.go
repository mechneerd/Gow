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
