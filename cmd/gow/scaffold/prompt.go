package scaffold

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptResult holds the user's answers from the interactive wizard.
type PromptResult struct {
	StarterKit string // "minimal", "api", "web", "auth"
	Database   string // "sqlite", "mysql", "postgres"
	ModulePath string
}

// RunInteractiveWizard asks the user questions when no flags are provided.
func RunInteractiveWizard(defaultAppName string) (PromptResult, error) {
	reader := bufio.NewReader(os.Stdin)
	result := PromptResult{}

	fmt.Println("\nNo starter kit specified. Let's create your project interactively.")

	fmt.Println("\n🎯 Which starter kit would you like?")
	fmt.Println("  1) Minimal")
	fmt.Println("  2) API (with Sanctum)")
	fmt.Println("  3) Web (Blade + views)")
	fmt.Println("  4) Web + Auth (recommended)")

	fmt.Print("Enter choice [1-4] (default: 4): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		result.StarterKit = "minimal"
	case "2":
		result.StarterKit = "api"
	case "3":
		result.StarterKit = "web"
	case "4", "":
		result.StarterKit = "auth"
	default:
		result.StarterKit = "auth"
	}

	for {
		fmt.Println("\n🗄️  Which database driver?")
		fmt.Println("  1) SQLite (default)")
		fmt.Println("  2) MySQL")
		fmt.Println("  3) PostgreSQL")

		fmt.Print("Enter choice [1-3] (default: 1): ")
		dbChoice, _ := reader.ReadString('\n')
		dbChoice = strings.TrimSpace(dbChoice)

		switch dbChoice {
		case "2":
			result.Database = "mysql"
		case "3":
			result.Database = "postgres"
		default:
			result.Database = "sqlite"
		}

		if IsValidDatabase(result.Database) {
			break
		}
		fmt.Println("Invalid choice. Please enter 1, 2, or 3.")
	}

	fmt.Printf("\n📦 Module path for go.mod? (leave empty to use '%s')\n", defaultAppName)
	fmt.Print("Module path: ")
	module, _ := reader.ReadString('\n')
	result.ModulePath = strings.TrimSpace(module)

	return result, nil
}

