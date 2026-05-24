package broadcasting

// WebSocketDriver implements the Broadcaster interface to push messages
// into the native WebSocket server Hub.
type WebSocketDriver struct {
	hub *Hub
}

// NewWebSocketDriver creates a new WebSocket driver.
func NewWebSocketDriver(hub *Hub) *WebSocketDriver {
	return &WebSocketDriver{
		hub: hub,
	}
}

// Broadcast sends the event payload to the specified channels via the Hub.
func (d *WebSocketDriver) Broadcast(channels []string, eventName string, payload map[string]any) error {
	for _, channel := range channels {
		msg := BroadcastMessage{
			Event:   eventName,
			Channel: channel,
			Data:    payload,
		}
		
		// Push message to the hub's broadcast channel to be routed to all subscribers
		d.hub.Broadcast <- msg
	}
	
	return nil
}

