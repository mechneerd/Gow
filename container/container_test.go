package container

import (
	"errors"
	"testing"
)

type Mailer interface {
	Send(msg string) string
}

type SmtpMailer struct {
	Host string
}

func (s *SmtpMailer) Send(msg string) string {
	return "Sent " + msg + " via " + s.Host
}

func TestContainerBindAndMake(t *testing.T) {
	app := New()

	app.Bind((*Mailer)(nil), func() Mailer {
		return &SmtpMailer{Host: "smtp.example.com"}
	})

	mailer, err := Make[Mailer](app)
	if err != nil {
		t.Fatalf("Failed to make Mailer: %v", err)
	}

	result := mailer.Send("Hello")
	if result != "Sent Hello via smtp.example.com" {
		t.Errorf("Unexpected result: %s", result)
	}

	// Verify it's not a singleton (each Make calls factory again)
	mailer2, _ := Make[Mailer](app)
	if mailer == mailer2 {
		t.Errorf("Bind should return a new instance each time, got same pointer")
	}
}

func TestContainerSingleton(t *testing.T) {
	app := New()

	app.Singleton((*Mailer)(nil), func() Mailer {
		return &SmtpMailer{Host: "smtp.example.com"}
	})

	mailer1, _ := Make[Mailer](app)
	mailer2, _ := Make[Mailer](app)

	if mailer1 != mailer2 {
		t.Errorf("Singleton should return the exact same instance, got different pointers")
	}
}

func TestContainerInstance(t *testing.T) {
	app := New()
	existing := &SmtpMailer{Host: "already.existing"}

	err := app.Instance((*Mailer)(nil), existing)
	if err != nil {
		t.Fatalf("Instance failed: %v", err)
	}

	resolved, _ := Make[Mailer](app)
	if resolved != existing {
		t.Errorf("Expected existing instance to be resolved")
	}
}

func TestContainerFreeze(t *testing.T) {
	app := New()
	app.Freeze()

	err := app.Bind((*Mailer)(nil), func() Mailer { return &SmtpMailer{} })
	if !errors.Is(err, ErrFrozen) {
		t.Errorf("Expected ErrFrozen, got %v", err)
	}
}

func TestContainerUnbound(t *testing.T) {
	app := New()

	_, err := Make[Mailer](app)
	if !errors.Is(err, ErrBindingNotFound) {
		t.Errorf("Expected ErrBindingNotFound, got %v", err)
	}
}

type DependentService struct {
	Mailer Mailer
}

func TestContainerDependencyInjection(t *testing.T) {
	app := New()

	app.Bind((*Mailer)(nil), func() Mailer {
		return &SmtpMailer{Host: "smtp.example.com"}
	})

	// Factory that takes dependencies automatically resolved by the container
	app.Bind((*DependentService)(nil), func(m Mailer) *DependentService {
		return &DependentService{Mailer: m}
	})

	service, err := Make[*DependentService](app)
	if err != nil {
		t.Fatalf("Failed to make DependentService: %v", err)
	}

	if service.Mailer == nil {
		t.Errorf("Expected Mailer dependency to be injected, got nil")
	}
}
