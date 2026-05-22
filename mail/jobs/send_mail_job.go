package mail

import (
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
	// Could log or notify about failed email
	// For now, just a placeholder
}
