package request

import (
	"bytes"
	"context"
	"github.com/mechneerd/gow/routing"
	"io"
	"net/http/httptest"
	"testing"
)

func TestRequestHelpers(t *testing.T) {
	t.Run("Param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), routing.ParamsKey, map[string]string{"id": "42"})
		req = req.WithContext(ctx)

		if Param(req, "id") != "42" {
			t.Errorf("expected Param to return 42")
		}
	})

	t.Run("Query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?search=gow", nil)
		if Query(req, "search") != "gow" {
			t.Errorf("expected Query to return gow")
		}
	})

	t.Run("Input Fallback (Form vs Query)", func(t *testing.T) {
		body := bytes.NewBufferString("name=Alice")
		req := httptest.NewRequest("POST", "/?name=Bob", body)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Form should win over Query
		if Input(req, "name") != "Alice" {
			t.Errorf("expected Input to return Alice, got %s", Input(req, "name"))
		}
	})

	t.Run("Input JSON Caching and Priority", func(t *testing.T) {
		bodyStr := `{"role": "admin", "age": 30}`
		body := bytes.NewBufferString(bodyStr)
		req := httptest.NewRequest("POST", "/?role=user", body)
		req.Header.Set("Content-Type", "application/json")

		// JSON should win over Query
		if Input(req, "role") != "admin" {
			t.Errorf("expected JSON role 'admin', got %s", Input(req, "role"))
		}

		// Ensure body was restored
		restoredBody, _ := io.ReadAll(req.Body)
		if string(restoredBody) != bodyStr {
			t.Errorf("expected body to be restored, got %s", string(restoredBody))
		}

		// Ensure age is returned as string representation for generic Input
		// (json.Unmarshal parses numbers as float64, so it stringifies to "30")
		if Input(req, "age") != "30" {
			t.Errorf("expected JSON age '30', got %s", Input(req, "age"))
		}
	})

	t.Run("Only", func(t *testing.T) {
		bodyStr := `{"a": 1, "b": 2, "c": 3}`
		req := httptest.NewRequest("POST", "/", bytes.NewBufferString(bodyStr))
		req.Header.Set("Content-Type", "application/json")

		filtered := Only(req, "a", "c", "d")
		if len(filtered) != 2 {
			t.Errorf("expected 2 items, got %d", len(filtered))
		}
		if filtered["a"] != float64(1) || filtered["c"] != float64(3) {
			t.Errorf("unexpected filtered values: %v", filtered)
		}
	})

	t.Run("Has", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?exists=1", nil)
		if !Has(req, "exists") {
			t.Error("expected Has to be true")
		}
		if Has(req, "missing") {
			t.Error("expected Has to be false")
		}
	})
}

