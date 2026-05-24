# GoW Framework — Complete Laravel Gap Analysis (v2)

> **Methodology**: Every Go source file was read. "Implemented" means working code with real logic. "Skeletal" means an interface/struct exists but the body is a stub, a TODO, or trivially incomplete. "Missing" means no code exists at all.
>
> **v2 update**: Added ~80 additional core framework gaps and 17 first-party ecosystem packages not covered in v1.

---

## Executive Summary

| Category | Implemented | Skeletal / Stub | Missing (no code) |
|---|---|---|---|
| Core Architecture (Container, App, Kernel, Config, Logging, CLI) | ✅ Solid | ⚠️ Advanced container features | — |
| Routing | ✅ Basic | ⚠️ Many gaps | — |
| Middleware | ✅ Good set | ⚠️ Middleware params, trusted proxies | — |
| ORM / Database | ✅ Foundation | ⚠️ Heavy gaps | — |
| Views / Templates | ✅ Basic compiler | ⚠️ Advanced directives | — |
| Validation | ⚠️ Very thin (2 rules) | — | — |
| Auth & Security | ✅ Interfaces | ⚠️ Session guard body | — |
| Cache | ✅ Memory driver | ⚠️ Only 1 driver | — |
| Session | ✅ 2 drivers | ⚠️ Missing encryption/flash | — |
| Queue | ✅ Sync driver | ⚠️ No DB/Redis driver | — |
| Events | ✅ Sync dispatch | ⚠️ No queued/stoppable | — |
| Mail / Notifications / Broadcasting | ⚠️ Interfaces only | — | — |
| Sections 35–50 (API Resources, Collections, etc.) | ❌ None | ❌ None | ❌ All 16 sections |
| Additional Core Gaps (v2) | ❌ None | ❌ None | ❌ ~80 items |
| First-Party Ecosystem (v2) | ❌ None | ❌ None | ❌ All 17 packages |

**Total gap count: ~350+ individual missing features/items.**

---

## PART I — CORE FRAMEWORK GAPS

---

## 1. Routing — [router.go](file:///d:/Go%20framework/routing/router.go) (228 lines)

### What's implemented
- Radix-tree route matching (GET, POST only via helpers)
- Route groups with prefix
- Named routes
- Path parameters (`{id}`)
- Middleware application per-route
- `ServeHTTP` implementing `http.Handler`

### What's missing vs Laravel

| Laravel Feature | Status | Notes |
|---|---|---|
| PUT, PATCH, DELETE, OPTIONS, HEAD helpers | ❌ Missing | Only `Get()` and `Post()` exist |
| `Any()` / `Match()` multi-method routes | ❌ Missing | |
| Optional parameters `{id?}` | ❌ Missing | |
| Regex constraints on params | ❌ Missing | `where('id', '[0-9]+')` |
| Route model binding (implicit/explicit) | ❌ Missing | Spec §39 — no code |
| Implicit enum binding | ❌ Missing | Bind Go enum types to route params |
| Resource/API resource routes | ❌ Missing | `router.Resource("users", ctrl)` |
| Singleton resource routes | ❌ Missing | `Route::singleton()` |
| Redirect routes / View routes | ❌ Missing | `Route::redirect()`, `Route::view()` |
| Fallback route (`Route::fallback`) | ❌ Missing | |
| Subdomain routing | ❌ Missing | |
| Route caching / compilation | ❌ Missing | |
| Route listing CLI (`route:list`) | ❌ Missing | |
| URL/Redirect full API | ❌ Missing | `URL::current()`, `redirect()->route()` |
| URL generation from named routes | ❌ Missing | Named routes stored but no `Generate()` |
| Rate limiting per route (named limiters) | ❌ Missing | `RateLimiter::for('api', ...)` |
| Signed URL routes | ❌ Missing | Spec §38 — no code |
| Controller middleware with `except`/`only` | ❌ Missing | |
| Route `missing()` callback | ❌ Missing | Custom 404 per binding |

---

## 2. Container & Service Providers

### Container — [container.go](file:///d:/Go%20framework/container/container.go)

| Feature | Status |
|---|---|
| Bind / Singleton / Instance / Make | ✅ |
| Generics `Make[T]` | ✅ |
| Freeze (prevent mutation) | ✅ |
| Contextual binding | ❌ Missing | `when(PhotoController)->needs(Filesystem)->give(S3Filesystem)` |
| Container tagging | ❌ Missing | `tag([Report1, Report2], 'reports')` |
| Container extending | ❌ Missing | `extend(Service, fn($service) => decorate($service))` |
| Rebinding callbacks | ❌ Missing | |
| Method injection (resolve method params) | ❌ Missing | |
| Facades pattern | ❌ Missing | Static-like access to container services (e.g., `Cache::get()`) |
| Real-time facades | ❌ Missing | Auto-generate facade from any class |
| Scoped bindings (request-scoped) | ❌ Missing | |

### Service Providers — [foundation/provider.go](file:///d:/Go%20framework/foundation/provider.go)

| Feature | Status |
|---|---|
| Register / Boot lifecycle | ✅ |
| Deferred providers | ❌ Missing | `IsDeferred()` + `Provides()` |
| Provider manifest (cached) | ❌ Missing | |
| Package service providers | ❌ Missing | |
| `Publishable` assets/configs/migrations | ❌ Missing | |

---

## 3. ORM / Database Layer

### Query Builder — [builder.go](file:///d:/Go%20framework/database/query/builder.go) (122 lines)

