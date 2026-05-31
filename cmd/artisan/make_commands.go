package artisan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

// titleCase capitalizes the first letter of each word (replaces deprecated titleCase).
func titleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if i == 0 || runes[i-1] == ' ' || runes[i-1] == '_' || runes[i-1] == '-' {
			runes[i] = unicode.ToUpper(runes[i])
		}
	}
	return string(runes)
}

// generateFile is a helper for scaffold commands.
func generateFile(path, content string) {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("File already exists: %s\n", path)
		return
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %s\n", path, err)
		return
	}
	fmt.Printf("Created: %s\n", path)
}

// MakeControllerCmd scaffolds a new controller.
var MakeControllerCmd = &cobra.Command{
	Use:   "make:controller [name]",
	Short: "Create a new controller class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		isResource, _ := cmd.Flags().GetBool("resource")
		isApi, _ := cmd.Flags().GetBool("api")
		path := fmt.Sprintf("app/Http/Controllers/%s.go", name) // Standardized to app/Http/Controllers

		stub := `package controllers

import (
	"net/http"

	"github.com/mechneerd/gow/http/request"
)

type ` + name + ` struct{}
`

		if isApi {
			stub += `
// Index returns a listing of the resource.
func (c *` + name + `) Index(w http.ResponseWriter, r *http.Request) {}

// Show returns the specified resource.
func (c *` + name + `) Show(w http.ResponseWriter, r *http.Request) {}

// Store stores a newly created resource.
func (c *` + name + `) Store(req *request.FormRequest, w http.ResponseWriter, r *http.Request) {}

// Update updates the specified resource.
func (c *` + name + `) Update(req *request.FormRequest, w http.ResponseWriter, r *http.Request) {}

// Destroy removes the specified resource.
func (c *` + name + `) Destroy(w http.ResponseWriter, r *http.Request) {}
`
		} else if isResource {
			stub += `
// Index displays a listing of the resource.
func (c *` + name + `) Index(w http.ResponseWriter, r *http.Request) {
	// 
}

// Show displays the specified resource.
func (c *` + name + `) Show(w http.ResponseWriter, r *http.Request) {}

// Store stores a newly created resource in storage.
func (c *` + name + `) Store(req *request.FormRequest, w http.ResponseWriter, r *http.Request) {}

// Update updates the specified resource in storage.
func (c *` + name + `) Update(req *request.FormRequest, w http.ResponseWriter, r *http.Request) {}

// Destroy removes the specified resource from storage.
func (c *` + name + `) Destroy(w http.ResponseWriter, r *http.Request) {}
`
		} else {
			stub += `
// Index handles the default action.
func (c *` + name + `) Index(w http.ResponseWriter, r *http.Request) {
	// 
}
`
		}

		generateFile(path, stub)
	},
}

// MakeModelCmd scaffolds a new model.
var MakeModelCmd = &cobra.Command{
	Use:   "make:model [name]",
	Short: "Create a new model using the ORM",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		withMigration, _ := cmd.Flags().GetBool("migration")
		path := fmt.Sprintf("app/Models/%s.go", name) // Standardized to app/Models

		stub := `package Models

import "github.com/mechneerd/gow/database/orm"

type {Name} struct {
	orm.Model
	// Add your fields here
	// Name string
}

func ({Name}) TableName() string {
	return "{table}"
}
`
		content := strings.ReplaceAll(stub, "{Name}", name)
		tableName := strings.ToLower(name) + "s"
		content = strings.ReplaceAll(content, "{table}", tableName)

		generateFile(path, content)

		if withMigration {
			// Automatically generate a basic migration file
			migrationName := "create_" + strings.ToLower(name) + "s_table"
			timestamp := time.Now().Format("2006_01_02_150405")
			className := "Create" + titleCase(strings.ReplaceAll(migrationName, "_", "")) + "Table"
			filename := fmt.Sprintf("%s_%s.go", timestamp, migrationName)
			migPath := fmt.Sprintf("database/migrations/%s", filename)

			migStub := `package migrations

import (
	"github.com/mechneerd/gow/database/migration"
	"github.com/mechneerd/gow/database/schema"
)

type ` + className + ` struct{}

func (` + className + `) Up(m *schema.Builder) error {
	return m.Create("` + strings.ToLower(name) + `s", func(table *schema.Blueprint) {
		table.ID()
		// Add your columns here (only use currently supported methods)
		// table.String("name", 255)
		// table.String("email", 255).Unique()
		// table.Integer("age")
		// table.Boolean("is_active").Default(true)
		// table.Text("description").Nullable()
		table.Timestamps()
	})
}

func (` + className + `) Down(m *schema.Builder) error {
	return m.Drop("` + strings.ToLower(name) + `s")
}

func init() {
	migration.Register("` + timestamp + `_` + migrationName + `", ` + className + `{})
}
`
			generateFile(migPath, migStub)
			fmt.Printf("  → Created migration: %s\n", filename)
		}
	},
}

