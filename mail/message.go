package mail

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Mailable represents an email message that can be sent by the mailer.
type Mailable interface {
	Build() *Message
}

type Attachment struct {
	Name        string
	Data        []byte
	ContentType string
	Inline      bool
}

// Message contains the payload data for the email.
type Message struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	ReplyTo     string
	Subject     string
	Text        string
	HTML        string
	Attachments []Attachment
	Headers     map[string]string
}

// BaseMailable provides a convenient builder pattern for constructing messages.
type BaseMailable struct {
	message *Message
}

func NewMailable() *BaseMailable {
	return &BaseMailable{
		message: &Message{
			Headers: make(map[string]string),
		},
	}
}

func (m *BaseMailable) From(from string) *BaseMailable {
	m.message.From = from
	return m
}

func (m *BaseMailable) To(to string) *BaseMailable {
	m.message.To = append(m.message.To, to)
	return m
}

func (m *BaseMailable) Cc(cc string) *BaseMailable {
	m.message.Cc = append(m.message.Cc, cc)
	return m
}

func (m *BaseMailable) Bcc(bcc string) *BaseMailable {
	m.message.Bcc = append(m.message.Bcc, bcc)
	return m
}

func (m *BaseMailable) ReplyTo(replyTo string) *BaseMailable {
	m.message.ReplyTo = replyTo
	return m
}

func (m *BaseMailable) Attach(name string, data []byte) *BaseMailable {
	m.message.Attachments = append(m.message.Attachments, Attachment{
		Name:        name,
		Data:        data,
		ContentType: detectContentType(name),
	})
	return m
}

func (m *BaseMailable) AttachFile(path string) *BaseMailable {
	data, err := os.ReadFile(path)
	if err != nil {
		return m
	}
	name := filepath.Base(path)
	m.message.Attachments = append(m.message.Attachments, Attachment{
		Name:        name,
		Data:        data,
		ContentType: detectContentType(name),
	})
	return m
}

func (m *BaseMailable) AttachRaw(name string, data []byte, contentType string) *BaseMailable {
	m.message.Attachments = append(m.message.Attachments, Attachment{
		Name:        name,
		Data:        data,
		ContentType: contentType,
	})
	return m
}

func (m *BaseMailable) Header(key, value string) *BaseMailable {
	if m.message.Headers == nil {
		m.message.Headers = make(map[string]string)
	}
	m.message.Headers[key] = value
	return m
}

func (m *BaseMailable) Subject(subject string) *BaseMailable {
	m.message.Subject = subject
	return m
}

func (m *BaseMailable) Text(text string) *BaseMailable {
	m.message.Text = text
	return m
}

func (m *BaseMailable) HTML(html string) *BaseMailable {
	m.message.HTML = html
	return m
}

// Markdown sets the email content from Markdown.
// It automatically renders to HTML and keeps the original as Text fallback.
func (m *BaseMailable) Markdown(md string) *BaseMailable {
	htmlContent, err := RenderMarkdown(md)
	if err != nil {
		// Fallback to raw markdown if rendering fails
		m.message.HTML = md
		m.message.Text = md
		return m
	}

	m.message.HTML = htmlContent
	m.message.Text = StripForText(htmlContent) // or keep original md if preferred
	return m
}

func (m *BaseMailable) Build() *Message {
	return m.message
}

// String returns a simple text representation of the message for debugging.
func (m *Message) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\n", m.From))
	sb.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ", ")))
	if len(m.Cc) > 0 {
		sb.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.Cc, ", ")))
	}
	if len(m.Bcc) > 0 {
		sb.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.Bcc, ", ")))
	}
	if m.ReplyTo != "" {
		sb.WriteString(fmt.Sprintf("Reply-To: %s\n", m.ReplyTo))
	}
	sb.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	if m.Text != "" {
		sb.WriteString(fmt.Sprintf("\n%s", m.Text))
	}
	if len(m.Attachments) > 0 {
		sb.WriteString(fmt.Sprintf("\n[%d attachment(s)]", len(m.Attachments)))
	}
	return sb.String()
}

// AttachmentsBase64 returns all attachments encoded as base64 strings (for SMTP).
func (m *Message) AttachmentsBase64() []string {
	result := make([]string, len(m.Attachments))
	for i, att := range m.Attachments {
		result[i] = base64.StdEncoding.EncodeToString(att.Data)
	}
	return result
}

func detectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html":
		return "text/html"
	case ".json":
		return "application/json"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

