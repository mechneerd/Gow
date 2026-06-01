package localization

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Translator handles translation of language lines.
type Translator struct {
	path           string
	locale         string
	fallbackLocale string
	lines          map[string]map[string]string // [locale][key] = translation
}

// NewTranslator initializes a new translator.
func NewTranslator(path, defaultLocale string) *Translator {
	return &Translator{
		path:           path,
		locale:         defaultLocale,
		fallbackLocale: "en",
		lines:          make(map[string]map[string]string),
	}
}

// SetLocale sets the active locale.
func (t *Translator) SetLocale(locale string) {
	t.locale = locale
}

// SetFallbackLocale sets the fallback locale.
func (t *Translator) SetFallbackLocale(locale string) {
	t.fallbackLocale = locale
}

// GetLocale returns the current locale.
func (t *Translator) GetLocale() string {
	return t.locale
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
// Falls back to fallbackLocale if key not found in current locale.
func (t *Translator) Translate(key string, replace map[string]string) string {
	t.Load(t.locale) // Ensure loaded

	line, ok := t.lines[t.locale][key]
	if !ok {
		// Try fallback locale
		if t.fallbackLocale != t.locale {
			t.Load(t.fallbackLocale)
			line, ok = t.lines[t.fallbackLocale][key]
		}
		if !ok {
			return key
		}
	}

	for k, v := range replace {
		line = strings.ReplaceAll(line, ":"+k, v)
	}

	return line
}

// TransChoice returns the translated string for a key based on count (pluralization).
// Supports simple CLDR-style rules:
//   - "apples" => "{0} No apples|{1} One apple|[2,*] :count apples"
//   - Rules are pipe-separated, each prefixed with a count condition
//   - {0} = zero, {1} = one, [2,*] = two or more, [2,10] = range
func (t *Translator) TransChoice(key string, count int, replace map[string]string) string {
	t.Load(t.locale)

	line, ok := t.lines[t.locale][key]
	if !ok {
		if t.fallbackLocale != t.locale {
			t.Load(t.fallbackLocale)
			line, ok = t.lines[t.fallbackLocale][key]
		}
		if !ok {
			return key
		}
	}

	// Parse plural rules
	choice := t.choosePluralForm(line, count)

	// Replace placeholders
	choice = strings.ReplaceAll(choice, ":count", formatCount(count))

	for k, v := range replace {
		choice = strings.ReplaceAll(choice, ":"+k, v)
	}

	return choice
}

// choosePluralForm selects the correct plural form based on CLDR-like rules.
func (t *Translator) choosePluralForm(line string, count int) string {
	// Split on pipe
	parts := strings.Split(line, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check for condition prefix
		if idx := strings.Index(part, "}"); idx != -1 {
			condition := part[:idx+1]
			text := part[idx+1:]
			text = strings.TrimSpace(text)

			if t.matchesCondition(condition, count) {
				return text
			}
		} else if idx := strings.Index(part, "]"); idx != -1 {
			condition := part[:idx+1]
			text := part[idx+1:]
			text = strings.TrimSpace(text)

			if t.matchesCondition(condition, count) {
				return text
			}
		}
	}

	// Default: return the last part or the whole string
	if len(parts) > 0 {
		return strings.TrimSpace(parts[len(parts)-1])
	}
	return line
}

// matchesCondition checks if a count matches a CLDR-like condition.
func (t *Translator) matchesCondition(condition string, count int) bool {
	condition = strings.TrimSpace(condition)

	// {0} = exact zero
	if condition == "{0}" {
		return count == 0
	}

	// {1} = exact one
	if condition == "{1}" {
		return count == 1
	}

	// {n} = exact number
	if strings.HasPrefix(condition, "{") && strings.HasSuffix(condition, "}") {
		num := condition[1 : len(condition)-1]
		if n, err := parseConditionNumber(num); err == nil {
			return count == n
		}
	}

	// [min,max] = range
	if strings.HasPrefix(condition, "[") && strings.HasSuffix(condition, "]") {
		inner := condition[1 : len(condition)-1]
		return t.matchesRange(inner, count)
	}

	return false
}

// matchesRange checks if a count falls within a range like "2,*" or "2,10".
func (t *Translator) matchesRange(inner string, count int) bool {
	parts := strings.SplitN(inner, ",", 2)
	if len(parts) != 2 {
		return false
	}

	minStr := strings.TrimSpace(parts[0])
	maxStr := strings.TrimSpace(parts[1])

	min, minErr := parseConditionNumber(minStr)
	max, maxErr := parseConditionNumber(maxStr)

	if minErr != nil {
		return false
	}

	// Wildcard max
	if maxStr == "*" || maxStr == "Inf" {
		return count >= min
	}

	if maxErr != nil {
		return false
	}

	return count >= min && count <= max
}

func parseConditionNumber(s string) (int, error) {
	s = strings.TrimSpace(s)
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return 0, &strconv.NumError{Func: "parseConditionNumber", Num: s, Err: strconv.ErrSyntax}
		}
	}
	return n, nil
}

func formatCount(count int) string {
	return strings.TrimSpace(strings.Repeat(" ", 0) + strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					formatInt(count), " ", "", -1,
				), "\n", "", -1,
			), "\r", "", -1,
		), "\t", "", -1,
	))
}

func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + formatInt(-n)
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
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

// TransChoice is the global pluralization helper.
// Usage: TransChoice("apples", 5, map[string]string{})
func TransChoice(key string, count int, replaces ...map[string]string) string {
	if defaultTranslator == nil {
		return key
	}
	var replace map[string]string
	if len(replaces) > 0 {
		replace = replaces[0]
	}
	return defaultTranslator.TransChoice(key, count, replace)
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

// SetFallbackLocale sets the fallback locale globally.
func SetFallbackLocale(locale string) {
	if defaultTranslator != nil {
		defaultTranslator.SetFallbackLocale(locale)
	}
}

// DetectLocaleFromRequest detects the locale from the HTTP request's Accept-Language header.
func DetectLocaleFromRequest(r *http.Request) string {
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang == "" {
		return "en"
	}

	// Parse Accept-Language header (e.g., "en-US,en;q=0.9,fr;q=0.8")
	parts := strings.Split(acceptLang, ",")
	if len(parts) == 0 {
		return "en"
	}

	// Get the primary language (first entry)
	primary := strings.TrimSpace(parts[0])
	// Remove quality factor if present
	if idx := strings.Index(primary, ";"); idx != -1 {
		primary = strings.TrimSpace(primary[:idx])
	}

	// Normalize: en-US -> en, pt-BR -> pt_BR
	primary = strings.Replace(primary, "-", "_", 2)
	primary = strings.ToLower(primary)

	return primary
}

// DetectLocaleAndSet detects locale from request and sets it globally.
func DetectLocaleAndSet(r *http.Request) string {
	locale := DetectLocaleFromRequest(r)
	SetLocale(locale)
	return locale
}


