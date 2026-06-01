# GoW Framework — Feature Implementation Tracker

> **Created**: 2026-06-01 12:15
> **Purpose**: Track implementation progress of missing Laravel features
> **Last Updated**: 2026-06-01 22:00

---

## Summary Dashboard

| Category | Total | Completed | In Progress | Pending |
|----------|-------|-----------|-------------|---------|
| Critical Bugs | 6 | 6 | 0 | 0 |
| Critical Missing | 30 | 30 | 0 | 0 |
| Important Missing | 130 | 79 | 0 | 51 |
| Nice to Have | 84 | 38 | 0 | 46 |
| **Total** | **250** | **153** | **0** | **97** |

---

## Phase 1: ORM Completeness ✅ COMPLETED

| Feature | Status | Date Completed | Notes |
|---------|--------|----------------|-------|
| Model Factories (Make[T], MakeMany[T], Seed[T]) | ✅ Done | 2026-05-22 | database/orm/query.go |
| WhereHas / WhereDoesntHave / WhereRelation | ✅ Done | 2026-05-22 | database/orm/query.go |
| UUID primary key generation | ✅ Done | 2026-05-22 | UUID(), UUIDTrait, UUIDModel |
| Lazy loading prevention | ✅ Done | 2026-05-22 | SetPreventLazyLoading() |

---

## Phase 2: Validation & Forms ✅ COMPLETED

| Feature | Status | Date Completed | Notes |
|---------|--------|----------------|-------|
| Rule objects (10 rules) | ✅ Done | 2026-05-22 | validation/validator.go |
| ValidateWithRules() | ✅ Done | 2026-05-22 | validation/validator.go |
| ValidateWithHooks() (before/after) | ✅ Done | 2026-05-22 | validation/validator.go |
| ValidateSometimes() + When() | ✅ Done | 2026-05-22 | validation/validator.go |
| SetAttributeName() / GetAttributeName() | ✅ Done | 2026-05-22 | validation/validator.go |

---

## Phase 3: Queue & Jobs ✅ COMPLETED

| Feature | Status | Date Completed | Notes |
|---------|--------|----------------|-------|
| Failed jobs tracking | ✅ Done | 2026-05-22 | queue/batch.go |
| Job batching with callbacks | ✅ Done | 2026-05-22 | JobBatch with Then/Catch |
| Job chaining | ✅ Done | 2026-05-22 | JobChain with Then/Catch |
| Job middleware (3 types) | ✅ Done | 2026-05-22 | Retry/Timeout/Logging |
| Delayed dispatch | ✅ Done | 2026-05-22 | Bus.Delay() |

---

## Phase 4: Routing & Middleware ✅ COMPLETED

| Feature | Status | Date Completed | Notes |
|---------|--------|----------------|-------|
| Redirect routes | ✅ Done | 2026-05-22 | routing/router.go |
| Fallback routes | ✅ Done | 2026-05-22 | routing/router.go |

---

## Phase 5: Views & Components ✅ COMPLETED

| Feature | Status | Date Completed | Notes |
|---------|--------|----------------|-------|
| @json directive | ✅ Done | 2026-05-22 | view/engine.go |

---

## Phase 6: Drivers & Infrastructure ✅ COMPLETED

| Feature | Status | Date Completed | Notes |
|---------|--------|----------------|-------|
| Session database driver | ✅ Done | 2026-05-22 | session/driver_database.go |
| Cache file driver | ✅ Done | 2026-05-22 | cache/driver_file.go |
| Cache tags | ✅ Done | 2026-05-22 | cache/tags.go |
| Mailgun driver | ✅ Done | 2026-05-22 | mail/mailer.go |
| Postmark driver | ✅ Done | 2026-05-22 | mail/mailer.go |
| SES driver | ✅ Done | 2026-05-22 | mail/mailer.go |

---

## Phase 7: Auth & Socialite ✅ COMPLETED

