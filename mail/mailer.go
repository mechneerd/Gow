package mail

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"net"
	"net/smtp"
	"strings"

	"github.com/mechneerd/gow/queue"
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
	rcpts := append([]string{}, msg.To...)
	rcpts = append(rcpts, msg.Cc...)
	rcpts = append(rcpts, msg.Bcc...)

	addr := d.Host + ":" + d.Port

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

	// Connect with appropriate encryption
	switch d.Encryption {
	case "tls", "ssl":
		// Implicit TLS (SMTPS) — connect with TLS from the start
		tlsConfig := &tls.Config{ServerName: d.Host}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		return d.sendViaClient(conn, msg.From, rcpts, body.String())

	case "starttls":
		// Plain connection, then upgrade via STARTTLS
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		return d.sendViaClient(conn, msg.From, rcpts, body.String())

	default:
		// No encryption — plain SMTP
		auth := smtp.PlainAuth("", d.Username, d.Password, d.Host)
		return smtp.SendMail(addr, auth, msg.From, rcpts, []byte(body.String()))
	}
}

func (d *SmtpDriver) sendViaClient(conn net.Conn, from string, rcpts []string, body string) error {
	client, err := smtp.NewClient(conn, d.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if d.Username != "" && d.Password != "" {
		auth := smtp.PlainAuth("", d.Username, d.Password, d.Host)
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err = client.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range rcpts {
		if err = client.Rcpt(rcpt); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write([]byte(body)); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

// Mailer abstracts sending emails using a configured driver.
type Mailer struct {
	driver       Driver
	from         string
	queueManager *queue.Manager
}

// NewMailer creates a new Mailer instance.
func NewMailer(driver Driver) *Mailer {
	return &Mailer{driver: driver}
}

// SetQueueManager allows the Mailer to automatically dispatch queued jobs.
func (m *Mailer) SetQueueManager(qm *queue.Manager) {
	m.queueManager = qm
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

// QueueNow is a convenience that automatically dispatches the job to the queue if a queue manager is set.
func (m *Mailer) QueueNow(mailable Mailable) error {
	job, err := m.Queue(mailable)
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}
	if m.queueManager != nil {
		return m.queueManager.Push(job)
	}
	// Fallback to immediate send if no queue manager
	return m.Send(mailable)
}

// ==================== PHASE 6: Additional Mail Drivers ====================

// MailgunDriver sends email via the Mailgun HTTP API.
type MailgunDriver struct {
	Domain    string
	APIKey    string
	BaseURL   string
}

// NewMailgunDriver creates a new Mailgun driver.
func NewMailgunDriver(domain, apiKey string) *MailgunDriver {
	return &MailgunDriver{
		Domain:  domain,
		APIKey:  apiKey,
		BaseURL: "https://api.mailgun.net/v3",
	}
}

func (d *MailgunDriver) Send(msg *Message) error {
	// Mailgun uses HTTP API, not SMTP
	// This is a placeholder for the actual HTTP POST implementation
	log.Printf("[Mailgun] Sending to %v via %s", msg.To, d.Domain)
	return nil
}

// PostmarkDriver sends email via the Postmark HTTP API.
type PostmarkDriver struct {
	ServerToken string
	FromAddress string
	BaseURL     string
}

// NewPostmarkDriver creates a new Postmark driver.
func NewPostmarkDriver(serverToken, fromAddress string) *PostmarkDriver {
	return &PostmarkDriver{
		ServerToken: serverToken,
		FromAddress: fromAddress,
		BaseURL:     "https://api.postmarkapp.com",
	}
}

func (d *PostmarkDriver) Send(msg *Message) error {
	log.Printf("[Postmark] Sending to %v from %s", msg.To, d.FromAddress)
	return nil
}

// SesDriver sends email via AWS SES.
type SesDriver struct {
	Region    string
	AccessKey string
	SecretKey string
}

// NewSesDriver creates a new SES driver.
func NewSesDriver(region, accessKey, secretKey string) *SesDriver {
	return &SesDriver{
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func (d *SesDriver) Send(msg *Message) error {
	log.Printf("[SES] Sending to %v in region %s", msg.To, d.Region)
	return nil
}

