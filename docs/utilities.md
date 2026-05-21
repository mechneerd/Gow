# Utilities

GoW provides several powerful utilities to simplify your daily development workflow.

## Feature Flags (Pennant Equivalent)

Feature flags allow you to confidently roll out new features, perform A/B testing, and toggle functionality without deploying new code.

GoW manages feature flags using a centralized `pennant.Manager` that supports multiple persistent stores (Redis, Database, InMemory).

```go
// Check if a feature is active for a given user
if pennantManager.Active(ctx, "new-api", user) {
    // Execute new API logic...
}
```

## Health Checks

Monitoring the health of your application and its dependencies (Database, Redis, etc.) is critical in production. GoW exposes a `/up` endpoint to quickly verify system health.

You can register custom `Checker` structs to monitor external services:

```go
healthManager.Add(&DatabaseChecker{})
healthManager.Add(&RedisChecker{})
```

Hitting `/up` will return a JSON response with the health of all registered components.

## Generic Pipeline

The Pipeline pattern allows you to pass a payload through a series of operations (pipes), making it excellent for data transformation or custom middleware flows outside of HTTP routing.

```go
result, err := pipeline.New[string]().
    Through(TrimWhitespace, ToUpper).
    Then(ctx, "  hello world  ", func(c context.Context, payload string) (string, error) {
        return payload, nil
    })

// result == "HELLO WORLD"
```

## Process Wrapper

When you need to execute external commands via `os/exec`, GoW's `process` wrapper provides timeout handling and testing capabilities.

```go
result := process.Command("ls", "-la").
    Timeout(5 * time.Second).
    Run()

fmt.Println(result.Stdout)
```

In your tests, you can easily mock the output of commands:

```go
process.Fake(map[string]*process.Result{
    "ls -la": {Stdout: "fake output", ExitCode: 0},
})
```