| Feature | Status | Date Completed | Notes |
|---------|--------|----------------|-------|
| Socialite Facebook provider | ✅ Done | 2026-05-22 | auth/socialite/providers_extra.go |
| Socialite Twitter/X provider | ✅ Done | 2026-05-22 | auth/socialite/providers_extra.go |
| Socialite LinkedIn provider | ✅ Done | 2026-05-22 | auth/socialite/providers_extra.go |

---

## 🔴 CRITICAL BUGS — All Verified Fixed ✅

| # | Bug | File | Status | Assigned | Target Date | Date Fixed |
|---|-----|------|--------|----------|-------------|------------|
| 1 | Eager Loading Returns Nil | database/orm/relation.go | ✅ Done | — | — | 2026-06-01 |
| 2 | Model Observers Never Called | database/orm/query.go | ✅ Done | — | — | 2026-06-01 |
| 3 | Global Scopes Never Applied | database/orm/model.go | ✅ Done | — | — | 2026-06-01 |
| 4 | Collection.Chunk() Returns Nil | support/collection/collection.go | ✅ Done | — | — | 2026-06-01 |
| 5 | @while/@once Compile to Comments | view/compiler.go | ✅ Done | — | — | 2026-06-01 |
| 6 | S3Driver Methods Are No-Ops | storage/s3_driver.go | ✅ Done | — | — | 2026-06-01 |

> All 6 bugs verified fixed via `go test`. Tests pass: TestEagerLoading, TestORMEvents, TestGlobalScopes, TestChunk, TestCompilerDirectives/while_directive, TestCompilerDirectives/once_directive. S3Driver fully implemented with real AWS SDK v2 calls.

---

## 🔴 CRITICAL MISSING — Required for Any Real Application

### Routing

| # | Feature | Laravel Equivalent | Status | Assigned | Target Date | Date Fixed | Notes |
|---|---------|-------------------|--------|----------|-------------|------------|-------|
| 1 | Route Model Binding (implicit) | `Route::get('/users/{user}', ...)` | ✅ Done | — | — | 2026-05-22 | routing/dispatcher.go |
| 2 | Route Model Binding (explicit) | `Route::bind('user', fn($value) => ...)` | ✅ Done | — | — | 2026-06-01 | router.go Bind() |
| 3 | Optional Parameters | `{id?}` with default value | ✅ Done | — | — | 2026-06-01 | Route.Default() |
| 4 | Regex Constraints | `->where('id', '[0-9]+')` | ✅ Done | — | — | 2026-06-01 | Route.With() |
| 5 | URL Generation | `route('user.show', ['id' => 1])` | ✅ Done | — | — | 2026-06-01 | Router.Route() |
| 6 | Fallback Route | `Route::fallback(fn() => ...)` | ✅ Done | — | — | 2026-05-22 | routing/router.go |
| 7 | Redirect Routes | `Route::redirect('/here', '/there')` | ✅ Done | — | — | 2026-05-22 | routing/router.go |
| 8 | View Routes | `Route::view('/welcome', 'welcome')` | ✅ Done | — | — | 2026-06-01 | Router.View() |

### Database / ORM

| # | Feature | Laravel Equivalent | Status | Assigned | Target Date | Date Fixed | Notes |
|---|---------|-------------------|--------|----------|-------------|------------|-------|
| 1 | Transactions | `DB.Transaction(func(tx) error)` | ✅ Done | — | — | 2026-05-22 | query.go Transaction() |
| 2 | Soft Deletes | `WithTrashed()`, `OnlyTrashed()`, `Restore()` | ✅ Done | — | — | 2026-05-22 | query.go |
| 3 | Working Relationships | HasOne, HasMany, BelongsTo | ✅ Done | — | — | 2026-05-22 | relation.go |
| 4 | BelongsToMany (pivot) | attach, detach, sync, toggle | ✅ Done | — | — | 2026-05-22 | relation.go |
| 5 | Raw Expressions | `DB.Raw("COALESCE(col, 0)")` | ✅ Done | — | — | 2026-05-22 | query builder WhereRaw |
| 6 | GROUP BY / HAVING | Group aggregation queries | ✅ Done | — | — | 2026-05-22 | builder.go |
| 7 | Distinct | `SELECT DISTINCT` | ✅ Done | — | — | 2026-06-01 | builder.go Distinct() |
| 8 | Find / FindOrFail | `User.Find(1)` with 404 | ✅ Done | — | — | 2026-06-01 | query.go FindOrFail() |
| 9 | FirstOrCreate / FirstOrNew | Find-or-create patterns | ✅ Done | — | — | 2026-06-01 | query.go |
| 10 | Mass Assignment (fillable/guarded) | Enforce during Insert/Update | ✅ Done | — | — | 2026-05-22 | model.go |
| 11 | toJSON / toArray serialization | Model serialization | ✅ Done | — | — | 2026-06-01 | query.go ToJSON/ToArray |
| 12 | Query Builder integration | `builder.Paginate(15)` | ✅ Done | — | — | 2026-05-22 | query.go Paginate() |
| 13 | JSON response meta/links | Standard API pagination | ✅ Done | — | — | 2026-05-22 | pagination package |

