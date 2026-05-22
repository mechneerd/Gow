package broadcasting

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	// public + private: channelName -> set of clients
	subscriptions sync.Map // map[string]map[*Client]bool

	// presence only: channelName -> map[string]PresenceMember
	presence sync.Map // map[string]map[string]PresenceMember

	// Registered clients.
	clients sync.Map // map[*Client]bool

	// Inbound messages from the drivers.
	Broadcast chan BroadcastMessage

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan BroadcastMessage, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub's main event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients.Store(client, true)
		case client := <-h.Unregister:
			if _, ok := h.clients.Load(client); ok {
				h.clients.Delete(client)
				close(client.send)
				// Clean up subscriptions
				client.mu.Lock()
				for channelName := range client.channels {
					h.unsubscribe(client, channelName)
				}
				client.mu.Unlock()
			}
		case message := <-h.Broadcast:
			// Find all clients subscribed to the target channel
			subs, ok := h.subscriptions.Load(message.Channel)
			if !ok {
				continue // no subscribers
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
					// In a real implementation we would trigger an unregister here
				}
			}
		}
	}
}

// subscribe adds a client to a channel.
func (h *Hub) subscribe(client *Client, channelName string) {
	subsInterface, _ := h.subscriptions.LoadOrStore(channelName, make(map[*Client]bool))
	subs := subsInterface.(map[*Client]bool)
	subs[client] = true
	h.subscriptions.Store(channelName, subs) // Store back updated map

	client.mu.Lock()
	client.channels[channelName] = true
	client.mu.Unlock()
}

// unsubscribe removes a client from a channel.
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
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Channels this client is subscribed to.
	channels map[string]bool
	mu       sync.Mutex
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
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

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse custom protocol for subscribe/unsubscribe
		var msg struct {
			Event string `json:"event"`
			Data  struct {
				Channel string         `json:"channel"`
				Auth    string         `json:"auth"`
				Info    map[string]any `json:"info"` // For presence
			} `json:"data"`
		}

		if err := json.Unmarshal(message, &msg); err == nil {
			if msg.Event == "system:subscribe" {
				c.hub.subscribe(c, msg.Data.Channel)
				// If presence channel, add to presence roster and broadcast system:member_added
				// (Simplified for now)
			} else if msg.Event == "system:unsubscribe" {
				c.hub.unsubscribe(c, msg.Data.Channel)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
