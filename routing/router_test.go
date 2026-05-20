package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterBasic(t *testing.T) {
	router := NewRouter()
	
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home"))
	})
	
	router.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("about"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Body.String() != "home" {
		t.Errorf("Expected home, got %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/about", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Body.String() != "about" {
		t.Errorf("Expected about, got %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/notfound", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestRouterParameters(t *testing.T) {
	router := NewRouter()
	
	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(ParamsKey).(map[string]string)
		w.Write([]byte("user:" + params["id"]))
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Body.String() != "user:123" {
		t.Errorf("Expected user:123, got %s", w.Body.String())
	}
}

func TestRouterMiddlewareAndGroups(t *testing.T) {
	router := NewRouter()
	
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Test", "1")
			next.ServeHTTP(w, r)
		})
	}

	router.Group("/api", func(r *Router) {
		r.Use(mw)
		r.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("pong"))
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Body.String() != "pong" {
		t.Errorf("Expected pong, got %s", w.Body.String())
	}
	if w.Header().Get("X-Test") != "1" {
		t.Errorf("Expected X-Test: 1, got %s", w.Header().Get("X-Test"))
	}
}
