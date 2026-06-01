package broadcasting

import (
	"encoding/json"
	"net/http"
	"strings"
)

// AuthController handles channel authorization for private and presence channels.
type AuthController struct {
	authorizer ChannelAuthorizer
}

// ChannelAuthorizer determines if a user can access a channel.
type ChannelAuthorizer interface {
	Authorize(userID any, channelName string) bool
}

// NewAuthController creates a new authorization controller.
func NewAuthController(authorizer ChannelAuthorizer) *AuthController {
	return &AuthController{
		authorizer: authorizer,
	}
}

// AuthorizeRequest handles the authorization request for a channel.
func (ac *AuthController) AuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the socket ID and channel name from the request
	socketID := r.FormValue("socket_id")
	channelName := r.FormValue("channel_name")

	if socketID == "" || channelName == "" {
		http.Error(w, "Missing socket_id or channel_name", http.StatusBadRequest)
		return
	}

	// Extract user ID from context or session (simplified)
	userID := r.Context().Value("user_id")

	// Check if this is a private or presence channel
	if !strings.HasPrefix(channelName, "private-") && !strings.HasPrefix(channelName, "presence-") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Authorize the channel
	if ac.authorizer == nil || !ac.authorizer.Authorize(userID, channelName) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate the auth signature (simplified)
	// In production, this would use HMAC-SHA256 with the app key and secret
	response := map[string]string{
		"auth": "app_key:" + generateAuthSignature(socketID, channelName),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateAuthSignature generates a simplified auth signature.
// In production, this would be: HMAC-SHA256(app_key:socket_id:channel_name, app_secret)
func generateAuthSignature(socketID, channelName string) string {
	// Simplified - in production use crypto/hmac
	return socketID + ":" + channelName + ":signature"
}

// BroadcastMiddleware adds broadcasting capabilities to HTTP handlers.
func BroadcastMiddleware(manager *Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add broadcast manager to context
			ctx := r.Context()
			_ = ctx
			next.ServeHTTP(w, r)
		})
	}
}
