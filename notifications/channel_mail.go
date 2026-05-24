package notifications

import (
	"fmt"

	"github.com/mechneerd/gow/mail"
)

// MailNotification is implemented by notifications that can be sent via email.
type MailNotification interface {
	ToMail(notifiable Notifiable) *mail.Message
}

// MailChannel sends notifications via the Mail system.
type MailChannel struct {
	mailer *mail.Mailer
}

// NewMailChannel creates a new MailChannel.
func NewMailChannel(mailer *mail.Mailer) *MailChannel {
	return &MailChannel{mailer: mailer}
}

// Send sends the notification as an email.
func (c *MailChannel) Send(notifiable Notifiable, notification Notification) error {
	mailNotif, ok := notification.(MailNotification)
	if !ok {
		return fmt.Errorf("notification does not implement MailNotification")
	}

	message := mailNotif.ToMail(notifiable)
	if message == nil {
		return nil // notification chose not to send an email
	}

	// Use BaseMailable and set the built message
	mailable := mail.NewMailable().
		From(message.From).
		To(message.To[0]) // simplify for now

	if message.Subject != "" {
		mailable.Subject(message.Subject)
	}
	if message.HTML != "" {
		mailable.HTML(message.HTML)
	} else {
		mailable.Text(message.Text)
	}

	return c.mailer.Send(mailable)
}

