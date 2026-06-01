package mail

import (
	"sync"
)

// MailFake is a fake mailer for testing purposes.
// It captures all sent messages for assertion.
type MailFake struct {
	mu       sync.RWMutex
	sent     []*Message
	sentTo   map[string]bool
	queuing  bool
}

// NewMailFake creates a new fake mailer for testing.
func NewMailFake() *MailFake {
	return &MailFake{
		sent:    make([]*Message, 0),
		sentTo:  make(map[string]bool),
		queuing: false,
	}
}

// Send captures the message instead of sending it.
func (f *MailFake) Send(mailable Mailable) error {
	msg := mailable.Build()
	f.mu.Lock()
	defer f.mu.Unlock()

	f.sent = append(f.sent, msg)
	for _, to := range msg.To {
		f.sentTo[to] = true
	}
	return nil
}

// Queue captures the message (pretends to queue).
func (f *MailFake) Queue(mailable Mailable) error {
	return f.Send(mailable)
}

// QueueNow captures the message (pretends to queue immediately).
func (f *MailFake) QueueNow(mailable Mailable) error {
	return f.Send(mailable)
}

// Sent returns all messages that were sent.
func (f *MailFake) Sent() []*Message {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make([]*Message, len(f.sent))
	copy(result, f.sent)
	return result
}

// SentTo checks if a message was sent to the given address.
func (f *MailFake) SentTo(email string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.sentTo[email]
}

// SentCount returns the number of messages sent.
func (f *MailFake) SentCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return len(f.sent)
}

// Last returns the last message that was sent, or nil if no messages.
func (f *MailFake) Last() *Message {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if len(f.sent) == 0 {
		return nil
	}
	return f.sent[len(f.sent)-1]
}

// HasSent checks if any message was sent.
func (f *MailFake) HasSent() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return len(f.sent) > 0
}

// AssertSent asserts that at least one message was sent.
func (f *MailFake) AssertSent(mailable Mailable) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	expected := mailable.Build()
	for _, sent := range f.sent {
		if sent.Subject == expected.Subject && len(sent.To) == len(expected.To) {
			return true
		}
	}
	return false
}

// AssertSentTo asserts that a message was sent to the given address.
func (f *MailFake) AssertSentTo(email string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.sentTo[email]
}

// AssertNotSent asserts that no messages were sent.
func (f *MailFake) AssertNotSent() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return len(f.sent) == 0
}

// AssertSentCount asserts the exact number of messages sent.
func (f *MailFake) AssertSentCount(count int) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return len(f.sent) == count
}

// Clear resets the fake mailer, removing all captured messages.
func (f *MailFake) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.sent = make([]*Message, 0)
	f.sentTo = make(map[string]bool)
}

// Fake replaces the mailer's driver with a MailFake for testing.
// Returns a MailFake that can be used to assert sent messages.
func Fake(mailer *Mailer) *MailFake {
	fake := NewMailFake()
	mailer.mu.Lock()
	defer mailer.mu.Unlock()
	mailer.fakes = append(mailer.fakes, fake)
	return fake
}
