# GoW Feature Parity Release – One-Page Migration Checklist

**Release**: May 23, 2026 (Feature Parity)  
**Goal**: Quickly check what you should update after upgrading to the new version.

Print this page or keep it open while upgrading.

---

## Pre-Upgrade

- [ ] Backup your project (especially `app/`, `database/`, `routes/`, `config/`)
- [ ] Run `go get gow@latest && go mod tidy`
- [ ] Run `go build ./...` (note any errors before changes)
- [ ] Read the full [Upgrade Guide](UPGRADE.md)

---

## ORM Changes (Strongly Recommended)

- [ ] Add `Casts()` method to models that need automatic type conversion
- [ ] Add Accessors (`GetXxxAttribute`) and Mutators (`SetXxxAttribute`) where useful
- [ ] Update relationship tags to use new wrappers (`MorphMany`, `HasManyThrough`, etc.)
- [ ] Register polymorphic types in bootstrap:

  ```go
  orm.RegisterMorph("post", models.Post{})
  ```

- [ ] Replace large `Get()` calls with `Chunk()` where appropriate
- [ ] Add `Touches()` method on models that should update parent timestamps

---

## Authentication Changes

- [ ] (Optional but recommended) Migrate custom OAuth code to `socialite` package
- [ ] Add 2FA support using `auth/fortify/two_factor.go` (if needed)
- [ ] Implement Account Lockout in your login controller
- [ ] Update login calls to support Remember Me: `guard.Attempt(creds, true)`

---

## Routing & HTTP

- [ ] Replace manual signed URL logic with `router.SignedRoute()` / `TemporarySignedRoute()`
- [ ] Protect sensitive routes with `middleware.ValidateSignature()`
- [ ] Upgrade rate limiting to named version (`ThrottleNamed`) where useful
- [ ] Enhance existing Form Requests with `Messages()`, `Attributes()`, `PrepareForValidation()`, etc.

---

## Testing

- [ ] Replace manual auth simulation with `tc.ActingAs(user)`
- [ ] Add file upload tests using `tc.Upload(...)`
- [ ] Use new JSON assertions (`AssertJson`, `AssertJsonStructure`)

---

## Artisan / CLI

- [ ] Update any CI or scripts that depended on previously missing commands
- [ ] Start using new generators (`make:mail`, `make:policy`, `make:resource`, etc.)

---

## Documentation & Learning

- [ ] Review new [Socialite Guide](guide/socialite.md)
- [ ] Review new [Testing Guide](guide/testing.md)
- [ ] Re-read the heavily updated [ORM Guide](guide/orm.md)
- [ ] Read the full [Upgrade Guide](UPGRADE.md) (this checklist is a summary only)

---

## Post-Upgrade Verification

- [ ] `go build ./...`
- [ ] `go test ./...`
- [ ] `./artisan route:list`
- [ ] `./artisan schedule:list` (if using scheduler)
- [ ] Manual test of login, OAuth (if used), and key user flows

---

## Optional but Nice

- [ ] Add Slack notifications using the new `SlackChannel`
- [ ] Start experimenting with Inertia or Livewire foundations
- [ ] Add Prometheus metrics if you need observability

---

**Tip**: You don’t need to do everything at once.  
Adopt the new features gradually — most old code will continue to work.

---

**Need help?**  
Full release notes: `docs/whats_new/2026-05-23.md`  
Upgrade Guide: `docs/UPGRADE.md`

Print this checklist and check items off as you go.
