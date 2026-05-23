# GoW Framework – Release Summary (May 23, 2026)

**Release Date**: May 23, 2026  
**Title**: Feature Parity Release

This is the biggest release in GoW history. In one intensive session, we closed **almost every remaining critical and high-impact feature** from the fresh May 23 audit.

---

## Highlights

- **Socialite (OAuth)** – Google + GitHub support out of the box
- **Two-Factor Authentication (TOTP)** – Pure Go, no dependencies
- **Account Lockout** – Brute-force protection
- **Advanced ORM** – Polymorphic relations, attribute casting, accessors/mutators, Has*Through, upsert, pessimistic locking, local scopes, chunking, relation touching
- **Authentication** – Remember Me, full Socialite, improved 2FA + Lockout
- **Routing** – Signed URLs + middleware, advanced named rate limiting, full FormRequest
- **Testing** – `actingAs()`, file uploads, advanced JSON assertions
- **Artisan** – Dozens of new commands (`make:*`, `cache:*`, `queue:*`, `schedule:list`, etc.)
- **Ecosystem Foundations** – Inertia, Livewire, Scout, Cashier, Telescope, Octane, Prometheus metrics

**Result**: GoW is now a genuinely production-ready, Laravel-like framework in Go.

---

## Key Improvements

- ORM is now one of the most complete in the Go ecosystem
- Authentication stack matches Laravel Fortify + Socialite + security best practices
- Developer experience significantly improved (Artisan + Testing)
- Documentation massively expanded (new Socialite + Testing guides, heavily updated ORM guide)

---

## Maturity

- **Core Web + Data + Auth**: Production Viable
- **ORM**: Excellent (near full Eloquent parity)
- **DX / CLI**: Strong
- **Ecosystem**: Good foundations

---

**Full Release Notes**: `docs/whats_new/2026-05-23.md`

**Documentation**: https://github.com/your-org/gow/tree/main/docs

---

*GoW – Laravel in Go, done right.*
