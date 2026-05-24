package orm

import (
	"database/sql"
	"errors"
	"github.com/mechneerd/gow/database/dialect"
	"github.com/mechneerd/gow/database/query"
	"testing"
)

type EventUser struct {
	ID   int    `db:"id" gow:"primaryKey"`
	Name string `db:"name"`

	beforeSaveCalled   bool
	beforeCreateCalled bool
	afterCreateCalled  bool
}

func (EventUser) TableName() string { return "event_users" }

func (u *EventUser) BeforeSave() error {
	u.beforeSaveCalled = true
	if u.Name == "CancelSave" {
		return errors.New("save cancelled")
	}
	return nil
}

func (u *EventUser) BeforeCreate() error {
	u.beforeCreateCalled = true
	return nil
}

func (u *EventUser) AfterCreate() error {
	u.afterCreateCalled = true
	return nil
}

func TestORMEvents(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Exec(`CREATE TABLE event_users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	db := &DB{
		Conn:    conn,
		Builder: query.NewBuilder(conn, &dialect.SQLiteDialect{}),
	}

	u1 := &EventUser{Name: "Alice"}
	err = NewQuery[EventUser](db).Insert(u1)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if !u1.beforeSaveCalled || !u1.beforeCreateCalled || !u1.afterCreateCalled {
		t.Errorf("Expected all create/save hooks to be called")
	}

	u2 := &EventUser{Name: "CancelSave"}
	err = NewQuery[EventUser](db).Insert(u2)
	if err != ErrCancelled {
		t.Errorf("Expected ErrCancelled, got %v", err)
	}
}

