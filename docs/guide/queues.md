> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Queues

> **Status**: ✅ Implemented

GoW provides a goroutine-safe, driver-based queue system for deferring time-consuming tasks out of the HTTP request cycle.

## Drivers

GoW ships with two built-in queue drivers:

| Driver | Description |
|---|---|
| `memory` | Buffered Go channel — zero dependencies, perfect for development and single-process apps |
| `database` | Backed by a `jobs` SQL table — persists across restarts |

The default driver is `memory` with a 10,000-job buffer, automatically registered when you call `queue.NewManager()`.

## Defining Jobs

A job is any type implementing the `Job` interface:

```go
type Job interface {
    Handle() error
    Failed(err error)
}
```

```go
type SendWelcomeEmail struct {
    UserID int
    Email  string
}

func (j *SendWelcomeEmail) Handle() error {
    // send the email...
    return nil
}

func (j *SendWelcomeEmail) Failed(err error) {
    log.Printf("job failed for user %d: %v", j.UserID, err)
}
```

## Dispatching Jobs

Push a job onto the queue via the manager:

```go
queueManager := queue.NewManager("memory")

err := queueManager.Push(&SendWelcomeEmail{
    UserID: user.ID,
    Email:  user.Email,
})
```

## Consuming Jobs

### Blocking Pop

Blocks until a job is available — ideal for a dedicated worker goroutine:

```go
driver := queueManager.Connection("memory")

for {
    job, err := driver.Pop()
    if err != nil {
        log.Printf("queue error: %v", err)
        continue
    }
    if err := job.Handle(); err != nil {
        job.Failed(err)
    }
}
```

### Non-Blocking TryPop

Returns `nil` immediately if the queue is empty:

```go
if job := driver.(*queue.MemoryDriver).TryPop(); job != nil {
    job.Handle()
}
```

## Queue Depth

Use `Len()` for health checks and monitoring:

```go
depth := driver.(*queue.MemoryDriver).Len()
fmt.Printf("Jobs waiting: %d\n", depth)
```

## Using a Custom Driver

Register any driver satisfying the `Driver` interface:

```go
queueManager.AddDriver("redis", myRedisDriver)
```
