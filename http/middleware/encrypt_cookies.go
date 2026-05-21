package middleware

import (
	"gow/cookie"
	"net/http"
)

// EncryptCookies middleware automatically decrypts incoming cookies
// and provides a hook (or relies on the cookie manager) for outgoing cookies.
func EncryptCookies(manager *cookie.Manager, except []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Decrypt incoming cookies
			for _, c := range r.Cookies() {
				// Check exemption list
				exempt := false
				for _, e := range except {
					if c.Name == e {
						exempt = true
						break
					}
				}

				if !exempt {
					decrypted, err := manager.Decrypt(c.Value)
					if err == nil {
						c.Value = decrypted
						// Update the request with the decrypted cookie so downstream handlers see plaintext
						r.AddCookie(c) 
					}
				}
			}

			// In Go, intercepting ResponseWriter to encrypt outgoing cookies 
			// requires a custom ResponseWriter wrapper, but for simplicity
			// controllers should use the cookie.Manager.Set() method directly.
			
			next.ServeHTTP(w, r)
		})
	}
}
