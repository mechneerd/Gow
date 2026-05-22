package queue

// MemoryDriver is a goroutine-safe queue driver backed by a buffered channel.
type MemoryDriver struct {
	jobs chan Job
}

// NewMemoryDriver initializes a new MemoryDriver with a specified buffer capacity.
func NewMemoryDriver(capacity int) *MemoryDriver {
	return &MemoryDriver{
		jobs: make(chan Job, capacity),
	}
}

// Push adds a new job to the channel. Blocks if the buffer is full.
func (d *MemoryDriver) Push(job Job) error {
	d.jobs <- job
	return nil
}

// Pop blocks until a job is available and returns it.
func (d *MemoryDriver) Pop() (Job, error) {
	job, ok := <-d.jobs
	if !ok {
		return nil, nil // channel closed
	}
	return job, nil
}

// TryPop attempts to read a job from the channel without blocking.
// Returns nil if the queue is empty.
func (d *MemoryDriver) TryPop() Job {
	select {
	case job, ok := <-d.jobs:
		if !ok {
			return nil
		}
		return job
	default:
		return nil
	}
}

// Len returns the current number of jobs waiting in the queue.
func (d *MemoryDriver) Len() int {
	return len(d.jobs)
}
