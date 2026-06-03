package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/mechneerd/gow/session"
)

// MockUser implements Authenticatable for testing.
type MockUser struct {
	ID       string
	Password string
}

func (u *MockUser) GetAuthIdentifier() string { return u.ID }
func (u *MockUser) GetAuthPassword() string   { return u.Password }

// MockProvider implements UserProvider for testing.
type MockProvider struct {
	users    map[string]*MockUser
	validate func(user any, credentials map[string]any) bool
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		users: make(map[string]*MockUser),
		validate: func(user any, credentials map[string]any) bool {
			if u, ok := user.(*MockUser); ok {
				if pw, ok := credentials["password"].(string); ok {
					return u.Password == pw
				}
			}
			return false
		},
	}
}

func (p *MockProvider) AddUser(user *MockUser) {
	p.users[user.ID] = user
}

func (p *MockProvider) RetrieveByID(identifier string) any {
	return p.users[identifier]
}

func (p *MockProvider) RetrieveByToken(identifier string, token string) any {
	return p.users[identifier]
}

func (p *MockProvider) UpdateRememberToken(user any, token string) {}

func (p *MockProvider) RetrieveByCredentials(credentials map[string]any) any {
	if email, ok := credentials["email"].(string); ok {
		return p.users[email]
	}
	return nil
}

func (p *MockProvider) ValidateCredentials(user any, credentials map[string]any) bool {
	return p.validate(user, credentials)
}

func newTestSessionManager(t *testing.T) *session.Manager {
	t.Helper()
	store := session.NewArrayDriver()
	mgr := session.NewManager(store, "test-session-id")
	if err := mgr.Start(); err != nil {
		t.Fatalf("failed to start session: %v", err)
	}
	return mgr
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if len(m.guards) != 0 {
		t.Errorf("expected 0 guards, got %d", len(m.guards))
	}
}

func TestAddGuard(t *testing.T) {
	m := NewManager()
	guard := &SessionGuard{name: "web"}
	m.AddGuard("web", guard)
	if m.Guard("web") == nil {
		t.Error("expected guard 'web' to exist")
	}
}

func TestGuardReturnsNilForUnknown(t *testing.T) {
	m := NewManager()
	if g := m.Guard("unknown"); g != nil {
		t.Errorf("expected nil, got %v", g)
	}
}

func TestSessionGuard_Check(t *testing.T) {
	provider := NewMockProvider()
	provider.AddUser(&MockUser{ID: "1", Password: "secret"})

	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	if guard.Check() {
		t.Error("expected Check() to be false for guest")
	}

	guard.Login(&MockUser{ID: "1", Password: "secret"})
	if !guard.Check() {
		t.Error("expected Check() to be true after login")
	}
}

func TestSessionGuard_Guest(t *testing.T) {
	provider := NewMockProvider()
	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	if !guard.Guest() {
		t.Error("expected Guest() to be true for unauthenticated user")
	}

	guard.Login(&MockUser{ID: "1", Password: "secret"})
	if guard.Guest() {
		t.Error("expected Guest() to be false after login")
	}
}

func TestSessionGuard_User(t *testing.T) {
	provider := NewMockProvider()
	user := &MockUser{ID: "1", Password: "secret"}
	provider.AddUser(user)

	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	if guard.User() != nil {
		t.Error("expected User() to be nil for guest")
	}

	guard.Login(user)
	if guard.User() == nil {
		t.Error("expected User() to return user after login")
	}
}

func TestSessionGuard_ID(t *testing.T) {
	provider := NewMockProvider()
	user := &MockUser{ID: "1", Password: "secret"}
	provider.AddUser(user)

	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	if id := guard.ID(); id != "" {
		t.Errorf("expected empty ID, got %q", id)
	}

	guard.Login(user)
	if id := guard.ID(); id != "1" {
		t.Errorf("expected ID '1', got %q", id)
	}
}

func TestSessionGuard_Attempt(t *testing.T) {
	provider := NewMockProvider()
	provider.AddUser(&MockUser{ID: "1", Password: "secret"})

	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	// Wrong credentials
	if guard.Attempt(map[string]any{"email": "1", "password": "wrong"}) {
		t.Error("expected Attempt to fail with wrong password")
	}

	// Correct credentials
	if !guard.Attempt(map[string]any{"email": "1", "password": "secret"}) {
		t.Error("expected Attempt to succeed with correct password")
	}

	if !guard.Check() {
		t.Error("expected user to be logged in after successful Attempt")
	}
}

func TestSessionGuard_Logout(t *testing.T) {
	provider := NewMockProvider()
	provider.AddUser(&MockUser{ID: "1", Password: "secret"})

	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	guard.Login(&MockUser{ID: "1", Password: "secret"})
	if !guard.Check() {
		t.Fatal("expected user to be logged in")
	}

	guard.Logout()
	if guard.Check() {
		t.Error("expected user to be logged out")
	}
}

func TestSessionGuard_LoginByID(t *testing.T) {
	provider := NewMockProvider()
	provider.AddUser(&MockUser{ID: "42", Password: "test"})

	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	guard.LoginByID("42")
	if id := guard.ID(); id != "42" {
		t.Errorf("expected ID '42', got %q", id)
	}
}

func TestSessionGuard_SetUser(t *testing.T) {
	provider := NewMockProvider()
	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	user := &MockUser{ID: "99", Password: "x"}
	guard.SetUser(user)
	if guard.User() != user {
		t.Error("expected SetUser to set the user")
	}
}

func TestSessionGuard_RegenerateSession(t *testing.T) {
	provider := NewMockProvider()
	sess := newTestSessionManager(t)
	guard := NewSessionGuard("web", provider, sess)

	// Should not panic
	guard.RegenerateSession()
}

func TestUserFromContext(t *testing.T) {
	user := &MockUser{ID: "1", Password: "secret"}
	req := (&http.Request{}).WithContext(context.WithValue(context.Background(), UserContextKey, user))

	got := User(req)
	if got != user {
		t.Error("expected User() to return user from context")
	}
}

func TestUserFromNilContext(t *testing.T) {
	req := &http.Request{}
	req = req.WithContext(context.Background())

	if got := User(req); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestUserIDFromContext(t *testing.T) {
	user := &MockUser{ID: "1", Password: "secret"}
	req := (&http.Request{}).WithContext(context.WithValue(context.Background(), UserContextKey, user))

	if id := UserID(req); id != "1" {
		t.Errorf("expected '1', got %q", id)
	}
}

func TestUserIDFromNonAuthUser(t *testing.T) {
	req := (&http.Request{}).WithContext(context.WithValue(context.Background(), UserContextKey, "not-a-user"))

	if id := UserID(req); id != "" {
		t.Errorf("expected empty, got %q", id)
	}
}

func TestGenerateRememberToken(t *testing.T) {
	token := generateRememberToken()
	if token == "" {
		t.Error("expected non-empty token")
	}
	if len(token) < 32 {
		t.Error("expected token to be at least 32 chars")
	}

	// Two tokens should be different
	token2 := generateRememberToken()
	if token == token2 {
		t.Error("expected two tokens to be different")
	}
}

func TestAuthenticatableInterface(t *testing.T) {
	user := &MockUser{ID: "1", Password: "secret"}
	var _ Authenticatable = user
}
