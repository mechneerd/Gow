package queue

import (
	"fmt"
	"sync"
	"time"
)

// FailedJob represents a failed job record.
type FailedJob struct {
	ID        int64
	Queue     string
	Payload   string
	Exception string
	FailedAt  time.Time
}

// FailedJobStore manages failed jobs.
type FailedJobStore struct {
	jobs []FailedJob
	mu   sync.RWMutex
	next int64
}

var globalFailedStore = &FailedJobStore{}

// RecordFailedJob stores a failed job.
func RecordFailedJob(job Job, err error) {
	globalFailedStore.mu.Lock()
	defer globalFailedStore.mu.Unlock()

	globalFailedStore.next++
	globalFailedStore.jobs = append(globalFailedStore.jobs, FailedJob{
		ID:        globalFailedStore.next,
		Queue:     "default",
		Payload:   fmt.Sprintf("%T", job),
		Exception: err.Error(),
		FailedAt:  time.Now(),
	})
}

// GetFailedJobs returns all failed jobs.
func GetFailedJobs() []FailedJob {
	globalFailedStore.mu.RLock()
	defer globalFailedStore.mu.RUnlock()
	result := make([]FailedJob, len(globalFailedStore.jobs))
	copy(result, globalFailedStore.jobs)
	return result
}

// RetryFailedJob retries a failed job by ID.
func RetryFailedJob(id int64) error {
	globalFailedStore.mu.Lock()
	defer globalFailedStore.mu.Unlock()

	for i, job := range globalFailedStore.jobs {
		if job.ID == id {
			globalFailedStore.jobs = append(globalFailedStore.jobs[:i], globalFailedStore.jobs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("failed job %d not found", id)
}

// FlushFailedJobs clears all failed jobs.
func FlushFailedJobs() {
	globalFailedStore.mu.Lock()
	defer globalFailedStore.mu.Unlock()
	globalFailedStore.jobs = nil
}

// JobBatch represents a batch of jobs.
type JobBatch struct {
	ID        int64
	Jobs      []Job
	Pending   int
	Completed int
	Failed    int
	Then      func(batch *JobBatch) // callback when all jobs complete
	Catch     func(batch *JobBatch, err error) // callback on failure
}

var (
	batches   = make(map[int64]*JobBatch)
	batchMu   sync.RWMutex
	batchNext int64
)

// Bus dispatches jobs and manages batches.
type Bus struct {
	manager *Manager
}

// NewBus creates a new job bus.
func NewBus(manager *Manager) *Bus {
	return &Bus{manager: manager}
}

// Dispatch pushes a job to the queue.
func (b *Bus) Dispatch(job Job) error {
	return b.manager.Push(job)
}

// Batch creates a batch of jobs and dispatches them.
func (b *Bus) Batch(jobs []Job) *BatchBuilder {
	return &BatchBuilder{bus: b, jobs: jobs}
}

// BatchBuilder configures a batch of jobs.
type BatchBuilder struct {
	bus   *Bus
	jobs  []Job
	then  func(*JobBatch)
	catch func(*JobBatch, error)
}

// Then sets the callback for when all jobs complete.
func (bb *BatchBuilder) Then(fn func(*JobBatch)) *BatchBuilder {
	bb.then = fn
	return bb
}

// Catch sets the callback for when a job fails.
func (bb *BatchBuilder) Catch(fn func(*JobBatch, error)) *BatchBuilder {
	bb.catch = fn
	return bb
}

// Dispatch dispatches all jobs in the batch.
func (bb *BatchBuilder) Dispatch() (*JobBatch, error) {
	batchMu.Lock()
	batchNext++
	batch := &JobBatch{
		ID:      batchNext,
		Jobs:    bb.jobs,
		Pending: len(bb.jobs),
		Then:    bb.then,
		Catch:   bb.catch,
	}
	batches[batch.ID] = batch
	batchMu.Unlock()

	for _, job := range bb.jobs {
		if err := bb.bus.Dispatch(job); err != nil {
			batch.Failed++
			batch.Pending--
			if batch.Catch != nil {
				batch.Catch(batch, err)
			}
			return batch, err
		}
	}

	return batch, nil
}

// Chain dispatches jobs sequentially.
func (b *Bus) Chain(jobs []Job) *ChainBuilder {
	return &ChainBuilder{bus: b, jobs: jobs}
}

// ChainBuilder configures sequential job execution.
type ChainBuilder struct {
	bus   *Bus
	jobs  []Job
	then  func()
	catch func(error)
}

// Then sets the callback for when all jobs complete.
func (cb *ChainBuilder) Then(fn func()) *ChainBuilder {
	cb.then = fn
	return cb
}

// Catch sets the callback for when a job fails.
func (cb *ChainBuilder) Catch(fn func(error)) *ChainBuilder {
	cb.catch = fn
	return cb
}

// Dispatch dispatches jobs sequentially.
func (cb *ChainBuilder) Dispatch() error {
	for _, job := range cb.jobs {
		if err := cb.bus.Dispatch(job); err != nil {
			if cb.catch != nil {
				cb.catch(err)
			}
			return err
		}
	}
	if cb.then != nil {
		cb.then()
	}
	return nil
}

// Delay dispatches a job with a delay.
func (b *Bus) Delay(job Job, delay time.Duration) error {
	go func() {
		time.Sleep(delay)
		b.manager.Push(job)
	}()
	return nil
}
