package queue

import (
	"sync"
	"testing"
)

type dummyJob struct {
	id int
}

func (j *dummyJob) Handle() error { return nil }
func (j *dummyJob) Failed(err error) {}

func TestMemoryDriverConcurrentPushPop(t *testing.T) {
	driver := NewMemoryDriver(1000)
	var wg sync.WaitGroup

	numJobs := 500
	wg.Add(numJobs)

	// Concurrent Push
	for i := 0; i < numJobs; i++ {
		go func(id int) {
			err := driver.Push(&dummyJob{id: id})
			if err != nil {
				t.Errorf("Unexpected error on push: %v", err)
			}
		}(i)
	}

	// Concurrent Pop
	for i := 0; i < numJobs; i++ {
		go func() {
			defer wg.Done()
			job, err := driver.Pop()
			if err != nil {
				t.Errorf("Unexpected error on pop: %v", err)
			}
			if job == nil {
				t.Errorf("Expected job, got nil")
			}
		}()
	}

	wg.Wait()

	if driver.Len() != 0 {
		t.Errorf("Expected queue length to be 0, got %d", driver.Len())
	}
}

func TestMemoryDriverTryPop(t *testing.T) {
	driver := NewMemoryDriver(10)

	if driver.TryPop() != nil {
		t.Errorf("Expected nil on empty queue")
	}

	driver.Push(&dummyJob{id: 1})

	if driver.Len() != 1 {
		t.Errorf("Expected queue length 1, got %d", driver.Len())
	}

	job := driver.TryPop()
	if job == nil {
		t.Errorf("Expected to pop a job")
	}
	
	if driver.Len() != 0 {
		t.Errorf("Expected queue length 0, got %d", driver.Len())
	}
}

