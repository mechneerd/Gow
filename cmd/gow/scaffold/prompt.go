package scaffold

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PromptResult holds the user's answers from the interactive wizard.
type PromptResult struct {
	StarterKit string // "minimal", "minimal-api", "api", "web", "auth", "full", etc.
	Database   string // "sqlite", "mysql", "postgres"
	ModulePath string
}

// RunInteractiveWizard asks the user questions when no flags are provided.
func RunInteractiveWizard(defaultAppName string) (PromptResult, error) {
	reader := bufio.NewReader(os.Stdin)
	result := PromptResult{}

	fmt.Println("\nNo starter kit specified. Let's create your project interactively.")

	kits := GetStarterKits()
	maxReady := 0
	fmt.Println("\n🎯 Which starter kit would you like?")
	for i, kit := range kits {
		num := i + 1
		if kit.Ready {
			maxReady = num
			marker := ""
			if kit.Key == "auth" {
				marker = " (recommended)"
			}
			fmt.Printf("  %d) %-22s — %s%s\n", num, kit.Name, kit.Description, marker)
		} else {
			fmt.Printf("  %d) %-22s — %s (planned)\n", num, kit.Name, kit.Description)
		}
	}

	defaultChoice := "4" // Web + Auth
	fmt.Printf("Enter choice [1-%d] (default: %s): ", len(kits), defaultChoice)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "" {
		choice = defaultChoice
	}

choiceLoop:
	for {
		idx, err := strconv.Atoi(choice)
		if err == nil && idx >= 1 && idx <= len(kits) {
			selected := kits[idx-1]
			if selected.Ready {
				result.StarterKit = selected.Key
				break choiceLoop
			}
			fmt.Printf("   ⚠️  %s is not yet available. Please choose a ready starter kit [1-%d]: ", selected.Name, maxReady)
			choice, _ = reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			continue choiceLoop
		}
		// Invalid input — default to auth
		result.StarterKit = "auth"
		break choiceLoop
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

