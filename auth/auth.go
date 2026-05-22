package auth

import "net/http"

type ContextKey string

const UserContextKey ContextKey = "user"

// User retrieves the authenticated user from the HTTP request context.
// It returns nil if no user is authenticated.
func User(r *http.Request) any {
	return r.Context().Value(UserContextKey)
}

// UserID retrieves the ID of the authenticated user from the HTTP request context.
// It assumes the user implements the Authenticatable interface.
func UserID(r *http.Request) string {
	user := User(r)
	if authUser, ok := user.(Authenticatable); ok {
		return authUser.GetAuthIdentifier()
	}
	return ""
}
