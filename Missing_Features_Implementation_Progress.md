# Missing Features Implementation Progress

> **Date**: 2026-05-22  
> **Based on**: `MISSING_FEATURES.md` (482 lines, ~250 items)  
> **Current Status**: Top 10 + Wave 2 completed. Most critical bugs from the original audit have been fixed.  
> **Legend**: 🔴 Critical | 🟡 Important | 🟢 Nice to Have

---

## Executive Summary

The `MISSING_FEATURES.md` document reveals a very large surface of missing functionality (~250 items).

After completing the original Top 10 and Wave 2, the framework has a solid foundation. The next phase must focus on **making the framework usable for real-world applications**.

This document defines a phased implementation roadmap (Wave 3 → Wave 8) to systematically close the gaps.

---

## Verification of Critical Bugs Section

The first 6 "🔴 CRITICAL BUGS" listed in `MISSING_FEATURES.md` were largely resolved during earlier phases:

| Bug | Status | Notes |
|-----|--------|-------|
| Eager Loading Returns Nil | ✅ Fixed | `loadRelation` now executes queries |
| Observers Never Called | ✅ Fixed | `DispatchModelEvent` calls exist in `query.go` |
| Global Scopes Never Applied | ✅ Fixed | Applied in `NewQuery()` |
| Collection.Chunk() Returns Nil | ✅ Fixed | `Chunked()` renamed + fixed |
| @while/@once Compile to Comments | ✅ Fixed | Custom template funcs implemented |
| S3Driver is All No-Ops | ⚠️ Partial | Still needs real aws-sdk-go-v2 implementation |

**Action**: A quick audit pass will be done in Wave 3 to confirm these are truly resolved.

---

## Implementation Waves

### Wave 3 — Foundation Hardening (Critical Data Layer)
**Goal**: Make basic CRUD + relationships + transactions reliable.

**Priority Items** (from MISSING_FEATURES.md):

| # | Feature | Category | Priority | Files to Touch / New |
|---|---------|----------|----------|----------------------|
| 1 | Transactions (`DB.Transaction`, Begin/Commit/Rollback, Savepoints) | Database | 🔴 | `database/manager.go`, `database/connection.go`, `database/orm/query.go` |
| 2 | Soft Deletes (full: WithTrashed, OnlyTrashed, Restore, ForceDelete) | ORM | 🔴 | `database/orm/softdelete.go` (expand) |
| 3 | Working Eager Loading + BelongsToMany (pivot) | ORM | 🔴 | `database/orm/relation.go`, new `belongs_to_many.go` |
| 4 | Full Pagination (`Paginate()`, `SimplePaginate()`, Cursor) | Database | 🔴 | `database/pagination/paginator.go` + builder integration |
| 5 | Mass Assignment (`$fillable`, `$guarded`) + `toJSON`/`toArray` | ORM | 🔴 | `database/orm/model.go` |
| 6 | FirstOrCreate / FirstOrNew / UpdateOrCreate | ORM | 🔴 | `database/orm/query.go` |
| 7 | Raw Expressions + GROUP BY / HAVING / Distinct | Query Builder | 🔴 | `database/query/builder.go` + dialects |
| 8 | Verify & close any remaining Critical Bugs | Audit | 🔴 | All ORM + Collection files |

**Estimated Effort**: 3–4 weeks

---

### Wave 4 — Routing & HTTP Ergonomics
**Goal**: Make routing and request handling feel complete (Laravel-like DX).

**Key Items**:

| Feature | Priority | Implementation Notes |
|---------|----------|----------------------|
| Route Model Binding (implicit + explicit) | 🔴 | `routing/router.go`, new `binding.go` |
| Optional params `{id?}`, Regex constraints | 🔴 | Router parameter handling |
| URL Generation (`route()`, `url()`) | 🔴 | New `routing/url.go` + named routes |
| Request helpers (`Input()`, `Only()`, `Except()`, `Collect()`) | 🔴 | `http/request/request.go` |
| Redirect with flash + old input + named routes | 🔴 | `http/response/redirect.go` + session |
| Fallback routes + View routes + Redirect routes | 🟡 | Router sugar methods |
| Content negotiation (`WantsJson()`, `Accepts()`) | 🔴 | `http/request/request.go` |

**Estimated Effort**: 2–3 weeks

---

### Wave 5 — Views & Validation Completion
**Goal**: Make Blade templates and validation production-ready.

**Key Items**:

