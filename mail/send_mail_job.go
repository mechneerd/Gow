package mail

// SendMailJob is the queue job for sending emails when a Mailable implements ShouldQueue.
type SendMailJob struct {
	Mailable Mailable
	Mailer   *Mailer
}

func (j *SendMailJob) Handle() error {
	if j.Mailer != nil {
		return j.Mailer.Send(j.Mailable)
	}
	return (&LogDriver{}).Send(j.Mailable.Build())
}

func (j *SendMailJob) Failed(err error) {
	// Consumers can override or log via their preferred logger.
}

