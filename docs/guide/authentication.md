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
