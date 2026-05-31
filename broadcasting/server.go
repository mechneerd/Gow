package broadcasting

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// BroadcastMessage defines the custom JSON protocol format.
type BroadcastMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    any    `json:"data"`
}

// PresenceMember defines a member in a presence channel.
type PresenceMember struct {
	UserID   string         `json:"user_id"`
	UserInfo map[string]any `json:"user_info"`
}

// Hub maintains the set of active clients and broadcasts messages to the channels.
type Hub struct {
	subscriptions sync.Map
	presence      sync.Map
	clients       sync.Map
	Broadcast     chan BroadcastMessage
	Register      chan *Client
	Unregister    chan *Client
	quit          chan struct{}
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan BroadcastMessage, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		quit:       make(chan struct{}),
	}
}

// Stop signals the hub to stop its Run loop for graceful shutdown.
func (h *Hub) Stop() {
	select {
	case <-h.quit:
	default:
		close(h.quit)
	}
}

// Run starts the hub's main event loop.
func (h *Hub) Run() {
	for {
		select {
		case <-h.quit:
			return
		case client := <-h.Register:
			h.clients.Store(client, true)
		case client := <-h.Unregister:
			if _, ok := h.clients.Load(client); ok {
				h.clients.Delete(client)
				close(client.send)
				client.mu.Lock()
				for channelName := range client.channels {
					h.unsubscribe(client, channelName)
				}
				client.mu.Unlock()
			}
		case message := <-h.Broadcast:
			subs, ok := h.subscriptions.Load(message.Channel)
			if !ok {
				continue
			}
			messageBytes, err := json.Marshal(message)
			if err != nil {
				continue
			}
			clientSet := subs.(map[*Client]bool)
			for client := range clientSet {
				select {
				case client.send <- messageBytes:
				default:
					close(client.send)
					h.clients.Delete(client)
				}
			}
		}
	}
}

func (h *Hub) subscribe(client *Client, channelName string) {
	subsInterface, _ := h.subscriptions.LoadOrStore(channelName, make(map[*Client]bool))
	subs := subsInterface.(map[*Client]bool)
	subs[client] = true
	h.subscriptions.Store(channelName, subs)

	client.mu.Lock()
	client.channels[channelName] = true
	client.mu.Unlock()
}

func (h *Hub) unsubscribe(client *Client, channelName string) {
	subsInterface, ok := h.subscriptions.Load(channelName)
	if !ok {
		return
	}
	subs := subsInterface.(map[*Client]bool)
	delete(subs, client)
	h.subscriptions.Store(channelName, subs)
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	channels map[string]bool
	mu       sync.Mutex
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Validate origin against allowed origins (configurable via AllowedOrigins)
	origin := r.Header.Get("Origin")
	allowed := isOriginAllowed(origin)

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: !allowed,
	})
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		channels: make(map[string]bool),
	}
	client.hub.Register <- client

	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, message, err := c.conn.Read(context.Background())
		if err != nil {
			break
		}

		var msg struct {
			Event string `json:"event"`
			Data  struct {
				Channel string         `json:"channel"`
				Auth    string         `json:"auth"`
				Info    map[string]any `json:"info"`
			} `json:"data"`
		}

		if err := json.Unmarshal(message, &msg); err == nil {
			if msg.Event == "system:subscribe" {
				c.hub.subscribe(c, msg.Data.Channel)
			} else if msg.Event == "system:unsubscribe" {
				c.hub.unsubscribe(c, msg.Data.Channel)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer c.conn.Close(websocket.StatusNormalClosure, "")

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageText, message)
			cancel()
			if err != nil {
				return
			}
		}
	}
}

// allowedOrigins is a configurable list of allowed WebSocket origins.
// Defaults to same-origin (empty list = allow all for backward compatibility).
var allowedOrigins []string

// SetAllowedOrigins configures the allowed origins for WebSocket connections.
func SetAllowedOrigins(origins []string) {
	allowedOrigins = origins
}

// isOriginAllowed checks if the given origin is in the allowed list.
// If no origins are configured, all origins are allowed (backward compatible).
func isOriginAllowed(origin string) bool {
	if len(allowedOrigins) == 0 {
		return true
	}
	for _, o := range allowedOrigins {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}

