# GoW Framework v0.9.0 – Feature Parity Release

**Release Date**: May 23, 2026  
**Codename**: Feature Parity

This is the **largest and most important release** in GoW history.

After an intensive implementation wave, GoW has closed nearly every remaining critical and high-impact gap identified in the May 23 audit. The framework has reached a new level of maturity and is now genuinely production-ready for real-world applications.

---

## 🚀 Highlights

- **Socialite (OAuth)** – First-class support for Google, GitHub, and custom providers
- **Two-Factor Authentication** – Pure-Go TOTP with QR codes (no external dependencies)
- **Account Lockout** – Built-in brute-force protection
- **Advanced ORM** – Polymorphic relations, attribute casting, accessors/mutators, Has*Through, upsert, pessimistic locking, local scopes, chunking, and relation touching
- **Authentication** – Remember Me, full Socialite integration, improved 2FA + Lockout
- **Routing & HTTP** – Signed URLs + middleware, named rate limiting, complete FormRequest support
- **Testing** – `actingAs()`, file upload helpers, advanced JSON assertions
- **Artisan** – Massive expansion of `make:*` and operational commands
- **Ecosystem Foundations** – Inertia, Livewire, Scout, Cashier, Telescope, Octane, Prometheus metrics

**Result**: GoW now offers near-complete Laravel feature parity in Go.

---

## ✨ New Features

### Authentication & Security
- Full **Socialite** (OAuth) system with Google and GitHub providers
- **Two-Factor Authentication** (TOTP) – `GenerateSecret`, `VerifyTOTP`, QR URL generation
- **Account Lockout** after repeated failed login attempts
- **Remember Me** tokens support

### ORM (Goquent)
- Polymorphic relations (`MorphOne`, `MorphMany`, `MorphTo`)
- `HasOneThrough` / `HasManyThrough`
- Full **Attribute Casting** system (6 built-in casts)
- **Accessors & Mutators** with magic methods
- **Relation Touching** (`$touches`)
- **Upsert** support
- **Pessimistic Locking** (`LockForUpdate`, `SharedLock`)
- Local scope auto-discovery
- `Chunk()` for processing large datasets

### HTTP & Routing
- Signed and Temporary Signed URLs + `ValidateSignature` middleware
- Advanced named rate limiting (`ThrottleNamed`, per-user throttling)
- Complete **FormRequest** with all hooks (`Messages`, `Attributes`, `PrepareForValidation`, etc.)

### Testing
- `TestCase.ActingAs(user)`
- File upload testing
- Advanced JSON structure assertions

### CLI (Artisan)
- Many new `make:*` generators (`make:mail`, `make:event`, `make:policy`, `make:resource`, `make:job`, `make:notification`, `make:test`, etc.)
- `cache:*`, `config:*`, `view:*` commands
- `queue:retry`, `queue:failed`, `queue:flush`
- `schedule:list`

### Other
- Slack notification channel
- Foundations for Inertia, Livewire, Scout, Cashier, Telescope, Octane, and Prometheus metrics

---

## 📚 Documentation

- New guide: [Socialite (OAuth)](docs/guide/socialite.md)
- New guide: [Testing](docs/guide/testing.md)
- Major expansions to [ORM](docs/guide/orm.md), [Authentication](docs/guide/authentication.md), and [Routing](docs/guide/routing.md)
- Deep-dive code examples added throughout
- Roadmap and documentation index updated to reflect current reality

---

## 📦 Installation / Upgrade

```bash
go get gow@latest
```

Full release notes and migration guidance:  
[docs/whats_new/2026-05-23.md](docs/whats_new/2026-05-23.md)

---

## 🙏 Thank You

This release represents an enormous amount of focused work. GoW is now one of the most complete Laravel-like frameworks available in Go.

We look forward to seeing what you build with it.

---

**Full Changelog**: [CHANGELOG.md](CHANGELOG.md)  
**Documentation**: [docs/](docs/)  
**Release Notes**: [docs/whats_new/2026-05-23.md](docs/whats_new/2026-05-23.md)

---

*GoW – Laravel in Go, done right.*
