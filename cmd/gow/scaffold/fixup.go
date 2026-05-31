package scaffold

import (
	"os"
	"path/filepath"
	"strings"
)

// FixSkeletonBugs applies targeted fixes for known bugs in the gow-skeleton templates.
// These are structural issues that simple placeholder replacement cannot handle.
func FixSkeletonBugs(projectDir string, moduleName string) error {
	return filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		original := string(content)
		updated := fixFileContent(original, path, moduleName)

		if updated != original {
			return os.WriteFile(path, []byte(updated), info.Mode())
		}
		return nil
	})
}

func fixFileContent(content string, filePath string, moduleName string) string {
	result := content

	// Fix 1: config/auth.go — remove duplicate getEnv function
	// The skeleton has getEnv in both config/app.go and config/auth.go
	if strings.HasSuffix(filePath, filepath.Join("config", "auth.go")) ||
		strings.HasSuffix(filePath, "config\\auth.go") || strings.HasSuffix(filePath, "config/auth.go") {
		result = removeDuplicateGetEnv(result)
	}

	// Fix 2: bootstrap/app.go — remove unused "routes" import
	if strings.HasSuffix(filePath, filepath.Join("bootstrap", "app.go")) ||
		strings.HasSuffix(filePath, "bootstrap\\app.go") || strings.HasSuffix(filePath, "bootstrap/app.go") {
		result = removeUnusedRoutesImport(result)
	}

	// Fix 3: app/Livewire/Counter.go — add livewire import, fix BaseComponent reference
	if strings.Contains(filePath, "Livewire") && strings.HasSuffix(filePath, "Counter.go") {
		result = fixLivewireCounter(result)
	}

	// Fix 4: app/Models/User.go — fix broken SetDB method and missing sql import
	if strings.HasSuffix(filePath, filepath.Join("app", "Models", "User.go")) ||
		strings.HasSuffix(filePath, "app\\Models\\User.go") || strings.HasSuffix(filePath, "app/Models/User.go") {
		result = fixUserModel(result)
	}

	// Fix 5: database/seeders/RoleSeeder.go — remove unused Models import
	if strings.Contains(filePath, "seeders") && strings.HasSuffix(filePath, "RoleSeeder.go") {
		result = removeUnusedModelsImport(result)
	}

	return result
}

func removeDuplicateGetEnv(content string) string {
	// If this file has its own getEnv, remove it (it's duplicated from config/app.go)
	if strings.Contains(content, "func getEnv(key, fallback string) string") {
		// Remove the import "os" if present (no longer needed without getEnv)
		content = strings.Replace(content, "import \"os\"\n\n", "", 1)
		// Remove the getEnv function
		idx := strings.Index(content, "\nfunc getEnv(key, fallback string) string")
		if idx == -1 {
			idx = strings.Index(content, "func getEnv(key, fallback string) string")
		}
		if idx != -1 {
			// Find the end of the function (next closing brace at column 0 or end of file)
			endIdx := strings.Index(content[idx+1:], "\n}\n")
			if endIdx != -1 {
				content = content[:idx] + content[idx+1+endIdx+3:]
			}
		}
	}
	return content
}

func removeUnusedRoutesImport(content string) string {
	// Remove "demo/routes" or "<module>/routes" import if unused
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	inImport := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImport = true
		}
		if inImport && strings.Contains(line, "/routes\"") {
			continue // skip the unused import
		}
		if trimmed == ")" && inImport {
			inImport = false
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func fixLivewireCounter(content string) string {
	// Add the livewire import if missing
	if !strings.Contains(content, "github.com/mechneerd/gow/http/livewire") &&
		strings.Contains(content, "BaseComponent") {
		content = strings.Replace(content,
			`import "fmt"`,
			`import (
	"fmt"

	"github.com/mechneerd/gow/http/livewire"
)`, 1)
	}

	// Fix BaseComponent reference to livewire.BaseComponent
	if strings.Contains(content, "\tBaseComponent\n") {
		content = strings.Replace(content, "\tBaseComponent\n", "\tlivewire.BaseComponent\n", 1)
	}

	return content
}

func fixUserModel(content string) string {
	// Remove the broken SetDB method that references unexported field
	idx := strings.Index(content, "\n// SetDB wires")
	if idx != -1 {
		endIdx := strings.Index(content[idx+1:], "\n}\n")
		if endIdx != -1 {
			content = content[:idx] + "\n"
		}
	}

	// Remove unused "database/sql" import
	content = strings.Replace(content, "\t\"database/sql\"\n", "", 1)

	return content
}

func removeUnusedModelsImport(content string) string {
	// Remove the unused Models import
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	inImport := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImport = true
		}
		if inImport && strings.Contains(line, "/Models\"") {
			continue
		}
		if trimmed == ")" && inImport {
			inImport = false
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
