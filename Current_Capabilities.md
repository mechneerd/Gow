# GoW Framework — Current Capabilities (vs Laravel)

> **Date**: 2026-05-23 (Major Update - Feature Parity Release)  
> **Framework Version**: Post Wave 3 + Wave 4 (Feature Parity)  
> **Total Go source files**: 180+ (massive expansion in ORM, Auth, CLI, and ecosystem foundations)

**May 23, 2026 Feature Parity Release** — This document has been significantly updated after the largest implementation wave in GoW history. Most remaining critical gaps have been closed. See `docs/whats_new/2026-05-23.md` for full release details.

This document provides a clear, honest comparison of what **GoW currently has implemented** versus Laravel equivalents.

**Legend**:
- ✅ **Implemented** — Feature works and is tested/usable
- 🟡 **Partial** — Core exists but incomplete or has limitations
- ❌ **Missing** — Not implemented yet

---

## 1. Foundation & Architecture

**Status: ✅ Fully Implemented**

### Implementation Delivered (2026-05-22)
- Created canonical `bootstrap/app.go` — standard way to create and boot an Application with all core providers.
- Wired real `vendor:publish` command (`cmd/artisan/vendor_commands.go`) to use `PublishableProvider` + `ProviderRegistry` + `PublishAssets()`.
- Updated Console Kernel to properly inject `*foundation.Application` into cobra context (so publish and future commands can access the booted app).
- Refactored `artisan.go` to use the new bootstrap helper.
- All Foundation interfaces (`ServiceProvider`, `PublishableProvider`, `ProviderRegistry`, discovery & publishing) are now production-wired.

| Feature                        | Laravel Equivalent          | GoW Status | GoW Package / Notes |
|--------------------------------|-----------------------------|------------|---------------------|
| Dependency Injection Container | Service Container           | ✅ Full    | `container/container.go` — Bind, Singleton, Make, Resolve, contextual bindings, Freeze |
| Application Bootstrap          | Application                 | ✅ Full    | `foundation/application.go` + Boot lifecycle |
| Service Providers              | Service Providers           | ✅ Full    | `foundation/provider.go` — Base + `PublishableProvider` interface |
| Provider Discovery & Publishing| `vendor:publish` + auto-discovery | ✅ Full | `foundation/discovery.go` — `ProviderRegistry`, `PublishAssets()`, `DiscoverProviders()`, `PublishProvider()` on Application (explicit Go-style registration) |
| Configuration                  | Config                      | ✅ Full    | `config/repository.go` + `.env` loading + cached repository |

---

## 2. HTTP Layer

**Status: ✅ Significantly Improved (2026-05-22)**

### Implementation Delivered
- Extended `http/request/request.go` with many Laravel-style helpers: `Except()`, `Boolean()`, `Integer()`, `Float()`, `Collect()`, `ExpectsJson()`, `WantsJson()`, `Accepts()`, `Old()`.
- Added **Middleware Groups & Aliases** support to the router (`Alias()`, `GroupMiddleware()`, `Middleware()`).
- Implemented **Route Model Binding** (implicit + explicit):
  - `router.Bind("user", resolver)`
  - `routing.Model[T](req, "user")` generic helper
  - Automatic resolution during request dispatch with 404 on binding failure
- Content negotiation methods are now fully available.

| Feature                          | Laravel Equivalent               | GoW Status | Notes |
|----------------------------------|----------------------------------|------------|-------|
| Router (Radix tree)              | Route                            | ✅ Strong  | All HTTP verbs, Groups, Resource routes, Named routes |
| Route Model Binding              | Implicit + Explicit              | ✅ Good    | `Bind()`, `ResolveBinding()`, `routing.Model[T]()`, `routing.Binding()` — full implicit + explicit support via registered resolvers |
| Middleware (global + route)      | Middleware                       | ✅ Good    | 14+ middlewares implemented |
| Middleware Groups & Aliases      | `middlewareGroups`, aliases      | ✅ Good    | `Alias()`, `GroupMiddleware()`, `Middleware()` supported |
| Kernel + Graceful Shutdown       | `app/Http/Kernel.php`            | ✅ Full    | `http/kernel.go` — `Serve()` + `OnShutdown()` hooks (Wave 2) |
| Request Helpers (`Input`, `Only`, `old()`) | `$request->input()`, `old()` | ✅ Good    | `Input`, `Only`, `Except`, `Boolean`, `Integer`, `Float`, `Collect`, `Old` + more |
| Form Request                     | FormRequest                      | ✅ Good    | `http/request/form_request.go` |
| API Resources                    | API Resources                    | ✅ Good    | `http/resources/resource.go` |
| Exception Handling               | Exception Handler                | ✅ Good    | `http/exception` + `ErrorHandler` middleware |
| Content Negotiation              | `wantsJson()`, `accepts()`       | ✅ Good    | `ExpectsJson`, `WantsJson`, `Accepts` implemented |

