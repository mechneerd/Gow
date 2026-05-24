package middleware

import (
	"context"
	"github.com/mechneerd/gow/session"
	"net/http"
	"time"
)

type contextKey string

const SessionKey = contextKey("session")

// StartSession initializes the session for the request.
func StartSession(store session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var sessionID string
			cookie, err := r.Cookie("gow_session")
			if err == nil {
				sessionID = cookie.Value
			}

			manager := session.NewManager(store, sessionID)
			manager.Start()

			// Add manager to context
			ctx := context.WithValue(r.Context(), SessionKey, manager)
			r = r.WithContext(ctx)

			// Execute next handler
			next.ServeHTTP(w, r)

			// Save session and write cookie
			manager.Save()

			http.SetCookie(w, &http.Cookie{
				Name:     "gow_session",
				Value:    manager.ID(),
				Path:     "/",
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(24 * time.Hour),
			})
		})
	}
}

// GetSession is a helper to extract the session from the request context.
func GetSession(r *http.Request) *session.Manager {
	if manager, ok := r.Context().Value(SessionKey).(*session.Manager); ok {
		return manager
	}
	return nil
}

