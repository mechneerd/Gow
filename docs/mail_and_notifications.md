# Mail & Notifications

GoW provides a clean, simple API over the popular `gomail` library for sending emails, as well as a robust notification system for sending short messages across various channels.

## Mail

You can configure your SMTP credentials in `config/mail.go`.

### Sending Mail

```go
mailer := mail.NewMailer()

message := mail.NewMessage().
    To("recipient@example.com").
    Subject("Welcome to GoW").
    Body("text/html", "<h1>Hello World!</h1>")

// Add attachments or CC/BCC
message.Attach("/path/to/file.pdf")
message.Cc("cc@example.com")

err := mailer.Send(message)
```

## Notifications

In addition to email, GoW provides support for sending notifications across a variety of delivery channels, including mail, SMS (via Nexmo/Twilio), Slack, and storing them in a database so they may be displayed in your web interface.

### Defining a Notification

A notification is represented by a struct that implements the `Notification` interface.

```go
type InvoicePaid struct {
    InvoiceID int
}

func (n *InvoicePaid) Via() []string {
    return []string{"mail", "database"}
}

func (n *InvoicePaid) ToMail() *mail.Message {
    return mail.NewMessage().
        Subject("Invoice Paid").
        Body("text/plain", "Your invoice has been paid!")
}

func (n *InvoicePaid) ToDatabase() map[string]any {
    return map[string]any{
        "invoice_id": n.InvoiceID,
    }
}
```

### Sending Notifications

You can send notifications using the `NotificationManager`.

```go
notifications.Send(user, &InvoicePaid{InvoiceID: 123})
```
