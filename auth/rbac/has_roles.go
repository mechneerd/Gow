package rbac

import "gow/database/orm"

// HasRoles provides RBAC methods for User models
type HasRoles struct {
	orm.Model
}

// HasRole checks if the user has a specific role
func (h *HasRoles) HasRole(roleName string) bool {
	// Real implementation will query the role_user pivot
	// This is a placeholder that always returns false until full wiring
	return false
}

// HasPermission checks if the user has a specific permission (via roles)
func (h *HasRoles) HasPermission(permissionName string) bool {
	return false
}

// AssignRole attaches a role to the user
func (h *HasRoles) AssignRole(roleName string) {
	// Implementation pending full pivot logic
}

// Can is an alias for HasPermission
func (h *HasRoles) Can(permission string) bool {
	return h.HasPermission(permission)
}
