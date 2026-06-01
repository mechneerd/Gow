package notifications

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mechneerd/gow/database/orm"
	"github.com/mechneerd/gow/database/query"
)

// DatabaseNotification defines the interface for notifications that can be stored in the database.
type DatabaseNotification interface {
	ToDatabase(notifiable Notifiable) map[string]any
}

// DatabaseNotificationModel represents a notification stored in the database.
type DatabaseNotificationModel struct {
	ID             string
	Type           string
	NotifiableType string
	NotifiableID   string
	Data           map[string]any
	ReadAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// IsRead returns true if the notification has been read.
func (n *DatabaseNotificationModel) IsRead() bool {
	return n.ReadAt != nil
}

// DatabaseChannel stores notifications in a database table.
type DatabaseChannel struct {
	db *orm.DB
}

// NewDatabaseChannel creates a new DatabaseChannel.
func NewDatabaseChannel(db *orm.DB) *DatabaseChannel {
	return &DatabaseChannel{db: db}
}

// Send stores the notification in the database.
func (c *DatabaseChannel) Send(notifiable Notifiable, notification Notification) error {
	dbNotification, ok := notification.(DatabaseNotification)
	if !ok {
		return fmt.Errorf("notification does not implement DatabaseNotification")
	}

	data := dbNotification.ToDatabase(notifiable)

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	notifiableType, notifiableID := getNotifiableKey(notifiable)

	d, err := c.db.Dialect()
	if err != nil {
		return err
	}

	builder := query.NewBuilder(c.db.RawDB(), d)
	builder.Table("notifications")

	_, err = builder.Insert(map[string]any{
		"type":            fmt.Sprintf("%T", notification),
		"notifiable_type": notifiableType,
		"notifiable_id":   notifiableID,
		"data":            string(bytes),
		"read_at":         nil,
		"created_at":      time.Now().Unix(),
		"updated_at":      time.Now().Unix(),
	})

	return err
}

// MarkAsRead marks a notification as read.
func (c *DatabaseChannel) MarkAsRead(notifiable Notifiable, notificationID string) error {
	notifiableType, notifiableID := getNotifiableKey(notifiable)

	d, err := c.db.Dialect()
	if err != nil {
		return err
	}

	builder := query.NewBuilder(c.db.RawDB(), d)
	builder.Table("notifications")
	builder.Where("id", "=", notificationID)
	builder.Where("notifiable_type", "=", notifiableType)
	builder.Where("notifiable_id", "=", notifiableID)

	_, err = builder.Update(map[string]any{
		"read_at":    time.Now(),
		"updated_at": time.Now().Unix(),
	})
	return err
}

// MarkAsUnread marks a notification as unread (clears read_at).
func (c *DatabaseChannel) MarkAsUnread(notifiable Notifiable, notificationID string) error {
	notifiableType, notifiableID := getNotifiableKey(notifiable)

	d, err := c.db.Dialect()
	if err != nil {
		return err
	}

	builder := query.NewBuilder(c.db.RawDB(), d)
	builder.Table("notifications")
	builder.Where("id", "=", notificationID)
	builder.Where("notifiable_type", "=", notifiableType)
	builder.Where("notifiable_id", "=", notifiableID)

	_, err = builder.Update(map[string]any{
		"read_at":    nil,
		"updated_at": time.Now().Unix(),
	})
	return err
}

// MarkAllAsRead marks all notifications for a notifiable as read.
func (c *DatabaseChannel) MarkAllAsRead(notifiable Notifiable) error {
	notifiableType, notifiableID := getNotifiableKey(notifiable)

	d, err := c.db.Dialect()
	if err != nil {
		return err
	}

	builder := query.NewBuilder(c.db.RawDB(), d)
	builder.Table("notifications")
	builder.Where("notifiable_type", "=", notifiableType)
	builder.Where("notifiable_id", "=", notifiableID)
	builder.Where("read_at", "IS", nil)

	_, err = builder.Update(map[string]any{
		"read_at":    time.Now(),
		"updated_at": time.Now().Unix(),
	})
	return err
}

// GetNotifications returns all notifications for a notifiable.
func (c *DatabaseChannel) GetNotifications(notifiable Notifiable) ([]*DatabaseNotificationModel, error) {
	notifiableType, notifiableID := getNotifiableKey(notifiable)

	d, err := c.db.Dialect()
	if err != nil {
		return nil, err
	}

	builder := query.NewBuilder(c.db.RawDB(), d)
	builder.Table("notifications")
	builder.Where("notifiable_type", "=", notifiableType)
	builder.Where("notifiable_id", "=", notifiableID)
	builder.OrderBy("created_at", "DESC")

	rows, err := builder.Get()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*DatabaseNotificationModel
	for rows.Next() {
		n := &DatabaseNotificationModel{}
		var dataStr string
		var readAt *time.Time
		var createdAt, updatedAt int64

		if err := rows.Scan(&n.ID, &n.Type, &n.NotifiableType, &n.NotifiableID, &dataStr, &readAt, &createdAt, &updatedAt); err != nil {
			continue
		}

		json.Unmarshal([]byte(dataStr), &n.Data)
		n.ReadAt = readAt
		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// GetUnreadNotifications returns only unread notifications.
func (c *DatabaseChannel) GetUnreadNotifications(notifiable Notifiable) ([]*DatabaseNotificationModel, error) {
	all, err := c.GetNotifications(notifiable)
	if err != nil {
		return nil, err
	}

	var unread []*DatabaseNotificationModel
	for _, n := range all {
		if !n.IsRead() {
			unread = append(unread, n)
		}
	}
	return unread, nil
}

// GetReadNotifications returns only read notifications.
func (c *DatabaseChannel) GetReadNotifications(notifiable Notifiable) ([]*DatabaseNotificationModel, error) {
	all, err := c.GetNotifications(notifiable)
	if err != nil {
		return nil, err
	}

	var read []*DatabaseNotificationModel
	for _, n := range all {
		if n.IsRead() {
			read = append(read, n)
		}
	}
	return read, nil
}

// UnreadCount returns the number of unread notifications.
func (c *DatabaseChannel) UnreadCount(notifiable Notifiable) (int, error) {
	unread, err := c.GetUnreadNotifications(notifiable)
	if err != nil {
		return 0, err
	}
	return len(unread), nil
}

// DeleteNotification removes a specific notification.
func (c *DatabaseChannel) DeleteNotification(notifiable Notifiable, notificationID string) error {
	notifiableType, notifiableID := getNotifiableKey(notifiable)

	d, err := c.db.Dialect()
	if err != nil {
		return err
	}

	builder := query.NewBuilder(c.db.RawDB(), d)
	builder.Table("notifications")
	builder.Where("id", "=", notificationID)
	builder.Where("notifiable_type", "=", notifiableType)
	builder.Where("notifiable_id", "=", notifiableID)

	_, err = builder.Delete()
	return err
}

// DeleteAllNotifications removes all notifications for a notifiable.
func (c *DatabaseChannel) DeleteAllNotifications(notifiable Notifiable) error {
	notifiableType, notifiableID := getNotifiableKey(notifiable)

	d, err := c.db.Dialect()
	if err != nil {
		return err
	}

	builder := query.NewBuilder(c.db.RawDB(), d)
	builder.Table("notifications")
	builder.Where("notifiable_type", "=", notifiableType)
	builder.Where("notifiable_id", "=", notifiableID)

	_, err = builder.Delete()
	return err
}

// getNotifiableKey is a helper to extract the polymorphic type/id of a notifiable.
func getNotifiableKey(notifiable Notifiable) (string, string) {
	id := notifiable.RouteNotificationFor("database")
	typeName := fmt.Sprintf("%T", notifiable)
	return typeName, id
}

