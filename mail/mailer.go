package mail

import (
	"log"
	"net/smtp"
	"strings"
)

// Driver handles the actual transport of the email.
type Driver interface {
	Send(msg *Message) error
}

// LogDriver outputs the email to standard log for local development.
type LogDriver struct{}

func (d *LogDriver) Send(msg *Message) error {
	log.Printf("--- EMAIL ---\nFrom: %s\nTo: %v\nSubject: %s\nBody: %s\n-------------\n",
		msg.From, msg.To, msg.Subject, msg.Text)
	return nil
}

// SmtpDriver connects to an SMTP server to send email.
type SmtpDriver struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (d *SmtpDriver) Send(msg *Message) error {
	auth := smtp.PlainAuth("", d.Username, d.Password, d.Host)

	body := "To: " + strings.Join(msg.To, ",") + "\r\n" +
		"Subject: " + msg.Subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		msg.HTML

	if msg.HTML == "" {
		body = "To: " + strings.Join(msg.To, ",") + "\r\n" +
			"Subject: " + msg.Subject + "\r\n" +
			"\r\n" + msg.Text
	}

	addr := d.Host + ":" + d.Port
	return smtp.SendMail(addr, auth, msg.From, msg.To, []byte(body))
}

// Mailer abstracts sending emails using a configured driver.
type Mailer struct {
	driver Driver
}

// NewMailer creates a new Mailer instance.
func NewMailer(driver Driver) *Mailer {
	return &Mailer{driver: driver}
}

// Send dispatches a mailable using the active driver.
func (m *Mailer) Send(mailable Mailable) error {
	msg := mailable.Build()
	return m.driver.Send(msg)
}