### Views / Template Engine

| # | Feature | Laravel Equivalent | Status | Assigned | Target Date | Date Fixed | Notes |
|---|---------|-------------------|--------|----------|-------------|------------|-------|
| 1 | @extends layout resolution | Engine resolves parent template | ✅ Done | — | — | 2026-05-22 | engine.go |
| 2 | $loop variable | `$loop.index`, `$loop.first`, `$loop.last` | ✅ Done | — | — | 2026-05-22 | loop.go |
| 3 | {!! !!} unescaped output | Raw HTML output | ✅ Done | — | — | 2026-05-22 | compiler.go |
| 4 | @method | Form method spoofing | ✅ Done | — | — | 2026-06-01 | compiler.go |
| 5 | @error directive | Display validation errors | ✅ Done | — | — | 2026-06-01 | compiler.go + engine.go |

### Authentication & Authorization

| # | Feature | Laravel Equivalent | Status | Assigned | Target Date | Date Fixed | Notes |
|---|---------|-------------------|--------|----------|-------------|------------|-------|
| 1 | SessionGuard implementation | Session-based auth persistence | ✅ Done | — | — | 2026-06-01 | auth/manager.go Login/Logout/User/ID |
| 2 | Password Reset flow | Tokens, email, broker, reset endpoint | ✅ Done | — | — | 2026-06-01 | auth/password/broker.go Reset() |
| 3 | Email Verification | Signed URLs, verification middleware | ✅ Done | — | — | 2026-06-01 | auth/verification/verification.go MarkAsVerified() |
| 4 | Session Fixation Protection | Regenerate session ID on login | ✅ Done | — | — | 2026-06-01 | auth/middleware.go SessionFixationMiddleware() + RegenerateSession() |
| 5 | Sanctum token expiration check | Currently accepts expired tokens | ✅ Done | — | — | 2026-06-01 | auth/sanctum/token.go IsExpired() + middleware check |
| 6 | Sanctum token abilities verification | `tokenCan('server:update')` | ✅ Done | — | — | 2026-06-01 | auth/sanctum/token.go HasAbility/Can/Cannot/RequireAbilities |

### HTTP Request / Response

| # | Feature | Laravel Equivalent | Status | Assigned | Target Date | Date Fixed | Notes |
|---|---------|-------------------|--------|----------|-------------|------------|-------|
| 1 | Input(), Query(), Post() | Convenient request data access | ✅ Done | — | — | 2026-06-01 | request.go Post() |
| 2 | Only(), Except() | Filter request data | ✅ Done | — | — | 2026-05-22 | request.go |
| 3 | HasFile(), File() | File upload handling | ✅ Done | — | — | 2026-06-01 | request.go HasFile()/File() |
| 4 | ExpectsJson(), WantsJson() | Content negotiation | ✅ Done | — | — | 2026-05-22 | request.go |
| 5 | Redirect with flash data | `redirect()->with('status', 'saved')` | ✅ Done | — | — | 2026-06-01 | response/redirect.go Flash() |
| 6 | Redirect to named route | `redirect()->route('dashboard')` | ✅ Done | — | — | 2026-06-01 | response/redirect.go RedirectToRoute() |
| 7 | Redirect back | `redirect()->back()` | ✅ Done | — | — | 2026-06-01 | response/redirect.go RedirectBack() |
| 8 | Old Input | `old('email')` after validation failure | ✅ Done | — | — | 2026-05-22 | request.go Old() |

