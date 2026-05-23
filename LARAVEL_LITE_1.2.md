# GoW Target: Laravel-lite 1.2

**Date**: 2026-05-23  
**Goal**: Deliver a **solid, reliable, production-usable "Laravel-lite"** framework by focusing ruthlessly on the core three pillars.

This document replaces the overly optimistic "Feature Parity" narrative. GoW will now target **Laravel-lite 1.2** — a well-tested, well-documented subset that actually works end-to-end.

## Definition of Laravel-lite 1.2

Laravel-lite means we deliver the **80% of Laravel that 80% of developers use every day**, done properly, instead of chasing 100% parity with half-working features.

### In Scope for 1.2 (The Core Three)

| # | Pillar                          | Must Be Production Ready                          |
|---|---------------------------------|---------------------------------------------------|
| 1 | Foundation & Architecture       | Container, Bootstrap, Service Providers, Config, .env, Publishing |
| 2 | HTTP Layer                      | Router (with groups, resources, named routes, signed URLs, model binding), Middleware (global + route + groups + aliases), Request helpers, FormRequest, Resources, Exception handling, Content negotiation |
| 3 | Database & ORM                  | Query Builder + full Eloquent-style ORM (Models, Relations, Eager loading, Soft deletes, Scopes, Observers, Mass assignment, Transactions, Pagination, Migrations, Seeders) |

### Explicitly Out of Scope for 1.2

- Livewire / Inertia (move to 2.x)
- Advanced Broadcasting (Reverb-like) — keep basic WebSocket for now
- Telescope, Horizon, Pulse
- Full Socialite (OAuth) — keep basic if time allows, otherwise defer
- Rich notification channels beyond Mail + Database + Slack
- Production-grade Mail drivers (SMTP + DKIM + retries) — basic only
- Most "ecosystem" tools and starter kits

## Success Criteria for 1.2

1. `go build ./...` succeeds with zero errors on the entire project
2. `go test ./... -race` passes cleanly for packages in pillars 1-2-3
3. Test coverage ≥ 65% on `foundation/`, `http/`, `routing/`, `database/`
4. Real integration test app (in `examples/`) that demonstrates:
   - Auth (login/register/logout)
   - CRUD with Eloquent models + relations
   - Middleware + FormRequest validation
   - Queued jobs + basic mail
5. Documentation in `docs/guide/` covers the three pillars with working examples
6. No more "placeholder" or "TODO" code in the core three pillars

## Philosophy

- **Boring is good.** Make the core features work reliably and predictably.
- **Honest documentation.** If something is 🟡 Partial or ❌ Missing, say so clearly.
- **Ship what works.** Defer everything else.

## Versioning

- **1.0** — Foundation + basic HTTP + basic DB (already largely achieved)
- **1.2** — Complete the three pillars + high test coverage + clean build (current target)
- **2.0** — Re-evaluate full ecosystem features (Livewire, advanced tools, etc.)

---

**Transition complete**:
- `Current_Capabilities.md` has been added to `.gitignore` and removed from git tracking.
- All future development and documentation should reference this `LARAVEL_LITE_1.2.md` document as the source of truth for project scope.
