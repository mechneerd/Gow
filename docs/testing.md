# Testing

Testing is a first-class citizen in GoW. By integrating deeply with Go's standard `testing` package and popular assertion libraries like `testify/assert`, GoW makes writing feature and unit tests a breeze.

## HTTP Testing

GoW provides a fluent testing wrapper around your `Router` utilizing the `httptest` package. This allows you to simulate HTTP requests without starting a real web server.

```go
func TestWelcomeScreen(t *testing.T) {
    // Initialize your application router
    router := setupRouter()
    
    // Create a new TestCase
    tc := testing.NewTestCase(t, router)
    
    // Perform a request and assert the response
    tc.Get("/welcome").
       AssertStatus(200).
       AssertSee("Welcome to GoW!")
}
```

### JSON APIs

Testing JSON APIs is simple using the `AssertJson` method.

```go
func TestUserProfile(t *testing.T) {
    tc := testing.NewTestCase(t, router)
    
    tc.Get("/api/user/1").
       AssertStatus(200).
       AssertJson(map[string]any{
           "id": 1,
           "name": "Taylor",
       })
}
```

## Database Assertions

When running tests, you frequently need to verify that certain data was inserted or deleted from the database. GoW provides helpers like `AssertDatabaseHas` and `AssertDatabaseMissing`.

```go
func TestUserRegistration(t *testing.T) {
    tc := testing.NewTestCase(t, router)
    
    tc.Post("/register", map[string]string{
        "email": "test@example.com",
        "password": "password123",
    }).AssertStatus(201)
    
    // Verify the user was created in the database
    tc.AssertDatabaseHas("users", map[string]any{
        "email": "test@example.com",
    })
}
```

## Time Travel

When testing features that depend on time (e.g. expiring tokens), you can "travel" to a specific point in time so your assertions behave predictably.

```go
func TestTokenExpiration(t *testing.T) {
    // Travel exactly 2 days into the future
    testing.Travel(48 * time.Hour)
    
    // Test logic...
    
    // Reset time back to normal
    testing.TravelBack()
}
```
