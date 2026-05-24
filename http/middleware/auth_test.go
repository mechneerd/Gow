package middleware

import (
	"github.com/mechneerd/gow/auth"
	"github.com/mechneerd/gow/auth/access"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockGuard for testing
type MockGuard struct {
	authenticated bool
	user          any
}

func (m *MockGuard) Check() bool { return m.authenticated }
func (m *MockGuard) Guest() bool { return !m.authenticated }
func (m *MockGuard) User() any   { return m.user }
func (m *MockGuard) ID() string  { return "1" }
func (m *MockGuard) Attempt(map[string]any) bool { return false }
func (m *MockGuard) Login(any) {}
func (m *MockGuard) Logout() {}

type TestUser struct {
	ID   int
	Role string
}

func TestAuthenticateMiddleware(t *testing.T) {
	manager := auth.NewManager()

	// 1. Test Unauthenticated
	manager.AddGuard("web", &MockGuard{authenticated: false})
	mw := Authenticate(manager, "web")

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be reached for unauthenticated user")
	}))

	req := httptest.NewRequest("GET", "/dashboard", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected redirect 302, got %d", w.Code)
	}

	// 2. Test Authenticated
	manager.AddGuard("web", &MockGuard{authenticated: true, user: TestUser{ID: 1}})
	mwAuth := Authenticate(manager, "web")

	handlerAuth := mwAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := auth.User(r)
		if user == nil {
			t.Error("User should be present in context")
		} else if u, ok := user.(TestUser); !ok || u.ID != 1 {
			t.Errorf("Expected TestUser with ID 1, got %v", user)
		}
		w.WriteHeader(http.StatusOK)
	}))

	reqAuth := httptest.NewRequest("GET", "/dashboard", nil)
	wAuth := httptest.NewRecorder()
	handlerAuth.ServeHTTP(wAuth, reqAuth)

	if wAuth.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", wAuth.Code)
	}
}

func TestAuthorizeMiddleware(t *testing.T) {
	gate := access.NewGate()
	gate.Define("edit-post", func(user any, args ...any) bool {
		u := user.(TestUser)
		return u.Role == "admin"
	})

	mw := Authorize(gate, "edit-post")

	// 1. Test Without User in Context
	reqNoUser := httptest.NewRequest("GET", "/post/1/edit", nil)
	wNoUser := httptest.NewRecorder()
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(wNoUser, reqNoUser)

	if wNoUser.Code != http.StatusFound { // defaults to redirect to login
		t.Errorf("Expected redirect for missing user, got %d", wNoUser.Code)
	}

	// 2. Test Unauthorized User
	manager := auth.NewManager()
	manager.AddGuard("web", &MockGuard{authenticated: true, user: TestUser{ID: 1, Role: "user"}})
	
	handlerUnauthorized := Authenticate(manager, "web")(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be reached for unauthorized user")
	})))

	reqUnauth := httptest.NewRequest("GET", "/post/1/edit", nil)
	wUnauth := httptest.NewRecorder()
	
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic (HTTP Abort) for unauthorized user, got nothing")
			}
		}()
		handlerUnauthorized.ServeHTTP(wUnauth, reqUnauth)
	}()

	// 3. Test Authorized User
	manager.AddGuard("web", &MockGuard{authenticated: true, user: TestUser{ID: 1, Role: "admin"}})
	
	handlerAuthorized := Authenticate(manager, "web")(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	reqAuth := httptest.NewRequest("GET", "/post/1/edit", nil)
	wAuth := httptest.NewRecorder()
	handlerAuthorized.ServeHTTP(wAuth, reqAuth)

	if wAuth.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", wAuth.Code)
	}
}

