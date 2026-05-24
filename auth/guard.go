package auth

// Guard defines the interface for authenticating users per request.
type Guard interface {
	// Check determines if the current user is authenticated.
	Check() bool
	// Guest determines if the current user is a guest.
	Guest() bool
	// User returns the currently authenticated user.
	User() any
	// ID returns the ID of the currently authenticated user.
	ID() string
	// Attempt tries to authenticate a user using the given credentials.
	Attempt(credentials map[string]any) bool
	// Login logs a user into the application.
	Login(user any)
	// Logout logs the user out of the application.
	Logout()
}

