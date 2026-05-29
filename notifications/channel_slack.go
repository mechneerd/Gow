package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("slack webhook returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

