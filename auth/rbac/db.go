package rbac

import "database/sql"

var defaultDB *sql.DB

// SetDefaultDB allows setting a global database connection that HasRoles can use.
// This is a pragmatic helper for generated projects until a better DI solution is in place.
func SetDefaultDB(db *sql.DB) {
	defaultDB = db
}

// getDefaultDB returns the globally set DB (if any).
func getDefaultDB() *sql.DB {
	return defaultDB
}

