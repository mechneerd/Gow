package queue

import "sync"

// SyncDriver is an in-memory channel-based queue driver for local development.
type SyncDriver struct {
	jobs chan Job
	wg   sync.WaitGroup
}

// NewSyncDriver creates a new SyncDriver with the given buffer size.
func NewSyncDriver(bufferSize int) *SyncDriver {
	return &SyncDriver{
		jobs: make(chan Job, bufferSize),
	}
}

// Push adds a job to the channel.
func (d *SyncDriver) Push(job Job) error {
	d.jobs <- job
	return nil
}

// Pop retrieves the next job from the channel (blocking).
// In a real channel based worker, the worker just ranges over the channel.
func (d *SyncDriver) Pop() (Job, error) {
	job := <-d.jobs
	return job, nil
}

// Channel returns the underlying Go channel.
func (d *SyncDriver) Channel() <-chan Job {
	return d.jobs
}

