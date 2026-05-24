package fakes

import (
	"reflect"
	"testing"

	"gow/mail"
	"gow/queue"

	"github.com/stretchr/testify/assert"
)

// MailFake tracks sent emails for testing.
type MailFake struct {
	sent []mail.Mailable
}

func (f *MailFake) Send(mailable mail.Mailable) error {
	f.sent = append(f.sent, mailable)
	return nil
}

func (f *MailFake) AssertSent(t *testing.T, mailableType any, count int) {
	actualCount := 0
	targetType := reflect.TypeOf(mailableType)
	
	for _, m := range f.sent {
		if reflect.TypeOf(m) == targetType {
			actualCount++
		}
	}
	
	assert.Equal(t, count, actualCount, "Expected %d mails of type %v to be sent, found %d", count, targetType, actualCount)
}

// QueueFake tracks jobs pushed to the queue.
type QueueFake struct {
	pushed []queue.Job
}

func (f *QueueFake) Push(job queue.Job) error {
	f.pushed = append(f.pushed, job)
	return nil
}

func (f *QueueFake) Pop() (queue.Job, error) {
	if len(f.pushed) == 0 {
		return nil, nil
	}
	job := f.pushed[0]
	f.pushed = f.pushed[1:]
	return job, nil
}

func (f *QueueFake) AssertPushed(t *testing.T, jobType any, count int) {
	actualCount := 0
	targetType := reflect.TypeOf(jobType)
	
	for _, j := range f.pushed {
		if reflect.TypeOf(j) == targetType {
			actualCount++
		}
	}
	
	assert.Equal(t, count, actualCount, "Expected %d jobs of type %v to be pushed, found %d", count, targetType, actualCount)
}

// EventFake tracks dispatched events.
type EventFake struct {
	dispatched []any
}

func (f *EventFake) Dispatch(event any) {
	f.dispatched = append(f.dispatched, event)
}

func (f *EventFake) AssertDispatched(t *testing.T, eventType any, count int) {
	actualCount := 0
	targetType := reflect.TypeOf(eventType)
	
	for _, e := range f.dispatched {
		if reflect.TypeOf(e) == targetType {
			actualCount++
		}
	}
	
	assert.Equal(t, count, actualCount, "Expected %d events of type %v to be dispatched, found %d", count, targetType, actualCount)
}
