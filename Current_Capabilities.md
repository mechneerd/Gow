# GoW Framework — Current Capabilities (vs Laravel)

> **Date**: 2026-05-22 (Updated)  
> **Framework Version**: Post Wave 2 + Foundation Hardening  
> **Total Go source files**: ~157 (new bootstrap package)

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

| Feature                          | Laravel Equivalent               | GoW Status | Notes |
|----------------------------------|----------------------------------|------------|-------|
| Router (Radix tree)              | Route                            | ✅ Strong  | All HTTP verbs, Groups, Resource routes, Named routes |
| Route Model Binding              | Implicit + Explicit              | ❌ Missing | — |
| Middleware (global + route)      | Middleware                       | ✅ Good    | 14+ middlewares implemented |
| Middleware Groups & Aliases      | `middlewareGroups`, aliases      | ❌ Missing | — |
| Kernel + Graceful Shutdown       | `app/Http/Kernel.php`            | ✅ Full    | `http/kernel.go` — `Serve()` + `OnShutdown()` hooks (Wave 2) |
| Request Helpers (`Input`, `Only`, `old()`) | `$request->input()`, `old()` | 🟡 Partial | Basic request object; many helpers still missing |
| Form Request                     | FormRequest                      | ✅ Good    | `http/request/form_request.go` |
| API Resources                    | API Resources                    | ✅ Good    | `http/resources/resource.go` |
| Exception Handling               | Exception Handler                | ✅ Good    | `http/exception` + `ErrorHandler` middleware |
| Content Negotiation              | `wantsJson()`, `accepts()`       | ❌ Missing | — |

---

## 3. Database & ORM (One of GoW's strongest areas)

| Feature                        | Laravel Equivalent                  | GoW Status | Notes |
|--------------------------------|-------------------------------------|------------|-------|
| Query Builder                  | Query Builder                       | ✅ Very Strong | Dialect-aware, all major methods |
| Eloquent ORM                   | Eloquent                            | ✅ Strong  | Models, relations, soft deletes, scopes, observers |
| Relationships (HasOne/HasMany/BelongsTo) | Eloquent Relations | ✅ Working | Eager loading fixed in earlier phase |
| BelongsToMany (pivot)          | BelongsToMany                       | ❌ Missing | — |
| Transactions                   | `DB::transaction()`                 | ❌ Missing | High priority for Wave 3 |
| Soft Deletes                   | SoftDeletes trait                   | ✅ Full    | Complete implementation |
| Global Scopes + Observers      | Global scopes + Observers           | ✅ Working | Fully functional |
| Raw Expressions / GROUP BY / HAVING | Raw + aggregates | 🟡 Partial | Basic raw exists, advanced missing |
| Schema Builder / Migrations    | Schema + Migrations                 | ✅ Good    | `database/schema`, migrator |
| Seeders & Factories            | Seeders + Factories                 | ✅ Present | Basic factory support |
| Pagination                     | Paginator                           | 🟡 Partial | Basic `Paginate()` works; cursor/simple missing |
| Database Dialects              | MySQL, Postgres, SQLite, SQL Server | ✅ Good    | SQLite + MySQL + PostgreSQL (Wave 2) |
| Connection Pooling             | Via config                          | ✅ Good    | Auto-config from env (Wave 2) |

---

## 4. Views & Templating

| Feature                     | Laravel Equivalent     | GoW Status | Notes |
|-----------------------------|------------------------|------------|-------|
| Blade-like Compiler         | Blade                  | ✅ Good    | `@if`, `@foreach`, `@extends`, `@include`, custom directives |
| Layouts (`@extends`)        | Layouts                | ✅ Working | Fixed in earlier phase |
| `$loop` variable            | `$loop`                | 🟡 Partial | Not fully implemented |
| Components & Slots          | Blade Components       | 🟡 Partial | Basic support exists |
| `{!! !!}` unescaped output  | `{!! !!}`              | ❌ Missing | — |
| `@auth` / `@can` directives | Blade auth directives  | 🟡 Partial | `@can` exists, full auth directives incomplete |
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
| Session-based Auth         | Session Guard                   | ❌ Missing | Only token auth currently works |
| Password Reset             | Password Broker                 | ❌ Missing | — |
| Email Verification         | Email Verification              | ❌ Missing | — |
| Policies & Gate            | Policies + Gate                 | ❌ Missing | — |
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
| Mailer / Sending           | Mail facade              | 🟡 Partial | Interface exists, real SMTP driver missing |
| Notifications              | Notifications            | 🟡 Partial | Manager + Database channel started |
| Mail Queueing              | ShouldQueue              | ❌ Missing | — |

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
| Route listing              | route:list               | ❌ Missing | — |
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

## Production Readiness Assessment (Post Wave 2)

| Area                        | Readiness | Comments |
|-----------------------------|-----------|----------|
| **Core HTTP + Routing**     | High      | Very usable with graceful shutdown |
| **Database / ORM**          | High      | One of the best parts of GoW right now |
| **Queue + Broadcasting**    | Medium-High | Native Go advantages |
| **Auth**                    | Low-Medium | Only API tokens work reliably |
| **Mail / Notifications**    | Low       | Needs real drivers |
| **Views + Validation**      | Medium    | Usable but incomplete |
| **Transactions & Advanced ORM** | Low    | Blocker for many real apps |
| **CLI & DX**                | Medium    | Good generators, missing many artisan commands |
| **Overall Production**      | Medium    | Can be used for APIs and background jobs today. Not yet for full-stack apps with auth + mail. |

---

## Key Strengths (What GoW Does Well Today)

- Excellent generic container + service provider system
- Strong, dialect-aware Query Builder + ORM with real relationships
- Graceful shutdown + connection pooling (production polish)
- Native WebSocket broadcasting (unique Go advantage)
- Multiple queue drivers + worker
- Clean middleware and routing architecture
- Good testing foundation

---

## Biggest Gaps Blocking Real-World Use

1. No real **transactions**
2. No **session-based authentication**
3. No **Route Model Binding**
4. Missing **Mail** transport (SMTP)
5. Many advanced ORM features (BelongsToMany, Mass Assignment enforcement, etc.)
6. Incomplete **Views** (`{!! !!}`, `$loop`, components)

---

## Recommended Next Focus (Wave 3)

From the assessment above, the highest-impact areas for the next wave are:

1. Database Transactions + Savepoints
2. Session-based Auth (SessionGuard)
3. Route Model Binding
4. Real Mail SMTP driver
5. BelongsToMany relationships + pivot tables
6. Complete Mass Assignment protection

---

**Bottom Line**

GoW has moved from "promising skeleton" to **"actually usable framework for APIs and background processing"** after completing Top 10 + Wave 2.

It is not yet a drop-in Laravel replacement for full-stack applications, but the foundation is now solid enough that building real features on top of it makes sense.

This document will be kept up to date after each implementation wave.
