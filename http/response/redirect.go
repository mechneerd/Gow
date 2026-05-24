package response

import (
	"net/http"
)

// Redirect sends an HTTP redirect response.
func Redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	http.Redirect(w, r, url, status)
}

