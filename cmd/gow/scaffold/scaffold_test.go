package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidDatabase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"sqlite", true},
		{"SQLite", true},
		{"mysql", true},
		{"postgres", true},
		{"mongodb", false},
		{"", false},
		{"  Postgres  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidDatabase(tt.input); got != tt.expected {
				t.Errorf("IsValidDatabase(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalizeDatabase(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"sqlite", "sqlite", false},
		{"MySQL", "mysql", false},
		{"POSTGRES", "postgres", false},
		{"mariadb", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := NormalizeDatabase(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("NormalizeDatabase(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("NormalizeDatabase(%q) unexpected error: %v", tt.input, err)
				}
				if got != tt.expected {
					t.Errorf("NormalizeDatabase(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			}
		})
	}
}

func TestReplacePlaceholders(t *testing.T) {
	// Create a temporary directory with a file containing placeholders
	tmpDir, err := os.MkdirTemp("", "gow-test-replacer-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a sample file with placeholders
	content := `module {{ .ModulePath }}

APP_NAME={{ .AppName }}
DB={{ .DatabaseDriver }}
YEAR={{ .Year }}
`
	filePath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Also create a .template file to test renaming
	templateFile := filepath.Join(tmpDir, "go.mod.template")
	if err := os.WriteFile(templateFile, []byte("module {{ .ModulePath }}"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := ReplaceContext{
		AppName:        "MyCoolApp",
		ModulePath:     "github.com/test/myapp",
		DatabaseDriver: "postgres",
		Year:           "2026",
	}

	if err := ReplacePlaceholders(tmpDir, ctx); err != nil {
		t.Fatalf("ReplacePlaceholders failed: %v", err)
	}

	// Check main file
	updated, _ := os.ReadFile(filePath)
	expected := `module github.com/test/myapp

APP_NAME=MyCoolApp
DB=postgres
YEAR=2026
`
	if string(updated) != expected {
		t.Errorf("Replacement mismatch.\nGot:\n%s\nWant:\n%s", updated, expected)
	}

	// Check that .template file was renamed
	newGoMod := filepath.Join(tmpDir, "go.mod")
	if _, err := os.Stat(newGoMod); os.IsNotExist(err) {
		t.Error("Expected go.mod to be created from go.mod.template")
	}
}

func TestDefaultReplaceContext(t *testing.T) {
	ctx := DefaultReplaceContext("TestApp")

	if ctx.AppName != "TestApp" {
		t.Errorf("Expected AppName to be 'TestApp', got %s", ctx.AppName)
	}
	if ctx.DatabaseDriver != "sqlite" {
		t.Errorf("Expected default DatabaseDriver to be 'sqlite'")
	}
	if ctx.Year == "" {
		t.Error("Expected Year to be set")
	}
}

// Note: Remote cloning tests are intentionally limited to avoid network dependency in CI.
// Local path testing is covered in TestPrepareSkeleton_LocalPath.

func TestPrepareSkeleton_LocalPath(t *testing.T) {
	// Create a fake local skeleton
	tmpSkeleton, err := os.MkdirTemp("", "fake-skeleton-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpSkeleton)

	// Add a minimal structure
	os.MkdirAll(filepath.Join(tmpSkeleton, "templates/web"), 0755)
	os.WriteFile(filepath.Join(tmpSkeleton, "templates/web/go.mod.template"), []byte("module {{ .ModulePath }}"), 0644)

	// Use PrepareSkeleton with local path
	resultPath, err := PrepareSkeleton(tmpSkeleton)
	if err != nil {
		t.Fatalf("PrepareSkeleton failed for local path: %v", err)
	}
	defer CleanupTemp(resultPath)

	// Verify the structure was copied
	if _, err := os.Stat(filepath.Join(resultPath, "templates/web/go.mod.template")); os.IsNotExist(err) {
		t.Error("Expected skeleton structure to be copied")
	}
}

func TestInjectRBACBootstrapExamples(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "rbac-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	bootstrapDir := filepath.Join(tmpDir, "bootstrap")
	os.MkdirAll(bootstrapDir, 0755)

	original := `package main

func main() { }
`
	os.WriteFile(filepath.Join(bootstrapDir, "app.go"), []byte(original), 0644)

	err = InjectRBACBootstrapExamples(tmpDir)
	if err != nil {
		t.Fatalf("InjectRBACBootstrapExamples failed: %v", err)
	}

	content, _ := os.ReadFile(filepath.Join(bootstrapDir, "app.go"))
	if !strings.Contains(string(content), "RBAC + Auth Middleware") {
		t.Error("Expected RBAC examples to be injected into bootstrap/app.go")
	}
}