| Laravel Feature | Status | Notes |
|---|---|---|
| SELECT, INSERT, UPDATE, DELETE | ✅ | Basic working versions |
| WHERE, OR WHERE | ✅ | |
| ORDER BY, LIMIT, OFFSET | ✅ | |
| JOIN (inner, left, right, cross) | ❌ Missing | |
| WHERE IN, WHERE NOT IN, WHERE NULL | ❌ Missing | |
| WHERE BETWEEN, WHERE EXISTS | ❌ Missing | |
| WHERE date/day/month/year/time | ❌ Missing | |
| WHERE JSON (`->`) column access | ❌ Missing | |
| GROUP BY / HAVING | ❌ Missing | |
| UNION / UNION ALL | ❌ Missing | |
| Subqueries (where, select, from) | ❌ Missing | |
| Aggregates (count, sum, avg, min, max) | ❌ Missing | |
| Distinct | ❌ Missing | |
| Chunk / chunkById | ❌ Missing | |
| Lazy / Cursor iteration | ❌ Missing | |
| Increment / Decrement | ❌ Missing | |
| Upsert / insertOrIgnore / insertGetId | ❌ Missing | |
| Raw expressions (`DB::raw`) | ❌ Missing | |
| Transaction support (begin/commit/rollback) | ❌ Missing | |
| Savepoints | ❌ Missing | |
| Database locking (shared/exclusive `lockForUpdate`) | ❌ Missing | |
| Multiple database connections | ⚠️ Skeletal | Manager exists but no real switching |
| Database sticky connections | ❌ Missing | Read replica but write to master |
| Read/write connection splitting | ❌ Missing | |
| `when()` conditional clauses | ❌ Missing | |
| `explain()` | ❌ Missing | |
| Fulltext indexes | ❌ Missing | |
| Spatial indexes | ❌ Missing | |
| JSON columns support | ❌ Missing | |
| Virtual/Stored generated columns | ❌ Missing | |

### ORM (Goquent) — [query.go](file:///d:/Go%20framework/database/orm/query.go) (200 lines)

| Laravel Feature | Status | Notes |
|---|---|---|
| Model struct with tags | ✅ | `db:` and `gow:` tags |
| Auto timestamps | ✅ | `autoCreateTime` / `autoUpdateTime` |
| Generic query `ModelQuery[T]` | ✅ | |
| `First()`, `Get()`, `Insert()` | ✅ | |
| `Find()` / `FindOrFail()` by ID | ❌ Missing | |
| `Update()` on model | ❌ Missing | |
| `Delete()` / `Destroy()` | ❌ Missing | |
| `Save()` (insert or update) | ❌ Missing | |
| `FirstOrCreate()` / `FirstOrNew()` / `UpdateOrCreate()` | ❌ Missing | |
| Soft deletes (`deleted_at`) | ❌ Missing | Tag recognized but no query filtering |
| `withTrashed()` / `onlyTrashed()` / `restore()` | ❌ Missing | |
| Mass assignment (fillable/guarded) | ❌ Missing | Tags exist but not enforced |
| Mutators & Accessors (new `Attribute` class style) | ❌ Missing | |
| Attribute casting (dates, JSON, enum, encrypted, custom) | ❌ Missing | |
| Enum casting | ❌ Missing | `AsStringable`, `AsCollection`, `AsArrayObject` |
| Custom casts (`CastsAttributes` interface) | ❌ Missing | |
| Encrypted casts | ❌ Missing | |
| Scopes (local & global) | ❌ Missing | |
| Model events/observers (creating, created, updating, etc.) | ❌ Missing | |
| Hidden/visible attributes | ❌ Missing | `makeVisible()`, `makeHidden()` |
| Appended attributes (`$appends`) | ❌ Missing | |
| Default attribute values | ❌ Missing | `$attributes = [...]` |
| Dirty tracking (`isDirty()`, `wasChanged()`) | ❌ Missing | |
| Touching parent timestamps | ❌ Missing | |
| `toJSON()` / `toArray()` serialization | ❌ Missing | |
| Model replicating (`replicate()`) | ❌ Missing | |
| Model comparing (`is()` / `isNot()`) | ❌ Missing | |
| `HasUuids` / `HasUlids` traits | ❌ Missing | UUID/ULID as primary key |
| Prevent lazy loading (N+1 protection) | ❌ Missing | `Model::preventLazyLoading()` |
| Model strictness (`preventSilentlyDiscardingAttributes`, `preventAccessingMissingAttributes`) | ❌ Missing | |
| Model factories (real implementation) | ⚠️ Skeletal | Interface only, no faker integration |
| Factory sequences, states, relationships | ❌ Missing | |
| Model pruning (`Prunable` / `MassPrunable`) | ❌ Missing | Spec §42 — no code |
| Broadcasting model events | ❌ Missing | |

### Relationships — [relation.go](file:///d:/Go%20framework/database/orm/relation.go) (45 lines)

| Laravel Feature | Status | Notes |
|---|---|---|
| HasOne | ❌ **Stub** | `loadRelation()` returns nil |
| HasMany | ❌ **Stub** | Same |
| BelongsTo | ❌ **Stub** | Same |
| BelongsToMany (pivot) | ❌ Missing | |
| Pivot table operations (attach, detach, sync, toggle, updateExistingPivot) | ❌ Missing | |
| HasOneThrough / HasManyThrough | ❌ Missing | |
| MorphOne / MorphMany / MorphTo / MorphToMany / MorphedByMany | ❌ Missing | |
| Default models for relationships | ❌ Missing | `withDefault()` |
| Relationship counting (`withCount`, `withSum`, `withAvg`, `withMin`, `withMax`) | ❌ Missing | |
| Relationship existence queries (`has`, `whereHas`, `doesntHave`) | ❌ Missing | |
| Constraining eager loads | ❌ Missing | `with(['posts' => fn($q) => $q->where()])` |
| Lazy eager loading (`load()`, `loadMissing()`) | ❌ Missing | |
| Nested eager loading (`with('author.contacts')`) | ❌ Missing | |

> [!CAUTION]
> Eager loading (`With()`) calls `eagerLoadRelationships` which **does nothing** — it's a stub that returns `nil`. This is the single biggest gap in the ORM.

### Pagination — exists in `database/pagination/`

| Feature | Status |
|---|---|
| Length-aware paginator | ⚠️ Basic struct |
| Simple pagination (no total count) | ❌ Missing |
| Cursor pagination | ❌ Missing |
| Manual paginator creation | ❌ Missing |
| View integration for HTML links | ❌ Missing |
| JSON response meta/links | ❌ Missing |

