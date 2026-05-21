# GoW Framework Implementation Plan: Phases B & C

This document serves as the master implementation roadmap for both Phase B (Production Readiness) and Phase C (Ecosystem & Packages) of the GoW framework.

---

## Phase B: Production Readiness

Transitioning the framework from a core foundation to a production-ready system with advanced infrastructure support (Redis, SMTP, Queues), APIs (Sanctum, Resources), and observability (Logging, Error Handling).

### Track 1: Core Infrastructure & Drivers

#### [MODIFY] `cache/manager.go`, `session/manager.go`, `queue/manager.go`
- **Feature**: Redis Drivers.
- Implement Redis stores for Cache, Session, and Queue to allow distributed scaling.

#### [NEW] `queue/database_driver.go` & `cmd/artisan/queue_commands.go`
- **Feature**: Database Queue System.
- Build a driver to store jobs in a `jobs` table.
- Implement a `failed_jobs` table and retry mechanism.
- Add `queue:work` and `queue:retry` commands.

#### [MODIFY] `mail/mailer.go` & `notifications/channel_database.go`
- **Feature**: SMTP & Notifications.
- Add standard SMTP delivery support (attachments, CC, BCC).
- Add a database notification channel to store notifications persistently.

#### [MODIFY] `logging/manager.go`
- **Feature**: Advanced Logging.
- Implement daily log rotation (io.MultiWriter + cron-based file rotation), stack channels, and contextual enrichment.

#### [MODIFY] `database/query/builder.go`
- **Feature**: Query Logging & Profiling.
- Add `QueryEvent` dispatcher to the Builder to capture SQL, execution time, and caller stack for N+1 detection.

### Track 2: Advanced Web & Middleware

#### [MODIFY] `routing/router.go`
- **Feature**: Resource Routes & Macros.
- Implement `.Resource()` and `.ApiResource()` to automatically register CRUD routes.
- Add macro support for custom Request/Response methods.

#### [NEW] `app/http/middleware/trusted_proxies.go`
- **Feature**: Trusted Proxies.
- Middleware to handle X-Forwarded-For headers securely.
- Add `TrimStrings` middleware.

#### [MODIFY] `foundation/application.go` & `config/`
- **Feature**: Error Handling & Config.
- Custom error pages (404, 500), maintenance mode (`artisan down`), and detailed developer error pages.
- Nested dot-notation config caching (`artisan config:cache` generating a compiled Go binary step for zero deserialization overhead).

#### [MODIFY] `http/response.go`
- **Feature**: Streaming Responses.
- Add capabilities to stream large file responses smoothly.

### Track 3: APIs & Authentication

#### [NEW] `http/resources/resource.go`
- **Feature**: API Resources / Transformers.
- Build a data transformation layer to convert ORM models to JSON API payloads.

#### [NEW] `auth/sanctum/`
- **Feature**: Sanctum Equivalent.
- Implement an API token authentication package with abilities/scopes, storing tokens securely.

#### [NEW] `support/httpclient/client.go`
- **Feature**: HTTP Client Wrapper.
- A fluent HTTP client with retry, fake (for testing), and connection pool capabilities.

### Track 4: Tooling & Diagnostics

#### [NEW] `testing/`
- **Feature**: Testing Utilities.
- Provide testing helpers: database assertions (e.g., `AssertDatabaseHas` wrapping Query Builder's Count/Where methods), time travel, and authentication mock helpers.

---

## Phase C: Ecosystem & Packages

Transitioning the framework to a fully-featured ecosystem with robust Authentication/Authorization, advanced ORM features, ecosystem utilities (Health, Feature Flags), and real-time broadcasting capabilities.

### Track 1: Complete Authentication & Authorization

#### [NEW] `auth/fortify/`
- **Feature**: Headless Auth Backend (PC.11, PC.12).
- Build JSON API endpoints: `POST /register`, `POST /login`, `POST /forgot-password`, `POST /reset-password`, `GET /email/verify/{id}/{hash}`, `POST /two-factor/enable`.
- Response format follows standard JSON payloads.

#### [NEW] `auth/authorization/`
- **Feature**: Policies and Gates (PC.13).
- Implement a `Gate` facade to define authorization closures.
- Map controller actions to standard policy methods (`view`, `create`, `update`, `delete`).

#### [MODIFY] `view/goblade/compiler.go`
- **Feature**: Authorization Directives.
- Add `@can` and `@cannot` directives to compile down to Gate checks.

### Track 2: ORM Mastery & Data Management

#### [MODIFY] `database/orm/model.go` & `database/orm/builder.go`
- **Feature**: Advanced ORM Features (PC.9, PC.10).
- Add support for Global/Local Scopes.
- Implement `Observer` interface + event dispatcher hooks (`creating`, `updating`, `deleting`).
- Add `Cast` interface for JSON/encrypted/date attribute casting, plus Accessors/Mutators and Default Values.
- Add ORM strictness (`PreventLazyLoading`, `HasUuids`).

#### [NEW] `cmd/artisan/schema_commands.go`
- **Feature**: Schema Dumping & Model Pruning (PC.7, PC.8).
- `schema:dump` to export SQL schema files.
- `model:prune` to automatically delete obsolete records based on a `Prunable` interface.

### Track 3.0: Core Completeness

#### [NEW] `support/pipeline/`
- **Feature**: Generic Pipeline Class (PB.17).
- Standardize middleware-like behavior via `func(T) (T, error)` chains.

#### [NEW] `storage/`
- **Feature**: File Storage Drivers (PB.18).
- Complete the multi-disk spec with Local driver, S3 driver, and Storage manager.

### Track 3: Ecosystem Utilities

#### [NEW] `routing/signed.go`
- **Feature**: Signed URLs (PC.1).
- Generate and verify cryptographically signed URLs to prevent tampering (e.g., for email verification).

#### [NEW] `support/health/`
- **Feature**: Health Checks (PC.2).
- Build a configurable `/up` endpoint to monitor Database, Redis, and overall system health.

#### [NEW] `support/pennant/`
- **Feature**: Feature Flags (PC.3).
- Implement a `Manager` with `Store` interface. Default to Redis store (`HSET/HGET`), with Database and InMemory fallbacks.

#### [NEW] `support/process/`
- **Feature**: Process Execution Wrapper (PC.6).
- Build a fluent wrapper around `os/exec` that supports fake testing and timeout handling.

### Track 4: Frontend & Real-time

#### [MODIFY] `view/goblade/compiler.go`
- **Feature**: Advanced Blade Directives (PC.4).
- Implement `$loop` variable for iterations, `@once`, `@class`, and `@checked`.
- Add basic support for Blade Components and Slots.

#### [NEW] `broadcasting/` & `events/`
- **Feature**: Broadcasting Server & Advanced Events (PC.14, PC.15).
- Build `broadcasting.Manager` with `PusherDriver` and `RedisDriver` to bridge to external WebSocket servers (e.g., Reverb).
- Provide `routes/channels.go` authorization middleware for private/presence channels.
- Upgrade the Event dispatcher to support queued listeners and Event Subscribers.

#### [NEW] `support/publishing/`
- **Feature**: Package Auto-Discovery & Publishing (PC.5).
- Build `vendor:publish` logic to extract configs, migrations, and assets from third-party Go packages.
