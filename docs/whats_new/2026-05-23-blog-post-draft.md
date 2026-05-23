# GoW Framework Reaches Feature Parity — May 23, 2026 Release

**Published**: May 23, 2026  
**Author**: GoW Team

Today we are incredibly excited to announce the **May 23, 2026 Feature Parity Release** of the GoW Framework — the largest and most significant update since the project began.

This release marks a turning point. After an intensive, focused implementation wave, GoW has closed the vast majority of remaining critical gaps. The framework is no longer “promising” — it is now genuinely **production-ready** for building real-world applications in Go with a Laravel-like experience.

---

## What We Delivered

This release focused on eliminating the last major pain points that teams encountered when trying to adopt GoW for serious projects.

### Authentication & Security

- **Socialite (OAuth)** — Official support for logging in with Google, GitHub, and any custom provider. Full redirect + callback flow with an extensible architecture.
- **Two-Factor Authentication** — A complete pure-Go TOTP implementation (no external dependencies) including secret generation, verification with clock skew, and QR code URLs.
- **Account Lockout** — Automatic protection against brute-force login attempts.
- **Remember Me** tokens — Long-lived “remember me” support in the session guard.

### ORM (Goquent) — Massive Leap Forward

The ORM received one of its biggest upgrades yet:

- Polymorphic relations (`MorphOne`, `MorphMany`, `MorphTo`)
- `HasOneThrough` and `HasManyThrough`
- Full **Attribute Casting** system with 6 built-in casts (`datetime`, `json`, `bool`, `int`, `float`, `string`)
- **Accessors & Mutators** using familiar magic method naming (`GetXxxAttribute`, `SetXxxAttribute`)
- Relation Touching (`$touches`)
- Upsert support
- Pessimistic locking (`LockForUpdate` / `SharedLock`)
- Local scope auto-discovery
- `Chunk()` for safe processing of large datasets

Goquent is now one of the most complete and powerful ORMs available in the Go ecosystem.

### Developer Experience

- Dozens of new and completed Artisan commands (`make:mail`, `make:policy`, `make:resource`, `make:job`, `make:event`, `make:listener`, `cache:*`, `config:*`, `view:*`, `queue:retry`, `schedule:list`, and many more).
- Signed and temporary signed URLs with middleware protection.
- Advanced named rate limiting with per-user support.
- Complete Form Request validation with all Laravel-style hooks.
- Significantly improved testing tools (`actingAs()`, file upload testing, advanced JSON assertions).

### Ecosystem Foundations

We also laid the groundwork for many advanced features:

- Inertia.js support
- Livewire-like component system
- Scout (full-text search)
- Cashier (billing & subscriptions)
- Telescope / Pulse style debugging
- Octane-style high-performance server
- Prometheus metrics exporter

### Documentation

We published several new and heavily expanded guides, including:

- New [Socialite Guide](guide/socialite.md)
- New [Testing Guide](guide/testing.md)
- Major updates to the [ORM](guide/orm.md), [Authentication](guide/authentication.md), and [Routing](guide/routing.md) guides
- A comprehensive [Upgrade Guide](UPGRADE.md) for existing users

---

## Why This Release Matters

Before May 23, teams building production applications in Go with Laravel-like expectations would eventually hit hard walls around OAuth, 2FA, advanced relationships, casting, Artisan tooling, and testing ergonomics.

With this release, **most of those walls are gone**.

GoW is now a framework you can confidently choose for:
- SaaS products
- Admin panels and internal tools
- REST and GraphQL APIs
- Full-stack applications

---

## Upgrade Path

We’ve published a detailed [Upgrade Guide](UPGRADE.md) that walks through the recommended changes and new patterns. Most applications can upgrade with minimal friction, and we’ve made the new features optional where possible so you can adopt them gradually.

---

## What’s Next?

While this release brings GoW very close to feature parity, our work is not finished. Future focus areas include:

- Production-grade Mail drivers and more notification channels
- Richer error pages rendered through the view system
- First-class Inertia + Livewire experiences
- Official starter kits and scaffolding
- Continued performance and DX improvements

We’re also very open to community contributions in these areas.

---

## Thank You

This release would not have been possible without the incredible focus and effort put in over the last few days. We’re proud of how far GoW has come in such a short time.

To everyone who has been following the project, testing it, giving feedback, or simply watching — thank you.

GoW is now ready for the next chapter.

---

**Links**

- Full Release Notes: [docs/whats_new/2026-05-23.md](https://github.com/your-org/gow/blob/main/docs/whats_new/2026-05-23.md)
- GitHub Release: [2026-05-23 GitHub Release Notes](https://github.com/your-org/gow/releases)
- Upgrade Guide: [docs/UPGRADE.md](https://github.com/your-org/gow/blob/main/docs/UPGRADE.md)
- Documentation: [docs/](https://github.com/your-org/gow/tree/main/docs)

---

*GoW – Laravel in Go, done right.*

We can’t wait to see what you build.
