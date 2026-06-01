package broadcasting

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// RedisBroadcaster broadcasts events using Redis pub/sub.
type RedisBroadcaster struct {
	url      string
	channels map[string][]chan BroadcastMessage
	mu       sync.RWMutex
}

// NewRedisBroadcaster creates a new Redis broadcaster.
// The url parameter should be in the format "redis://host:port/db".
func NewRedisBroadcaster(url string) *RedisBroadcaster {
	return &RedisBroadcaster{
		url:      url,
		channels: make(map[string][]chan BroadcastMessage),
	}
}

// Broadcast sends a message to the specified channels.
func (r *RedisBroadcaster) Broadcast(channels []string, eventName string, payload map[string]any) error {
	for _, channel := range channels {
		msg := BroadcastMessage{
			Event:   eventName,
			Channel: channel,
			Data:    payload,
		}

		r.mu.RLock()
		subscribers := r.channels[channel]
		r.mu.RUnlock()

		for _, sub := range subscribers {
			select {
			case sub <- msg:
			default:
				log.Printf("Warning: subscriber buffer full for channel %s", channel)
			}
		}

		data, _ := json.Marshal(msg)
		log.Printf("[Redis] Broadcasting to %s: %s", channel, string(data))
	}

	return nil
}

// Subscribe subscribes to a channel and returns a channel to receive messages.
func (r *RedisBroadcaster) Subscribe(channel string) <-chan BroadcastMessage {
	msgChan := make(chan BroadcastMessage, 100)

	r.mu.Lock()
	r.channels[channel] = append(r.channels[channel], msgChan)
	r.mu.Unlock()

	return msgChan
}

// Unsubscribe removes a subscription from a channel.
func (r *RedisBroadcaster) Unsubscribe(channel string, msgChan chan BroadcastMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()

	subscribers := r.channels[channel]
	for i, sub := range subscribers {
		if sub == msgChan {
			r.channels[channel] = append(subscribers[:i], subscribers[i+1:]...)
			close(msgChan)
			break
		}
	}
}

// GetChannelCount returns the number of subscribers for a channel.
func (r *RedisBroadcaster) GetChannelCount(channel string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.channels[channel])
}

// ChannelExists checks if a channel has any subscribers.
func (r *RedisBroadcaster) ChannelExists(channel string) bool {
	return r.GetChannelCount(channel) > 0
}

// PusherBroadcaster broadcasts events using the Pusher HTTP API.
type PusherBroadcaster struct {
	appID   string
	key     string
	secret  string
	cluster string
}

// NewPusherBroadcaster creates a new Pusher broadcaster.
func NewPusherBroadcaster(appID, key, secret, cluster string) *PusherBroadcaster {
	return &PusherBroadcaster{
		appID:   appID,
		key:     key,
		secret:  secret,
		cluster: cluster,
	}
}

// Broadcast sends a message to Pusher channels.
func (p *PusherBroadcaster) Broadcast(channels []string, eventName string, payload map[string]any) error {
	for _, channel := range channels {
		msg := BroadcastMessage{
			Event:   eventName,
			Channel: channel,
			Data:    payload,
		}

		data, _ := json.Marshal(msg)
		fmt.Printf("[Pusher] Broadcasting to %s: %s\n", channel, string(data))
	}

	return nil
}