### Validation

| # | Rule | Status | Assigned | Target Date | Date Fixed |
|---|------|--------|----------|-------------|------------|
| 1 | boolean, integer, array | ✅ Done | — | — | 2026-05-22 |
| 2 | date, before, after, date_format | ✅ Done | — | — | 2026-06-01 |
| 3 | image, mimes, mimetypes, file | ✅ Done | — | — | 2026-06-01 | validation/validator.go image/mimes/file rules |
| 4 | url, active_url, ip, ipv4, ipv6 | ✅ Done | — | — | 2026-06-01 |
| 5 | json, uuid, ulid | ✅ Done | — | — | 2026-05-22 |
| 6 | starts_with, ends_with, contains | ✅ Done | — | — | 2026-06-01 |
| 7 | alpha, alpha_num, alpha_dash | ✅ Done | — | — | 2026-05-22 |
| 8 | digits, digits_between | ✅ Done | — | — | 2026-06-01 |
| 9 | gt, gte, lt, lte (compare fields) | ✅ Done | — | — | 2026-05-22 |
| 10 | same, different (compare fields) | ✅ Done | — | — | 2026-05-22 |
| 11 | required_if, required_unless, required_with, required_without | ✅ Done | — | — | 2026-06-01 |
| 12 | nullable, sometimes, bail | ✅ Done | — | — | 2026-05-22 |
| 13 | Nested array validation (items.*.name) | ✅ Done | — | — | 2026-06-01 | validation/validator.go expandWildcard() |

---

## 🟡 IMPORTANT MISSING — Production Readiness

### Cache

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Atomic locks (distributed) | ✅ Done | — | — | 2026-06-01 | cache/helpers.go AcquireLock/NewLock |
| 2 | Remember() / RememberForever() | ✅ Done | — | — | 2026-06-01 | cache/helpers.go Remember/RememberForever |
| 3 | Many() / PutMany() bulk ops | ✅ Done | — | — | 2026-06-01 | cache/helpers.go Many/PutMany/PutManyForever |
| 4 | Cache events (hit/miss/write/delete) | ✅ Done | — | — | 2026-06-01 | cache/store.go EventType/CacheEvent/EventDispatcher |
| 5 | Cache prefix support | ✅ Done | — | — | 2026-06-01 | All drivers: Has(), prefix in RedisStore |

### Session

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Array driver (testing) | ✅ Done | — | — | 2026-06-01 | session/driver_array.go ArrayDriver |
| 2 | Session encryption | ✅ Done | — | — | 2026-06-01 | session/encrypted.go EncryptedStore/NewEncryptedStore |
| 3 | Previous URL storage (intended()) | ✅ Done | — | — | 2026-06-01 | session/manager.go PreviousURL/SetPreviousURL/Intended/SetIntendedURL |

### Queue

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Unique jobs (ShouldBeUnique) | ✅ Done | — | — | 2026-06-01 | queue/job.go ShouldBeUnique/ShouldBeUniqueUntilProcessing |
| 2 | Job timeout / max attempts / backoff | ✅ Done | — | — | 2026-06-01 | queue/job.go HasTimeout/HasMaxRetries/HasBackoff + worker.go |
| 3 | Job rate limiting | ✅ Done | — | — | 2026-06-01 | queue/middleware.go RateLimitMiddleware/RateLimitByKeyMiddleware/EnsureUnique |
| 4 | Queue priorities | ✅ Done | — | — | 2026-06-01 | queue/job.go HasPriority interface |
| 5 | SQS / RabbitMQ drivers | ✅ Done | — | — | 2026-06-01 | queue/driver_sqs_rabbit.go SQSDriver/RabbitMQDriver |

