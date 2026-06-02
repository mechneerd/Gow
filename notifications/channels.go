package notifications

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// BroadcastNotification represents a notification sent via broadcast channels.
type BroadcastNotification struct {
	Channel string
	Data    map[string]any
}

// BroadcastChannel sends notifications via WebSocket/broadcasting.
type BroadcastChannel struct {
	broadcaster interface {
		Broadcast(channel string, data any)
	}
}

// NewBroadcastChannel creates a new broadcast notification channel.
func NewBroadcastChannel(broadcaster interface {
	Broadcast(channel string, data any)
}) *BroadcastChannel {
	return &BroadcastChannel{broadcaster: broadcaster}
}

// Send sends a broadcast notification.
func (bc *BroadcastChannel) Send(notifiable Notifiable, notification Notification) error {
	channel := notifiable.RouteNotificationFor("broadcast")
	if channel == "" {
		return fmt.Errorf("no broadcast channel specified")
	}

	data := map[string]any{
		"type": fmt.Sprintf("%T", notification),
	}

	bc.broadcaster.Broadcast(channel, data)
	return nil
}

// SMSNotification represents a notification sent via SMS.
type SMSNotification struct {
	To      string
	Message string
}

// SMSChannel sends notifications via SMS provider.
type SMSChannel struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

// NewSMSChannel creates a new SMS notification channel.
func NewSMSChannel(apiKey, apiURL string) *SMSChannel {
	return &SMSChannel{
		apiKey:     apiKey,
		apiURL:     apiURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// SendSMS sends an SMS notification.
func (sc *SMSChannel) SendSMS(to string, message string) error {
	payload := fmt.Sprintf(`{"to":"%s","message":"%s"}`, to, message)
	req, err := http.NewRequest("POST", sc.apiURL, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sc.apiKey)

	resp, err := sc.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SMS API returned status %d", resp.StatusCode)
	}
	return nil
}

// NotificationBell provides an in-app notification bell.
type NotificationBell struct {
	store interface {
		GetByUserID(ctx context.Context, userID string) ([]DatabaseNotificationModel, error)
		Update(ctx context.Context, notification *DatabaseNotificationModel) error
	}
}

// NewNotificationBell creates a new in-app notification bell.
func NewNotificationBell(store interface {
	GetByUserID(ctx context.Context, userID string) ([]DatabaseNotificationModel, error)
	Update(ctx context.Context, notification *DatabaseNotificationModel) error
}) *NotificationBell {
	return &NotificationBell{store: store}
}

// UnreadCount returns the count of unread notifications for a user.
func (nb *NotificationBell) UnreadCount(userID string) (int, error) {
	notifications, err := nb.store.GetByUserID(context.Background(), userID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, n := range notifications {
		if n.ReadAt == nil {
			count++
		}
	}
	return count, nil
}

// MarkAllRead marks all notifications as read for a user.
func (nb *NotificationBell) MarkAllRead(userID string) error {
	notifications, err := nb.store.GetByUserID(context.Background(), userID)
	if err != nil {
		return err
	}
	now := time.Now()
	for _, n := range notifications {
		if n.ReadAt == nil {
			n.ReadAt = &now
			nb.store.Update(context.Background(), &n)
		}
	}
	return nil
}
