package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	scaffoldpkg "gow/cmd/gow/scaffold"
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

		yes := cmd.Flag("yes").Changed

		hasStarterKit := cmd.Flag("minimal").Changed ||
			cmd.Flag("api").Changed ||
			cmd.Flag("auth").Changed

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
				"minimal": cmd.Flag("minimal").Changed,
				"api":     cmd.Flag("api").Changed,
				"auth":    cmd.Flag("auth").Changed,
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
	newCmd.Flags().Bool("minimal", false, "Create a minimal project")
	newCmd.Flags().Bool("api", false, "Create an API-only project (with Sanctum)")
	newCmd.Flags().Bool("auth", false, "Create a full web app with authentication")
	newCmd.Flags().String("module", "", "Module path for go.mod (e.g. github.com/username/myapp)")
	newCmd.Flags().String("db", "sqlite", "Database driver: sqlite, mysql, postgres")
	newCmd.Flags().Bool("force", false, "Overwrite existing directory if it already exists")
	newCmd.Flags().Bool("no-git", false, "Skip git repository initialization")
	newCmd.Flags().String("skeleton", "", "Custom skeleton source (remote URL or local path). Useful for custom templates.")
	newCmd.Flags().Bool("yes", false, "Accept all defaults and skip interactive prompts (useful for CI/scripts)")

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

// getStarterKitFromFlags converts flag map into starter kit name
func getStarterKitFromFlags(flags map[string]bool) string {
	if flags["minimal"] {
		return "minimal"
	}
	if flags["api"] {
		return "api"
	}
	if flags["auth"] {
		return "auth"
	}
	return "web"
}

// scaffoldWithOptions is the new entry point that uses the gow-skeleton repository.
// The skeletonURL parameter allows using a custom repository (experimental, for future release).
func scaffoldWithOptions(name string, result scaffoldpkg.PromptResult, force bool, noGit bool, skeletonURL string) error {
	// Convert prompt result to flags for selector
	flags := map[string]bool{
		"minimal": result.StarterKit == "minimal",
		"api":     result.StarterKit == "api",
		"auth":    result.StarterKit == "auth",
	}

	templateRelative := scaffoldpkg.SelectTemplate(flags)

	fmt.Printf("🚀  Creating GoW project \"%s\"...\n\n", name)
	fmt.Println("→ Fetching template...")

	clonedPath, err := scaffoldpkg.PrepareSkeleton(skeletonURL)
	if err != nil {
		return err
	}
	defer scaffoldpkg.CleanupTemp(clonedPath)
	fmt.Println("   ✓ Template ready")

	targetDir := name

	if _, err := os.Stat(targetDir); err == nil {
		if !force {
			return fmt.Errorf("directory %q already exists. Use --force to overwrite", targetDir)
		}
		fmt.Printf("→ Removing existing directory: %s\n", targetDir)
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
		fmt.Println("   ✓ Directory cleared")
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	fullTemplatePath := scaffoldpkg.GetTemplateFullPath(clonedPath, templateRelative)
	if err := scaffoldpkg.CopyTemplate(fullTemplatePath, targetDir); err != nil {
		return err
	}

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

	err = scaffoldpkg.RunPostInstall(targetDir, scaffoldpkg.PostInstallOptions{
		RunGoModTidy: true,
		CopyEnv:      true,
	})
	if err != nil {
		return err
	}

	// 5. Initialize git (unless --no-git is passed)
	if !noGit {
		fmt.Println("→ Initializing git repository...")
		cmd := exec.Command("git", "init")
		cmd.Dir = targetDir
		if err := cmd.Run(); err != nil {
			fmt.Printf("   ⚠️  Git init skipped (git not found or failed)\n")
		} else {
			fmt.Println("   ✓ Git repository initialized")
		}
	}

	scaffoldpkg.PrintNextSteps(name)

	return nil
}