### Events

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Queued listeners (ShouldQueue) | ✅ Done | — | — | 2026-06-01 | events/manager.go ShouldQueue interface/QueueListener |
| 2 | Event subscribers | ✅ Done | — | — | 2026-05-22 | events/manager.go Subscriber interface |
| 3 | Stoppable events | ✅ Done | — | — | 2026-06-01 | events/manager.go StoppableEvent/StopPropagation |
| 4 | Event fake for testing | ✅ Done | — | — | 2026-06-01 | events/fake.go Fake |

### Mail

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | SMTP driver (actual email sending) | ✅ Done | — | — | 2026-05-22 | mail/mailer.go SmtpDriver with STARTTLS/TLS |
| 2 | Markdown mail support | ✅ Done | — | — | 2026-05-22 | mail/markdown.go RenderMarkdown |
| 3 | Attachments (files, raw data, from storage) | ✅ Done | — | — | 2026-06-01 | mail/message.go Attach/AttachFile/AttachRaw |
| 4 | CC, BCC, Reply-To | ✅ Done | — | — | 2026-06-01 | mail/message.go Cc/Bcc/ReplyTo |
| 5 | Mail queuing (ShouldQueue) | ✅ Done | — | — | 2026-05-22 | mail/mailer.go Queue/QueueNow |
| 6 | Mail preview in browser | ✅ Done | — | — | 2026-06-01 | mail/mailer.go Preview method |
| 7 | Mail fake for testing | ✅ Done | — | — | 2026-06-01 | mail/fake.go MailFake |

### Notifications

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Mail channel implementation | ✅ Done | — | — | 2026-05-22 | notifications/channel_mail.go |
| 2 | Slack channel | ✅ Done | — | — | 2026-05-22 | notifications/channel_slack.go |
| 3 | On-demand notifications | ✅ Done | — | — | 2026-06-01 | notifications/manager.go OnDemand/SendNow |
| 4 | Mark as read/unread | ✅ Done | — | — | 2026-06-01 | notifications/channel_database.go MarkAsRead/MarkAsUnread/MarkAllAsRead |
| 5 | Queued notifications | ✅ Done | — | — | 2026-06-01 | notifications/manager.go SendQueued/NotifyQueued/QueueManager interface |

### Broadcasting

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | WebSocket server (Reverb equivalent) | ✅ Done | — | — | 2026-06-01 | broadcasting/server.go WebSocketHub |
| 2 | Pusher driver | ✅ Done | — | — | 2026-06-01 | broadcasting/redis.go PusherBroadcaster |
| 3 | Redis pub/sub driver | ✅ Done | — | — | 2026-06-01 | broadcasting/redis.go RedisBroadcaster |
| 4 | Channel authorization routes | ✅ Done | — | — | 2026-06-01 | broadcasting/auth.go AuthController/ChannelAuthorizer |

### File Storage

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Real S3 driver (using aws-sdk-go-v2) | ✅ Done | — | — | 2026-06-01 | storage/s3_driver.go S3Driver with full FilesystemWithExtras |
| 2 | UploadedFile wrapper with Store() / StoreAs() | ✅ Done | — | — | 2026-06-01 | storage/storage.go UploadedFile with Store/StoreAs/Contents/Move/TemporaryURL |
| 3 | Temporary URLs | ✅ Done | — | — | 2026-06-01 | storage/storage.go LocalDriver.TemporaryURL + UploadedFile.TemporaryURL |
| 4 | Visibility (public/private) | ✅ Done | — | — | 2026-06-01 | storage/storage.go LocalDriver.SetVisibility/GetVisibility |
| 5 | storage:link command | ✅ Done | — | — | 2026-06-01 | cmd/artisan/maintenance_commands.go StorageLinkCmd |

### Error Handling

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Custom error pages (404/500/503) | ✅ Done | — | — | 2026-05-22 |
| 2 | Dev error page (Whoops-like) | ✅ Done | — | — | 2026-06-01 | http/exception/exception.go APP_DEBUG debug page |
| 3 | Renderable exceptions | ✅ Done | — | — | 2026-06-01 | http/exception/exception.go Renderable interface/Handler |
| 4 | Reportable exceptions | ✅ Done | — | — | 2026-06-01 | http/exception/exception.go Reportable interface/WithReport |

