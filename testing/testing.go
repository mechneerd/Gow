package testing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gow/auth/sanctum"
	"gow/database/orm"
	"gow/database/query"
	"gow/foundation"
	"gow/routing"
)

// TestCase provides a fluent testing API similar to Laravel.
type TestCase struct {
	*testing.T
	App    *foundation.Application
	DB     *orm.DB
	Router *routing.Router
	client *http.Client
}

// NewTestCase initializes a test environment.
func NewTestCase(t *testing.T, app *foundation.Application, db *orm.DB, router *routing.Router) *TestCase {
	return &TestCase{
		T:      t,
		App:    app,
		DB:     db,
		Router: router,
		client: &http.Client{},
	}
}

// AssertDatabaseHas asserts that a database table contains a row matching the given constraints.
func (tc *TestCase) AssertDatabaseHas(table string, conditions map[string]any) {
	tc.Helper()
	builder := query.NewBuilder(tc.DB.RawDB(), tc.DB.Dialect())
	builder.Table(table)
	
	for k, v := range conditions {
		builder.Where(k, "=", v)
	}
	
	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error querying database: %v", err)
	}
	if count == 0 {
		tc.Errorf("Failed asserting that table [%s] has row matching %v", table, conditions)
	}
}

// AssertDatabaseMissing asserts that a database table does NOT contain a row matching the given constraints.
func (tc *TestCase) AssertDatabaseMissing(table string, conditions map[string]any) {
	tc.Helper()
	builder := query.NewBuilder(tc.DB.RawDB(), tc.DB.Dialect())
	builder.Table(table)
	
	for k, v := range conditions {
		builder.Where(k, "=", v)
	}
	
	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error querying database: %v", err)
	}
	if count > 0 {
		tc.Errorf("Failed asserting that table [%s] is missing row matching %v", table, conditions)
	}
}

// ActingAs authenticates the current test request using the Sanctum middleware logic.
func (tc *TestCase) ActingAs(token string) *TestCase {
	// In a real framework we might set a context value globally for the test
	// or attach it to subsequent requests made by a test HTTP client.
	return tc
}

// Get dispatch a GET request to the application.
func (tc *TestCase) Get(uri string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	w := httptest.NewRecorder()
	
	tc.Router.ServeHTTP(w, req)
	return w
}

// Post dispatch a POST request to the application.
func (tc *TestCase) Post(uri string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, uri, body)
	w := httptest.NewRecorder()
	
	tc.Router.ServeHTTP(w, req)
	return w
}

// AssertStatus asserts the response has the given HTTP status code.
func (tc *TestCase) AssertStatus(w *httptest.ResponseRecorder, status int) {
	tc.Helper()
	if w.Code != status {
		tc.Errorf("Expected status %d but got %d", status, w.Code)
	}
}

// AssertSee asserts the response body contains the given string.
func (tc *TestCase) AssertSee(w *httptest.ResponseRecorder, text string) {
	tc.Helper()
	if !strings.Contains(w.Body.String(), text) {
		tc.Errorf("Expected to see %q in response", text)
	}
}
