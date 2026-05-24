package middleware

import (
	"net/http"

	"gow/auth"
)

// RoleMiddleware checks if the authenticated user has the required role.
func RoleMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := auth.User() // Assumes auth.User() is available in context or session

			if user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Type assert to our RBAC-enabled user if possible
			if rbacUser, ok := user.(interface{ HasRole(string) bool }); ok {
				if !rbacUser.HasRole(role) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			} else {
				// Fallback if not using HasRoles
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// PermissionMiddleware checks if the user has a specific permission.
func PermissionMiddleware(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := auth.User()
			if user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if rbacUser, ok := user.(interface{ Can(string) bool }); ok {
				if !rbacUser.Can(permission) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
