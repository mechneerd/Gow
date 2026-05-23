package routing

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRouterBasicVerbs(t *testing.T) {
	router := NewRouter()

	handler := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.Method + " OK"))
		return nil
	}

	router.Get("/get", handler)
	router.Post("/post", handler)
	router.Put("/put", handler)
	router.Delete("/delete", handler)
	router.Patch("/patch", handler)
	router.Options("/options", handler)

	tests := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}

	for _, method := range tests {
		req := httptest.NewRequest(method, "/"+strings.ToLower(method), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK for %s, got %d", method, res.StatusCode)
		}

		body, _ := io.ReadAll(res.Body)
		if string(body) != method+" OK" {
			t.Errorf("Expected body '%s OK', got '%s'", method, body)
		}
	}
}

func TestRouterParameters(t *testing.T) {
	router := NewRouter()

	router.Get("/users/{id}/posts/{post_id}", func(w http.ResponseWriter, r *http.Request) error {
		params := r.Context().Value(ParamsKey).(map[string]string)
		w.Write([]byte(fmt.Sprintf("User %s Post %s", params["id"], params["post_id"])))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123/posts/456", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if string(body) != "User 123 Post 456" {
		t.Errorf("Expected 'User 123 Post 456', got '%s'", body)
	}
}

func TestRouterGroupsAndPrefixes(t *testing.T) {
	router := NewRouter()

	router.Group("/api", func(r *Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) error {
			w.Write([]byte("pong"))
			return nil
		})

		r.Group("/v1", func(r2 *Router) {
			r2.Get("/status", func(w http.ResponseWriter, r *http.Request) error {
				w.Write([]byte("v1 ok"))
				return nil
			})
		})
	})

	tests := []struct {
		Path         string
		ExpectedBody string
		ExpectedCode int
	}{
		{"/api/ping", "pong", http.StatusOK},
		{"/api/v1/status", "v1 ok", http.StatusOK},
		{"/ping", "404 page not found\n", http.StatusNotFound},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, tc.Path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		if res.StatusCode != tc.ExpectedCode {
			t.Errorf("Path %s expected code %d, got %d", tc.Path, tc.ExpectedCode, res.StatusCode)
		}

		if string(body) != tc.ExpectedBody {
			t.Errorf("Path %s expected body '%s', got '%s'", tc.Path, tc.ExpectedBody, body)
		}
	}
}

func TestRouterMiddlewareOrder(t *testing.T) {
	router := NewRouter()

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("MW1 Start|"))
			next.ServeHTTP(w, r)
			w.Write([]byte("|MW1 End"))
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("MW2 Start|"))
			next.ServeHTTP(w, r)
			w.Write([]byte("|MW2 End"))
		})
	}

	router.Use(mw1)
	router.Use(mw2)

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("Handler"))
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	expected := "MW1 Start|MW2 Start|Handler|MW2 End|MW1 End"
	if string(body) != expected {
		t.Errorf("Expected middleware execution order '%s', got '%s'", expected, string(body))
	}
}

func TestRouteNaming(t *testing.T) {
	router := NewRouter()

	router.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) error { return nil }).SetRouteName(router, "dashboard")

	if route, exists := router.namedRoutes["dashboard"]; !exists || route.Path != "/dashboard" {
		t.Errorf("Expected named route 'dashboard' to point to '/dashboard'")
	}
}

// Mock controller for resource tests
type UserController struct{}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("index"))
	return nil
}
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("create"))
	return nil
}
func (c *UserController) Store(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("store"))
	return nil
}
func (c *UserController) Show(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("show"))
	return nil
}
func (c *UserController) Edit(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("edit"))
	return nil
}
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("update"))
	return nil
}
func (c *UserController) Destroy(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("destroy"))
	return nil
}

func TestResourceRegistration(t *testing.T) {
	router := NewRouter()
	ctrl := &UserController{}
	router.Resource("users", ctrl)

	expected := []struct {
		method string
		path   string
		name   string
	}{
		{"GET", "/users", "users.index"},
		{"GET", "/users/create", "users.create"},
		{"POST", "/users", "users.store"},
		{"GET", "/users/{user}", "users.show"},
		{"GET", "/users/{user}/edit", "users.edit"},
		{"PUT", "/users/{user}", "users.update"},
		{"DELETE", "/users/{user}", "users.destroy"},
	}

	for _, exp := range expected {
		found := false
		for _, route := range router.allRoutes {
			if route.Method == exp.method && route.Path == exp.path && route.Name == exp.name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected resource route %s %s (%s) not found", exp.method, exp.path, exp.name)
		}
	}
}

func TestApiResourceRegistration(t *testing.T) {
	router := NewRouter()
	ctrl := &UserController{}
	router.ApiResource("posts", ctrl)

	// Should skip Create and Edit
	for _, route := range router.allRoutes {
		if strings.Contains(route.Path, "/create") || strings.Contains(route.Path, "/edit") {
			t.Errorf("ApiResource should not register create/edit routes, got %s", route.Path)
		}
	}
}

func TestSignedRoutes(t *testing.T) {
	router := NewRouter()
	router.Get("/verify/{id}", func(w http.ResponseWriter, r *http.Request) error { return nil }).SetRouteName(router, "verify")

	gen := NewURLGenerator(router, "super-secret-key")
	expires := time.Now().Add(1 * time.Hour)

	url, err := gen.SignedRoute("verify", map[string]string{"id": "123"}, expires)
	if err != nil {
		t.Fatalf("SignedRoute error: %v", err)
	}
	if !strings.Contains(url, "signature=") || !strings.Contains(url, "expires=") {
		t.Errorf("Signed URL missing signature or expires: %s", url)
	}

	// Simulate request
	req := httptest.NewRequest("GET", url, nil)
	if !gen.HasValidSignature(req) {
		t.Error("Expected valid signature for freshly generated URL")
	}
}