---

## 3. Database & ORM (One of GoW's strongest areas)

| Feature                        | Laravel Equivalent                  | GoW Status | Notes |
|--------------------------------|-------------------------------------|------------|-------|
| Query Builder                  | Query Builder                       | ✅ Very Strong | Dialect-aware, all major methods |
| Eloquent ORM                   | Eloquent                            | ✅ Strong  | Models, relations, soft deletes, scopes, observers |
| Relationships (HasOne/HasMany/BelongsTo) | Eloquent Relations | ✅ Working | Eager loading fixed in earlier phase |
| BelongsToMany (pivot)          | BelongsToMany                       | ✅ Good    | Eager loading (`With()`) works. Full helpers: Attach, Detach, Sync, Toggle implemented. Pivot table support complete |
| Transactions                   | `DB::transaction()`                 | ✅ Good    | `db.Transaction(ctx, fn)`, `db.Begin(ctx)` / `BeginTx()`, `Commit()`, `Rollback()`, `InTransaction()` — full manual + callback support |
| Soft Deletes                   | SoftDeletes trait                   | ✅ Full    | Complete implementation |
| Global Scopes + Observers      | Global scopes + Observers           | ✅ Working | Fully functional |
| Mass Assignment protection     | $fillable / $guarded                | ✅ Good    | Full protection via MassAssignable interface (Fillable/Guarded methods) enforced in Insert/Update |
| Raw Expressions / GROUP BY / HAVING | Raw + aggregates | ✅ Good    | GroupBy(), Having(), OrHaving(), SelectRaw(), WhereRaw(), OrWhereRaw() supported |
| Schema Builder / Migrations    | Schema + Migrations                 | ✅ Good    | `database/schema`, migrator |
| Seeders & Factories            | Seeders + Factories                 | ✅ Present | Basic factory support |
| Pagination                     | Paginator                           | ✅ Good    | Full implementation: LengthAwarePaginator, SimplePaginator, CursorPaginator + Paginate(), SimplePaginate(), CursorPaginate() on ModelQuery |
| Database Dialects              | MySQL, Postgres, SQLite, SQL Server | ✅ Good    | SQLite + MySQL + PostgreSQL (Wave 2) |
| Connection Pooling             | Via config                          | ✅ Good    | Auto-config from env (Wave 2) |

---

## 4. Views & Templating

| Feature                     | Laravel Equivalent     | GoW Status | Notes |
|-----------------------------|------------------------|------------|-------|
| Blade-like Compiler         | Blade                  | ✅ Good    | `@if`, `@foreach`, `@extends`, `@include`, custom directives |
| Layouts (`@extends`)        | Layouts                | ✅ Working | Fixed in earlier phase |
| `$loop` variable            | `$loop`                | ✅ Good    | Full support: $loop.Index, $loop.Iteration, $loop.First, $loop.Last, $loop.Even, $loop.Odd, $loop.Remaining via mkloop helper |
| Components & Slots          | Blade Components       | ✅ Good    | Full support for <x-foo attr="val">slot content</x-foo> with attributes + slots passed to component views via "attributes" and "slot" |
| `{!! !!}` unescaped output  | `{!! !!}`              | ✅ Good    | Implemented via "raw" template func + compiler support |
| `@auth` / `@can` directives | Blade auth directives  | ✅ Good    | Full support: @auth / @guest / @can / @cannot / @canany with proper end tags and template functions |
| View Composers              | View Composers         | ✅ Present | `view/composer.go` |

---

## 5. Validation

| Feature                        | Laravel Equivalent          | GoW Status | Notes |
|--------------------------------|-----------------------------|------------|-------|
| Validator                      | Validator                   | ✅ Good    | ~14–16 rules + DB rules |
| Form Request Validation        | FormRequest                 | ✅ Good    | Integrated |
| Missing Rules (date, uuid, required_if, etc.) | 20+ rules | 🟡 Partial | Major progress: boolean, integer, array, url, uuid, alpha*, required_if, same/different, gt/gte/lt/lte, nullable, bail, date, before, after, json added |

---


## 6. Authentication & Authorization

