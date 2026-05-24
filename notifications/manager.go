package notifications

import (
	"fmt"
	"sync"
)

// Channel represents an interface for sending notifications.
type Channel interface {
	Send(notifiable Notifiable, notification Notification) error
}

// Manager resolves channels and dispatches notifications.
type Manager struct {
	channels map[string]Channel
	mu       sync.RWMutex
}

// NewManager creates a new Notification Manager.
func NewManager() *Manager {
	return &Manager{
		channels: make(map[string]Channel),
	}
}

// Extend registers a custom channel driver.
func (m *Manager) Extend(name string, channel Channel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.channels[name] = channel
}

// Send dispatches the notification to the given notifiables.
func (m *Manager) Send(notifiables []Notifiable, notification Notification) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, notifiable := range notifiables {
		channels := notification.Via(notifiable)
		for _, channelName := range channels {
			channel, exists := m.channels[channelName]
			if !exists {
				return fmt.Errorf("notification channel [%s] not found", channelName)
			}

			if err := channel.Send(notifiable, notification); err != nil {
				return err
			}
		}
	}
	return nil
}

// Notify is a convenient helper to send a notification to one or more notifiables.
// It uses the default global notification manager if registered.
func Notify(notifiables []Notifiable, notification Notification) error {
	if defaultManager == nil {
		return fmt.Errorf("no default notification manager registered")
	}
	return defaultManager.Send(notifiables, notification)
}

// defaultManager is the global notification manager used by the Notify helper.
var defaultManager *Manager
var once sync.Once

// SetDefaultManager sets the global notification manager (usually called from a service provider).
func SetDefaultManager(m *Manager) {
	once.Do(func() {
		defaultManager = m
	})
}