| Category | High Priority Items | Files |
|----------|---------------------|-------|
| Blade | `$loop` variable, `{!! !!}`, `@method`, `@error`, `@auth/@guest` | `view/compiler.go`, `view/engine.go` |
| Validation | 15+ missing rules (`boolean`, `date`, `required_if`, `uuid`, etc.) + nested arrays | `validation/validator.go` |
| View Sharing + Components | Global data, custom components/slots | `view/composer.go` + new component system |

**Estimated Effort**: 2–3 weeks

---

### Wave 6 — Authentication & Authorization
**Goal**: Complete the auth system (currently only Sanctum tokens exist).

**Key Items**:

| Feature | Priority | Notes |
|---------|----------|-------|
| SessionGuard (real session auth) | 🔴 | New `auth/session_guard.go` |
| Password Reset flow | 🔴 | Tokens, emails, broker |
| Email Verification | 🔴 | Signed URLs + middleware |
| Sanctum fixes (expiry + abilities) | 🔴 | `auth/sanctum/` |
| Policies + Gate | 🟡 | `auth/policy.go` + `Gate` facade alternative |

**Estimated Effort**: 3 weeks

---

### Wave 7 — Communication Layer (Mail, Queue, Notifications, Broadcasting)
**Goal**: Make background jobs, emails, and real-time features actually work.

**Key Items**:

| Area | Must-Have | Files |
|------|-----------|-------|
| Queue | Delayed jobs, Failed jobs table, Retry, Backoff, Job middleware | `queue/` package |
| Mail | Real SMTP driver, Attachments, CC/BCC, Markdown mails | `mail/` + new `mail/smtp.go` |
| Notifications | Mail channel + Database channel | `notifications/` |
| Broadcasting | Full WebSocket + Pusher/Redis drivers + channel auth | Already started in Wave 2 — complete it |

**Estimated Effort**: 4 weeks

---

### Wave 8 — Production Polish & Ecosystem
**Goal**: CLI commands, middleware, testing utilities, storage, logging, config, etc.

**Major Buckets**:

- **CLI / Artisan**: `migrate:fresh`, `route:list`, `serve`, `key:generate`, `tinker`, `make:*` generators
- **Middleware**: Groups, aliases, parameters, more built-in middleware
- **Storage**: Real S3 driver, `UploadedFile`, `storage:link`
- **Testing**: `RefreshDatabase`, `ActingAs`, `AssertDatabaseHas`, Mail/Queue fakes
- **Cache/Session**: File driver, Database driver, `Remember()`, Previous URL
- **Error Handling & Logging**: Custom error pages, Daily rotation, Structured logging
- **Ecosystem**: First-party packages (Socialite, Scout, Telescope, etc.) — lower priority

**Estimated Effort**: 6–8 weeks (can be parallelized)

---

## Tracking Table (Overall Progress)

| Wave | Focus Area | Items | Status | Target Completion |
|------|------------|-------|--------|-------------------|
| Wave 3 | Core Data Layer (Transactions, ORM, Pagination) | ~25 | 📋 Planned | — |
| Wave 4 | Routing + HTTP | ~20 | 📋 Planned | — |
| Wave 5 | Views + Validation | ~30 | 📋 Planned | — |
| Wave 6 | Authentication | ~15 | 📋 Planned | — |
| Wave 7 | Mail / Queue / Notifications / Broadcasting | ~35 | 📋 Planned | — |
| Wave 8 | CLI + Middleware + Storage + Testing + Polish | ~125 | 📋 Planned | — |

---

## Recommended Approach

1. **Start with Wave 3** — it unblocks almost everything else.
2. For every feature, follow the pattern used in previous waves:
   - Add or extend interfaces first
   - Implement in the relevant package
   - Write tests immediately
   - Update documentation
3. Keep a running `kilocode_changes.md` or similar for each wave.
4. After each wave, update this file and `SUGGESTIONS.md`.

---

## Verification Commands (Per Wave)

```bash
go build ./...
go test ./... -count=1
go test -race ./...
```

For database-heavy waves, also run with SQLite in-memory + real MySQL/Postgres when possible.

---

## Next Steps

- [ ] Confirm status of the original 6 Critical Bugs (quick audit)
- [ ] Finalize Wave 3 detailed implementation plan (similar to `Wave_2_Implementation_Plan.md`)
- [ ] Begin Wave 3

---

**Bottom Line**

`MISSING_FEATURES.md` is no longer just a scary list — it is now the **official backlog**. With a structured wave-based approach, GoW can realistically reach feature parity with Laravel in the critical areas within 4–6 months of focused development.

This document will be the single source of truth for implementation progress.