| Feature                    | Laravel Equivalent              | GoW Status | Notes |
|----------------------------|---------------------------------|------------|-------|
| Sanctum (API tokens)       | Laravel Sanctum                 | ✅ Good    | Token creation/verification (expiry & abilities partial) |
| Session-based Auth         | Session Guard                   | ✅ Good    | Full SessionGuard + UserProvider + Authenticatable + ready-to-use Middleware (auth.Middleware) and GuestMiddleware. Login/Logout/Attempt fully functional |
| Password Reset             | Password Broker                 | ✅ Good    | Full flow: token generation, Markdown email, validation & reset via Fortify + password.Broker |
| Email Verification         | Email Verification              | ✅ Good    | Signed URL verification + Fortify handler + RequireVerified middleware implemented |
| Policies & Gate            | Policies + Gate                 | ✅ Good    | Full Gate with Define, Policy, Before/After hooks, Allows/Denies available in auth/access |
| `@can` directive           | `@can`                          | 🟡 Partial | Exists in views but limited |

---

## 7. Session, Cache, Queue, Broadcasting

| Feature                    | Laravel Equivalent         | GoW Status | Drivers |
|----------------------------|----------------------------|------------|---------|
| Session                    | Session                    | ✅ Good    | Cookie, File, Redis |
| Cache                      | Cache                      | ✅ Good    | Memory, Redis |
| Queue                      | Queue                      | ✅ Strong  | Sync, Memory, Redis, Database + Worker |
| Broadcasting               | Broadcasting + Reverb      | ✅ Good    | Native WebSocket Hub + driver (Wave 2) |
| Events                     | Events + Listeners         | ✅ Good    | Manager + listeners |

---

## 8. Communication (Mail, Notifications)

| Feature                    | Laravel Equivalent       | GoW Status | Notes |
|----------------------------|--------------------------|------------|-------|
| Mail Message Builder       | Mailable                 | ✅ Good    | `mail/message.go` |
| Mailer / Sending           | Mail facade              | ✅ Good    | Improved SmtpDriver (base64 attachments, encryption modes documented), better QueueNow helper, Markdown support, ServiceProvider |
| Notifications              | Notifications            | ✅ Good    | Manager + DatabaseChannel + new MailChannel + Notify() helper + ServiceProvider |
| Mail Queueing              | ShouldQueue              | ✅ Good    | Full ShouldQueue support via SendMailJob + Queue/QueueNow. Mailer supports automatic dispatch when queue manager is attached. |

---

## 9. Other Core Features

| Feature                    | Laravel Equivalent          | GoW Status | Notes |
|----------------------------|-----------------------------|------------|-------|
| Encryption                 | Encrypter                   | ✅ Good    | AES-256 |
| Hashing                    | Hash                        | ✅ Good    | Bcrypt + Argon2 |
| Logging                    | Logging                     | ✅ Good    | slog-based with levels |
| Storage (Filesystem)       | Storage                     | ✅ Good    | Local + Real S3 driver implemented using aws-sdk-go-v2 (Put, Get, Delete, Exists, URL, PresignedURL) |
| HTTP Client                | Http Client                 | ✅ Good    | Fluent client with Get/Post/Put/Patch/Delete, query params, tokens, retries, JSON/form, and basic faking |
| Process / Shell            | Process                     | ✅ Good    | `support/process` |
| Health Checks              | Health route                | ✅ Good    | `/up` endpoint + checkers (Wave 2) |
| Pipeline                   | Pipeline                    | ✅ Good    | `support/pipeline` |

---

## 10. CLI & Tooling

| Feature                    | Laravel Equivalent       | GoW Status | Notes |
|----------------------------|--------------------------|------------|-------|
| Artisan-like CLI           | Artisan                  | ✅ Good    | `cmd/artisan` + `gow` binary |
| Generators (`make:*`)      | make:* commands          | ✅ Several | Controller, Model, Migration, Middleware, Command, etc. |
| Migration Commands         | migrate, migrate:fresh   | ✅ Good    | Full support: migrate, migrate:fresh, migrate:refresh implemented on Migrator + artisan commands |
| Route listing              | route:list               | ✅ Good    | Full `route:list` command implemented (lists method, URI, name). Router exposes GetAllRoutes(). |
| Tinker / REPL              | tinker                   | ✅ Good    | Full REPL using Yaegi. Preloads App, DB, Config. Users can import any gow package at runtime. |

---

## 11. Testing Utilities

| Feature                    | Laravel Equivalent          | GoW Status | Notes |
|----------------------------|-----------------------------|------------|-------|
| TestCase base              | TestCase                    | ✅ Good    | `testing/testcase.go` |
| Database testing helpers   | RefreshDatabase, etc.       | ✅ Good    | RefreshDatabase + transaction support + AssertDatabaseHas/Missing/Count/Exactly implemented |
| Fakes                      | Mail::fake, Queue::fake     | ✅ Good    | MailFake, QueueFake, EventFake with AssertSent / AssertPushed / AssertDispatched |
| Assertions                 | AssertDatabaseHas, etc.     | ✅ Good    | Rich set: AssertDatabaseHas, Missing, Count, Exactly, HasNoRecords + HTTP assertions |

