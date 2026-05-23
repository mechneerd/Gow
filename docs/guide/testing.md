> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Testing

> **Status**: ✅ Good

GoW provides a rich testing experience inspired by Laravel.

## Base TestCase

```go
func TestUserCreation(t *testing.T) {
    tc := testing.NewTestCase(t, router)

    response := tc.Post("/users", map[string]any{
        "name": "John",
        "email": "john@example.com",
    })

    response.AssertStatus(201)
}
```

## Acting As (Authenticated User)

```go
tc.ActingAs(user).Get("/dashboard").AssertOk()
```

## File Upload Testing

```go
tc.Upload("/upload", "avatar", "photo.jpg", "fake-image-data")
```

## Advanced Assertions

```go
response.AssertJson("success", true)
response.AssertJsonStructure("data", "message")
```

## Database Assertions

```go
tc.AssertDatabaseHas("users", map[string]any{"email": "john@example.com"})
tc.AssertDatabaseMissing("users", map[string]any{"email": "ghost@example.com"})
tc.AssertDatabaseCount("posts", 42)
```

## Fakes

- `Mail::fake()`
- `Queue::fake()`
- `Event::fake()`

With rich assertion helpers (`AssertSent`, `AssertPushed`, etc.).

## RefreshDatabase Trait

Use `testing.RefreshDatabase` in your test base class for clean database state between tests.
