package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gow",
	Short: "GoW framework CLI",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve command not yet implemented. Use 'go run cmd/app/main.go' for now.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the framework version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoW Framework version 1.0.0")
	},
}

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new GoW application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := scaffold(name)
		if err != nil {
			fmt.Printf("Error scaffolding project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully created GoW project: %s\n", name)
		fmt.Printf("Run 'cd %s' and 'go run cmd/app/main.go' to start.\n", name)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// scaffold creates the directory structure and base files.
func scaffold(name string) error {
	dirs := []string{
		"cmd/app",
		"routes",
		"app/http/controllers",
		"resources/views",
		"storage",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(name, dir), 0755); err != nil {
			return err
		}
	}

	files := map[string]string{
		"go.mod": fmt.Sprintf(`module %s

go 1.20

require github.com/yourname/gow v1.0.0 // Replace with actual module path
`, name),
		"cmd/app/main.go": `package main

import (
	"log"

	"github.com/yourname/gow/foundation"
	_ "` + name + `/routes" // trigger init() to register routes
)

func main() {
	app := foundation.NewApplication()
	
	// Boot the framework and start the HTTP server
	if err := app.Serve(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
`,
		"routes/web.go": `package routes

import (
	"github.com/yourname/gow/http/router"
	"` + name + `/app/http/controllers"
)

func init() {
	// Web routes
	router.Get("/", controllers.HomeController)
}
`,
		"app/http/controllers/home_controller.go": `package controllers

import (
	"net/http"

	"github.com/yourname/gow/view"
)

func HomeController(w http.ResponseWriter, r *http.Request) error {
	return view.Render(w, "welcome", map[string]any{
		"Title": "Welcome to GoW",
	})
}
`,
		"resources/views/welcome.html": `<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <h1>{{ .Title }}</h1>
    <p>Your GoW application is running.</p>
</body>
</html>
`,
		"storage/.gitkeep": "",
		".env.example": `APP_NAME=GoWApp
APP_ENV=local
APP_KEY=
APP_DEBUG=true
APP_URL=http://localhost
APP_PORT=8080

DB_CONNECTION=sqlite
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=database.sqlite
DB_USERNAME=root
DB_PASSWORD=
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(name, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}
