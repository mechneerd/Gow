package notifications

// Notifiable defines the interface for an entity that can receive notifications.
type Notifiable interface {
	RouteNotificationFor(channel string) string
}

// Notification represents a message to be sent via multiple channels.
type Notification interface {
	// Via determines which channels the notification should be sent through.
	Via(notifiable Notifiable) []string
}
