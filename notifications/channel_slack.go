package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackChannel sends notifications to Slack.
type SlackChannel struct {
	WebhookURL string
}

func NewSlackChannel(webhookURL string) *SlackChannel {
	return &SlackChannel{WebhookURL: webhookURL}
}

func (s *SlackChannel) Send(notifiable any, notification Notification) error {
	message := map[string]string{
		"text": fmt.Sprintf("Notification for %v: %T", notifiable, notification),
	}

	body, _ := json.Marshal(message)
	http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(body))
	return nil
}
