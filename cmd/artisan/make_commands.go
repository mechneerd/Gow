package artisan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

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
		path := fmt.Sprintf("app/http/controllers/%s.go", name)
		
		stub := `package controllers

import (
	"net/http"
)

type {Name} struct {}

func (c *{Name}) Index(w http.ResponseWriter, r *http.Request) {
	// 
}
`
		if isResource {
			stub += `
func (c *{Name}) Show(w http.ResponseWriter, r *http.Request) {}
func (c *{Name}) Store(w http.ResponseWriter, r *http.Request) {}
func (c *{Name}) Update(w http.ResponseWriter, r *http.Request) {}
func (c *{Name}) Destroy(w http.ResponseWriter, r *http.Request) {}
`
		}

		content := strings.ReplaceAll(stub, "{Name}", name)
		generateFile(path, content)
	},
}

// MakeModelCmd scaffolds a new model.
var MakeModelCmd = &cobra.Command{
	Use:   "make:model [name]",
	Short: "Create a new Goquent model class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		path := fmt.Sprintf("app/models/%s.go", name)
		
		stub := `package models

type {Name} struct {
	ID        int    ` + "`db:\"id\"`" + `
	CreatedAt string ` + "`db:\"created_at\"`" + `
	UpdatedAt string ` + "`db:\"updated_at\"`" + `
}
`
		content := strings.ReplaceAll(stub, "{Name}", name)
		generateFile(path, content)
	},
}

// MakeMigrationCmd scaffolds a new migration.
var MakeMigrationCmd = &cobra.Command{
	Use:   "make:migration [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		timestamp := time.Now().Format("20060102150405")
		filename := fmt.Sprintf("%s_%s.go", timestamp, name)
		path := fmt.Sprintf("database/migrations/%s", filename)
		
		stub := `package migrations

import "gow/database/schema"

func init() {
	schema.Register(
		// Up
		func(b *schema.Blueprint) {
			// b.Create("table_name", func(table *schema.Table) {
			// 	table.Increments("id")
			// 	table.Timestamps()
			// })
		},
		// Down
		func(b *schema.Blueprint) {
			// b.DropIfExists("table_name")
		},
	)
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
		path := fmt.Sprintf("app/http/middleware/%s.go", name)
		
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
		path := fmt.Sprintf("app/jobs/%s.go", name)
		
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

func init() {
	MakeControllerCmd.Flags().Bool("resource", false, "Generate a resource controller class")
}
