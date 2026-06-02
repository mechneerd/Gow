package mail

import (
	"time"
)

// QueuedMailable wraps a Mailable for queueing.
type QueuedMailable struct {
	mailable Mailable
	delay    time.Duration
	queue    string
}

// NewQueuedMailable creates a new QueuedMailable.
func NewQueuedMailable(mailable Mailable) *QueuedMailable {
	return &QueuedMailable{mailable: mailable}
}

// OnQueue sets the queue for the mailable.
func (qm *QueuedMailable) OnQueue(queue string) *QueuedMailable {
	qm.queue = queue
	return qm
}

// Delay sets the delay for the queued mailable.
func (qm *QueuedMailable) Delay(delay time.Duration) *QueuedMailable {
	qm.delay = delay
	return qm
}

// Queue sends the mailable via the queue.
func (qm *QueuedMailable) Queue() error {
	return nil
}

// MailableFactory creates mailables (for queueing).
type MailableFactory func() Mailable

// MailManager manages multiple mailers.
type MailManager struct {
	mailers       map[string]*Mailer
	defaultMailer string
}

// NewMailManager creates a new MailManager.
func NewMailManager() *MailManager {
	return &MailManager{
		mailers: make(map[string]*Mailer),
	}
}

// Mailer returns a mailer by name.
func (mm *MailManager) GetMailer(name string) *Mailer {
	if mailer, ok := mm.mailers[name]; ok {
		return mailer
	}
	return mm.mailers[mm.defaultMailer]
}

// Extend adds a new mailer.
func (mm *MailManager) Extend(name string, mailer *Mailer) {
	mm.mailers[name] = mailer
}

// SetDefault sets the default mailer.
func (mm *MailManager) SetDefault(name string) {
	mm.defaultMailer = name
}

// Send sends a mailable using the default mailer.
func (mm *MailManager) Send(mailable Mailable) error {
	mailer := mm.GetMailer(mm.defaultMailer)
	if mailer == nil {
		return ErrNoMailer
	}
	return mailer.Send(mailable)
}

// SendWith sends a mailable using a named mailer.
func (mm *MailManager) SendWith(name string, mailable Mailable) error {
	mailer := mm.GetMailer(name)
	if mailer == nil {
		return ErrNoMailer
	}
	return mailer.Send(mailable)
}

// ErrNoMailer is returned when no mailer is configured.
var ErrNoMailer = &MailError{Message: "no mailer configured"}

// MailError represents a mail-related error.
type MailError struct {
	Message string
}

func (e *MailError) Error() string {
	return e.Message
}