// MakeMigrationCmd scaffolds a new migration.
var MakeMigrationCmd = &cobra.Command{
	Use:   "make:migration [name]",
	Short: "Create a new database migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		timestamp := time.Now().Format("2006_01_02_150405")
		className := "Create" + titleCase(strings.ReplaceAll(name, "_", "")) + "Table"
		filename := fmt.Sprintf("%s_%s.go", timestamp, name)
		path := fmt.Sprintf("database/migrations/%s", filename) // Standardized to database/migrations (lowercase)

		stub := `package migrations

import (
	"github.com/mechneerd/gow/database/migration"
	"github.com/mechneerd/gow/database/schema"
)

type ` + className + ` struct{}

func (` + className + `) Up(m *schema.Builder) error {
	return m.Create("` + strings.ToLower(name) + `s", func(table *schema.Blueprint) {
		table.ID()
		// Add columns here (only use currently supported methods):
		// table.String("name", 255)
		// table.String("email", 255).Unique()
		// table.Integer("age")
		// table.Boolean("is_active").Default(true)
		// table.Text("description").Nullable()
		table.Timestamps()
	})
}

func (` + className + `) Down(m *schema.Builder) error {
	return m.Drop("` + strings.ToLower(name) + `s")
}

func init() {
	migration.Register("` + timestamp + `_` + name + `", ` + className + `{})
}
`
		generateFile(path, stub)
	},
}

// MakeMiddlewareCmd scaffolds a new middleware.
var MakeMiddlewareCmd = &cobra.Command{
	Use:   "make:middleware [name]",
	Short: "Create a new middleware class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := fmt.Sprintf("app/Http/Middleware/%s.go", name) // Standardized to app/Http/Middleware
		
		stub := `package middleware

import (
	"net/http"
)

func {Name}() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pre-middleware logic
			
			next.ServeHTTP(w, r)
			
			// Post-middleware logic
		})
	}
}
`
		content := strings.ReplaceAll(stub, "{Name}", name)
		generateFile(path, content)
	},
}

// MakeJobCmd scaffolds a new queue job.
var MakeJobCmd = &cobra.Command{
	Use:   "make:job [name]",
	Short: "Create a new job class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := fmt.Sprintf("app/Jobs/%s.go", name)
		
		stub := `package jobs

type {Name} struct {
	// payload fields
}

func (j *{Name}) Handle() error {
	// Job logic here
	return nil
}

func (j *{Name}) Failed(err error) {
	// Handle failure
}
`
		content := strings.ReplaceAll(stub, "{Name}", name)
		generateFile(path, content)
	},
}

// MakeCommandCmd scaffolds a new custom artisan command.
var MakeCommandCmd = &cobra.Command{
	Use:   "make:command [name]",
	Short: "Create a new artisan command (meta)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !strings.HasSuffix(name, "Cmd") {
			name += "Cmd"
		}

		path := fmt.Sprintf("app/Console/Commands/%s.go", name)

		stub := `package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ` + name + ` = &cobra.Command{
	Use:   "` + strings.ToLower(strings.TrimSuffix(name, "Cmd")) + `",
	Short: "Description of ` + name + `",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("` + name + ` executed successfully!")
	},
}

// Register this command in your console kernel (artisan.go or console/kernel.go)
`
		generateFile(path, stub)
	},
}

// MakeViewCmd scaffolds a new view file using the preferred .goblade extension.
var MakeViewCmd = &cobra.Command{
	Use:   "make:view [name]",
	Short: "Create a new view (uses .goblade by default for Go-native feel)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		// Convert dot notation to path, e.g. auth.login → auth/login.goblade
		relPath := strings.ReplaceAll(name, ".", "/")
		path := fmt.Sprintf("resources/views/%s.goblade", relPath)

		stub := `{{/* resources/views/` + relPath + `.goblade */}}

<div>
    <h1>{{ .Title }}</h1>
    {{-- Add your GoBlade content here --}}
</div>
`
		generateFile(path, stub)
	},
}

func init() {
	MakeControllerCmd.Flags().Bool("resource", false, "Generate a resource controller class")
	MakeControllerCmd.Flags().Bool("api", false, "Generate an API controller class")
	MakeModelCmd.Flags().Bool("migration", false, "Also create a migration for the model")
}

