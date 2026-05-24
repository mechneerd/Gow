package view

import (
	"regexp"
	"strings"
)

// parseAttributes turns a raw attribute string like `type="error" :title="$name"`
// into a map. This is a simplified parser.
func parseAttributes(attrString string) map[string]any {
	result := make(map[string]any)
	if strings.TrimSpace(attrString) == "" {
		return result
	}

	// Very naive parser — good enough for most cases
	re := regexp.MustCompile(`([:\w-]+)\s*=\s*["']([^"']*)["']`)
	matches := re.FindAllStringSubmatch(attrString, -1)

	for _, m := range matches {
		key := strings.TrimPrefix(m[1], ":")
		value := m[2]

		// If it starts with :, we treat it as an expression reference (for now just store the string)
		// In a more advanced system we would evaluate it against the current context.
		result[key] = value
	}

	return result
}

