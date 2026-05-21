package auth

import "gow/session"

// SessionGuard implements the Guard interface using standard HTTP sessions.
type SessionGuard struct {
	name     string
	provider UserProvider
	session  *session.Manager
	user     any
}

// NewSessionGuard creates a new session guard instance.
func NewSessionGuard(name string, provider UserProvider, session *session.Manager) *SessionGuard {
	return &SessionGuard{
		name:     name,
		provider: provider,
		session:  session,
	}
}

func (g *SessionGuard) Check() bool {
	return g.User() != nil
}

func (g *SessionGuard) Guest() bool {
	return !g.Check()
}

func (g *SessionGuard) User() any {
	if g.user != nil {
		return g.user
	}

	id := g.session.Get(g.getName())
	if id != nil {
		g.user = g.provider.RetrieveByID(id.(string))
	}

	return g.user
}

func (g *SessionGuard) ID() string {
	if g.Check() {
		return g.session.Get(g.getName()).(string)
	}
	return ""
}

func (g *SessionGuard) Attempt(credentials map[string]any) bool {
	user := g.provider.RetrieveByCredentials(credentials)
	if user == nil {
		return false
	}

	if g.provider.ValidateCredentials(user, credentials) {
		g.Login(user)
		return true
	}

	return false
}

func (g *SessionGuard) Login(user any) {
	// Typically we need a way to extract the ID from the generic 'user' type.
	// For simplicity, we assume the provider gives us a struct we can reflect on,
	// or the User implements an Authenticatable interface.
	// As a placeholder, we use a generic type assertion here if possible.
	if authUser, ok := user.(Authenticatable); ok {
		g.session.Put(g.getName(), authUser.GetAuthIdentifier())
		g.session.Regenerate() // Prevent session fixation
		g.user = user
	}
}

func (g *SessionGuard) Logout() {
	g.user = nil
	g.session.Put(g.getName(), nil) // clear from session
	g.session.Regenerate()
}

func (g *SessionGuard) getName() string {
	return "login_" + g.name + "_id"
}

// Authenticatable is an interface that users must implement to be logged in via SessionGuard.
type Authenticatable interface {
	GetAuthIdentifier() string
	GetAuthPassword() string
}

// Manager resolves the guards configured for the application.
type Manager struct {
	guards map[string]Guard
}

// NewManager creates a new Auth Manager.
func NewManager() *Manager {
	return &Manager{
		guards: make(map[string]Guard),
	}
}

// AddGuard registers a new guard.
func (m *Manager) AddGuard(name string, guard Guard) {
	m.guards[name] = guard
}

// Guard retrieves a guard instance by name.
func (m *Manager) Guard(name string) Guard {
	return m.guards[name]
}
