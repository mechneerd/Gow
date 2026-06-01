package response

import (
	"net/http"
	"strings"
)

// Redirect sends an HTTP redirect response.
func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	http.Redirect(w, r, url, status)
}

// RedirectBack redirects to the previous URL (from Referer header) or fallback.
func RedirectBack(w http.ResponseWriter, r *http.Request, fallback string, status int) {
	if referer := r.Header.Get("Referer"); referer != "" {
		// Only redirect to same host to prevent open redirect
		if strings.HasPrefix(referer, r.Host) || strings.Contains(referer, r.Host) {
			http.Redirect(w, r, referer, status)
			return
		}
	}
	http.Redirect(w, r, fallback, status)
}

// RedirectToRoute redirects to a named route.
// Usage: response.RedirectToRoute(w, r, "user.show", map[string]string{"id": "1"}, http.StatusFound)
func RedirectToRoute(w http.ResponseWriter, r *http.Request, name string, params map[string]string, status int) {
	// Get the router from context (if available)
	if router, ok := r.Context().Value("router").(routeRouter); ok {
		url := router.Route(name, params)
		if url != "" {
			http.Redirect(w, r, url, status)
			return
		}
	}
	// Fallback: redirect to root
	http.Redirect(w, r, "/", status)
}

// Flash stores a flash message in the session.
func Flash(w http.ResponseWriter, r *http.Request, key string, value any) {
	// Flash data is stored in session for next request
	// This is a simplified implementation - full version would use session manager
	if session, ok := r.Context().Value("session").(interface{ Set(string, any) }); ok {
		session.Set("flash."+key, value)
	}
}

// FlashInput stores all input data in flash for next request.
func FlashInput(w http.ResponseWriter, r *http.Request) {
	// Store current input in session flash
	if session, ok := r.Context().Value("session").(interface{ Set(string, any) }); ok {
		input := make(map[string]any)
		for k, v := range r.URL.Query() {
			if len(v) > 0 {
				input[k] = v[0]
			}
		}
		session.Set("flash.old_input", input)
	}
}

// routeRouter is an interface to avoid circular imports
type routeRouter interface {
	Route(name string, params ...map[string]string) string
}

