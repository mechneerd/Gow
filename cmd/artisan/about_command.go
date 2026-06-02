package artisan

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// AboutCmd displays information about the GoW installation.
var AboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Display basic information about the GoW installation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(AboutCommand())
	},
}

// AboutCommand displays information about the GoW installation.
func AboutCommand() string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("        ___           ___         \n")
	sb.WriteString("       /   \\         /   \\        \n")
	sb.WriteString("      /     \\       /     \\       \n")
	sb.WriteString("     /  /\\   \\     /  /\\   \\      \n")
	sb.WriteString("    /  /  \\   \\   /  /  \\   \\     \n")
	sb.WriteString("   /  /    \\   \\ /  /    \\   \\    \n")
	sb.WriteString("  /  /      \\   V  /      \\   \\   \n")
	sb.WriteString(" /__/        \\____/        \\__\\  \n")
	sb.WriteString("\n")
	sb.WriteString("GoW Framework — Laravel for Go\n")
	sb.WriteString("──────────────────────────────\n\n")

	sb.WriteString(fmt.Sprintf("  Go Version:    %s\n", runtime.Version()))
	sb.WriteString(fmt.Sprintf("  GOOS/GOARCH:   %s/%s\n", runtime.GOOS, runtime.GOARCH))
	sb.WriteString(fmt.Sprintf("  CPU Cores:     %d\n", runtime.NumCPU()))
	sb.WriteString(fmt.Sprintf("  Goroutines:    %d\n", runtime.NumGoroutine()))
	sb.WriteString(fmt.Sprintf("  Module:        github.com/mechneerd/gow\n"))
	sb.WriteString(fmt.Sprintf("  Current Time:  %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("\n")

	sb.WriteString("  Packages:\n")
	sb.WriteString("    ├── auth       — Authentication & Authorization\n")
	sb.WriteString("    ├── billing    — Payment Processing (Cashier)\n")
	sb.WriteString("    ├── broadcasting — WebSocket & Real-time\n")
	sb.WriteString("    ├── cache      — Caching Drivers & Tags\n")
	sb.WriteString("    ├── config     — Configuration Management\n")
	sb.WriteString("    ├── container  — IoC Container\n")
	sb.WriteString("    ├── contracts  — Interface Definitions\n")
	sb.WriteString("    ├── database   — ORM & Query Builder\n")
	sb.WriteString("    ├── encryption — Encryption & Signing\n")
	sb.WriteString("    ├── events     — Event Dispatcher\n")
	sb.WriteString("    ├── hashing    — Password Hashing\n")
	sb.WriteString("    ├── horizon    — Queue Dashboard\n")
	sb.WriteString("    ├── http       — HTTP Layer\n")
	sb.WriteString("    ├── localization — Translations\n")
	sb.WriteString("    ├── logging    — Logging\n")
	sb.WriteString("    ├── mail       — Email Sending\n")
	sb.WriteString("    ├── notifications — Notifications\n")
	sb.WriteString("    ├── nova       — Admin Panel\n")
	sb.WriteString("    ├── pulse      — Health Monitoring\n")
	sb.WriteString("    ├── queue      — Job Queue\n")
	sb.WriteString("    ├── routing    — HTTP Routing\n")
	sb.WriteString("    ├── session    — Session Management\n")
	sb.WriteString("    ├── storage    — File Storage\n")
	sb.WriteString("    ├── support    — Helpers & Utilities\n")
	sb.WriteString("    ├── testing    — Test Utilities\n")
	sb.WriteString("    ├── validation — Form Validation\n")
	sb.WriteString("    └── view       — Template Engine\n\n")

	sb.WriteString("  Features:\n")
	sb.WriteString("    ✓ 250+ Laravel features implemented\n")
	sb.WriteString("    ✓ Generic ORM with relationships\n")
	sb.WriteString("    ✓ Blade template engine\n")
	sb.WriteString("    ✓ Queue workers with batching\n")
	sb.WriteString("    ✓ Sanctum API authentication\n")
	sb.WriteString("    ✓ Socialite OAuth providers\n")
	sb.WriteString("    ✓ Horizon queue dashboard\n")
	sb.WriteString("    ✓ Nova admin panel\n")
	sb.WriteString("    ✓ Pulse health monitoring\n")
	sb.WriteString("    ✓ Sail Docker environment\n")
	sb.WriteString("    ✓ Telescope debug dashboard\n")
	sb.WriteString("    ✓ 50+ Artisan commands\n\n")

	return sb.String()
}

// AboutDetailed returns detailed package information.
func AboutDetailed() map[string]any {
	return map[string]any{
		"framework": "GoW",
		"version":   "1.0.0",
		"go_version": runtime.Version(),
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"cpus":      runtime.NumCPU(),
		"goroutines": runtime.NumGoroutine(),
		"modules": []string{
			"auth", "billing", "broadcasting", "cache", "config",
			"container", "contracts", "database", "encryption",
			"events", "hashing", "horizon", "http", "localization",
			"logging", "mail", "notifications", "nova", "pulse",
			"queue", "routing", "session", "storage", "support",
			"testing", "validation", "view",
		},
		"features": []string{
			"ORM with relationships", "Blade templates", "Queue workers",
			"Sanctum auth", "Socialite OAuth", "Horizon dashboard",
			"Nova admin panel", "Pulse health", "Sail Docker",
			"Telescope debug", "50+ Artisan commands",
		},
	}
}