---

## Production Readiness Assessment (May 23, 2026 - Post Feature Parity Release)

| Area                        | Readiness      | Comments |
|-----------------------------|----------------|----------|
| **Core HTTP + Routing**     | High           | Excellent (Signed URLs, advanced rate limiting, full FormRequest, content negotiation) |
| **Database / ORM**          | Excellent      | One of the strongest Laravel-like ORMs in Go (Polymorphic, Casting, Accessors, Through, Locking, Chunk, Upsert, etc.) |
| **Queue + Broadcasting**    | High           | Native WebSocket + multiple queue drivers + worker |
| **Auth**                    | High           | Production-ready (Session + Sanctum + Gate + 2FA + Socialite + Lockout + Remember Me) |
| **Mail / Notifications**    | Medium         | Markdown + queueing + Slack channel; real SMTP still basic |
| **Views + Validation**      | Good           | Strong Blade-like experience + many validation rules |
| **Transactions & Advanced ORM** | Excellent   | Full feature set now present |
| **CLI & DX**                | High           | Very rich Artisan surface + excellent generators + REPL |
| **Testing**                 | Good           | Strong utilities (`actingAs`, uploads, assertions, fakes) |
| **Documentation**           | Good           | Significantly expanded with new guides and deep examples |
| **Overall Production**      | **High**       | Genuinely usable for real-world full-stack and API applications. Most major gaps closed. |

---

## Key Strengths (What GoW Does Well Today - May 23, 2026)

- Excellent generic container + service provider system
- **World-class ORM** — Polymorphic, Casting, Accessors/Mutators, Has*Through, Upsert, Pessimistic Locking, Chunk, Relation Touching, full BelongsToMany
- Full database transactions + advanced querying
- **Production-ready Authentication** — Session + Sanctum + Gate + 2FA + Socialite + Account Lockout + Remember Me + Email Verification
- Rich HTTP layer (Signed URLs, advanced rate limiting, full FormRequest, content negotiation)
- Native WebSocket broadcasting (unique advantage)
- Very strong CLI (dozens of Artisan commands + generators)
- Good testing experience (`actingAs`, file uploads, assertions)
- Foundation publishing + clean architecture
- Significantly improved documentation with deep examples

---

## Biggest Gaps Blocking Real-World Use (Post May 23, 2026)

Most major blockers have been eliminated. Remaining gaps are mostly polish or ecosystem:

1. Real production-grade **Mail** drivers (SMTP + better retry/DKIM)
2. More notification channels (beyond DB, Mail, Slack)
3. Richer error pages rendered from Blade views (stack traces in templates)
4. Full Inertia + Livewire experience (foundations exist)
5. Advanced ecosystem tools (Telescope UI, Horizon dashboard, etc.)
6. More complete starter kits / scaffolding

Most "day-to-day" features developers expect from a Laravel-like framework are now present and well documented.

---

## Recommended Next Focus (Post May 23, 2026 Feature Parity Release)

**Most high-impact items are now complete.** The May 23 wave closed the majority of the critical gaps.

**Completed in the May 23 Feature Parity wave:**
- Socialite (OAuth)
- 2FA + Account Lockout + Remember Me
- Advanced ORM (Polymorphic, Casting, Accessors, Through relations, Locking, Upsert, Chunk, Touching, etc.)
- Full Artisan command surface + generators
- Signed URLs, advanced rate limiting, complete FormRequest
- Testing improvements (`actingAs`, uploads, assertions)
- Documentation overhaul + deep examples

**Recommended next focus (polish & ecosystem):**
1. Production-grade Mail drivers + better retry handling
2. More notification channels (SMS, etc.)
3. Richer Blade error pages with stack traces
4. Full Inertia + Livewire experience
5. Advanced debugging / monitoring tools (Telescope UI, Horizon)
6. Official starter kits / scaffolding

---

**Bottom Line (May 23, 2026)**

After the massive Feature Parity release, GoW has reached a new level of maturity.

The framework now includes:
- One of the most complete ORMs available in Go (near full Eloquent parity)
- Production-ready Authentication stack (Socialite + 2FA + Lockout + Sanctum + Gate)
- Very strong CLI and Developer Experience
- Solid foundations for most ecosystem features

**GoW is now genuinely production-viable** for building real-world applications (SaaS, admin panels, APIs, full-stack apps).

This document is kept up to date after each major implementation wave.

**See also**: `docs/whats_new/2026-05-23.md` for the full release notes.