### Logging

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Daily rotation | ✅ Done | — | — | 2026-06-01 | logging/logger.go dailyWriter with MaxAge cleanup |
| 2 | Stack channels | ✅ Done | — | — | 2026-05-22 | logging/logger.go io.MultiWriter stack |
| 3 | Contextual data enrichment | ✅ Done | — | — | 2026-06-01 | logging/logger.go With/WithContext/LogRequest/LogError |
| 4 | Request ID injection | ✅ Done | — | — | 2026-06-01 | logging/logger.go RequestMiddleware + X-Request-ID |
| 5 | Structured JSON logging | ✅ Done | — | — | 2026-05-22 | logging/logger.go slog.NewJSONHandler |

### Config

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Config caching (config:cache) | ✅ Done | — | — | 2026-05-22 | config/cached.go CachedValues |
| 2 | Nested dot-notation Get() | ✅ Done | — | — | 2026-06-01 | config/repository.go Get() with dot-notation |
| 3 | Environment-specific overrides | ✅ Done | — | — | 2026-05-22 | config/repository.go .env loading |
| 4 | Typed access helpers (GetInt, GetBool, GetFloat) | ✅ Done | — | — | 2026-06-01 | config/repository.go GetInt/GetBool/GetFloat |

### Localization

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Pluralization rules (CLDR) | ✅ Done | — | — | 2026-06-01 | localization/translator.go TransChoice/choosePluralForm |
| 2 | Locale detection middleware | ✅ Done | — | — | 2026-06-01 | localization/translator.go DetectLocaleFromRequest/DetectLocaleAndSet |
| 3 | TransChoice() | ✅ Done | — | — | 2026-06-01 | localization/translator.go TransChoice |
| 4 | JSON translation files | ✅ Done | — | — | 2026-05-22 | localization/translator.go Load |

### Middleware

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | ValidatePostSize | ✅ Done | — | — | 2026-06-01 | http/middleware/post_size_cache.go |
| 2 | EnsureEmailIsVerified | ✅ Done | — | — | 2026-06-01 | auth/require_verified.go RequireVerified |
| 3 | SetCacheHeaders | ✅ Done | — | — | 2026-06-01 | http/middleware/post_size_cache.go SetCacheHeaders/CacheFor/CachePrivate |
| 4 | ConvertEmptyStringsToNull | ✅ Done | — | — | 2026-06-01 | http/middleware/post_size_cache.go |
| 5 | Middleware parameters | ✅ Done | — | — | 2026-05-22 | routing/router.go |
| 6 | Middleware groups (web, api) | ✅ Done | — | — | 2026-05-22 | routing/router.go GroupMiddleware |
| 7 | Middleware aliases | ✅ Done | — | — | 2026-05-22 | routing/router.go Alias |

### Missing CLI / Artisan Commands

