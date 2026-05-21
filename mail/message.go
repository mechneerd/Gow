package mail

// Mailable represents an email message that can be sent by the mailer.
type Mailable interface {
	Build() *Message
}

// Message contains the payload data for the email.
type Message struct {
	From    string
	To      []string
	Subject string
	Text    string
	HTML    string
	// Attachments could be added here
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

func (m *BaseMailable) Build() *Message {
	return m.message
}
