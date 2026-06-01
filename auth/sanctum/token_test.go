package sanctum

import (
	"database/sql"
	"github.com/mechneerd/gow/database/dialect"
	"github.com/mechneerd/gow/database/orm"
	"github.com/mechneerd/gow/database/query"
	"net/http"
	"net/http/httptest"
	"testing"
	
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *orm.DB {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open memory db: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE personal_access_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tokenable_type TEXT,
			tokenable_id TEXT,
			name TEXT,
			token TEXT,
			abilities TEXT,
			last_used_at DATETIME,
			created_at INTEGER,
			updated_at INTEGER,
			expires_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return &orm.DB{
		Conn:    conn,
		Builder: query.NewBuilder(conn, &dialect.SQLiteDialect{}),
	}
}

func TestTokenManager(t *testing.T) {
	db := setupTestDB(t)
	manager := NewTokenManager(db)

	plainToken, err := manager.CreateToken("User", "1", "test-token", []string{"read", "write"})
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}
	if len(plainToken) != 40 {
		t.Errorf("Expected token length 40, got %d", len(plainToken))
	}

	// Test Middleware Authentication
	mw := manager.Middleware()
	
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(TokenKey).(*PersonalAccessToken)
		if token == nil {
			t.Errorf("Expected token in context, got nil")
		}
		w.WriteHeader(http.StatusOK)
	}))

	// Missing Header
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 without header, got %d", w.Result().StatusCode)
	}

	// Invalid Token
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 with invalid token, got %d", w.Result().StatusCode)
	}

	// Valid Token
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+plainToken)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected 200 with valid token, got %d", w.Result().StatusCode)
	}
}

