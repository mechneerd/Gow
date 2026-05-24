package jobs

import "fmt"

// SendMailJob is a job for sending emails via queue.
// For full mail integration, wire this with your Mailable + Mailer in the app.
type SendMailJob struct {
	// Extend this struct with Mailable data as needed in your application.
}

func (j *SendMailJob) Handle() error {
	fmt.Println("[Mail] SendMailJob executed")
	return nil
}

func (j *SendMailJob) Failed(err error) {
	fmt.Printf("[Mail] SendMailJob failed: %v\n", err)
}
