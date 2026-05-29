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

```go
token, plainTextToken, err := sanctum.IssueToken(user.ID, "my-device", []string{"read", "write"})
```

### Protecting Routes

```go
api := router.Group("/api")
api.Use(authMiddleware.Guard("sanctum"))
```

## Auth Middleware Guidance

GoW provides two auth middleware options. Choose the right one for your use case:

### `auth.Middleware(guard)` — Recommended for most cases

This is the **primary auth middleware**. It stores the authenticated user in the request context, making it available via `auth.User(r)`.

```go
import "github.com/mechneerd/gow/auth"

guard := authManager.Guard("web")
router.Use(auth.Middleware(guard))
```

**Use when:**
- You need `auth.User(r)` to work in downstream handlers
- Building web applications with session-based auth
- You want the user available in template data

### `middleware.Authenticate(manager, guardName)` — Deprecated

This is a **legacy wrapper** kept for backward compatibility. It does NOT store the user in context — use `auth.Middleware` instead.

```go
import "github.com/mechneerd/gow/http/middleware"

// Deprecated: use auth.Middleware(guard) instead
router.Use(middleware.Authenticate(authManager, "web"))
```

### `middleware.Authorize(gate, ability)` — For authorization

Use this **after** `auth.Middleware` to check permissions via the Gate system:

```go
router.Use(auth.Middleware(guard))
router.Use(middleware.Authorize(gate, "edit-post"))
```

### `auth.GuestMiddleware(guard, redirectTo)` — Guest-only routes

Restricts routes to unauthenticated users (e.g., login/register pages):

```go
router.Use(auth.GuestMiddleware(guard, "/dashboard"))
```

---

## Socialite (OAuth)

GoW supports third-party login via the `socialite` package.

See the full guide: [Socialite (OAuth)](socialite.md)

Quick example:

```go
manager.Extend("google", socialite.NewGoogleProvider(clientID, secret, redirectURL))
url := manager.Driver("google").RedirectURL(state)
```

---

## Two-Factor Authentication (2FA)

GoW includes a pure-Go TOTP implementation in `auth/fortify/two_factor.go`.

### Generating a Secret

```go
secret, _ := fortify.GenerateSecret()
qrURL := fortify.GetQRCodeURL("MyApp", user.Email, secret)
```

### Verifying Codes

```go
valid := fortify.VerifyTOTP(secret, userInputCode, 1) // 1 = allow 30s clock skew
```

Enable 2FA via the Fortify routes (`/api/two-factor/enable`).

---

## Account Lockout

Prevent brute-force attacks:

```go
lockout := auth.NewLockout(5, 15) // 5 attempts → 15 min lock

if lockout.IsLocked(email) {
    return errors.New("account locked")
}

if !guard.Attempt(creds) {
    lockout.RecordFailure(email)
} else {
    lockout.Reset(email)
}
```

### Recommended Usage in Login Controller

```go
func Login(w, r) {
    email := request.Input(r, "email")

    if lockout.IsLocked(email) {
        return response.Error("Too many attempts. Try again later.")
    }

    if guard.Attempt(map[string]any{"email": email, "password": ...}) {
        lockout.Reset(email)
        // login success
    } else {
        lockout.RecordFailure(email)
    }
}
```

---

## Remember Me

```go
guard.Attempt(credentials, true) // second param = remember
```

The `SessionGuard` will call `UpdateRememberToken` on your user provider.

---

## Email Verification & Password Reset

Fully supported via Fortify + signed URLs (see earlier sections).
