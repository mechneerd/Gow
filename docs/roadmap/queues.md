# Queues

> **Status**: ?? Planned (Currently stubbed or not implemented)


While building your web application, you may have some tasks, such as parsing and storing an uploaded CSV file, that take too long to perform during a typical web request. GoW's queue behavior allows you to defer the processing of a time-consuming task until a later time.

## Drivers

GoW supports multiple queue drivers: `sync` (default), `database`, and `redis`.

To use the `database` driver, you will need a `jobs` table and a `failed_jobs` table.

## Dispatching Jobs

You can push a job onto the queue using the `queue` manager.

```go
// Define a job handler
queue.Register("ProcessPodcast", func(payload []byte) error {
    // Process the podcast...
    return nil
})

// Dispatch the job
queue.Push("ProcessPodcast", []byte(`{"podcast_id": 1}`))
```

## Running the Queue Worker

To process jobs in the queue, you must run the Artisan queue worker.

```bash
artisan queue:work
```

This worker will continue to run in the foreground, pulling jobs off the queue and executing them.

## Retrying Failed Jobs

If a job fails, it will be placed in the `failed_jobs` table (if using the database driver). You can retry all failed jobs using Artisan:

```bash
artisan queue:retry
```
