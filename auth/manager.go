package auth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/mechneerd/gow/session"
)

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
		val := g.session.Get(g.getName())
		if val == nil {
			return ""
		}
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

func (g *SessionGuard) Attempt(credentials map[string]any, remember ...bool) bool {
	user := g.provider.RetrieveByCredentials(credentials)
	if user == nil {
		return false
	}

	if g.provider.ValidateCredentials(user, credentials) {
		rem := false
		if len(remember) > 0 {
			rem = remember[0]
		}
		g.Login(user, rem)
		return true
	}

	return false
}

func (g *SessionGuard) Login(user any, remember ...bool) {
	if authUser, ok := user.(Authenticatable); ok {
		g.session.Put(g.getName(), authUser.GetAuthIdentifier())
		g.session.Regenerate()
		g.user = user

		// Remember Me support (Wave 4)
		if len(remember) > 0 && remember[0] {
			if provider, ok := g.provider.(interface {
				UpdateRememberToken(user any, token string)
			}); ok {
				token := generateRememberToken()
				provider.UpdateRememberToken(user, token)
				// Remember-me cookie support is available via the provider; full long-lived
				// cookie handling should be implemented in the application layer.
			}
		}
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

// generateRememberToken creates a secure random remember token.
func generateRememberToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate remember token: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)
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