### Migration System — [migrator.go](file:///d:/Go%20framework/database/migration/migrator.go) (4363 bytes)
- ✅ Runs Up/Down, tracks in `migrations` table
- ❌ Missing: `migrate:fresh`, `migrate:refresh`, `migrate:status`, `migrate:reset` commands
- ❌ Missing: Schema dumping (`schema:dump` / `schema:load`)
- ❌ Missing: Squashing migrations
- ❌ Missing: Anonymous migrations

---

## 4. Views / Template Engine

### Goblade Compiler — [compiler.go](file:///d:/Go%20framework/view/compiler.go) (100 lines)

| Feature | Status |
|---|---|
| `@if` / `@elseif` / `@else` / `@endif` | ✅ |
| `@foreach` / `@endforeach` | ✅ |
| `@yield` / `@section` / `@endsection` | ✅ |
| `@include` | ✅ |
| `@csrf` | ✅ |
| File compilation + SHA caching | ✅ |
| `@extends` (layout inheritance) | ❌ Missing | Only a comment — not actually resolved |
| `@for` / `@while` / `@switch` | ❌ Missing |
| `$loop` variable (index, first, last, etc.) | ❌ Missing |
| Components & Slots (anonymous + class-based) | ❌ Missing |
| Stacks (`@push` / `@stack` / `@prepend`) | ❌ Missing |
| Custom directives registration | ❌ Missing |
| `@auth` / `@guest` / `@can` / `@cannot` directives | ❌ Missing |
| `@once`, `@fragment`, `@aware`, `@inject` | ❌ Missing |
| `@checked`, `@selected`, `@disabled`, `@readonly` | ❌ Missing |
| `@class`, `@style` conditional helpers | ❌ Missing |
| `@env` / `@production` directives | ❌ Missing |
| `@method` (form method spoofing) | ❌ Missing |
| `@error` directive | ❌ Missing |
| `@props` directive | ❌ Missing |
| `{!! !!}` unescaped output | ❌ Missing |
| `{{-- comments --}}` | ❌ Missing |
| View sharing (global data) | ❌ Missing |
| View first (`view()->first()`) | ❌ Missing |
| Inline Blade templates | ❌ Missing |

### View Composers — [composer.go](file:///d:/Go%20framework/view/composer.go) (1059 bytes)
- ✅ Basic `ComposerRegistry` exists
- ❌ Missing: View creators (run on create, not render)

---

## 5. Validation — [validator.go](file:///d:/Go%20framework/validation/validator.go) (65 lines)

> [!WARNING]
> Only **2 rules** are implemented: `required` and `email` (via a naive `strings.Contains(@)` check).

| Laravel Rule | Status |
|---|---|
| required | ✅ (basic) |
| email | ⚠️ (just checks for `@`) |
| string, integer, numeric, boolean, array | ❌ Missing |
| min, max (string length, number value) | ❌ Missing |
| between | ❌ Missing |
| size | ❌ Missing |
| confirmed (field_confirmation) | ❌ Missing |
| unique / exists (database) | ❌ Missing |
| in / not_in | ❌ Missing |
| regex / not_regex | ❌ Missing |
| date, before, after, date_format, date_equals | ❌ Missing |
| image, mimes, mimetypes, file, dimensions | ❌ Missing |
| url, active_url, ip, ipv4, ipv6 | ❌ Missing |
| json, uuid, ulid | ❌ Missing |
| starts_with, ends_with, contains | ❌ Missing |
| alpha, alpha_num, alpha_dash | ❌ Missing |
| digits, digits_between | ❌ Missing |
| gt, gte, lt, lte | ❌ Missing |
| same, different | ❌ Missing |
| required_if, required_unless, required_with, required_without | ❌ Missing |
| nullable, sometimes, bail | ❌ Missing |
| Nested array validation (`items.*.name`) | ❌ Missing |
| Custom validation rules (class-based) | ❌ Missing |
| Conditional validation (`when`) | ❌ Missing |
| Validation error bags | ❌ Missing |
| After validation hooks | ❌ Missing |
| Localized error messages | ❌ Missing |
| Struct-tag-based validation | ❌ Missing |
| Precognition (live validation) | ❌ Missing |

---

## 6. Authentication & Authorization

### Auth Manager — [manager.go](file:///d:/Go%20framework/auth/manager.go) + [guard.go](file:///d:/Go%20framework/auth/guard.go)
- ✅ `Guard` interface (Check, Guest, User, Attempt, Login, Logout)
- ✅ `AuthManager` with multi-guard support
- ✅ `UserProvider` interface
- ❌ Missing: Actual `SessionGuard` implementation (session persistence)
- ❌ Missing: Token/API guard
- ❌ Missing: JWT support
- ❌ Missing: Password reset flow (tokens, email, broker)
- ❌ Missing: Email verification (signed URLs, middleware)
- ❌ Missing: Remember me (token + cookie)
- ❌ Missing: Two-factor authentication
- ❌ Missing: Password confirmation middleware
- ❌ Missing: Login throttling (beyond generic rate limiter)
- ❌ Missing: Session fixation protection (regenerate on login)
- ❌ Missing: Concurrent session limits
- ❌ Missing: Device tracking

### Authorization — [gate.go](file:///d:/Go%20framework/auth/access/gate.go) (1397 bytes)
- ✅ Gate `Define()` and `Allows()` / `Denies()`
- ❌ Missing: Policy classes mapped to models
- ❌ Missing: `@can` / `@cannot` / `@canany` template directives
- ❌ Missing: `can:` middleware parameter
- ❌ Missing: Gate `before()` / `after()` hooks
- ❌ Missing: Guest gate checks
- ❌ Missing: Policy auto-discovery

---

## 7. Cache — [driver_memory.go](file:///d:/Go%20framework/cache/driver_memory.go) (1750 bytes)

