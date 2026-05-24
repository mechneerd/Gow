package rbac

import (
	"database/sql"
	"fmt"
)

// HasRoles provides RBAC methods for User models.
// Embed this in your User model for role/permission support.
//
// Note: The embedding model should call SetID(user.ID) after loading
// so that HasRole / HasPermission / AssignRole work correctly.
type HasRoles struct {
	ID int
	db *sql.DB // Can be set per instance or via global SetDefaultDB
}

// HasRole checks if the user has a specific role by querying role_user table.
func (h *HasRoles) HasRole(roleName string) bool {
	db := h.db
	if db == nil {
		db = getDefaultDB()
	}
	if h.ID == 0 || db == nil {
		return false
	}

	var count int
	query := `
		SELECT COUNT(*) 
		FROM role_user ru
		JOIN roles r ON r.id = ru.role_id
		WHERE ru.user_id = ? AND r.name = ?
	`
	err := db.QueryRow(query, h.ID, roleName).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// HasPermission checks if the user has a permission through any of their roles.
func (h *HasRoles) HasPermission(permissionName string) bool {
	db := h.db
	if db == nil {
		db = getDefaultDB()
	}
	if h.ID == 0 || db == nil {
		return false
	}

	var count int
	query := `
		SELECT COUNT(*) 
		FROM role_user ru
		JOIN roles r ON r.id = ru.role_id
		JOIN permission_role pr ON pr.role_id = r.id
		JOIN permissions p ON p.id = pr.permission_id
		WHERE ru.user_id = ? AND p.name = ?
	`
	err := db.QueryRow(query, h.ID, permissionName).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// AssignRole assigns a role to the user (creates entry in role_user).
func (h *HasRoles) AssignRole(roleName string) error {
	if h.ID == 0 || h.db == nil {
		return fmt.Errorf("user ID or database not available")
	}

	// Find role ID
	var roleID int
	err := h.db.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("role not found: %s", roleName)
	}

	_, err = h.db.Exec(`
		INSERT OR IGNORE INTO role_user (role_id, user_id) 
		VALUES (?, ?)
	`, roleID, h.ID)

	return err
}

// Can is an alias for HasPermission (Laravel-style).
func (h *HasRoles) Can(permission string) bool {
	return h.HasPermission(permission)
}

// SetID sets the user ID for RBAC operations.
// Call this after loading the model (e.g. user.HasRoles.SetID(user.ID)).
func (h *HasRoles) SetID(id int) {
	h.ID = id
}
