> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Socialite (OAuth Authentication)

> **Status**: ✅ Implemented

GoW provides a clean Socialite-like package for authenticating users via third-party OAuth providers (Google, GitHub, etc.).

## Installation / Registration

Register the Socialite service provider in your `bootstrap/app.go` (or manually):

```go
app.RegisterProvider(&socialite.ServiceProvider{})
```

## Configuring Providers

You can register providers like this:

```go
manager := app.Make((*socialite.Manager)(nil)).(*socialite.Manager)

manager.Extend("google", socialite.NewGoogleProvider(
    os.Getenv("GOOGLE_CLIENT_ID"),
    os.Getenv("GOOGLE_CLIENT_SECRET"),
    "http://localhost:8080/auth/google/callback",
))

manager.Extend("github", socialite.NewGitHubProvider(...))
```

## Redirecting Users

```go
router.Get("/auth/google", func(w, r) {
    url := manager.Driver("google").RedirectURL("random-state")
    http.Redirect(w, r, url, 302)
})
```

## Handling the Callback

```go
router.Get("/auth/google/callback", func(w, r) {
    user, err := manager.Driver("google").User(r.Context(), r.URL.Query().Get("code"))
    if err != nil { ... }

    // user.Email, user.Name, user.ID, user.Avatar, etc. are available
    // Log the user in or create account
})
```

## Supported Providers (Built-in)

- Google
- GitHub

You can easily add more by implementing the `socialite.Provider` interface.

## User Object

The returned `*socialite.User` contains:

- `ID`
- `Email`
- `Name`
- `Nickname`
- `Avatar`
- `Token`
- `RefreshToken`

## Best Practices

- Store the OAuth user ID + provider name in your users table for future logins.
- Always verify the `state` parameter in production to prevent CSRF.

---

## Complete End-to-End Example

### 1. Register providers in `bootstrap/app.go`

```go
manager := app.Make((*socialite.Manager)(nil)).(*socialite.Manager)

manager.Extend("google", socialite.NewGoogleProvider(
    config.Get("services.google.client_id"),
    config.Get("services.google.client_secret"),
    "http://localhost:8080/auth/google/callback",
))
```

### 2. Routes (`routes/web.go`)

```go
router.Get("/auth/google", func(w http.ResponseWriter, r *http.Request) {
    state := uuid.New().String() // store in session for verification
    url := manager.Driver("google").RedirectURL(state)
    http.Redirect(w, r, url, http.StatusFound)
})

router.Get("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    oauthUser, err := manager.Driver("google").User(r.Context(), code)
    if err != nil {
        http.Error(w, "OAuth failed", 500)
        return
    }

    // Find or create user in your database
    user := findOrCreateUserFromOAuth(oauthUser, "google")

    // Log them in
    auth.Login(user)

    http.Redirect(w, r, "/dashboard", http.StatusFound)
})
```

### 3. Storing OAuth users

Recommended table structure:

```sql
ALTER TABLE users ADD COLUMN provider VARCHAR(50);
ALTER TABLE users ADD COLUMN provider_id VARCHAR(100);
ALTER TABLE users ADD UNIQUE KEY provider_unique (provider, provider_id);
```

This pattern gives you a complete, production-ready OAuth login flow.
