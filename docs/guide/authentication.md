> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Authentication

> **Status**: ✅ Implemented


GoW provides robust, built-in services for authenticating users out of the box, utilizing guards, providers, and token issuance strategies.

## Headless Authentication (Fortify Equivalent)

If you are building a Single Page Application (SPA) or a mobile application, GoW includes a headless authentication backend that provides standard JSON API endpoints for complete user management.

To use these routes, register the `Fortify` backend in your application router:

```go
fortifyBackend := fortify.New(authManager)
fortifyBackend.RegisterRoutes(router)
```

This will automatically expose the following endpoints:

- `POST /api/register`: Register a new user.
- `POST /api/login`: Authenticate the user and establish a session.
- `POST /api/forgot-password`: Send a password reset link.
- `POST /api/reset-password`: Reset the user's password.
- `GET /api/email/verify/{id}/{hash}`: Verify a user's email via signed URL.
- `POST /api/two-factor/enable`: Generate a 2FA QR code and recovery codes.

## Token Authentication (Sanctum Equivalent)

For API authentication, GoW includes a Sanctum equivalent that issues and verifies personal access tokens.

Tokens are cryptographically hashed using SHA-256 before being stored in your database, preventing database leaks from compromising active tokens.

### Issuing Tokens

You can issue a token using the `sanctum` package:

```go
token, plainTextToken, err := sanctum.IssueToken(user.ID, "my-device", []string{"read", "write"})
```

The `plainTextToken` should be returned to the client immediately, as it cannot be retrieved again.

### Protecting Routes

To protect your API routes, assign the `auth:sanctum` middleware to the route group.

```go
api := router.Group("/api")
api.Use(authMiddleware.Guard("sanctum"))

api.Get("/user", func(w http.ResponseWriter, r *http.Request) {
    // The user is authenticated!
})
```
