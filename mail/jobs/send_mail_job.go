package mail

import (
	"fmt"
	"gow/queue"
)

// SendMailJob is the queue job for sending emails.
type SendMailJob struct {
	Mailable Mailable
	Mailer   *Mailer // injected when creating the job
}

func (j *SendMailJob) Handle() error {
	if j.Mailer != nil {
		return j.Mailer.Send(j.Mailable)
	}
	// Fallback: create a new log mailer if no mailer provided
	return (&LogDriver{}).Send(j.Mailable.Build())
}

func (j *SendMailJob) Failed(err error) {
	fmt.Printf("[Mail] SendMailJob failed: %v\n", err)
}