- ✅ `Store` interface (Get, Put, Increment, Decrement, Forever, Forget, Flush)
- ✅ Memory driver with TTL expiration
- ✅ Rate limiter (token bucket)
- ❌ Missing: File cache driver
- ❌ Missing: Redis cache driver
- ❌ Missing: Database cache driver
- ❌ Missing: Cache tags
- ❌ Missing: Atomic locks (distributed)
- ❌ Missing: `Remember()` / `RememberForever()` helper
- ❌ Missing: `Many()` / `PutMany()` bulk operations
- ❌ Missing: Cache events (hit/miss/delete/write)
- ❌ Missing: Cache manager (multi-store switching)
- ❌ Missing: Cache prefix support

---

## 8. Session — 2 drivers (file, cookie)

- ✅ File driver, Cookie driver
- ✅ Session start/save middleware
- ❌ Missing: Redis driver
- ❌ Missing: Database driver
- ❌ Missing: DynamoDB driver
- ❌ Missing: Array driver (testing)
- ❌ Missing: Session encryption
- ❌ Missing: Flash data (`flash()`, `reflash()`, `keep()`)
- ❌ Missing: Session regeneration (`regenerate()`)
- ❌ Missing: Session invalidation (`invalidate()`)
- ❌ Missing: Previous URL storage
- ❌ Missing: Session blocking (concurrent request locking)

---

## 9. Queue — [worker.go](file:///d:/Go%20framework/queue/worker.go) + sync driver

- ✅ `Job` interface (Handle, Failed)
- ✅ Sync driver (immediate execution)
- ✅ Basic worker loop
- ❌ Missing: Database queue driver
- ❌ Missing: Redis queue driver
- ❌ Missing: SQS / RabbitMQ drivers
- ❌ Missing: Delayed jobs (`later()`)
- ❌ Missing: Job chaining (`Chain::of()`)
- ❌ Missing: Batch jobs (`Bus::batch()` with completion/failure callbacks)
- ❌ Missing: Failed jobs table + retry mechanism
- ❌ Missing: Job middleware
- ❌ Missing: Unique jobs (`ShouldBeUnique`)
- ❌ Missing: Job timeout / max attempts / max exceptions / backoff
- ❌ Missing: Job rate limiting
- ❌ Missing: `queue:work`, `queue:retry`, `queue:failed`, `queue:flush` CLI commands
- ❌ Missing: Queue priorities
- ❌ Missing: Job events (before/after process)
- ❌ Missing: Job batching progress tracking

---

## 10. Events — [dispatcher.go](file:///d:/Go%20framework/events/dispatcher.go) (54 lines)

- ✅ Type-based listener registration
- ✅ Wildcard listeners
- ✅ Synchronous dispatch
- ❌ Missing: Queued listeners (`ShouldQueue`)
- ❌ Missing: Event subscribers (class with `subscribe()`)
- ❌ Missing: Stoppable events (return `false` to halt)
- ❌ Missing: Event fake for testing
- ❌ Missing: Broadcasting integration (`ShouldBroadcast`)
- ❌ Missing: Event discovery (auto-register from method signatures)
- ❌ Missing: Model event integration (Eloquent observers)

---

## 11. Mail — [mailer.go](file:///d:/Go%20framework/mail/mailer.go) + [message.go](file:///d:/Go%20framework/mail/message.go)

- ✅ `Mailer` interface (Send)
- ✅ `LogMailer` driver (writes to stdout)
- ✅ `Mailable` with basic From/To/Subject/Body
- ❌ Missing: SMTP driver (actual email sending)
- ❌ Missing: Markdown mail support
- ❌ Missing: Attachments (files, raw data, from storage)
- ❌ Missing: CC, BCC, Reply-To
- ❌ Missing: Mail queuing (`ShouldQueue`)
- ❌ Missing: Multiple mailer drivers (SES, Mailgun, Postmark, Resend)
- ❌ Missing: `Mailable` with `envelope()` / `content()` / `attachments()`
- ❌ Missing: Mail preview in browser
- ❌ Missing: Mail localization
- ❌ Missing: Mail fake for testing (`Mail::assertSent`)

---

## 12. Notifications — [manager.go](file:///d:/Go%20framework/notifications/manager.go) + [channel.go](file:///d:/Go%20framework/notifications/channel.go)

- ✅ `Notification` interface with `Via()` method
- ✅ `Channel` interface
- ✅ Manager dispatches to channels
- ❌ Missing: Mail channel implementation
- ❌ Missing: Database channel (notifications table)
- ❌ Missing: Broadcast channel
- ❌ Missing: Slack channel
- ❌ Missing: SMS/Vonage channel
- ❌ Missing: On-demand notifications
- ❌ Missing: Mark as read/unread
- ❌ Missing: Queued notifications
- ❌ Missing: Notification fake for testing

---

## 13. Broadcasting — [manager.go](file:///d:/Go%20framework/broadcasting/manager.go) + [channel.go](file:///d:/Go%20framework/broadcasting/channel.go)

- ✅ `Broadcaster` interface
- ✅ `LogBroadcaster` driver
- ✅ Public/Private/Presence channel types
- ❌ Missing: Pusher driver
- ❌ Missing: Ably driver
- ❌ Missing: Redis pub/sub driver
- ❌ Missing: WebSocket server (Reverb equivalent)
- ❌ Missing: Channel authorization routes (`routes/channels.go`)
- ❌ Missing: Actual WebSocket upgrade handling
- ❌ Missing: Client-side event support (whisper)
- ❌ Missing: Echo JS library equivalent

---

## 14. File Storage — [filesystem.go](file:///d:/Go%20framework/storage/filesystem.go) (1306 bytes)

- ✅ `Filesystem` interface (Get, Put, Delete, Exists, etc.)
- ❌ Missing: Local driver implementation
- ❌ Missing: S3 driver (AWS SDK)
- ❌ Missing: SFTP driver
- ❌ Missing: GCS (Google Cloud Storage) driver
- ❌ Missing: Storage manager (multi-disk switching)
- ❌ Missing: `UploadedFile` wrapper with `store()` / `storeAs()`
- ❌ Missing: Temporary URLs
- ❌ Missing: Visibility (public/private)
- ❌ Missing: `storage:link` command (public symlink)
- ❌ Missing: Streaming uploads/downloads
- ❌ Missing: File fake for testing

---

