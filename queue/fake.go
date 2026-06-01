package queue

import (
	"sync"
)

// Fake is a fake queue dispatcher for testing.
type Fake struct {
	mu    sync.RWMutex
	jobs  []Job
	size  int
}

// NewFake creates a new fake queue dispatcher.
func NewFake() *Fake {
	return &Fake{
		jobs: make([]Job, 0),
	}
}

// Push captures the job instead of dispatching it.
func (f *Fake) Push(job Job) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.jobs = append(f.jobs, job)
	f.size++
	return nil
}

// Size returns the number of jobs in the fake queue.
func (f *Fake) Size() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.size
}

// GetJobs returns all captured jobs.
func (f *Fake) GetJobs() []Job {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make([]Job, len(f.jobs))
	copy(result, f.jobs)
	return result
}

// GetLastJob returns the last captured job, or nil if none.
func (f *Fake) GetLastJob() Job {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if len(f.jobs) == 0 {
		return nil
	}
	return f.jobs[len(f.jobs)-1]
}

// HasJob checks if a specific job type was queued.
func (f *Fake) HasJob(jobType any) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, job := range f.jobs {
		if job == jobType {
			return true
		}
	}
	return false
}

// GetJobCount returns the number of jobs queued.
func (f *Fake) GetJobCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.jobs)
}

// AssertPushed asserts that at least one job was pushed.
func (f *Fake) AssertPushed() bool {
	return f.GetJobCount() > 0
}

// AssertNotPushed asserts that no jobs were pushed.
func (f *Fake) AssertNotPushed() bool {
	return f.GetJobCount() == 0
}

// AssertPushedCount asserts the exact number of jobs pushed.
func (f *Fake) AssertPushedCount(count int) bool {
	return f.GetJobCount() == count
}

// AssertPushedJob asserts that a specific job type was pushed.
func (f *Fake) AssertPushedJob(jobType any) bool {
	return f.HasJob(jobType)
}

// AssertNotPushedJob asserts that a specific job type was NOT pushed.
func (f *Fake) AssertNotPushedJob(jobType any) bool {
	return !f.HasJob(jobType)
}

// Clear resets all captured jobs.
func (f *Fake) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.jobs = make([]Job, 0)
	f.size = 0
}

// Sync runs all captured jobs synchronously (useful in tests).
func (f *Fake) Sync() error {
	f.mu.Lock()
	jobs := make([]Job, len(f.jobs))
	copy(jobs, f.jobs)
	f.jobs = nil
	f.size = 0
	f.mu.Unlock()

	for _, job := range jobs {
		if err := job.Handle(); err != nil {
			return err
		}
	}
	return nil
}

// SyncJobs runs jobs of a specific type synchronously.
func (f *Fake) SyncJobs(jobType any) error {
	f.mu.Lock()
	var remaining []Job
	var toRun []Job
	for _, job := range f.jobs {
		if job == jobType {
			toRun = append(toRun, job)
		} else {
			remaining = append(remaining, job)
		}
	}
	f.jobs = remaining
	f.size = len(remaining)
	f.mu.Unlock()

	for _, job := range toRun {
		if err := job.Handle(); err != nil {
			return err
		}
	}
	return nil
}
