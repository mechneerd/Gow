package mail

// Mailable represents an email message that can be sent by the mailer.
type Mailable interface {
	Build() *Message
}

type Attachment struct {
	Name string
	Data []byte
}

// Message contains the payload data for the email.
type Message struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Text        string
	HTML        string
	Attachments []Attachment
}

// BaseMailable provides a convenient builder pattern for constructing messages.
type BaseMailable struct {
	message *Message
}

func NewMailable() *BaseMailable {
	return &BaseMailable{
		message: &Message{},
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

func (m *BaseMailable) Attach(name string, data []byte) *BaseMailable {
	m.message.Attachments = append(m.message.Attachments, Attachment{Name: name, Data: data})
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
