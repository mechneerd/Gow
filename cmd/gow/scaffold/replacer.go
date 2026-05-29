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

	// Fallback: replace hardcoded skeleton placeholder names with the actual module path.
	// Some skeleton templates use literal strings instead of {{ .ModulePath }} placeholders.
	if ctx.ModulePath != "" && ctx.ModulePath != "test-project" {
		result = strings.ReplaceAll(result, "test-project", ctx.ModulePath)
	}

	// Fix framework imports: skeleton templates may use <module>/auth/rbac, <module>/database/orm, etc.
	// instead of the full github.com/mechneerd/gow/... path.
	frameworkPrefixes := []string{
		"auth/rbac",
		"auth/sanctum",
		"auth/access",
		"auth/fortify",
		"auth/password",
		"auth/socialite",
		"auth/verification",
		"database/orm",
		"database/query",
		"database/migration",
		"database/schema",
		"database/factory",
		"database/seeder",
		"database/pagination",
		"container",
		"routing",
		"view",
		"validation",
		"encryption",
		"hashing",
		"session",
		"cache",
		"queue",
		"events",
		"mail",
		"notifications",
		"broadcasting",
		"logging",
		"config",
		"console",
		"cookie",
		"storage",
		"localization",
		"http/client",
		"http/middleware",
		"http/request",
		"http/response",
		"http/resources",
		"http/exception",
		"support/collection",
		"support/str",
		"support/arr",
		"support/pipeline",
		"support/telescope",
		"support/health",
		"support/metrics",
		"support/scout",
		"support/process",
		"support/pennant",
		"support/fakes",
		"support/httpclient",
		"testing",
		"foundation",
		"bootstrap",
	}
	for _, prefix := range frameworkPrefixes {
		// Replace "<module>/prefix" with "github.com/mechneerd/gow/prefix"
		oldImport := ctx.ModulePath + "/" + prefix
		newImport := "github.com/mechneerd/gow/" + prefix
		result = strings.ReplaceAll(result, oldImport, newImport)
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

