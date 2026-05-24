package jobs

import "fmt"

// SendMailJob is a placeholder job for sending emails via queue.
// The real implementation lives in the mail package.
type SendMailJob struct {
	// TODO: Wire with actual Mailable + Mailer from gow/mail
}

func (j *SendMailJob) Handle() error {
	fmt.Println("[Mail] SendMailJob executed (stub)")
	return nil
}

func (j *SendMailJob) Failed(err error) {
	fmt.Printf("[Mail] SendMailJob failed: %v\n", err)
}
