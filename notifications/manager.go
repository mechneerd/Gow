package notifications

import "fmt"

// Channel represents an interface for sending notifications.
type Channel interface {
	Send(notifiable Notifiable, notification Notification) error
}

// Manager resolves channels and dispatches notifications.
type Manager struct {
	channels map[string]Channel
}

// NewManager creates a new Notification Manager.
func NewManager() *Manager {
	return &Manager{
		channels: make(map[string]Channel),
	}
}

// Extend registers a custom channel driver.
func (m *Manager) Extend(name string, channel Channel) {
	m.channels[name] = channel
}

// Send dispatches the notification to the given notifiables.
func (m *Manager) Send(notifiables []Notifiable, notification Notification) error {
	for _, notifiable := range notifiables {
		channels := notification.Via(notifiable)
		for _, channelName := range channels {
			channel, exists := m.channels[channelName]
			if !exists {
				return fmt.Errorf("notification channel [%s] not found", channelName)
			}
			
			err := channel.Send(notifiable, notification)
			if err != nil {
				// We might want to log this and continue, or return immediately
				return err
			}
		}
	}
	return nil
}
