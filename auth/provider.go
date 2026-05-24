package auth

// UserProvider abstracts how users are retrieved from the storage mechanism.
type UserProvider interface {
	// RetrieveByID retrieves a user by their unique identifier.
	RetrieveByID(identifier string) any
	// RetrieveByToken retrieves a user by their unique identifier and "remember me" token.
	RetrieveByToken(identifier string, token string) any
	// UpdateRememberToken updates the "remember me" token for the given user.
	UpdateRememberToken(user any, token string)
	// RetrieveByCredentials retrieves a user by the given credentials.
	RetrieveByCredentials(credentials map[string]any) any
	// ValidateCredentials validates a user against the given credentials.
	ValidateCredentials(user any, credentials map[string]any) bool
}

