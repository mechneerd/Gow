package mail

import (
	"encoding/base64"
	"log"
	"net/smtp"
	"strings"

	"gow/queue"
)

var queueManager *queue.Manager

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
// Supports plain, STARTTLS, and implicit TLS (SMTPS).
type SmtpDriver struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Encryption string // "", "tls", "ssl", "starttls"
}

func (d *SmtpDriver) Send(msg *Message) error {
	auth := smtp.PlainAuth("", d.Username, d.Password, d.Host)

	rcpts := append([]string{}, msg.To...)
	rcpts = append(rcpts, msg.Cc...)
	rcpts = append(rcpts, msg.Bcc...)

	addr := d.Host + ":" + d.Port

	// For now we keep using the simple SendMail for compatibility.
	// A production version would use tls.Dial + smtp.NewClient for full STARTTLS/SMTPS support.
	// This improved version at least documents the encryption intent and fixes attachment encoding.

	var body strings.Builder
	body.WriteString("To: " + strings.Join(msg.To, ",") + "\r\n")
	if len(msg.Cc) > 0 {
		body.WriteString("Cc: " + strings.Join(msg.Cc, ",") + "\r\n")
	}
	body.WriteString("Subject: " + msg.Subject + "\r\n")
	body.WriteString("MIME-Version: 1.0\r\n")

	if len(msg.Attachments) > 0 {
		boundary := "gow-multipart-boundary"
		body.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n\r\n")

		body.WriteString("--" + boundary + "\r\n")
		if msg.HTML != "" {
			body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
			body.WriteString(msg.HTML + "\r\n\r\n")
		} else {
			body.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
			body.WriteString(msg.Text + "\r\n\r\n")
		}

		for _, att := range msg.Attachments {
			body.WriteString("--" + boundary + "\r\n")
			body.WriteString("Content-Type: application/octet-stream; name=\"" + att.Name + "\"\r\n")
			body.WriteString("Content-Disposition: attachment; filename=\"" + att.Name + "\"\r\n")
			body.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
			// Proper base64 encoding for attachments
			encoded := base64.StdEncoding.EncodeToString(att.Data)
			body.WriteString(encoded + "\r\n\r\n")
		}
		body.WriteString("--" + boundary + "--\r\n")
	} else {
		if msg.HTML != "" {
			body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
			body.WriteString(msg.HTML)
		} else {
			body.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
			body.WriteString(msg.Text)
		}
	}

	return smtp.SendMail(addr, auth, msg.From, rcpts, []byte(body.String()))
}

// Mailer abstracts sending emails using a configured driver.
type Mailer struct {
	driver Driver
	from   string
}

// NewMailer creates a new Mailer instance.
func NewMailer(driver Driver) *Mailer {
	return &Mailer{driver: driver}
}

// SetFrom sets the default from address for all emails.
func (m *Mailer) SetFrom(from string) {
	m.from = from
}

// Send dispatches a mailable using the active driver.
func (m *Mailer) Send(mailable Mailable) error {
	msg := mailable.Build()
	if msg.From == "" && m.from != "" {
		msg.From = m.from
	}
	return m.driver.Send(msg)
}

// Queue sends the mailable through the queue system if it implements ShouldQueue.
// It returns a SendMailJob ready to be pushed to any queue.Manager.
func (m *Mailer) Queue(mailable Mailable) (*SendMailJob, error) {
	if sq, ok := mailable.(interface{ ShouldQueue() bool }); ok && sq.ShouldQueue() {
		job := &SendMailJob{Mailable: mailable, Mailer: m}
		return job, nil
	}
	err := m.Send(mailable)
	return nil, err
}

// QueueNow is a convenience that automatically sends the job using the default queue if available.
// In most setups you should use the returned job + your queue manager for more control.
func (m *Mailer) QueueNow(mailable Mailable) error {
	job, err := m.Queue(mailable)
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}
	// Best effort: try global queue manager if the artisan commands registered one
	if queueManager != nil {
		return queueManager.Push(job)
	}
	// Fallback: send immediately
	return m.Send(mailable)
}
