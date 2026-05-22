# GoW Framework — Missing Features

> **Last Updated**: 2026-05-22
> **Based on**: Full audit of 115 Go source files (~7,624 lines)
> **Legend**: 🔴 Critical | 🟡 Important | 🟢 Nice to Have

---

## 🔴 CRITICAL BUGS (Existing Code That's Broken)

### 1. Eager Loading Returns Nil
- **File**: [relation.go:65](routing/../database/orm/relation.go#L65)
- `loadRelation()` collects parent IDs but returns `nil` — no SQL query fires
- Any `.With("relation").Get()` silently returns empty relation fields
- **Fix**: Execute `SELECT * FROM target WHERE fk IN (...)` and hydrate results back to parents

### 2. Model Observers Never Called
- **File**: [model.go](database/orm/model.go) defines `DispatchModelEvent()` and Observer interface
- **File**: [query.go](database/orm/query.go) — `Insert()`, `Update()`, `Delete()` never call it
- **Fix**: Add `DispatchModelEvent(model, "creating")` before DB calls, `DispatchModelEvent(model, "created")` after

### 3. Global Scopes Never Applied
- **File**: [model.go](database/orm/model.go) — `globalScopes` registry exists
- `NewQuery()` and `Get()` never apply registered scopes
- **Fix**: In `NewQuery()`, iterate `globalScopes[modelName]` and call `scope.Apply(builder)`

### 4. Collection.Chunk() Returns Nil
- **File**: [collection.go:77](support/collection/collection.go#L77)
- `Chunk()` has wrong return type `[][]*Collection[T]` and always returns nil
- `Chunked()` works correctly but `Chunk()` is broken dead code
- **Fix**: Remove broken `Chunk()`, rename `Chunked()` to `Chunk()`, fix signature

### 5. @while/@once Compile to Comments
- **File**: [compiler.go](view/compiler.go)
- `@while(condition)` compiles to `{{/* while $1 */}}` — a Go template comment (invisible)
- `@once` compiles to `{{/* @once start */}}` — also invisible
- **Fix**: Map @while to a custom template func, implement @once with a template variable guard

### 6. S3Driver Methods Are All No-Ops
- **File**: [storage.go:84-98](storage/storage.go)
- Every method returns nil/false — no actual S3 interaction
- Documented as working in docs

---

## 🔴 CRITICAL MISSING — Required for Any Real Application

### Routing

| Feature | Description | Laravel Equivalent |
|---------|-------------|--------------------|
| Route Model Binding (implicit) | Auto-resolve `{user}` param to User model | `Route::get('/users/{user}', ...)` |
| Route Model Binding (explicit) | Custom resolution logic | `Route::bind('user', fn($value) => ...)` |
| Optional Parameters | `{id?}` with default value | `Route::get('/users/{name?}', ...)` |
| Regex Constraints | Validate params against pattern | `->where('id', '[0-9]+')` |
| URL Generation | Generate URL from named route + params | `route('user.show', ['id' => 1])` |
| Fallback Route | Catch-all 404 handler | `Route::fallback(fn() => ...)` |
| Redirect Routes | Simple redirect without controller | `Route::redirect('/here', '/there')` |
| View Routes | Return a view without controller | `Route::view('/welcome', 'welcome')` |

### Database / ORM

| Feature | Description | Priority |
|---------|-------------|----------|
| **Transactions** | `DB.Transaction(func(tx) error)`, `BeginTransaction()`, `Commit()`, `Rollback()` | 🔴 Blocker |
| **Soft Deletes** | Auto-add `WHERE deleted_at IS NULL`, `WithTrashed()`, `OnlyTrashed()`, `Restore()` | 🔴 |
| **Working Relationships** | HasOne, HasMany, BelongsTo must actually execute SQL | 🔴 |
| BelongsToMany (pivot) | Many-to-many with pivot table (attach, detach, sync, toggle) | 🔴 |
| HasOneThrough / HasManyThrough | Indirect relationships | 🟡 |
| Morph relationships | MorphOne, MorphMany, MorphTo, MorphToMany | 🟡 |
| Raw Expressions | `DB.Raw("COALESCE(col, 0)")` inside select/where | 🔴 |
| Subqueries | Where subquery, select subquery, from subquery | 🟡 |
| GROUP BY / HAVING | Group aggregation queries | 🔴 |
| UNION / UNION ALL | Combine query results | 🟡 |
| Distinct | `SELECT DISTINCT` | 🔴 |
| Chunk / ChunkById | Process large datasets in batches | 🟡 |
| Increment / Decrement | `$user->increment('votes')` | 🟡 |
| Upsert | Insert or update on duplicate key | 🟡 |
| InsertGetId | Insert and return the auto-increment ID | 🟡 |
| Database Locking | `lockForUpdate()`, `sharedLock()` | 🟡 |
| Multiple DB Connections | Actually switch between connections at runtime | 🟡 |
| Savepoints | Nested transaction support | 🟢 |
| WHERE date/day/month/year/time | Date-specific where clauses | 🟢 |
| WHERE JSON (`->`) | JSON column access in queries | 🟢 |
| Explain | Query execution plan | 🟢 |

### ORM Model Features

| Feature | Description | Priority |
|---------|-------------|----------|
| Find / FindOrFail by ID | `User.Find(1)` with 404 on fail | 🔴 (Find exists, FindOrFail missing) |
| FirstOrCreate / FirstOrNew / UpdateOrCreate | Convenient find-or-create patterns | 🔴 |
| Mass Assignment (fillable/guarded) | Tags exist but not enforced during Insert/Update | 🔴 |
| Mutators & Accessors | Transform attributes on get/set | 🟡 |
| Attribute Casting | Dates, JSON, enum, encrypted, custom | 🟡 |
| Hidden / Visible attributes | Control serialization output | 🟡 |
| Dirty Tracking | `isDirty()`, `wasChanged()`, `getOriginal()` | 🟡 |
| Default Attribute Values | Pre-populate model defaults | 🟡 |
| toJSON / toArray serialization | Consistent model serialization | 🔴 |
| Model Replicating | `model.Replicate()` | 🟢 |
| HasUuids / HasUlids | UUID/ULID as primary key | 🟡 |
| Prevent Lazy Loading | N+1 protection in dev mode | 🟡 |
| Model Factories (real) | With faker integration, sequences, states | 🟡 |
| Model Pruning | Auto-delete old records on schedule | 🟢 |
| Touching Parent Timestamps | Update parent `updated_at` on child change | 🟢 |

### Pagination

| Feature | Description | Priority |
|---------|-------------|----------|
| Query Builder integration | `builder.Paginate(15)` returns paginator | 🔴 |
| Simple Pagination | No total count, just next/prev | 🟡 |
| Cursor Pagination | For infinite scroll / large datasets | 🟡 |
| JSON response meta/links | Standard API pagination format | 🔴 |
| View integration | HTML pagination links | 🟡 |

### Views / Template Engine

| Feature | Description | Priority |
|---------|-------------|----------|
| **`@extends` layout resolution** | Engine must resolve parent template, not just emit a comment | 🔴 |
| `$loop` variable | `$loop.index`, `$loop.first`, `$loop.last` in `@foreach` | 🔴 |
| Components & Slots | `<x-alert>`, `{{ $slot }}`, named slots | 🟡 |
| Stacks | `@push('scripts')` / `@stack('scripts')` / `@prepend` | 🟡 |
| `{!! !!}` unescaped output | Raw HTML output | 🔴 |
| `{{-- comments --}}` | Blade-style comments (not rendered) | 🟡 |
| `@auth` / `@guest` directives | Conditional rendering based on auth state | 🟡 |
| `@env` / `@production` | Conditional rendering based on environment | 🟡 |
| `@method` | Form method spoofing (`PUT`, `DELETE`) | 🔴 |
| `@error` directive | Display validation errors | 🔴 |
| `@props` directive | Declare component properties | 🟡 |
| Custom directive registration | `Compiler.Directive("datetime", func)` | 🟡 |
| View sharing (global data) | Share data across all views | 🟡 |
| View creators | Run callback when view is first created | 🟢 |

### Validation (14 rules exist, ~25+ missing)

| Rule | Priority |
|------|----------|
| `boolean`, `integer`, `array` | 🔴 |
| `date`, `before`, `after`, `date_format` | 🟡 |
| `image`, `mimes`, `mimetypes`, `file`, `dimensions` | 🟡 |
| `url`, `active_url`, `ip`, `ipv4`, `ipv6` | 🟡 |
| `json`, `uuid`, `ulid` | 🟡 |
| `starts_with`, `ends_with`, `contains` | 🟡 |
| `alpha`, `alpha_num`, `alpha_dash` | 🟡 |
| `digits`, `digits_between` | 🟡 |
| `gt`, `gte`, `lt`, `lte` (compare fields) | 🟡 |
| `same`, `different` (compare fields) | 🟡 |
| `required_if`, `required_unless`, `required_with`, `required_without` | 🟡 |
| `nullable`, `sometimes`, `bail` | 🔴 |
| Nested array validation (`items.*.name`) | 🟡 |
| Custom validation rules (class-based) | 🟡 |
| After validation hooks | 🟡 |
| Localized error messages | 🟡 |
| Struct-tag-based validation | 🟡 |

### Authentication & Authorization

| Feature | Description | Priority |
|---------|-------------|----------|
| **SessionGuard implementation** | Actual session-based auth persistence | 🔴 |
| **Password Reset flow** | Tokens, email, broker, reset endpoint | 🔴 |
| **Email Verification** | Signed URLs, verification middleware | 🔴 |
| Remember Me | Token + cookie for persistent login | 🟡 |
| Two-Factor Authentication | TOTP-based 2FA | 🟡 |
| Password Confirmation middleware | Re-confirm password for sensitive actions | 🟡 |
| Login Throttling | Per-user rate limiting on login | 🟡 |
| Session Fixation Protection | Regenerate session ID on login | 🔴 |
| Policy classes | Map policies to models | 🟡 |
| Gate before/after hooks | Pre/post authorization checks | 🟡 |
| `@can` / `@cannot` / `@canany` in templates | Authorization in views (partially done) | 🟡 |
| `can:` middleware parameter | Inline authorization in route definition | 🟡 |
| Sanctum — token expiration check | Currently accepts expired tokens | 🔴 |
| Sanctum — token abilities verification | `tokenCan('server:update')` | 🔴 |

### HTTP Request / Response

| Feature | Description | Priority |
|---------|-------------|----------|
| `Input()`, `Query()`, `Post()` | Convenient request data access | 🔴 |
| `Only()`, `Except()` | Filter request data | 🔴 |
| `HasFile()`, `File()` | File upload handling | 🔴 |
| `ExpectsJson()`, `WantsJson()`, `Accepts()` | Content negotiation | 🔴 |
| `Collect()` | Get input as Collection | 🟡 |
| `Boolean()`, `Integer()`, `Float()` | Type-casted input | 🟡 |
| `response()->download()` | File download response | 🟡 |
| `response()->file()` | File display response | 🟡 |
| `StreamedResponse` | Streaming responses | 🟡 |
| `StreamedJsonResponse` | Stream large JSON datasets | 🟡 |
| Redirect with flash data | `redirect()->with('status', 'saved')` | 🔴 |
| Redirect to named route | `redirect()->route('dashboard')` | 🔴 |
| Redirect back | `redirect()->back()` | 🔴 |
| Old Input | `old('email')` in views after validation failure | 🔴 |

---

## 🟡 IMPORTANT MISSING — Production Readiness

### Cache

| Feature | Priority |
|---------|----------|
| File cache driver | 🟡 |
| Database cache driver | 🟡 |
| Cache tags | 🟡 |
| Atomic locks (distributed) | 🟡 |
| `Remember()` / `RememberForever()` | 🔴 |
| `Many()` / `PutMany()` bulk ops | 🟡 |
| Cache events (hit/miss/write/delete) | 🟢 |
| Cache prefix support | 🟡 |

### Session

| Feature | Priority |
|---------|----------|
| Database driver | 🟡 |
| Array driver (testing) | 🟡 |
| Session encryption | 🟡 |
| Session blocking (concurrent request locking) | 🟢 |
| Previous URL storage (`intended()`) | 🔴 |

### Queue

| Feature | Priority |
|---------|----------|
| Delayed jobs (`Delay(5 * time.Minute)`) | 🔴 |
| Job chaining (`Chain(job1, job2, job3)`) | 🟡 |
| Batch jobs (with completion/failure callbacks) | 🟡 |
| Failed jobs table + retry mechanism | 🔴 |
| Job middleware | 🟡 |
| Unique jobs (`ShouldBeUnique`) | 🟡 |
| Job timeout / max attempts / backoff | 🔴 |
| Job rate limiting | 🟡 |
| Queue priorities | 🟡 |
| Job events (before/after process) | 🟡 |
| SQS / RabbitMQ drivers | 🟢 |

### Events

| Feature | Priority |
|---------|----------|
| Queued listeners (`ShouldQueue`) | 🟡 |
| Event subscribers | 🟡 |
| Stoppable events (return false to halt) | 🟡 |
| Event fake for testing | 🟡 |
| Broadcasting integration (`ShouldBroadcast`) | 🟡 |
| Model event integration | 🟡 |

### Mail

| Feature | Priority |
|---------|----------|
| **SMTP driver** (actual email sending) | 🔴 |
| Markdown mail support | 🟡 |
| Attachments (files, raw data, from storage) | 🔴 |
| CC, BCC, Reply-To | 🔴 |
| Mail queuing (`ShouldQueue`) | 🟡 |
| Multiple mailer drivers (SES, Mailgun, Postmark) | 🟢 |
| Mail preview in browser | 🟡 |
| Mail localization | 🟢 |
| Mail fake for testing (`Mail::assertSent`) | 🟡 |

### Notifications

| Feature | Priority |
|---------|----------|
| Mail channel implementation | 🔴 |
| Database channel (notifications table) | 🔴 |
| Broadcast channel | 🟡 |
| Slack channel | 🟡 |
| On-demand notifications | 🟡 |
| Mark as read/unread | 🟡 |
| Queued notifications | 🟡 |

### Broadcasting

| Feature | Priority |
|---------|----------|
| WebSocket server (Reverb equivalent) | 🟡 |
| Pusher driver | 🟡 |
| Redis pub/sub driver | 🟡 |
| Channel authorization routes | 🟡 |
| WebSocket upgrade handling | 🟡 |
| Client-side event support (whisper) | 🟢 |

### File Storage

| Feature | Priority |
|---------|----------|
| **Real S3 driver** (using aws-sdk-go-v2) | 🟡 |
| SFTP driver | 🟢 |
| GCS (Google Cloud Storage) driver | 🟢 |
| `UploadedFile` wrapper with `Store()` / `StoreAs()` | 🔴 |
| Temporary URLs | 🟡 |
| Visibility (public/private) | 🟡 |
| `storage:link` command | 🟡 |
| Streaming uploads/downloads | 🟡 |
| File fake for testing | 🟡 |

### Error Handling

| Feature | Priority |
|---------|----------|
| Custom error pages (404/500/503) | 🔴 |
| Dev error page (Whoops-like with stack trace) | 🟡 |
| Renderable exceptions | 🟡 |
| Reportable exceptions | 🟡 |
| Ignored exceptions list | 🟡 |
| Exception context | 🟡 |
| Sentry/Bugsnag integration hooks | 🟢 |

### Logging

| Feature | Priority |
|---------|----------|
| Daily rotation | 🔴 |
| Stack channels | 🟡 |
| Slack channel | 🟢 |
| Syslog / Papertrail | 🟢 |
| Contextual data enrichment | 🟡 |
| Request ID injection | 🟡 |
| Structured JSON logging | 🟡 |

### Config

| Feature | Priority |
|---------|----------|
| Config caching (`config:cache`) | 🟡 |
| Nested dot-notation `Get("database.connections.mysql.host")` | 🔴 |
| Environment-specific overrides | 🟡 |
| Typed access helpers (`GetInt()`, `GetBool()`, `GetDuration()`) | 🟡 |

### Localization

| Feature | Priority |
|---------|----------|
| Pluralization rules (CLDR) | 🟡 |
| Locale detection middleware | 🟡 |
| `TransChoice()` | 🟡 |
| JSON translation files | 🟡 |

---

## 🟡 MISSING CLI / Artisan Commands

### Migration Commands
| Command | Priority |
|---------|----------|
| `migrate:fresh` (drop all + migrate) | 🔴 |
| `migrate:refresh` (rollback + migrate) | 🟡 |
| `migrate:status` | 🟡 |
| `migrate:reset` | 🟡 |

### Missing Generators (beyond existing make:*)
| Command | Priority |
|---------|----------|
| `make:request` | 🟡 |
| `make:resource` | 🟡 |
| `make:event` | 🟡 |
| `make:listener` | 🟡 |
| `make:job` | 🟡 |
| `make:mail` | 🟡 |
| `make:notification` | 🟡 |
| `make:policy` | 🟡 |
| `make:rule` | 🟡 |
| `make:seeder` | 🟡 |
| `make:factory` | 🟡 |
| `make:test` | 🟡 |
| `make:observer` | 🟡 |
| `make:cast` | 🟢 |

### Utility Commands
| Command | Description | Priority |
|---------|-------------|----------|
| `serve` | Start dev server with hot reload | 🔴 |
| `route:list` | List all registered routes | 🔴 |
| `route:cache` / `route:clear` | Cache compiled routes | 🟡 |
| `key:generate` | Generate APP_KEY | 🔴 |
| `tinker` / REPL | Interactive Go console | 🟡 |
| `optimize` / `optimize:clear` | Cache config/routes/views | 🟡 |
| `down` / `up` | Maintenance mode (secret, retry-after) | 🟡 |
| `event:list` | List all events and listeners | 🟡 |
| `db:show` / `db:table` | Database info display | 🟡 |
| `storage:link` | Create public storage symlink | 🟡 |
| `about` | Framework/environment info | 🟢 |

---

## 🟡 MISSING MIDDLEWARE

| Middleware | Description | Priority |
|------------|-------------|----------|
| `PreventRequestsDuringMaintenance` | Block requests in maintenance mode | 🟡 (maintenance.go exists, verify behavior) |
| `ValidatePostSize` | Reject oversized POST bodies | 🟡 |
| `EnsureEmailIsVerified` | Block unverified users | 🟡 |
| `SetCacheHeaders` | Auto-set cache control headers | 🟡 |
| `ConvertEmptyStringsToNull` | Normalize input | 🟡 |
| Middleware parameters | `throttle:60,1` — pass args to middleware | 🔴 |
| Middleware groups | Named groups (`web`, `api`) | 🔴 |
| Middleware aliases | Short names for middleware classes | 🟡 |
| Middleware priority / sorting | Control execution order | 🟡 |

---

## 🟡 CROSS-CUTTING MISSING

| Feature | Description | Priority |
|---------|-------------|----------|
| HTTP Client — retry with backoff | Automatic retry on failure | 🟡 |
| HTTP Client — fake for testing | `Http.Fake()` for test stubs | 🟡 |
| HTTP Client — async requests | Concurrent HTTP calls | 🟡 |
| Encryption — key rotation | Rotate APP_KEY without breaking data | 🟡 |
| Pipeline class | Generic pipeline usable outside middleware | 🟡 |
| Streaming responses | `StreamedResponse`, `streamDownload()` | 🟡 |
| Vite / Asset Bundling | `@vite` directive, asset versioning | 🟡 |
| Contracts package | All interfaces as importable package | 🟡 |
| Content negotiation | `$request->expectsJson()`, `accepts()` | 🔴 |
| `Conditionable` / `Tappable` | `when()` / `unless()` / `tap()` chaining | 🟡 |
| `Once` helper | Memoize function results | 🟢 |
| `Defer` helper | Defer work until response sent | 🟢 |
| `Concurrency` helper | Run closures concurrently | 🟢 |
| `Context` propagation | Cross-layer context flowing to jobs/logs | 🟡 |

---

## 🟢 MISSING TESTING UTILITIES

| Feature | Priority |
|---------|----------|
| `AssertDatabaseHas` / `AssertDatabaseMissing` | 🔴 |
| `AssertDatabaseCount` | 🟡 |
| `AssertSoftDeleted` | 🟡 |
| `ActingAs()` auth helper | 🔴 |
| Time travel (`Travel()`, `FreezeTime()`) | 🟡 |
| `RefreshDatabase` (transaction rollback per test) | 🔴 |
| View assertions (`AssertViewHas`, `AssertViewIs`) | 🟡 |
| Command testing (`ExpectsOutput`) | 🟡 |
| `AssertRedirectToRoute` | 🟡 |
| `AssertValid` / `AssertInvalid` | 🟡 |
| Mail fake (`Mail.AssertSent`) | 🟡 |
| Notification fake | 🟡 |
| Event fake | 🟡 |
| Queue fake | 🟡 |
| Parallel testing | 🟢 |

---

## 🟢 MISSING ECOSYSTEM PACKAGES (First-Party)

| Package | Laravel Equivalent | Description | Priority |
|---------|-------------------|-------------|----------|
| Socialite | `laravel/socialite` | OAuth social login (Google, GitHub, etc.) | 🟡 |
| Scout | `laravel/scout` | Full-text search (Meilisearch, Algolia) | 🟡 |
| Telescope | `laravel/telescope` | Debug & monitoring dashboard | 🟡 |
| Horizon | `laravel/horizon` | Queue monitoring dashboard | 🟢 |
| Passport | `laravel/passport` | Full OAuth2 server | 🟢 |
| Cashier | `laravel/cashier` | Subscription billing (Stripe/Paddle) | 🟢 |
| Breeze | `laravel/breeze` | Auth starter kit scaffolding | 🟡 |
| Jetstream | `laravel/jetstream` | Advanced starter with teams | 🟢 |
| Dusk | `laravel/dusk` | Browser testing | 🟢 |
| Precognition | `laravel/precognition` | Live validation | 🟢 |
| Folio | `laravel/folio` | File-based routing | 🟢 |
| Pint | `laravel/pint` | Code style fixer (Go has `gofmt`) | 🟢 |
| Sail | `laravel/sail` | Docker dev environment | 🟢 |
| Pulse | `laravel/pulse` | APM monitoring | 🟢 |
| Reverb | `laravel/reverb` | First-party WebSocket server | 🟡 |
| Prompts | `laravel/prompts` | Rich CLI UX (spinners, multiselect) | 🟡 |
| Echo | `laravel/echo` | JS WebSocket client library | 🟡 |

---

## Summary Counts

| Priority | Count |
|----------|-------|
| 🔴 Critical (Bugs + Must-Have) | ~65 items |
| 🟡 Important (Production Ready) | ~130 items |
| 🟢 Nice to Have | ~55 items |
| **Total** | **~250 features** |

> **Note**: This excludes the ~17 ecosystem packages which would add another ~120+ features combined.
