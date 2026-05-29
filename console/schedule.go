package console

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

// Event represents a scheduled task.
type Event struct {
	command     string
	expression  string
	withoutOverlap bool
	isRunning   bool
	mu          sync.Mutex
}

// EveryMinute schedules the event to run every minute.
func (e *Event) EveryMinute() *Event {
	e.expression = "* * * * *"
	return e
}

// Hourly schedules the event to run hourly.
func (e *Event) Hourly() *Event {
	e.expression = "0 * * * *"
	return e
}

// Cron sets a custom cron expression.
func (e *Event) Cron(expression string) *Event {
	e.expression = expression
	return e
}

// WithoutOverlapping ensures the task does not overlap itself.
func (e *Event) WithoutOverlapping() *Event {
	e.withoutOverlap = true
	return e
}

func (e *Event) run() {
	if e.withoutOverlap {
		e.mu.Lock()
		if e.isRunning {
			e.mu.Unlock()
			return
		}
		e.isRunning = true
		e.mu.Unlock()

		defer func() {
			e.mu.Lock()
			e.isRunning = false
			e.mu.Unlock()
		}()
	}

	log.Printf("Running scheduled command: %s", e.command)

	// Parse command and args, then execute
	parts := strings.Fields(e.command)
	if len(parts) == 0 {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Scheduled command [%s] failed: %v", e.command, err)
	}
}

// Schedule manages registered tasks.
type Schedule struct {
	events []*Event
}

// NewSchedule creates a new Schedule instance.
func NewSchedule() *Schedule {
	return &Schedule{}
}

// Command registers a new scheduled command.
func (s *Schedule) Command(command string) *Event {
	event := &Event{command: command}
	s.events = append(s.events, event)
	return event
}

// Run executes the scheduler blocking the current thread.
// It stops when the context is cancelled.
func (s *Schedule) Run(ctx context.Context) {
	c := cron.New()
	
	for _, event := range s.events {
		if event.expression == "" {
			continue
		}

		ev := event 
		
		_, err := c.AddFunc(ev.expression, func() {
			ev.run()
		})
		
		if err != nil {
			log.Printf("Error scheduling command [%s]: %s", ev.command, err)
		}
	}

	log.Println("Scheduler started. Press CTRL+C to abort.")
	c.Start()
	
	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Scheduler shutting down...")
	c.Stop()
}

