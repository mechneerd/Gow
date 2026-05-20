package container

import (
	"errors"
	"testing"
)

type Database interface {
	Connect() string
}

type MySQLDatabase struct {
	ConnectionString string
}

func (m *MySQLDatabase) Connect() string {
	return "Connected to " + m.ConnectionString
}

type UserService struct {
	DB Database
}

func NewUserService(db Database) *UserService {
	return &UserService{DB: db}
}

func TestContainerBindingAndResolution(t *testing.T) {
	c := New()

	// Bind Database interface to MySQLDatabase factory
	err := c.Bind((*Database)(nil), func() Database {
		return &MySQLDatabase{ConnectionString: "mysql://localhost"}
	})
	if err != nil {
		t.Fatalf("Failed to bind: %v", err)
	}

	// Resolve the interface
	db, err := Make[Database](c)
	if err != nil {
		t.Fatalf("Failed to resolve: %v", err)
	}

	if db.Connect() != "Connected to mysql://localhost" {
		t.Errorf("Unexpected result: %s", db.Connect())
	}
}

func TestContainerDependencyInjection(t *testing.T) {
	c := New()

	c.Bind((*Database)(nil), func() Database {
		return &MySQLDatabase{ConnectionString: "mysql://localhost"}
	})

	c.Bind((*UserService)(nil), NewUserService)

	svc, err := Make[*UserService](c)
	if err != nil {
		t.Fatalf("Failed to resolve UserService: %v", err)
	}

	if svc.DB == nil {
		t.Fatal("Dependency was not injected")
	}

	if svc.DB.Connect() != "Connected to mysql://localhost" {
		t.Errorf("Unexpected result: %s", svc.DB.Connect())
	}
}

func TestContainerSingleton(t *testing.T) {
	c := New()

	callCount := 0
	c.Singleton((*Database)(nil), func() Database {
		callCount++
		return &MySQLDatabase{ConnectionString: "singleton://"}
	})

	db1, _ := Make[Database](c)
	db2, _ := Make[Database](c)

	if db1 != db2 {
		t.Error("Expected the same instance for singleton")
	}

	if callCount != 1 {
		t.Errorf("Expected factory to be called 1 time, got %d", callCount)
	}
}

func TestContainerFreeze(t *testing.T) {
	c := New()
	c.Freeze()

	err := c.Bind((*Database)(nil), func() Database {
		return &MySQLDatabase{}
	})

	if !errors.Is(err, ErrFrozen) {
		t.Errorf("Expected ErrFrozen, got %v", err)
	}
}
