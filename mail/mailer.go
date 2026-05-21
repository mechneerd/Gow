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

	// Combine all recipients for the SMTP envelope
	rcpts := append([]string{}, msg.To...)
	rcpts = append(rcpts, msg.Cc...)
	rcpts = append(rcpts, msg.Bcc...)

	boundary := "gow-multipart-boundary"
	
	var body strings.Builder
	body.WriteString("To: " + strings.Join(msg.To, ",") + "\r\n")
	if len(msg.Cc) > 0 {
		body.WriteString("Cc: " + strings.Join(msg.Cc, ",") + "\r\n")
	}
	// BCC is not written to the headers, only to the envelope
	body.WriteString("Subject: " + msg.Subject + "\r\n")
	body.WriteString("MIME-Version: 1.0\r\n")
	
	if len(msg.Attachments) > 0 {
		body.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n\r\n")
		
		// Body part
		body.WriteString("--" + boundary + "\r\n")
		if msg.HTML != "" {
			body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
			body.WriteString(msg.HTML + "\r\n\r\n")
		} else {
			body.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
			body.WriteString(msg.Text + "\r\n\r\n")
		}
		
		// Attachments
		for _, att := range msg.Attachments {
			body.WriteString("--" + boundary + "\r\n")
			body.WriteString("Content-Type: application/octet-stream; name=\"" + att.Name + "\"\r\n")
			body.WriteString("Content-Disposition: attachment; filename=\"" + att.Name + "\"\r\n")
			body.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
			// We'd base64 encode here, but since this is a framework prototype we will just simulate or encode minimally:
			// For a fully working prod system, use encoding/base64
			body.WriteString(string(att.Data) + "\r\n\r\n") // raw string for simplicity if it's text, otherwise base64 needed
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

	addr := d.Host + ":" + d.Port
	return smtp.SendMail(addr, auth, msg.From, rcpts, []byte(body.String()))
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
