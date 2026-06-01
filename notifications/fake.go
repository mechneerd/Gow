package notifications

import (
	"sync"
)

// Fake is a fake notification dispatcher for testing.
type Fake struct {
	mu            sync.RWMutex
	notifications []Notification
	recipients    map[string][]Notification
}

// NewFake creates a new fake notification dispatcher.
func NewFake() *Fake {
	return &Fake{
		notifications: make([]Notification, 0),
		recipients:    make(map[string][]Notification),
	}
}

// Send captures the notification instead of dispatching it.
func (f *Fake) Send(notifiable any, notification Notification) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.notifications = append(f.notifications, notification)
}

// SendTo captures notifications sent to a specific recipient.
func (f *Fake) SendTo(notifiable any, notification Notification, channel string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.notifications = append(f.notifications, notification)
	if n, ok := notifiable.(interface{ NotificationName() string }); ok {
		f.recipients[n.NotificationName()] = append(f.recipients[n.NotificationName()], notification)
	}
}

// GetNotifications returns all captured notifications.
func (f *Fake) GetNotifications() []Notification {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make([]Notification, len(f.notifications))
	copy(result, f.notifications)
	return result
}

// GetLastNotification returns the last captured notification, or nil if none.
func (f *Fake) GetLastNotification() Notification {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if len(f.notifications) == 0 {
		return nil
	}
	return f.notifications[len(f.notifications)-1]
}

// GetNotificationCount returns the number of notifications dispatched.
func (f *Fake) GetNotificationCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.notifications)
}

// HasSent checks if any notification was dispatched.
func (f *Fake) HasSent() bool {
	return f.GetNotificationCount() > 0
}

// HasSentNotification checks if a specific notification was dispatched.
func (f *Fake) HasSentNotification(notification any) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, n := range f.notifications {
		if n == notification {
			return true
		}
	}
	return false
}

// HasSentTo checks if a notification was sent to a specific recipient.
func (f *Fake) HasSentTo(notifiable any, notification any) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if n, ok := notifiable.(interface{ NotificationName() string }); ok {
		for _, notif := range f.recipients[n.NotificationName()] {
			if notif == notification {
				return true
			}
		}
	}
	return false
}

// AssertSent asserts that at least one notification was dispatched.
func (f *Fake) AssertSent() bool {
	return f.HasSent()
}

// AssertNotSent asserts that no notifications were dispatched.
func (f *Fake) AssertNotSent() bool {
	return f.GetNotificationCount() == 0
}

// AssertSentCount asserts the exact number of notifications dispatched.
func (f *Fake) AssertSentCount(count int) bool {
	return f.GetNotificationCount() == count
}

// AssertSentNotification asserts that a specific notification was dispatched.
func (f *Fake) AssertSentNotification(notification any) bool {
	return f.HasSentNotification(notification)
}

// AssertNotSentNotification asserts that a specific notification was NOT dispatched.
func (f *Fake) AssertNotSentNotification(notification any) bool {
	return !f.HasSentNotification(notification)
}

// AssertSentTo asserts a notification was sent to a specific recipient.
func (f *Fake) AssertSentTo(notifiable any, notification any) bool {
	return f.HasSentTo(notifiable, notification)
}

// Clear resets all captured notifications.
func (f *Fake) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.notifications = make([]Notification, 0)
	f.recipients = make(map[string][]Notification)
}
