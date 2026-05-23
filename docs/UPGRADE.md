# Upgrade Guide — GoW Feature Parity Release (May 23, 2026)

This guide helps you upgrade your existing GoW application to the **May 23, 2026 Feature Parity Release**.

This release is one of the largest in GoW history. It closes most remaining critical gaps and brings near-complete Laravel feature parity.

---

## Before You Upgrade

1. **Backup your project** (especially `app/`, `database/`, `routes/`, and `config/`).
2. Update your `go.mod`:

   ```bash
   go get gow@latest
   go mod tidy
   ```

3. Run `go build ./...` before and after to catch any compilation issues.

---

## Breaking Changes & Important Notes

This release is **mostly backward compatible**, but there are some recommended changes and new patterns you should adopt.

### 1. ORM Changes (Highly Recommended)

#### New Recommended Ways

- Use the new **Attribute Casting** system instead of manual JSON/datetime handling.
- Prefer **Accessors & Mutators** (`GetXxxAttribute` / `SetXxxAttribute`) for computed or transformed fields.
- Register polymorphic types early in bootstrap:

  ```go
  orm.RegisterMorph("post", models.Post{})
  orm.RegisterMorph("video", models.Video{})
  ```

- Use `ModelQuery.Chunk()` for large datasets instead of loading everything at once.

#### Relation Tags

New relation types are now supported via `gow` struct tags:

```go
Comments orm.MorphMany[Comment] `gow:"morphMany,morphType=commentable_type,morphId=commentable_id,type=Post"`
Posts    orm.HasManyThrough[Post] `gow:"hasManyThrough,through=users,foreignKey=country_id,relatedKey=user_id"`
```

Update your models to use the new relation wrappers where appropriate (`MorphOne`, `MorphMany`, `MorphTo`, `HasOneThrough`, `HasManyThrough`).

### 2. Authentication Changes

#### Socialite (OAuth)

If you were previously using custom OAuth code, consider migrating to the official `socialite` package.

New recommended flow:

```go
manager.Extend("google", socialite.NewGoogleProvider(...))
```

See the new [Socialite Guide](guide/socialite.md).

#### Two-Factor Authentication

2FA is now available out of the box via `auth/fortify/two_factor.go`.

Enable it through the existing Fortify routes (`POST /api/two-factor/enable`).

#### Account Lockout

Add brute-force protection in your login flow:

```go
lockout := auth.NewLockout(5, 15) // 5 attempts, 15 min lock
```

### 3. Routing & HTTP

#### Signed URLs

Replace any manual signed URL logic with the new helpers:

```go
url, _ := router.TemporarySignedRoute("verify.email", 30*time.Minute, params)
```

Protect routes with:

```go
router.Use(middleware.ValidateSignature(urlGenerator))
```

#### Rate Limiting

Replace basic `ThrottleRequests` calls with the new named version when you want per-user or named limiters:

```go
router.Use(middleware.ThrottleNamed("login", 5, 1))
```

#### Form Requests

Form Request structs can now implement additional methods:

- `Messages()`
- `Attributes()`
- `PrepareForValidation()`
- `PassedValidation()`
- `FailedValidation()`

Update your existing Form Requests to take advantage of these for better DX.

### 4. Testing

Replace manual auth header tricks with the new helper:

```go
tc.ActingAs(user).Get("/dashboard")
```

File upload testing is now available:

```go
tc.Upload("/upload", "file", "test.txt", "content")
```

### 5. Artisan Commands

Many new commands are now available. Update any custom scripts or CI that relied on missing commands.

Notable new commands:
- `make:mail`, `make:event`, `make:listener`, `make:policy`, `make:resource`, `make:job`, `make:notification`, `make:test`
- `cache:clear`, `config:cache`, `view:clear`
- `queue:retry`, `queue:failed`, `schedule:list`

---

## Recommended Adoption Steps

1. **ORM Modernization** (High priority)
   - Add `Casts()` to models that need casting.
   - Add accessors/mutators for computed fields.
   - Register polymorphic types.
   - Start using `Chunk()` for large data processing.

2. **Authentication Hardening** (Recommended)
   - Add Socialite for third-party login.
   - Enable 2FA for admin/user accounts.
   - Add Account Lockout to login flows.

3. **Routing Improvements**
   - Use Signed URLs for sensitive actions (email verification, password reset, etc.).
   - Switch to named rate limiters for better control.

4. **Testing**
   - Refactor test files to use `ActingAs()` and new assertion helpers.

5. **Documentation**
   - Review the updated guides, especially:
     - [ORM](guide/orm.md) (heavily expanded)
     - [Authentication](guide/authentication.md)
     - New [Socialite](guide/socialite.md) and [Testing](guide/testing.md) guides

---

## New Packages You Can Start Using

| Package | Purpose |
|---------|---------|
| `auth/socialite` | OAuth login (Google, GitHub, etc.) |
| `auth/fortify/two_factor` | TOTP 2FA |
| `auth/lockout` | Account lockout |
| `http/middleware/signed` | Validate signed URLs |
| `http/inertia` | Inertia.js responses |
| `http/livewire` | Livewire-like components |
| `support/scout` | Full-text search |
| `billing/cashier` | Subscriptions & billing |
| `support/telescope` | Debugging panel |
| `support/metrics` | Prometheus exporter |

---

## Deprecations & Removals

- No major deprecations in this release.
- Old manual implementations of features now covered by the framework (e.g., custom OAuth, manual casting) are still supported but **not recommended**.

---

## After Upgrading

Run the following to verify everything works:

```bash
go build ./...
go test ./...
./artisan route:list
./artisan schedule:list
```

Then gradually adopt the new patterns listed above.

---

## Getting Help

- Full release notes: `docs/whats_new/2026-05-23.md`
- GitHub Release notes: `docs/whats_new/2026-05-23-github-release.md`
- Updated capabilities: `Current_Capabilities.md`

If you run into issues during upgrade, please open an issue with your `go.mod` and the error message.

---

**Welcome to the new era of GoW.**

This release brings GoW very close to full Laravel-like completeness. Enjoy building with it!
