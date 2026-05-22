package resources

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type UserResource struct{}

func (UserResource) ToMap(req *http.Request, user User) map[string]any {
	return map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		// Password intentionally omitted
	}
}

func TestResourceTransform(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com", Password: "secret"}
	transformer := UserResource{}

	result := Transform(req, user, transformer)

	if result["id"] != 1 {
		t.Errorf("Expected id 1, got %v", result["id"])
	}
	if result["name"] != "Alice" {
		t.Errorf("Expected Alice, got %v", result["name"])
	}
	if _, exists := result["password"]; exists {
		t.Errorf("Expected password to be omitted")
	}
}

func TestResourceCollection(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	users := []User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}
	transformer := UserResource{}

	results := Collection(req, users, transformer)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	if results[0]["name"] != "Alice" {
		t.Errorf("Expected Alice, got %v", results[0]["name"])
	}
	if results[1]["name"] != "Bob" {
		t.Errorf("Expected Bob, got %v", results[1]["name"])
	}
}

func TestResourceRespond(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	transformer := UserResource{}

	Respond(w, req, user, transformer, http.StatusOK)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var payload ResourceResponse
	err := json.Unmarshal(w.Body.Bytes(), &payload)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	data, ok := payload.Data.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be a map, got %T", payload.Data)
	}

	if data["id"] != float64(1) { // JSON unmarshals numbers as float64
		t.Errorf("Expected id 1, got %v", data["id"])
	}
}

func TestRespondWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	
	meta := map[string]any{
		"current_page": 1,
		"last_page":    10,
	}
	links := map[string]any{
		"next": "http://example.com/api/users?page=2",
	}

	RespondWithMeta(w, []map[string]any{{"id": 1}}, meta, links, http.StatusOK)

	var payload ResourceResponse
	err := json.Unmarshal(w.Body.Bytes(), &payload)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if payload.Meta["current_page"] != float64(1) {
		t.Errorf("Expected meta current_page 1, got %v", payload.Meta["current_page"])
	}
	if payload.Links["next"] != "http://example.com/api/users?page=2" {
		t.Errorf("Expected next link, got %v", payload.Links["next"])
	}
}
