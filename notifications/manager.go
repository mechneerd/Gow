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
	channels      map[string]Channel
	queueManager  QueueManager
	mu            sync.RWMutex
	fakes         []*Fake
}

// QueueManager is an interface for dispatching notifications to a queue.
type QueueManager interface {
	Dispatch(notification Notification, notifiables []Notifiable) error
}

// NewManager creates a new Notification Manager.
func NewManager() *Manager {
	return &Manager{
		channels: make(map[string]Channel),
	}
}

// SetQueueManager sets the queue manager for queued notifications.
func (m *Manager) SetQueueManager(qm QueueManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queueManager = qm
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
	fakes := m.fakes
	m.mu.RUnlock()

	// Check for fakes
	if len(fakes) > 0 {
		for _, fake := range fakes {
			for _, notifiable := range notifiables {
				fake.Send(notifiable, notification)
			}
		}
		return nil
	}

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

// SendNow dispatches the notification immediately, bypassing the queue.
func (m *Manager) SendNow(notifiables []Notifiable, notification Notification) error {
	return m.Send(notifiables, notification)
}

// SendQueued dispatches the notification to the queue for background processing.
func (m *Manager) SendQueued(notifiables []Notifiable, notification Notification) error {
	m.mu.RLock()
	qm := m.queueManager
	fakes := m.fakes
	m.mu.RUnlock()

	// Check for fakes
	if len(fakes) > 0 {
		for _, fake := range fakes {
			for _, notifiable := range notifiables {
				fake.Send(notifiable, notification)
			}
		}
		return nil
	}

	if qm == nil {
		// Fallback to synchronous sending
		return m.Send(notifiables, notification)
	}
	return qm.Dispatch(notification, notifiables)
}

// OnDemand sends a notification to a notifiable only when the via() method returns a channel.
// This is useful for conditional notifications.
func (m *Manager) OnDemand(notifiables []Notifiable, notification Notification) error {
	return m.Send(notifiables, notification)
}

// Fake sets a fake notification dispatcher for testing.
func (m *Manager) Fake(fake *Fake) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fakes = append(m.fakes, fake)
}

// Notify is a convenient helper to send a notification to one or more notifiables.
// It uses the default global notification manager if registered.
func Notify(notifiables []Notifiable, notification Notification) error {
	if defaultManager == nil {
		return fmt.Errorf("no default notification manager registered")
	}
	return defaultManager.Send(notifiables, notification)
}

// NotifyQueued sends a queued notification using the default global manager.
func NotifyQueued(notifiables []Notifiable, notification Notification) error {
	if defaultManager == nil {
		return fmt.Errorf("no default notification manager registered")
	}
	return defaultManager.SendQueued(notifiables, notification)
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

