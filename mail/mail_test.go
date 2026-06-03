package mail

import (
	"testing"
)

func TestNewMailable(t *testing.T) {
	m := NewMailable()
	if m == nil {
		t.Fatal("NewMailable returned nil")
	}
}

func TestMailable_From(t *testing.T) {
	m := NewMailable().From("sender@example.com")
	if m.message.From != "sender@example.com" {
		t.Errorf("expected From 'sender@example.com', got %q", m.message.From)
	}
}

func TestMailable_To(t *testing.T) {
	m := NewMailable().To("recipient@example.com")
	if len(m.message.To) != 1 || m.message.To[0] != "recipient@example.com" {
		t.Errorf("expected To ['recipient@example.com'], got %v", m.message.To)
	}
}

func TestMailable_MultipleTo(t *testing.T) {
	m := NewMailable().
		To("one@example.com").
		To("two@example.com")

	if len(m.message.To) != 2 {
		t.Errorf("expected 2 To recipients, got %d", len(m.message.To))
	}
}

func TestMailable_Cc(t *testing.T) {
	m := NewMailable().Cc("cc@example.com")
	if len(m.message.Cc) != 1 || m.message.Cc[0] != "cc@example.com" {
		t.Errorf("expected Cc ['cc@example.com'], got %v", m.message.Cc)
	}
}

func TestMailable_Bcc(t *testing.T) {
	m := NewMailable().Bcc("bcc@example.com")
	if len(m.message.Bcc) != 1 || m.message.Bcc[0] != "bcc@example.com" {
		t.Errorf("expected Bcc ['bcc@example.com'], got %v", m.message.Bcc)
	}
}

func TestMailable_Subject(t *testing.T) {
	m := NewMailable().Subject("Test Subject")
	if m.message.Subject != "Test Subject" {
		t.Errorf("expected Subject 'Test Subject', got %q", m.message.Subject)
	}
}

func TestMailable_Text(t *testing.T) {
	m := NewMailable().Text("Plain text body")
	if m.message.Text != "Plain text body" {
		t.Errorf("expected Text 'Plain text body', got %q", m.message.Text)
	}
}

func TestMailable_HTML(t *testing.T) {
	m := NewMailable().HTML("<h1>Hello</h1>")
	if m.message.HTML != "<h1>Hello</h1>" {
		t.Errorf("expected HTML '<h1>Hello</h1>', got %q", m.message.HTML)
	}
}

func TestMailable_ReplyTo(t *testing.T) {
	m := NewMailable().ReplyTo("reply@example.com")
	if m.message.ReplyTo != "reply@example.com" {
		t.Errorf("expected ReplyTo 'reply@example.com', got %q", m.message.ReplyTo)
	}
}

func TestMailable_Header(t *testing.T) {
	m := NewMailable().Header("X-Custom", "value")
	if m.message.Headers["X-Custom"] != "value" {
		t.Errorf("expected header X-Custom 'value', got %q", m.message.Headers["X-Custom"])
	}
}

func TestMailable_Build(t *testing.T) {
	type WelcomeMail struct {
		BaseMailable
		Name string
	}

	m := &WelcomeMail{
		BaseMailable: *NewMailable(),
		Name:         "John",
	}

	m.From("sender@example.com").
		To("john@example.com").
		Subject("Welcome").
		HTML("<h1>Welcome " + m.Name + "</h1>")

	msg := m.BaseMailable.Build()

	if msg.From != "sender@example.com" {
		t.Errorf("expected From 'sender@example.com', got %q", msg.From)
	}
	if len(msg.To) != 1 || msg.To[0] != "john@example.com" {
		t.Errorf("expected To ['john@example.com'], got %v", msg.To)
	}
	if msg.Subject != "Welcome" {
		t.Errorf("expected Subject 'Welcome', got %q", msg.Subject)
	}
}

func TestMailManager_New(t *testing.T) {
	mm := NewMailManager()
	if mm == nil {
		t.Fatal("NewMailManager returned nil")
	}
}

func TestMailManager_SendWithoutMailer(t *testing.T) {
	mm := NewMailManager()
	mm.SetDefault("smtp")

	err := mm.Send(NewMailable().From("test@example.com").To("recv@example.com").Subject("Test"))
	if err == nil {
		t.Error("expected error when no mailer configured")
	}
}

func TestMailManager_Extend(t *testing.T) {
	mm := NewMailManager()
	mailer := NewMailer(&LogDriver{})
	mm.Extend("log", mailer)
	mm.SetDefault("log")

	// Should not panic
	mm.Send(NewMailable().From("test@example.com").To("recv@example.com").Subject("Test"))
}

func TestQueuedMailable(t *testing.T) {
	m := NewMailable().From("test@example.com").To("recv@example.com")
	qm := NewQueuedMailable(m)

	if qm.mailable == nil {
		t.Error("expected non-nil mailable")
	}

	// Test chaining
	qm.OnQueue("emails").Delay(0)
	if qm.queue != "emails" {
		t.Errorf("expected queue 'emails', got %q", qm.queue)
	}
}

func TestMailableInterface(t *testing.T) {
	type TestMail struct {
		BaseMailable
	}

	m := &TestMail{BaseMailable: *NewMailable()}
	var _ Mailable = m
}
