package routing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"gow/auth"
	"gow/broadcasting"
	"net/http"
	"strings"
)

// BroadcastRoutes registers the standard POST /broadcasting/auth endpoint.
func (r *Router) BroadcastRoutes(channelManager *broadcasting.ChannelManager, appSecret string) {
	r.Post("/broadcasting/auth", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		socketID := req.FormValue("socket_id")
		channelName := req.FormValue("channel_name")

		if socketID == "" || channelName == "" {
			http.Error(w, "socket_id and channel_name are required", http.StatusForbidden)
			return
		}

		user := auth.User(req)
		if user == nil {
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}

		authorized, channelData := channelManager.Authorize(user, channelName)
		if !authorized {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		isPresence := strings.HasPrefix(channelName, "presence-")

		var signatureString string
		var responseData map[string]any

		if isPresence {
			channelDataJSON, err := json.Marshal(channelData)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			signatureString = socketID + ":" + channelName + ":" + string(channelDataJSON)
			
			mac := hmac.New(sha256.New, []byte(appSecret))
			mac.Write([]byte(signatureString))
			authSignature := hex.EncodeToString(mac.Sum(nil))

			responseData = map[string]any{
				"auth":         authSignature,
				"channel_data": string(channelDataJSON),
			}
		} else {
			signatureString = socketID + ":" + channelName
			
			mac := hmac.New(sha256.New, []byte(appSecret))
			mac.Write([]byte(signatureString))
			authSignature := hex.EncodeToString(mac.Sum(nil))

			responseData = map[string]any{
				"auth": authSignature,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	})
}
