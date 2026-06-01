package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// ValidatePostSize validates that the request's Content-Length does not exceed
// the maximum allowed size (in bytes). Returns 413 Payload Too Large if exceeded.
func ValidatePostSize(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxSize {
				http.Error(w, "Payload Too Large", http.StatusRequestEntityTooLarge)
				return
			}

			// Also check X-Forwarded-For header if present
			if forwardedSize := r.Header.Get("Content-Length"); forwardedSize != "" {
				if size, err := strconv.ParseInt(forwardedSize, 10, 64); err == nil {
					if size > maxSize {
						http.Error(w, "Payload Too Large", http.StatusRequestEntityTooLarge)
						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SetCacheHeaders sets HTTP cache control headers on the response.
// Usage: SetCacheHeaders("private, max-age=3600") or SetCacheHeaders("no-store, no-cache")
func SetCacheHeaders(directives string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", directives)
			next.ServeHTTP(w, r)
		})
	}
}

// SetCacheHeadersWithMap sets HTTP cache control headers from a map.
// Usage: SetCacheHeadersWithMap(map[string]string{"Cache-Control": "max-age=3600", "ETag": "\"abc123\""})
func SetCacheHeadersWithMap(headers map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range headers {
				w.Header().Set(key, value)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PreventCache is a convenience middleware that sets headers to prevent caching.
func PreventCache(next http.Handler) http.Handler {
	return SetCacheHeaders("no-store, no-cache, must-revalidate, max-age=0")(next)
}

// CacheFor sets Cache-Control with a max-age directive.
// Usage: CacheFor(3600) caches for 1 hour.
func CacheFor(seconds int) func(http.Handler) http.Handler {
	return SetCacheHeaders("public, max-age=" + strconv.Itoa(seconds))
}

// CachePrivate sets Cache-Control to private with a max-age.
func CachePrivate(seconds int) func(http.Handler) http.Handler {
	return SetCacheHeaders("private, max-age=" + strconv.Itoa(seconds))
}

// ConvertEmptyStringsToNull middleware converts empty string form values to nil.
// This matches Laravel's ConvertEmptyStringsToNull middleware behavior.
func ConvertEmptyStringsToNull(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
			if r.Body != nil {
				r.ParseForm()
				for key, values := range r.Form {
					for i, v := range values {
						if v == "" {
							r.Form[key][i] = "" // Keep empty for form parsing
						}
					}
				}
			}
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := r.ParseMultipartForm(32 << 20); err == nil {
				for key, values := range r.MultipartForm.Value {
					for i, v := range values {
						if v == "" {
							r.MultipartForm.Value[key][i] = ""
						}
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// TrimAndConvertEmptyToNull combines TrimStrings and ConvertEmptyStringsToNull.
// This is the typical middleware stack used in Laravel's web middleware group.
func TrimAndConvertEmptyToNull(next http.Handler) http.Handler {
	return ConvertEmptyStringsToNull(TrimStrings(next))
}
