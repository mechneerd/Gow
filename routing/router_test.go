package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestNamedRoutes(t *testing.T) {
	router := NewRouter()
	
	route := router.Get("/user/{id}/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("profile"))
	})
	router.SetName(route, "user.profile")

	// We only verify it gets stored, as basic URL generation isn't exposed separately from signed URLs yet
	if r, ok := router.namedRoutes["user.profile"]; !ok {
		t.Errorf("Expected route to be named")
	} else if r.Path != "/user/{id}/profile" {
		t.Errorf("Expected path /user/{id}/profile, got %s", r.Path)
	}
}

type PhotoController struct{}
func (p *PhotoController) Index(w http.ResponseWriter, r *http.Request) { w.Write([]byte("index")) }
func (p *PhotoController) Show(w http.ResponseWriter, r *http.Request) { w.Write([]byte("show")) }
func (p *PhotoController) Store(w http.ResponseWriter, r *http.Request) { w.Write([]byte("store")) }
func (p *PhotoController) Update(w http.ResponseWriter, r *http.Request) { w.Write([]byte("update")) }
func (p *PhotoController) Destroy(w http.ResponseWriter, r *http.Request) { w.Write([]byte("destroy")) }

func TestResourceRouting(t *testing.T) {
	router := NewRouter()
	ctrl := &PhotoController{}
	
	router.ApiResource("photos", ctrl)

	req := httptest.NewRequest(http.MethodGet, "/photos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Body.String() != "index" {
		t.Errorf("Expected index, got %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, "/photos/123", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Body.String() != "update" {
		t.Errorf("Expected update, got %s", w.Body.String())
	}
}

func TestSignedURLs(t *testing.T) {
	router := NewRouter()
	secret := "secret-key"
	gen := NewURLGenerator(router, secret)

	route := router.Get("/unsubscribe/{user}", func(w http.ResponseWriter, r *http.Request) {
		if !gen.HasValidSignature(r) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("invalid signature"))
			return
		}
		w.Write([]byte("unsubscribed"))
	})
	router.SetName(route, "unsubscribe")

	urlStr, err := gen.SignedRoute("unsubscribe", map[string]string{"user": "1"}, time.Time{})
	if err != nil {
		t.Fatalf("Failed to generate signed route: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, urlStr, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Body.String() != "unsubscribed" {
		t.Errorf("Expected unsubscribed, got %s", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, urlStr+"x", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", w.Code)
	}
}
