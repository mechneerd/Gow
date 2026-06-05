package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/mechneerd/gow/bootstrap"
	gowhttp "github.com/mechneerd/gow/http"
	"github.com/mechneerd/gow/routing"

	"github.com/mechneerd/gow/cmd/artisan"
	scaffoldpkg "github.com/mechneerd/gow/cmd/gow/scaffold"

	// Database drivers (required for migrate/seed commands to work)
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// version is set via ldflags at build time by goreleaser.
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "gow",
	Short: "GoW framework CLI",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application on a local development server",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		// If the project has its own main.go, run it directly
		if _, err := os.Stat(filepath.Join(cwd, "main.go")); err == nil {
			fmt.Println("🚀  Starting application server...")
			fmt.Println("    Press Ctrl+C to stop.")

			serveCmd := exec.Command("go", "run", "main.go")
			serveCmd.Dir = cwd
			serveCmd.Stdout = os.Stdout
			serveCmd.Stderr = os.Stderr
			serveCmd.Stdin = os.Stdin

			if err := serveCmd.Run(); err != nil {
				fmt.Printf("Server error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Fallback: framework default server (no user main.go found)
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
		addr := ":" + port

		// Boot the framework application (loads config, providers, etc.)
		app := bootstrap.NewApplication(".")

		// Create router
		router := routing.NewRouter()

		// Default welcome route
		router.Get("/", func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<h1>Welcome to GoW Framework!</h1>
<p>Server is running on <strong>http://localhost%s</strong></p>
<p>Add your routes in <code>routes/web.go</code> or your <code>main.go</code>.</p>`, addr)
			return nil
		})

		// Create HTTP kernel with graceful shutdown support
		kernel := gowhttp.NewKernel(app, router)

		fmt.Printf("🚀  GoW development server started at http://localhost%s\n", addr)
		fmt.Println("    Press Ctrl+C to stop.")

		if err := kernel.Serve(addr); err != nil {
			fmt.Printf("Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the framework version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoW Framework version %s\n", version)
	},
}

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new GoW application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		yes := cmd.Flag("yes").Changed

		hasStarterKit := cmd.Flag("minimal").Changed ||
			cmd.Flag("minimal-api").Changed ||
			cmd.Flag("api").Changed ||
			cmd.Flag("web").Changed ||
			cmd.Flag("auth").Changed ||
			cmd.Flag("full").Changed ||
			cmd.Flag("admin-panel").Changed ||
			cmd.Flag("with-docker").Changed ||
			cmd.Flag("inertia-react").Changed ||
			cmd.Flag("inertia-vue").Changed ||
			cmd.Flag("livewire").Changed

		var result scaffoldpkg.PromptResult

		if !hasStarterKit && !yes {
			// Interactive mode (only when nothing was specified and --yes was not used)
			var err error
			result, err = scaffoldpkg.RunInteractiveWizard(name)
			if err != nil {
				fmt.Printf("Error during interactive setup: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Non-interactive mode (flags provided or --yes used)
			flags := map[string]bool{
				"minimal":       cmd.Flag("minimal").Changed,
				"minimal-api":   cmd.Flag("minimal-api").Changed,
				"api":           cmd.Flag("api").Changed,
				"web":           cmd.Flag("web").Changed,
				"auth":          cmd.Flag("auth").Changed,
				"full":          cmd.Flag("full").Changed,
				"admin-panel":   cmd.Flag("admin-panel").Changed,
				"with-docker":   cmd.Flag("with-docker").Changed,
				"inertia-react": cmd.Flag("inertia-react").Changed,
				"inertia-vue":   cmd.Flag("inertia-vue").Changed,
				"livewire":      cmd.Flag("livewire").Changed,
			}

			result.StarterKit = getStarterKitFromFlags(flags)

			// If --yes was used without a specific starter kit, default to the recommended one ("auth")
			if result.StarterKit == "" && yes {
				result.StarterKit = "auth"
				fmt.Println("→ Using default configuration (Web + Auth + SQLite)")
			}

			result.Database = cmd.Flag("db").Value.String()
			result.ModulePath = cmd.Flag("module").Value.String()

			// Validate database driver
			if result.Database != "" {
				normalized, err := scaffoldpkg.NormalizeDatabase(result.Database)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				result.Database = normalized
			}
		}

		force := cmd.Flag("force").Changed
		noGit := cmd.Flag("no-git").Changed
		skeletonURL := cmd.Flag("skeleton").Value.String()

		err := scaffoldWithOptions(name, result, force, noGit, skeletonURL)
		if err != nil {
			fmt.Printf("Error scaffolding project: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add flags to newCmd
	newCmd.Flags().Bool("minimal", false, "Create a minimal project (basic routing + views)")
	newCmd.Flags().Bool("minimal-api", false, "Create an ultra-light API project")
	newCmd.Flags().Bool("api", false, "Create an API-only project (with Sanctum)")
	newCmd.Flags().Bool("web", false, "Create a full web app with Blade views")
	newCmd.Flags().Bool("auth", false, "Create a full web app with authentication + RBAC")
	newCmd.Flags().Bool("full", false, "Create a full-stack project (web + API + auth + RBAC)")
	newCmd.Flags().Bool("admin-panel", false, "Create an admin panel project (planned)")
	newCmd.Flags().Bool("with-docker", false, "Create a Dockerized project (planned)")
	newCmd.Flags().Bool("inertia-react", false, "Create an Inertia.js + React project (planned)")
	newCmd.Flags().Bool("inertia-vue", false, "Create an Inertia.js + Vue 3 project (planned)")
	newCmd.Flags().Bool("livewire", false, "Create a Livewire-focused project (planned)")
	newCmd.Flags().String("module", "", "Module path for go.mod (e.g. github.com/username/myapp)")
	newCmd.Flags().String("db", "sqlite", "Database driver: sqlite, mysql, postgres")
	newCmd.Flags().Bool("force", false, "Overwrite existing directory if it already exists")
	newCmd.Flags().Bool("no-git", false, "Skip git repository initialization")
	newCmd.Flags().String("skeleton", "", "Custom skeleton source (remote URL or local path). Useful for custom templates.")
	newCmd.Flags().Bool("yes", false, "Accept all defaults and skip interactive prompts (useful for CI/scripts)")

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	// Wire important Artisan commands into the main gow CLI
	rootCmd.AddCommand(artisan.MigrateCmd)
	rootCmd.AddCommand(artisan.MigrateRollbackCmd)
	rootCmd.AddCommand(artisan.MigrateStatusCmd)
	rootCmd.AddCommand(artisan.MigrateRunCmd)
	rootCmd.AddCommand(artisan.MigrateFreshCmd)
	rootCmd.AddCommand(artisan.MigrateRefreshCmd)

	rootCmd.AddCommand(artisan.DbSeedCmd)
	rootCmd.AddCommand(artisan.KeyGenerateCmd)
	rootCmd.AddCommand(artisan.AboutCmd)
	rootCmd.AddCommand(artisan.ListCmd)

	// Make commands (most commonly used)
	rootCmd.AddCommand(artisan.MakeControllerCmd)
	rootCmd.AddCommand(artisan.MakeModelCmd)
	rootCmd.AddCommand(artisan.MakeMigrationCmd)
	rootCmd.AddCommand(artisan.MakeMiddlewareCmd)
	rootCmd.AddCommand(artisan.MakeRequestCmd)
	rootCmd.AddCommand(artisan.MakeSeederCmd)
	rootCmd.AddCommand(artisan.MakeCommandCmd)

	// Additional useful commands
	rootCmd.AddCommand(artisan.EnvCmd)
	rootCmd.AddCommand(artisan.RouteListCmd)
	rootCmd.AddCommand(artisan.CacheClearCmd)
	rootCmd.AddCommand(artisan.CacheForgetCmd)
	rootCmd.AddCommand(artisan.ViewClearCmd)
	rootCmd.AddCommand(artisan.ScheduleListCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// getStarterKitFromFlags converts flag map into starter kit name
func getStarterKitFromFlags(flags map[string]bool) string {
	if flags["minimal"] {
		return "minimal"
	}
	if flags["minimal-api"] {
		return "minimal-api"
	}
	if flags["api"] {
		return "api"
	}
	if flags["web"] {
		return "web"
	}
	if flags["auth"] {
		return "auth"
	}
	if flags["full"] {
		return "full"
	}
	if flags["admin-panel"] {
		return "admin-panel"
	}
	if flags["with-docker"] {
		return "with-docker"
	}
	if flags["inertia-react"] {
		return "inertia-react"
	}
	if flags["inertia-vue"] {
		return "inertia-vue"
	}
	if flags["livewire"] {
		return "livewire"
	}
	return "auth" // default
}

// scaffoldWithOptions is the new entry point that uses the gow-skeleton repository.
// The skeletonURL parameter allows using a custom repository (experimental, for future release).
func scaffoldWithOptions(name string, result scaffoldpkg.PromptResult, force bool, noGit bool, skeletonURL string) error {
	// Show GoW banner
	scaffoldpkg.ShowBanner()

	// Convert prompt result to flags for selector
	flags := map[string]bool{
		"minimal":       result.StarterKit == "minimal",
		"minimal-api":   result.StarterKit == "minimal-api",
		"api":           result.StarterKit == "api",
		"web":           result.StarterKit == "web",
		"auth":          result.StarterKit == "auth",
		"full":          result.StarterKit == "full",
		"admin-panel":   result.StarterKit == "admin-panel",
		"with-docker":   result.StarterKit == "with-docker",
		"inertia-react": result.StarterKit == "inertia-react",
		"inertia-vue":   result.StarterKit == "inertia-vue",
		"livewire":      result.StarterKit == "livewire",
	}

	templateRelative := scaffoldpkg.SelectTemplate(flags)

	totalSteps := 7
	completeStep := scaffoldpkg.ShowProgress(totalSteps)

	fmt.Printf("🚀  Creating GoW project \"%s\"...\n\n", name)

	// Step 1: Fetch template
	completeStep("Scaffolding project")

	clonedPath, err := scaffoldpkg.PrepareSkeleton(skeletonURL)
	if err != nil {
		return err
	}
	defer scaffoldpkg.CleanupTemp(clonedPath)
	completeStep("Fetching template")

	targetDir := name

	if _, err := os.Stat(targetDir); err == nil {
		if !force {
			return fmt.Errorf("directory %q already exists. Use --force to overwrite", targetDir)
		}
		fmt.Printf("→ Removing existing directory: %s\n", targetDir)
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	fullTemplatePath := scaffoldpkg.GetTemplateFullPath(clonedPath, templateRelative)
	if err := scaffoldpkg.CopyTemplate(fullTemplatePath, targetDir); err != nil {
		return err
	}
	completeStep("Configuring project")

	// 3. Replace placeholders and rename .template files
	ctx := scaffoldpkg.DefaultReplaceContext(name)
	if result.ModulePath != "" {
		ctx.ModulePath = result.ModulePath
	}
	if result.Database != "" {
		ctx.DatabaseDriver = result.Database
	}
	if err := scaffoldpkg.ReplacePlaceholders(targetDir, ctx); err != nil {
		return err
	}

	// Fix known skeleton template bugs (duplicate functions, broken imports, etc.)
	if err := scaffoldpkg.FixSkeletonBugs(targetDir, ctx.ModulePath); err != nil {
		fmt.Printf("   ⚠️  Skeleton fixup warning: %v\n", err)
	}

	// Inject RBAC middleware examples + DB wiring guidance into bootstrap/app.go
	// for auth-enabled kits (addresses high-priority gap in generated projects).
	if flags["auth"] || flags["full"] {
		_ = scaffoldpkg.InjectRBACBootstrapExamples(targetDir)
	}
	completeStep("Applying fixes")

	err = scaffoldpkg.RunPostInstall(targetDir, scaffoldpkg.PostInstallOptions{
		RunGoModTidy: true,
		CopyEnv:      true,
	})
	if err != nil {
		return err
	}
	completeStep("Installing dependencies")
	completeStep("Preparing environment")

	// 5. Initialize git (unless --no-git is passed)
	if !noGit {
		cmd := exec.Command("git", "init")
		cmd.Dir = targetDir
		if err := cmd.Run(); err != nil {
			completeStep("Initializing git (skipped)")
		} else {
			completeStep("Initializing git")
		}
	} else {
		completeStep("Initializing git (skipped)")
	}

	scaffoldpkg.PrintNextSteps(name)

	return nil
}


