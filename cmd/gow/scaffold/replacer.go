package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ReplaceContext holds the values to substitute in templates.
type ReplaceContext struct {
	AppName        string
	ModulePath     string
	DatabaseDriver string
	Year           string
}

// ReplacePlaceholders walks the project directory and replaces all placeholders.
// It also renames files ending with .template (e.g. go.mod.template → go.mod)
func ReplacePlaceholders(rootDir string, ctx ReplaceContext) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		original := string(content)
		updated := replaceInContent(original, ctx)

		if updated != original {
			if err := os.WriteFile(path, []byte(updated), info.Mode()); err != nil {
				return err
			}
		}

		// Rename *.template files
		if strings.HasSuffix(path, ".template") {
			newPath := strings.TrimSuffix(path, ".template")
			if err := os.Rename(path, newPath); err != nil {
				return err
			}
		}

		return nil
	})
}

func replaceInContent(content string, ctx ReplaceContext) string {
	replacements := map[string]string{
		"{{ .AppName }}":        ctx.AppName,
		"{{ .ModulePath }}":     ctx.ModulePath,
		"{{ .DatabaseDriver }}": ctx.DatabaseDriver,
		"{{ .Year }}":           ctx.Year,
	}

	result := content
	for placeholder, value := range replacements {
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// DefaultReplaceContext creates a context with sensible defaults.
func DefaultReplaceContext(appName string) ReplaceContext {
	return ReplaceContext{
		AppName:        appName,
		ModulePath:     appName,
		DatabaseDriver: "sqlite",
		Year:           time.Now().Format("2006"),
	}
}