## 15. Testing — [testcase.go](file:///d:/Go%20framework/testing/testcase.go) + [db.go](file:///d:/Go%20framework/testing/db.go)

- ✅ `TestCase` struct wrapping `httptest`
- ✅ `Get()` / `Post()` request helpers
- ✅ `AssertStatus()`, `AssertJson()`
- ✅ Basic fakes ([fakes.go](file:///d:/Go%20framework/support/fakes/fakes.go))
- ❌ Missing: `AssertDatabaseHas` / `AssertDatabaseMissing` / `AssertDatabaseCount`
- ❌ Missing: `AssertSoftDeleted` / `AssertNotSoftDeleted`
- ❌ Missing: `actingAs()` auth helper
- ❌ Missing: Time travel / freeze (`travel()`, `travelTo()`, `freezeTime()`)
- ❌ Missing: `RefreshDatabase` trait (transaction rollback per test)
- ❌ Missing: View assertions (`assertViewHas`, `assertViewIs`)
- ❌ Missing: Command testing (`expectsOutput`, `expectsConfirmation`)
- ❌ Missing: Browser testing (Dusk equivalent)
- ❌ Missing: `assertRedirectToRoute`
- ❌ Missing: `assertValid` / `assertInvalid`
- ❌ Missing: Parallel testing support
- ❌ Missing: Snapshot testing

---

## 16. CLI / Console — [console/kernel.go](file:///d:/Go%20framework/console/kernel.go)

- ✅ Cobra-based command framework
- ✅ Schedule support
- ❌ Missing: `make:controller`, `make:model`, `make:migration`, `make:middleware`, `make:request`, `make:resource`, `make:event`, `make:listener`, `make:job`, `make:mail`, `make:notification`, `make:policy`, `make:rule`, `make:seeder`, `make:factory`, `make:test`, `make:command`, `make:provider`, `make:observer`, `make:cast` generators
- ❌ Missing: `route:list`, `route:cache`, `route:clear`
- ❌ Missing: `config:cache`, `config:clear`
- ❌ Missing: `view:cache`, `view:clear`
- ❌ Missing: `cache:clear`, `cache:forget`
- ❌ Missing: `queue:work`, `queue:retry`, `queue:failed`, `queue:flush`, `queue:restart`
- ❌ Missing: `key:generate`
- ❌ Missing: `serve` command
- ❌ Missing: `tinker` / REPL
- ❌ Missing: `optimize` / `optimize:clear`
- ❌ Missing: `down` / `up` (maintenance mode with secret, redirect, retry-after, render, status)
- ❌ Missing: `event:list`, `event:cache`
- ❌ Missing: `schedule:list`, `schedule:test`
- ❌ Missing: `db:show`, `db:table`, `db:monitor`
- ❌ Missing: `vendor:publish`
- ❌ Missing: `storage:link`
- ❌ Missing: `model:prune`
- ❌ Missing: `schema:dump`
- ❌ Missing: Rich CLI UX (Prompts equivalent — spinners, multiselect, search, progress bars with ETA)
- ❌ Missing: `about` command (framework/environment info)

---

## 17. Middleware — Detailed Gaps

### Implemented middleware (8 files in `http/middleware/`)
- ✅ Auth, CORS, CSRF, EncryptCookies, Recovery (panic), SecurityHeaders, Session, Throttle

### Missing middleware
| Middleware | Status |
|---|---|
| Trusted Proxies (`TrustProxies`) | ❌ Missing | Sets `X-Forwarded-*` trusted IPs |
| Trusted Hosts (`TrustHosts`) | ❌ Missing | |
| `TrimStrings` | ❌ Missing | Auto-trim whitespace from input |
| `ConvertEmptyStringsToNull` | ❌ Missing | |
| `PreventRequestsDuringMaintenance` | ❌ Missing | |
| `ValidatePostSize` | ❌ Missing | |
| `EnsureEmailIsVerified` | ❌ Missing | |
| `SetCacheHeaders` | ❌ Missing | `cache.headers` middleware |
| Middleware parameters (`throttle:60,1`) | ❌ Missing | Middleware with arguments |
| Middleware priority / sorting | ❌ Missing | |
| Middleware groups (`web`, `api`) | ❌ Missing | Named middleware groups |
| Middleware aliases | ❌ Missing | Short names for middleware |

---

## 18. Cross-Cutting Core Gaps

| Area | Gap |
|---|---|
| **Error Handling** | Only basic `HttpException` + panic recovery. Missing: custom error pages (404/500/503), Whoops-like dev error pages, maintenance mode (`down`/`up`), renderable/reportable exceptions, ignored exceptions, exception context, Sentry integration hooks |
| **Localization** | Translator exists but missing: pluralization rules (CLDR), locale detection middleware, `trans_choice()`, JSON translation files, publishing language files |
| **HTTP Client** | Basic fluent wrapper. Missing: retry with backoff, connection pool, async requests, fake/mock for testing (`Http::fake()`), response sequence stubs, middleware (request/response), `throw()`, streaming |
| **Encryption** | AES-256-GCM exists. Missing: key rotation, encrypted model casts, deterministic encryption |
| **Hashing** | Bcrypt + Argon2id ✅ Solid |
| **Logging** | Basic slog wrapper. Missing: daily rotation, stack channels, Slack channel, Syslog, Papertrail, contextual data enrichment, log deprecation warnings, request ID injection |
| **Config** | .env loading works. Missing: config caching, nested dot-notation `Get("database.connections.mysql.host")`, environment-specific overrides, typed access helpers |
| **Pipeline** | ❌ No generic Pipeline class (only used implicitly in middleware). Laravel's Pipeline is usable standalone |
| **Macroable** | ❌ No Macroable trait equivalent — can't extend framework classes at runtime |
| **Conditionable / Tappable** | ❌ No `when()` / `unless()` / `tap()` chaining on builder/collection objects |
| **`Once` helper** | ❌ Memoize function results (Laravel 11+) |
| **`Defer` helper** | ❌ Defer work until response is sent (Laravel 11+) |
| **`Concurrency` helper** | ❌ Run closures concurrently and collect results (Laravel 11+) |
| **`Context` propagation** | ❌ Cross-layer context that flows to jobs/logs (Laravel 11+) |
| **Streaming responses** | ❌ Missing: `StreamedResponse`, `StreamedJsonResponse`, `streamDownload()` |
| **File downloads** | ❌ Missing: `response()->download()`, `response()->file()` |
| **JSONP responses** | ❌ Missing |
| **Vite / Asset Bundling** | ❌ Missing: `@vite` directive, `mix()` helper, asset versioning |
| **Contracts package** | ❌ Missing: All interfaces published as a separate importable package |
| **Content negotiation** | ❌ Missing: `$request->expectsJson()`, `wantsJson()`, `accepts()` |
| **Request convenience methods** | ❌ Missing: `only()`, `except()`, `input()`, `query()`, `collect()`, `file()`, `hasFile()` |

---

## 19. Sections 35–50 — ALL MISSING (specification only, zero code)

These 16 feature areas are **fully specified** in [GoW_SPEC.md](file:///d:/Go%20framework/GoW_SPEC.md) but have **no implementation code** anywhere in the project:

| Section | Feature | Priority |
|---|---|---|
| §35 | **API Resources / Transformers** | 🔴 High — essential for any API |
| §36 | **Collections & Lazy Collections** | 🔴 High — used everywhere |
| §37 | **Helpers (Str, Arr, Number, Date, Benchmark, Sleep)** | 🔴 High — DX foundation |
| §38 | **Signed URLs** | 🟡 Medium |
| §39 | **Route Model Binding** | 🔴 High — core Laravel DX |
| §40 | **Process / CLI Command Execution** | 🟡 Medium |
| §41 | **Schema Dumping** | 🟢 Low |
| §42 | **Model Pruning** | 🟢 Low |
| §43 | **Health Checks** | 🟡 Medium |
| §44 | **Feature Flags (Pennant equivalent)** | 🟡 Medium |
| §45 | **Advanced Blade Directives** | 🟡 Medium |
| §46 | **Request/Response Macros, Old Input, Flashing** | 🔴 High |
| §47 | **Query Logging & Database Profiling** | 🟡 Medium |
| §48 | **Package Auto-Discovery & Publishing** | 🟡 Medium |
| §49 | **Testing Utilities (DB assertions, time travel, auth)** | 🔴 High |
| §50 | **Modern Laravel Parity (Reverb, Pulse, Prompts, Context)** | 🟢 Low |

---

## PART II — FIRST-PARTY ECOSYSTEM PACKAGES (all missing from GoW)

> [!IMPORTANT]
> Laravel's power comes not just from its core framework but from its **first-party ecosystem**. None of these exist in GoW — not even as specifications. Each would need its own spec section, design, and implementation.

---

### 20. Sanctum — API Token Authentication
**What it does**: Lightweight API token auth for SPAs, mobile apps, and simple token-based APIs. Personal access tokens stored in database. SPA authentication via session cookies.

| Feature | Status |
|---|---|
| Personal access tokens (hashed in `personal_access_tokens` table) | ❌ Missing |
| Token abilities/scopes | ❌ Missing |
| SPA cookie-based authentication (CSRF + session) | ❌ Missing |
| Token expiration | ❌ Missing |
| `currentAccessToken()` / `tokenCan()` | ❌ Missing |
| Middleware: `auth:sanctum` | ❌ Missing |
| Token pruning | ❌ Missing |

**Priority**: 🔴 **High** — most Laravel APIs use Sanctum

---

### 21. Passport — Full OAuth2 Server
**What it does**: Full OAuth2 server implementation with authorization codes, personal access tokens, client credentials, password grants.

| Feature | Status |
|---|---|
| OAuth2 authorization code flow | ❌ Missing |
| Client credentials grant | ❌ Missing |
| Personal access tokens | ❌ Missing |
| Password grant | ❌ Missing |
| Token scopes | ❌ Missing |
| Token revocation | ❌ Missing |
| PKCE support | ❌ Missing |
| Vue.js management components | ❌ Missing |

**Priority**: 🟡 Medium — Sanctum covers most use cases

---

### 22. Socialite — OAuth Social Login
**What it does**: OAuth authentication with third-party providers (Google, Facebook, GitHub, Twitter, LinkedIn, GitLab, Bitbucket, Slack, etc.).

| Feature | Status |
|---|---|
| OAuth 1 & 2 provider abstraction | ❌ Missing |
| Redirect to provider | ❌ Missing |
| Handle callback with user info | ❌ Missing |
| Stateless authentication (API) | ❌ Missing |
| Access token retrieval | ❌ Missing |
| Scopes | ❌ Missing |
| Custom providers | ❌ Missing |

**Priority**: 🟡 Medium — common in user-facing apps

---

### 23. Scout — Full-Text Search
**What it does**: Driver-based full-text search for Eloquent models. Automatically syncs model data to search index.

| Feature | Status |
|---|---|
| Searchable trait on models | ❌ Missing |
| Meilisearch driver | ❌ Missing |
| Algolia driver | ❌ Missing |
| Database driver (LIKE/fulltext) | ❌ Missing |
| Typesense driver | ❌ Missing |
| Index management (`scout:import`, `scout:flush`) | ❌ Missing |
| Soft delete handling in index | ❌ Missing |
| Custom search builder | ❌ Missing |
| Pagination of search results | ❌ Missing |
| Conditional indexing (`shouldBeSearchable`) | ❌ Missing |

**Priority**: 🟡 Medium

---

### 24. Telescope — Debug & Monitoring Dashboard
**What it does**: Debug assistant providing insight into requests, exceptions, queries, jobs, mail, notifications, cache, schedule, dumps — all in a web dashboard.

| Feature | Status |
|---|---|
| Request watcher (log all HTTP requests) | ❌ Missing |
| Exception watcher | ❌ Missing |
| Query watcher (all SQL queries) | ❌ Missing |
| Job watcher | ❌ Missing |
| Mail watcher (preview sent emails) | ❌ Missing |
| Notification watcher | ❌ Missing |
| Cache watcher | ❌ Missing |
| Schedule watcher | ❌ Missing |
| Log watcher | ❌ Missing |
| Dump watcher (`dump()` to dashboard) | ❌ Missing |
| Model watcher (create/update/delete) | ❌ Missing |
| Gate watcher (auth checks) | ❌ Missing |
| View watcher | ❌ Missing |
| Web dashboard UI | ❌ Missing |
| Pruning old entries | ❌ Missing |
| Authorization (who can view dashboard) | ❌ Missing |

**Priority**: 🟡 Medium — huge DX improvement

---

### 25. Horizon — Queue Monitoring Dashboard
**What it does**: Dashboard and configuration for Redis-powered queues. Monitors queue throughput, runtime, failures, with auto-scaling workers.

| Feature | Status |
|---|---|
| Real-time queue metrics dashboard | ❌ Missing |
| Job/failure tracking | ❌ Missing |
| Auto-balancing workers across queues | ❌ Missing |
| Retry/delete failed jobs from UI | ❌ Missing |
| Queue wait time monitoring | ❌ Missing |
| Worker supervision | ❌ Missing |
| Tags and filtering | ❌ Missing |
| Notification on long wait times | ❌ Missing |

**Priority**: 🟢 Low — needs Redis queue driver first

---

### 26. Cashier — Subscription Billing
**What it does**: Expressive interface for Stripe and Paddle subscription billing services.

| Feature | Status |
|---|---|
| Subscription creation/management | ❌ Missing |
| Subscription plan swapping | ❌ Missing |
| Subscription quantities (per-seat) | ❌ Missing |
| Trial periods | ❌ Missing |
| Invoices (generation, PDF download) | ❌ Missing |
| Payment methods (cards) | ❌ Missing |
| Webhook handling (Stripe events) | ❌ Missing |
| Single charges | ❌ Missing |
| Customer portal redirect | ❌ Missing |
| Tax calculation | ❌ Missing |

**Priority**: 🟢 Low — niche use case

---

### 27. Fortify — Headless Auth Backend
**What it does**: Backend authentication implementation without UI. Registration, login, password reset, email verification, 2FA — all as API endpoints.

| Feature | Status |
|---|---|
| Registration flow with custom actions | ❌ Missing |
| Login with rate limiting | ❌ Missing |
| Password reset via email | ❌ Missing |
| Email verification | ❌ Missing |
| Two-factor authentication (TOTP) | ❌ Missing |
| Password confirmation | ❌ Missing |
| Profile update | ❌ Missing |
| Custom validation/creation actions | ❌ Missing |

**Priority**: 🔴 High — provides the auth actions GoW is missing

---

### 28. Breeze — Starter Kit
**What it does**: Minimal authentication scaffolding with login, registration, password reset, email verification, and profile management. Ships with Blade, Livewire, React, Vue, or API stacks.

| Feature | Status |
|---|---|
| `gow new --breeze` scaffolding command | ❌ Missing |
| Blade + Alpine.js stack | ❌ Missing |
| Livewire stack | ❌ Missing (Livewire N/A in Go) |
| Inertia + React stack | ❌ Missing |
| Inertia + Vue stack | ❌ Missing |
| API-only stack | ❌ Missing |
| Dark mode support | ❌ Missing |
| TypeScript support | ❌ Missing |

**Priority**: 🟡 Medium — great for onboarding

---

### 29. Jetstream — Advanced Starter Kit
**What it does**: Extends Breeze with team management, API tokens (via Sanctum), two-factor authentication, browser session management, profile photos.

| Feature | Status |
|---|---|
| Team management (create, invite, roles, permissions) | ❌ Missing |
| Team switching | ❌ Missing |
| API token management UI | ❌ Missing |
| Browser session management | ❌ Missing |
| Profile photo upload | ❌ Missing |
| Account deletion | ❌ Missing |

**Priority**: 🟢 Low — built on top of Breeze + Sanctum

---

### 30. Dusk — Browser Testing
**What it does**: Automated browser testing using ChromeDriver. Page objects, element interaction, screenshots, authentication helpers, component testing.

| Feature | Status |
|---|---|
| ChromeDriver integration | ❌ Missing |
| Browser automation API (`visit`, `type`, `click`, `press`) | ❌ Missing |
| Element assertions (`assertSee`, `assertVisible`, `assertInputValue`) | ❌ Missing |
| Page objects | ❌ Missing |
| Component testing | ❌ Missing |
| Screenshots on failure | ❌ Missing |
| Authentication helpers (`loginAs`) | ❌ Missing |
| Multiple browser windows | ❌ Missing |
| File upload testing | ❌ Missing |
| JavaScript dialog handling | ❌ Missing |

**Priority**: 🟢 Low — Go has other options (chromedp, rod)

---

### 31. Precognition — Live Validation
**What it does**: Validate form data against server-side rules in real-time before submission without full form request.

| Feature | Status |
|---|---|
| `Precognition` header detection | ❌ Missing |
| Partial validation (only changed fields) | ❌ Missing |
| `@precognitive` middleware | ❌ Missing |
| JS client integration (Axios/fetch wrapper) | ❌ Missing |
| Vue/React/Alpine.js integration | ❌ Missing |

**Priority**: 🟢 Low

---

### 32. Folio — Page-Based Routing
**What it does**: File-system based routing. Drop a Blade template in `resources/views/pages/` and it's automatically routed.

| Feature | Status |
|---|---|
| File-system route resolution | ❌ Missing |
| Dynamic segments via filenames (`[id].goblade`) | ❌ Missing |
| Route middleware in page files | ❌ Missing |
| Wildcard pages | ❌ Missing |

**Priority**: 🟢 Low — niche

---

### 33. Pint — Code Style Fixer
**What it does**: Opinionated code style fixer (wraps PHP CS Fixer). Go equivalent would wrap `gofmt`/`goimports` + custom rules.

| Feature | Status |
|---|---|
| `gow pint` / `gow fmt` command | ❌ Missing |
| Preset styles (laravel, psr-12) | N/A (Go has `gofmt`) |
| Custom rule configuration | ❌ Missing |

**Priority**: 🟢 Low — Go already has `gofmt`

---

### 34. Sail / Envoy — DevOps Tooling
**What it does**: Sail = Docker dev environment. Envoy = SSH task runner for deployment.

| Feature | Status |
|---|---|
| Docker Compose dev environment | ❌ Missing |
| `gow sail up` / `gow sail down` | ❌ Missing |
| Deployment task runner | ❌ Missing |

**Priority**: 🟢 Low — can use Docker directly

---

### 35. Livewire / Inertia.js / Echo — Frontend Integration
**What they do**: Livewire = full-stack reactive UI. Inertia = SPA without API. Echo = WebSocket event subscription (JS library).

| Feature | Status |
|---|---|
| Server-side reactive components (Livewire-like) | ❌ Missing |
| Inertia.js adapter (shared props, lazy loading) | ❌ Missing |
| Echo JS library for WebSocket subscription | ❌ Missing |
| Turbo/Hotwire integration (alternative to Livewire) | ❌ Missing |

**Priority**: 🟡 Medium — Echo is important for broadcasting

---

## PART III — UPDATED ROADMAP

---

### 🔴 Phase A — Make Core Usable (must-do first)

1. **Complete Query Builder** — JOINs, WHERE IN/NULL/BETWEEN, aggregates, transactions, subqueries, raw expressions
2. **Complete ORM Relationships** — actually implement `loadRelation()`, then pivot tables, polymorphic
3. **Finish Validation** — top 30+ rules
4. **Add missing HTTP verb helpers** — PUT, PATCH, DELETE, OPTIONS, HEAD
5. **Implement `@extends` properly** in view compiler + `@for`/`@while`/`@switch`
6. **Session flash data + old input** + session regeneration
7. **Collections** (§36) — `Collect[T]()` with map/filter/reduce/pluck/chunk
8. **Str/Arr/Number/Date helpers** (§37)
9. **Route Model Binding** (§39)
10. **All `make:*` CLI generators**
11. **ORM CRUD** — `Find()`, `Update()`, `Delete()`, `Save()`, soft delete query filtering
12. **Container** — contextual binding, tagging, scoped bindings

### 🟡 Phase B — Production Readiness

13. Redis drivers (cache, session, queue)
14. Database queue driver + failed jobs table + retry
15. SMTP mail driver + attachments + CC/BCC
16. Database notification channel
17. API Resources / Transformers (§35)
18. Request/Response helpers (old input, flashing, boolean/integer casting) (§46)
19. Query logging & N+1 detection (§47)
20. Testing utilities — `assertDatabaseHas`, `actingAs`, time travel (§49)
21. Resource / API resource routes + URL generation
22. Error handling — custom error pages, maintenance mode, dev error pages
23. Logging — daily rotation, stack channels, contextual enrichment
24. Config caching + nested dot-notation
25. Trusted Proxies + Trusted Hosts + TrimStrings middleware
26. Sanctum equivalent (API tokens) — **first ecosystem package**
27. HTTP Client — retry, fake, pool
28. Streaming responses + file downloads
29. Generic Pipeline class
30. File storage — local driver + S3 driver + storage manager

### 🟢 Phase C — Feature Completeness

31. Signed URLs (§38)
32. Health checks (§43)
33. Feature flags (§44)
34. Advanced Blade directives — `$loop`, `@once`, `@class`, `@checked`, components & slots (§45)
35. Package auto-discovery & publishing (§48)
36. Process execution wrapper (§40)
37. Schema dumping (§41)
38. Model pruning (§42)
39. ORM advanced — scopes, observers, mutators/accessors, attribute casting, enum casting, encrypted casts
40. ORM advanced — `HasUuids`/`HasUlids`, model strictness, prevent lazy loading
41. Fortify equivalent (headless auth backend)
42. Auth — password reset flow, email verification, 2FA, remember me
43. Authorization — policies, policy auto-discovery, `@can` directives
44. Broadcasting — real WebSocket server, channel authorization
45. Event — queued listeners, subscribers, model events

### 🔵 Phase D — Ecosystem Packages

46. Socialite equivalent (OAuth social login)
47. Scout equivalent (full-text search)
48. Telescope equivalent (debug dashboard)
49. Breeze equivalent (starter kit / scaffolding)
50. Passport equivalent (full OAuth2 server)
51. Cashier equivalent (subscription billing)
52. Horizon equivalent (queue dashboard)
53. Dusk equivalent (browser testing)
54. Echo JS equivalent (WebSocket client library)

### 🟣 Phase E — Aspirational / Modern Parity

55. Reverb equivalent (first-party WebSocket server)
56. Pulse equivalent (APM monitoring dashboard)
57. Prompts equivalent (rich CLI UX)
58. Context propagation (Laravel 11+)
59. `Once` / `Defer` / `Concurrency` helpers (Laravel 11+)
60. Precognition (live validation)
61. Folio (file-based routing)
62. Inertia.js adapter
63. Livewire-like server-side components (very ambitious in Go)
64. Jetstream equivalent (teams, advanced starter kit)

---

## Total Gap Count

| Category | Items |
|---|---|
| Core framework (Sections 1–18) | ~250 missing features |
| Spec sections 35–50 (no code) | ~80 features across 16 sections |
| First-party ecosystem (Sections 20–35) | ~120 features across 17 packages |
| **Grand Total** | **~450 missing features/items** |

> [!IMPORTANT]
> For comparison, Laravel's `Illuminate` source + first-party packages span **~600,000+ lines** of PHP. GoW currently has roughly **~1,500 lines of Go framework code**. Even accounting for Go's conciseness, the gap is **~100:1** in code volume and **~450 individual missing features**.

> [!TIP]
> A practical "Laravel-like Go framework" doesn't need 100% parity. Completing **Phase A + Phase B** (~30 items) would give you a **usable** framework. Completing through **Phase C** would make it **competitive** with existing Go frameworks (Gin, Echo, Fiber) while offering Laravel-like DX. The ecosystem packages (Phase D/E) are what would make it truly "Laravel for Go."
