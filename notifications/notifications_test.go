package notifications

import (
	"context"
	"testing"
	"time"
)

// MockNotifiable implements Notifiable for testing.
type MockNotifiable struct {
	Channel string
}

func (m *MockNotifiable) RouteNotificationFor(channel string) string {
	return m.Channel
}

// MockNotification implements Notification for testing.
type MockNotification struct{}

func (n *MockNotification) Via(notifiable Notifiable) []string {
	return []string{"mail", "database"}
}

func TestNewBroadcastChannel(t *testing.T) {
	bc := NewBroadcastChannel(nil)
	if bc == nil {
		t.Fatal("NewBroadcastChannel returned nil")
	}
}

func TestNewSMSChannel(t *testing.T) {
	sc := NewSMSChannel("api-key", "https://api.sms.com/send")
	if sc == nil {
		t.Fatal("NewSMSChannel returned nil")
	}
	if sc.apiKey != "api-key" {
		t.Errorf("expected apiKey 'api-key', got %q", sc.apiKey)
	}
}

func TestNewNotificationBell(t *testing.T) {
	store := &mockNotificationStore{
		notifications: []DatabaseNotificationModel{},
	}
	bell := NewNotificationBell(store)
	if bell == nil {
		t.Fatal("NewNotificationBell returned nil")
	}
}

func TestNotificationBell_UnreadCount(t *testing.T) {
	now := time.Now()
	store := &mockNotificationStore{
		notifications: []DatabaseNotificationModel{
			{ID: "1", ReadAt: nil},
			{ID: "2", ReadAt: &now},
			{ID: "3", ReadAt: nil},
		},
	}
	bell := NewNotificationBell(store)

	count, err := bell.UnreadCount("user-1")
	if err != nil {
		t.Fatalf("UnreadCount failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 unread, got %d", count)
	}
}

func TestNotificationBell_MarkAllRead(t *testing.T) {
	store := &mockNotificationStore{
		notifications: []DatabaseNotificationModel{
			{ID: "1", ReadAt: nil},
			{ID: "2", ReadAt: nil},
		},
	}
	bell := NewNotificationBell(store)

	err := bell.MarkAllRead("user-1")
	if err != nil {
		t.Fatalf("MarkAllRead failed: %v", err)
	}

	// Verify all notifications were marked as read
	for _, n := range store.notifications {
		if n.ReadAt == nil {
			t.Errorf("notification %s should have been marked as read", n.ID)
		}
	}
}

func TestMockNotifiable(t *testing.T) {
	n := &MockNotifiable{Channel: "mail"}
	if n.RouteNotificationFor("mail") != "mail" {
		t.Error("expected RouteNotificationFor to return correct channel")
	}
}

func TestMockNotification(t *testing.T) {
	n := &MockNotification{}
	channels := n.Via(&MockNotifiable{})
	if len(channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(channels))
	}
}

type mockNotificationStore struct {
	notifications []DatabaseNotificationModel
}

func (m *mockNotificationStore) GetByUserID(ctx context.Context, userID string) ([]DatabaseNotificationModel, error) {
	return m.notifications, nil
}

func (m *mockNotificationStore) Update(ctx context.Context, notification *DatabaseNotificationModel) error {
	for i, n := range m.notifications {
		if n.ID == notification.ID {
			m.notifications[i] = *notification
		}
	}
	return nil
}
