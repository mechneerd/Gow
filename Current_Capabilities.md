# GoW Framework — Current Capabilities (vs Laravel)

> **Date**: 2026-05-22 (Updated)  
> **Framework Version**: Post Wave 2 + Major Feature Waves (ORM, Auth, HTTP, Foundation)  
> **Total Go source files**: ~160+ (includes auth middleware, orm helpers, raw query improvements)

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
| Missing Rules (date, uuid, required_if, etc.) | 20+ rules | ❌ Many    | See `MISSING_FEATURES.md` |

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
| Notifications              | Notifications            | 🟡 Partial | Manager + Database channel started |
| Mail Queueing              | ShouldQueue              | ✅ Good    | Full ShouldQueue support via SendMailJob + Queue/QueueNow. Mailer supports automatic dispatch when queue manager is attached. |

---

## 9. Other Core Features

| Feature                    | Laravel Equivalent          | GoW Status | Notes |
|----------------------------|-----------------------------|------------|-------|
| Encryption                 | Encrypter                   | ✅ Good    | AES-256 |
| Hashing                    | Hash                        | ✅ Good    | Bcrypt + Argon2 |
| Logging                    | Logging                     | ✅ Good    | slog-based with levels |
| Storage (Filesystem)       | Storage                     | 🟡 Partial | Local works; S3 is still no-op |
| HTTP Client                | Http Client                 | 🟡 Partial | Basic client in `support/httpclient` |
| Process / Shell            | Process                     | ✅ Good    | `support/process` |
| Health Checks              | Health route                | ✅ Good    | `/up` endpoint + checkers (Wave 2) |
| Pipeline                   | Pipeline                    | ✅ Good    | `support/pipeline` |

---

## 10. CLI & Tooling

| Feature                    | Laravel Equivalent       | GoW Status | Notes |
|----------------------------|--------------------------|------------|-------|
| Artisan-like CLI           | Artisan                  | ✅ Good    | `cmd/artisan` + `gow` binary |
| Generators (`make:*`)      | make:* commands          | ✅ Several | Controller, Model, Migration, Middleware, Command, etc. |
| Migration Commands         | migrate, migrate:fresh   | 🟡 Partial | Basic migrate works; fresh/refresh missing |
| Route listing              | route:list               | ✅ Good    | Full `route:list` command implemented (lists method, URI, name). Router exposes GetAllRoutes(). |
| Tinker / REPL              | tinker                   | ❌ Missing | — |

---

## 11. Testing Utilities

| Feature                    | Laravel Equivalent          | GoW Status | Notes |
|----------------------------|-----------------------------|------------|-------|
| TestCase base              | TestCase                    | ✅ Good    | `testing/testcase.go` |
| Database testing helpers   | RefreshDatabase, etc.       | 🟡 Partial | Basic DB helpers exist |
| Fakes                      | Mail::fake, Queue::fake     | 🟡 Partial | Some fakes in `support/fakes` |
| Assertions                 | AssertDatabaseHas, etc.     | ❌ Missing | — |

---

## Production Readiness Assessment (May 2026)

| Area                        | Readiness | Comments |
|-----------------------------|-----------|----------|
| **Core HTTP + Routing**     | High      | Very usable with graceful shutdown + rich helpers |
| **Database / ORM**          | High      | One of the strongest areas (Transactions, BelongsToMany, multi-dialect) |
| **Queue + Broadcasting**    | Medium-High | Native Go advantages |
| **Auth**                    | Medium-High | Session Auth + Gate now functional + Sanctum |
| **Mail / Notifications**    | Low       | Needs real drivers |
| **Views + Validation**      | Medium    | Usable but incomplete |
| **Transactions & Advanced ORM** | High   | Now well supported |
| **CLI & DX**                | Medium    | Good generators, missing some artisan commands |
| **Overall Production**      | Medium-High | Can be used for APIs, background jobs, and many full-stack features. Mail is the main remaining gap for complete apps. |

---

## Key Strengths (What GoW Does Well Today)

- Excellent generic container + service provider system
- Strong, dialect-aware Query Builder + ORM (Transactions, BelongsToMany + pivot helpers, scopes, observers, soft deletes)
- Full database transactions with manual + callback API
- Session-based Authentication + Gate/Policies (production ready)
- Route Model Binding (implicit + explicit)
- Rich HTTP layer (Request helpers, Middleware groups/aliases, Content negotiation)
- Native WebSocket broadcasting (unique Go advantage)
- Multiple queue drivers + worker
- Graceful shutdown + connection pooling
- Clean middleware and routing architecture
- Good testing foundation
- Foundation publishing system (`vendor:publish`)

---

## Biggest Gaps Blocking Real-World Use

1. Missing **Mail** transport (real SMTP driver + Markdown mails)
2. Password Reset flow
3. Email Verification system
4. (Advanced Pagination now has solid foundation — LengthAware + Simple + Cursor)
5. Mass Assignment protection enforcement (`$fillable` / `$guarded`)
6. More complete View system (`{!! !!}` raw output, full `$loop` variable, proper components/slots)

---

## Recommended Next Focus (Post May 2026)

Many high-impact items from the original list have now been completed:

**Completed in recent waves:**
- Database Transactions (full Begin/Commit/Rollback + callback)
- Session-based Authentication (SessionGuard + Middleware)
- Route Model Binding (implicit + explicit with `Model[T]()`)
- BelongsToMany (pivot) + Attach/Detach/Sync/Toggle
- Raw Expressions, GROUP BY, HAVING, SelectRaw/WhereRaw
- Foundation & Architecture (bootstrap + vendor:publish)
- HTTP Layer improvements (Request helpers, Middleware groups, Content negotiation)

**Remaining high-priority items for next phase:**
1. Real Mail SMTP driver + full Markdown mail polish
2. More complete View components (`{!! !!}`, `$loop`, Slots)
3. CLI improvements (more artisan commands)

---

**Bottom Line**

GoW has matured significantly. After Top 10 + Wave 2 + multiple focused waves, the framework now includes:

- Full-featured ORM (Transactions, BelongsToMany + helpers, Mass Assignment protection, Pagination)
- Production-ready Authentication (Session Guard + Middleware, Password Reset, Email Verification, Gate)
- Solid Mail system (Markdown + basic queueing)
- Excellent HTTP layer and Foundation

The framework is now genuinely usable for building real-world applications.

This document is kept up to date after each implementation wave.
