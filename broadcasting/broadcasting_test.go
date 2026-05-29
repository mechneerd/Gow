package broadcasting

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

type TestUser struct {
	ID   string
	Name string
}

func TestChannelManager(t *testing.T) {
	cm := NewChannelManager()

	cm.Channel("private-chat.{id}", func(user any, channelName string) (bool, map[string]any) {
		u := user.(TestUser)
		parts := strings.Split(channelName, ".")
		id := parts[1]
		return u.ID == id, nil
	})

	cm.Channel("presence-room.*", func(user any, channelName string) (bool, map[string]any) {
		u := user.(TestUser)
		return true, map[string]any{"name": u.Name}
	})

	// Test private channel authorization
	user1 := TestUser{ID: "1", Name: "Alice"}
	user2 := TestUser{ID: "2", Name: "Bob"}

	auth1, _ := cm.Authorize(user1, "private-chat.1")
	if !auth1 {
		t.Error("User 1 should be authorized for private-chat.1")
	}

	auth2, _ := cm.Authorize(user2, "private-chat.1")
	if auth2 {
		t.Error("User 2 should NOT be authorized for private-chat.1")
	}

	// Test presence channel authorization
	auth3, data := cm.Authorize(user1, "presence-room.general")
	if !auth3 {
		t.Error("User 1 should be authorized for presence-room.general")
	}
	if data["name"] != "Alice" {
		t.Errorf("Expected presence data name Alice, got %v", data["name"])
	}
}

func TestWebSocketServer(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.Dial(context.Background(), wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to websocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Wait briefly for connection registration
	time.Sleep(50 * time.Millisecond)

	// Subscribe to a channel
	subscribeMsg := map[string]any{
		"event": "system:subscribe",
		"data": map[string]any{
			"channel": "chat.room",
			"auth":    "fake-auth",
		},
	}
	msgBytes, _ := json.Marshal(subscribeMsg)
	conn.Write(context.Background(), websocket.MessageText, msgBytes)

	// Wait briefly for subscription to process
	time.Sleep(50 * time.Millisecond)

	// Push a message to the hub
	hub.Broadcast <- BroadcastMessage{
		Event:   "chat.message",
		Channel: "chat.room",
		Data:    map[string]any{"text": "Hello World"},
	}

	// Read the broadcasted message
	readCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, rawMsg, err := conn.Read(readCtx)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var receivedMsg BroadcastMessage
	err = json.Unmarshal(rawMsg, &receivedMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if receivedMsg.Event != "chat.message" {
		t.Errorf("Expected event chat.message, got %s", receivedMsg.Event)
	}
	
	dataBytes, _ := json.Marshal(receivedMsg.Data)
	var dataMap map[string]string
	json.Unmarshal(dataBytes, &dataMap)
	
	if dataMap["text"] != "Hello World" {
		t.Errorf("Expected data text 'Hello World', got %v", dataMap["text"])
	}
}

func TestWebSocketDriver(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	driver := NewWebSocketDriver(hub)
	
	// Ensure we don't block if there are no subscribers
	err := driver.Broadcast([]string{"test.channel"}, "test.event", map[string]any{"foo": "bar"})
	if err != nil {
		t.Fatalf("Driver Broadcast returned error: %v", err)
	}
}

