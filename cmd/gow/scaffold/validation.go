package scaffold

import (
	"fmt"
	"strings"
)

// ValidDatabases defines the supported database drivers.
var ValidDatabases = []string{"sqlite", "mysql", "postgres"}

// IsValidDatabase checks if the given driver is supported.
func IsValidDatabase(driver string) bool {
	driver = strings.ToLower(strings.TrimSpace(driver))
	for _, v := range ValidDatabases {
		if v == driver {
			return true
		}
	}
	return false
}

// NormalizeDatabase returns the normalized (lowercase) driver name if valid.
func NormalizeDatabase(driver string) (string, error) {
	driver = strings.ToLower(strings.TrimSpace(driver))
	if IsValidDatabase(driver) {
		return driver, nil
	}
	return "", fmt.Errorf("invalid database driver '%s'. Valid options: %s", driver, strings.Join(ValidDatabases, ", "))
}