| # | Command | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | migrate:fresh | ✅ Done | — | — | 2026-05-22 | cmd/artisan/migrate_commands.go MigrateFreshCmd |
| 2 | migrate:refresh | ✅ Done | — | — | 2026-05-22 | cmd/artisan/migrate_commands.go MigrateRefreshCmd |
| 3 | migrate:status | ✅ Done | — | — | 2026-05-22 | cmd/artisan/migrate_commands.go MigrateStatusCmd |
| 4 | migrate:reset | ✅ Done | — | — | 2026-06-01 | cmd/artisan/migrate_commands.go MigrateResetCmd |
| 5 | make:request | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_request.go MakeRequestCmd |
| 6 | make:resource | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_more.go MakeResourceCmd |
| 7 | make:event | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_more.go MakeEventCmd |
| 8 | make:listener | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_more.go MakeListenerCmd |
| 9 | make:job | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_commands.go MakeJobCmd |
| 10 | make:mail | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_more.go MakeMailCmd |
| 11 | make:notification | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_more.go MakeNotificationCmd |
| 12 | make:policy | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_more.go MakePolicyCmd |
| 13 | make:rule | ✅ Done | — | — | 2026-06-01 | cmd/artisan/make_more.go MakeRuleCmd |
| 14 | make:seeder | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_seeder.go MakeSeederCmd |
| 15 | make:factory | ✅ Done | — | — | 2026-06-01 | cmd/artisan/make_more.go MakeFactoryCmd |
| 16 | make:test | ✅ Done | — | — | 2026-05-22 | cmd/artisan/make_more.go MakeTestCmd |
| 17 | make:observer | ✅ Done | — | — | 2026-06-01 | cmd/artisan/make_more.go MakeObserverCmd |
| 18 | serve (dev server) | ✅ Done | — | — | 2026-06-01 | cmd/artisan/serve_command.go ServeCmd |
| 19 | route:list | ✅ Done | — | — | 2026-05-22 | cmd/artisan/route_list.go RouteListCmd |
| 20 | route:cache / route:clear | ✅ Done | — | — | 2026-06-01 | cmd/artisan/route_cache.go RouteCacheCmd/RouteClearCmd |
| 21 | key:generate | ✅ Done | — | — | 2026-05-22 | cmd/artisan/key_generate.go KeyGenerateCmd |
| 22 | tinker / REPL | ✅ Done | — | — | 2026-05-22 | cmd/artisan/tinker.go |
| 23 | optimize / optimize:clear | ✅ Done | — | — | 2026-06-01 | cmd/artisan/optimize_commands.go OptimizeCmd/OptimizeClearCmd |
| 24 | down / up (maintenance mode) | ✅ Done | — | — | 2026-05-22 | cmd/artisan/maintenance_commands.go |
| 25 | event:list | ✅ Done | — | — | 2026-06-01 | cmd/artisan/maintenance_commands.go EventListCmd |
| 26 | db:show / db:table | ✅ Done | — | — | 2026-06-01 | cmd/artisan/db_commands.go DbShowCmd/DbTableCmd |
| 27 | storage:link | ✅ Done | — | — | 2026-06-01 | cmd/artisan/maintenance_commands.go StorageLinkCmd |
| 28 | horizon | ✅ Done | — | — | 2026-06-01 | cmd/artisan/dashboard_commands.go HorizonCmd |
| 29 | pulse | ✅ Done | — | — | 2026-06-01 | cmd/artisan/dashboard_commands.go PulseCmd |

---

## 🟡 MISSING TESTING UTILITIES

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | AssertDatabaseHas / AssertDatabaseMissing | ✅ Done | — | — | 2026-05-22 | testing/testing.go |
| 2 | AssertDatabaseCount | ✅ Done | — | — | 2026-05-22 | testing/testing.go AssertDatabaseCount/HasNoRecords/HasExactly |
| 3 | AssertSoftDeleted | ✅ Done | — | — | 2026-06-01 | testing/testing.go AssertSoftDeleted/AssertNotSoftDeleted/AssertDatabaseTable |
| 4 | ActingAs() auth helper | ✅ Done | — | — | 2026-05-22 | testing/testing.go ActingAs |
| 5 | Time travel (Travel, FreezeTime) | ✅ Done | — | — | 2026-06-01 | testing/time_travel.go TimeTravel/FreezeTime/TravelTime |
| 6 | RefreshDatabase (transaction rollback) | ✅ Done | — | — | 2026-05-22 | testing/db.go RefreshDatabase |
| 7 | View assertions (AssertViewHas, AssertViewIs) | ✅ Done | — | — | 2026-06-01 | testing/testing.go AssertViewIs/AssertViewHas/AssertViewHasValue |
| 8 | Command testing (ExpectsOutput) | ✅ Done | — | — | 2026-06-01 | testing/testing.go ArtisanTestCase/AssertOutputContains/AssertExitCode |
| 9 | Mail fake (Mail.AssertSent) | ✅ Done | — | — | 2026-06-01 | mail/fake.go MailFake with AssertSent |
| 10 | Notification fake | ✅ Done | — | — | 2026-06-01 | notifications/fake.go Fake/Send/AssertSent |
| 11 | Event fake | ✅ Done | — | — | 2026-06-01 | events/fake.go Fake/Dispatch/AssertDispatched |
| 12 | Queue fake | ✅ Done | — | — | 2026-06-01 | queue/fake.go Fake/Push/GetJobs/Sync |

