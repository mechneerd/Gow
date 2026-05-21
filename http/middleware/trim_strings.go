package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TrimStrings middleware trims whitespace from incoming string parameters.
// Currently, it handles URL query parameters and x-www-form-urlencoded body.
func TrimStrings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Trim Query Parameters
		query := r.URL.Query()
		for key, values := range query {
			for i, v := range values {
				query[key][i] = strings.TrimSpace(v)
			}
		}
		r.URL.RawQuery = query.Encode()

		// Trim Form Values
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				parsed, err := url.ParseQuery(string(bodyBytes))
				if err == nil {
					for key, values := range parsed {
						for i, v := range values {
							parsed[key][i] = strings.TrimSpace(v)
						}
					}
					newBody := parsed.Encode()
					r.Body = io.NopCloser(strings.NewReader(newBody))
					r.ContentLength = int64(len(newBody))
				} else {
					// Fallback
					r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
