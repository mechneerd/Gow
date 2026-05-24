package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InjectRBACBootstrapExamples appends ready-to-use commented examples
// for RBAC middleware and global DB wiring into the generated project's
// bootstrap/app.go (for auth-enabled kits).
//
// This ensures that even if the upstream skeleton's bootstrap/app.go
// is minimal, users get clear, copy-pasteable guidance immediately after
// `gow new --auth`.
func InjectRBACBootstrapExamples(projectDir string) error {
	bootstrapPath := filepath.Join(projectDir, "bootstrap", "app.go")

	content, err := os.ReadFile(bootstrapPath)
	if err != nil {
		// If no bootstrap/app.go exists yet, we skip silently (some minimal kits may not have it)
		return nil
	}

	original := string(content)

	// Avoid double-injection
	if strings.Contains(original, "RBAC / Auth Middleware - Ready-to-use examples") {
		return nil
	}

	injection := `

/*
 * ============================================================
 * RBAC + Auth Middleware — Ready-to-use setup (auth kits)
 * ============================================================
 * 1. After gow migrate && gow db:seed, wire the global DB once:
 *
 *      import "yourmodule/auth/rbac"
 *      rbac.SetDefaultDB(db)   // call this after opening your *sql.DB
 *
 * 2. Register / use the middleware (example protected routes):
 *
 *      import (
 *          "yourmodule/app/Http/Middleware"
 *          "yourmodule/auth/rbac"
 *      )
 *
 *      // In routes or bootstrap:
 *      router.Group(func(r *router.Router) {
 *          r.Use(Middleware.RoleMiddleware("super-admin"))
 *          // r.Use(Middleware.PermissionMiddleware("users.manage"))
 *
 *          r.Get("/admin/dashboard", adminDashboardHandler)
 *          r.Post("/admin/users", createUserHandler)
 *      })
 *
 * 3. Common roles created by the default RoleSeeder:
 *      - super-admin (full access)
 *      - admin
 *      - editor
 *      - user
 *
 * See docs/rbac.md or the generated README for full details.
 */
`

	// Append at the end of the file (before any final package-level closing if needed)
	newContent := original
	if !strings.HasSuffix(original, "\n") {
		newContent += "\n"
	}
	newContent += injection

	if err := os.WriteFile(bootstrapPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to inject RBAC examples into bootstrap/app.go: %w", err)
	}

	fmt.Println("   ✓ Added RBAC middleware examples to bootstrap/app.go")
	return nil
}