---

## 🟢 NICE TO HAVE — Ecosystem & DX

| # | Feature | Status | Assigned | Target Date | Date Fixed |
|---|---------|--------|----------|-------------|------------|
| 1 | Horizon (queue dashboard) | ✅ Done | — | — | 2026-06-01 | horizon/dashboard.go Dashboard/Store/InMemoryStore |
| 2 | Nova (admin panel) | ✅ Done | — | — | 2026-06-01 | nova/admin.go Admin/Resource/Panel |
| 3 | Sail (Docker dev environment) | ✅ Done | — | — | 2026-06-01 | sail/sail.go Config/Init/createDockerCompose |
| 4 | Pulse (health monitoring) | ✅ Done | — | — | 2026-06-01 | pulse/pulse.go Pulse/MetricProvider/CustomProvider |
| 5 | Socialite Apple provider | ✅ Done | — | — | 2026-06-01 | auth/socialite/providers_extra.go AppleProvider |
| 6 | Route group name prefixes | ✅ Done | — | — | 2026-06-01 | routing/router.go NameGroup/AsGroup |
| 7 | Controller middleware assignment | ✅ Done | — | — | 2026-05-22 | routing/router.go |
| 8 | Maintenance mode bypass IPs | ✅ Done | — | — | 2026-06-01 | cmd/artisan/maintenance_commands.go IsAllowedIP/GetRetryAfter/allow flag |
| 9 | TrustHosts middleware | ✅ Done | — | — | 2026-06-01 | http/middleware/trust_hosts.go TrustHostsMiddleware |
| 10 | Blade @vite directive | ✅ Done | — | — | 2026-06-01 | view/compiler.go @vite directive |
| 11 | HTTP Client retry with backoff | ✅ Done | — | — | 2026-06-01 | http/client/client.go WithRetry + exponential backoff |
| 12 | HTTP Client fake for testing | ✅ Done | — | — | 2026-06-01 | http/client/fake.go Fake/When/AssertSent |
| 13 | Pipeline class (generic) | ✅ Done | — | — | 2026-05-22 | support/pipeline/pipeline.go |
| 14 | Streaming responses | ✅ Done | — | — | 2026-06-01 | http/response/stream.go StreamLines/StreamSSE/StreamTimer |
| 15 | Contracts package (interfaces) | ✅ Done | — | — | 2026-06-01 | contracts/contracts.go |
| 16 | Conditionable / Tappable | ✅ Done | — | — | 2026-06-01 | support/helpers.go Conditionable/Tap |
| 17 | Once helper | ✅ Done | — | — | 2026-06-01 | support/helpers.go Once |
| 18 | Defer helper | ✅ Done | — | — | 2026-06-01 | support/helpers.go Defer |
| 19 | Concurrency helper | ✅ Done | — | — | 2026-06-01 | support/helpers.go Parallel/Collect/LimitConcurrency/Retry |
| 20 | Context propagation | ✅ Done | — | — | 2026-06-01 | support/context.go ContextPropagator/ContextWith* helpers |
| 21 | Parallel testing | ✅ Done | — | — | 2026-06-01 | testing/testing.go ParallelTestCase/ConcurrentTestCase |

---

## Implementation Notes

### How to Update This File
1. Change Status from `Pending` to `In Progress` or `Done`
2. Fill in `Assigned` with developer name/handle
3. Set `Target Date` in YYYY-MM-DD format
4. When complete, set `Date Fixed` and move to Completed section
5. Update the Summary Dashboard at top

### Priority Definitions
- **Critical (🔴)**: Blocks any real application from working
- **Important (🟡)**: Required for production readiness
- **Nice to Have (🟢)**: Improves developer experience, can wait

### File Naming Convention
This file is named `FEATURE_TRACKER_YYYY-MM-DD_HH-MM.md` to preserve the creation timestamp.
A new file should be created periodically (monthly or after major milestones) to keep the tracker fresh.
