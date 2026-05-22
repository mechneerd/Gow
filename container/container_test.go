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

type Filesystem interface {
	Disk() string
}

type LocalFilesystem struct{}

func (l *LocalFilesystem) Disk() string { return "local" }

type S3Filesystem struct{}

func (s *S3Filesystem) Disk() string { return "s3" }

type PhotoController struct {
	FS Filesystem
}

type VideoController struct {
	FS Filesystem
}

func TestContextualBinding(t *testing.T) {
	c := New()

	c.Bind((*PhotoController)(nil), func(fs Filesystem) *PhotoController {
		return &PhotoController{FS: fs}
	})

	c.Bind((*VideoController)(nil), func(fs Filesystem) *VideoController {
		return &VideoController{FS: fs}
	})

	c.When((*PhotoController)(nil)).
		Needs((*Filesystem)(nil)).
		Give(func() Filesystem { return &LocalFilesystem{} })

	c.When((*VideoController)(nil)).
		Needs((*Filesystem)(nil)).
		Give(func() Filesystem { return &S3Filesystem{} })

	pc, err := Make[*PhotoController](c)
	if err != nil {
		t.Fatalf("Failed to resolve PhotoController: %v", err)
	}
	if pc.FS.Disk() != "local" {
		t.Errorf("Expected local disk, got %s", pc.FS.Disk())
	}

	vc, err := Make[*VideoController](c)
	if err != nil {
		t.Fatalf("Failed to resolve VideoController: %v", err)
	}
	if vc.FS.Disk() != "s3" {
		t.Errorf("Expected s3 disk, got %s", vc.FS.Disk())
	}
}

func TestContainerInstances(t *testing.T) {
	c := New()
	fs := &LocalFilesystem{}

	err := c.Instance((*Filesystem)(nil), fs)
	if err != nil {
		t.Fatalf("Failed to register instance: %v", err)
	}

	resolved, err := Make[Filesystem](c)
	if err != nil {
		t.Fatalf("Failed to resolve instance: %v", err)
	}

	if resolved != fs {
		t.Error("Expected the exact instance to be resolved")
	}
}

type AutoWiredController struct {
	DB Database `inject:""`
}

func TestStructInjection(t *testing.T) {
	c := New()

	c.Instance((*Database)(nil), &MySQLDatabase{ConnectionString: "injected"})

	ctrl, err := Make[*AutoWiredController](c)
	if err != nil {
		t.Fatalf("Failed to resolve struct: %v", err)
	}

	if ctrl.DB == nil {
		t.Fatal("Expected DB to be injected")
	}

	if ctrl.DB.Connect() != "Connected to injected" {
		t.Errorf("Unexpected result: %s", ctrl.DB.Connect())
	}
}
