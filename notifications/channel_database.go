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

// getNotifiableKey is a helper to extract the polymorphic type/id of a notifiable.
func getNotifiableKey(notifiable Notifiable) (string, string) {
	// In a complete framework, this would use reflection or an interface method
	// to extract the underlying Model's table name and primary key.
	// We'll use the RouteNotificationFor method with a "database" channel hint for simplicity.
	id := notifiable.RouteNotificationFor("database")
	typeName := fmt.Sprintf("%T", notifiable)
	return typeName, id
}

