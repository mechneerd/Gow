# GoW Framework — Implementation Context Guide

> **Purpose**: This file gives ANY AI the full context needed to implement features without asking questions. It documents every implemented feature's exact API, every struct shape, every inter-package dependency, and the exact approach for implementing new features.
>
> **Last Updated**: 2026-05-22
> **Go Version**: 1.26.3
> **Module Name**: `gow`

---

## TABLE OF CONTENTS

1. [Module & Dependencies](#1-module--dependencies)
2. [Project Structure](#2-project-structure)
3. [Container (IoC)](#3-container-ioc)
4. [Foundation / Application](#4-foundation--application)
5. [Config](#5-config)
6. [Routing](#6-routing)
7. [HTTP Kernel & Middleware](#7-http-kernel--middleware)
8. [HTTP Request & Response](#8-http-request--response)
9. [Database Layer](#9-database-layer)
10. [ORM (Goquent)](#10-orm-goquent)
11. [Validation](#11-validation)
12. [Session](#12-session)
13. [Auth & Authorization](#13-auth--authorization)
14. [Cache](#14-cache)
15. [Queue](#15-queue)
16. [Events](#16-events)
17. [Mail](#17-mail)
18. [View Engine (Goblade)](#18-view-engine-goblade)
19. [Logging](#19-logging)
20. [Storage](#20-storage)
21. [Console / CLI](#21-console--cli)
22. [Support Packages](#22-support-packages)
23. [Testing](#23-testing)
24. [Implementation Rules](#24-implementation-rules)

---

## 1. Module & Dependencies

```
module gow
go 1.26.3

Direct deps:
- github.com/joho/godotenv v1.5.1    → .env loading
- github.com/spf13/cobra v1.10.2     → CLI framework
- github.com/robfig/cron/v3 v3.0.1   → Scheduled tasks
- github.com/stretchr/testify v1.11.1 → Test assertions
```

**Import Convention**: All internal imports use `gow/` prefix (e.g., `gow/routing`, `gow/database/orm`).

---

## 2. Project Structure

```
gow/
├── artisan.go              # CLI entry point (build tag: ignore)
├── go.mod / go.sum
│
├── container/              # IoC container
│   └── container.go        # Container, Make[T], Bind, Singleton, Instance, Freeze
│
├── foundation/             # App bootstrap
│   ├── application.go      # Application struct (embeds *Container)
│   └── provider.go         # ServiceProvider interface
│
├── config/                 # Configuration
│   ├── repository.go       # Repository (loads .env, Get/Set/GetBool/GetInt)
│   ├── cached.go           # CachedValues map[string]string
│   └── service_provider.go
│
├── routing/                # HTTP router
│   ├── router.go           # Router, Route, radix tree, all HTTP verbs
│   ├── resource.go         # Resource/ApiResource routes, Macros
│   ├── signed.go           # URLGenerator, SignedRoute, HasValidSignature
│   └── dispatcher.go       # Route Model Binding via reflection
│
├── http/
│   ├── kernel.go           # HTTP Kernel (wraps Router with global middleware)
│   ├── exception.go        # HttpException struct
│   ├── middleware/          # 12 middleware files
│   │   ├── auth.go, cors.go, csrf.go, encrypt_cookies.go
│   │   ├── error_handler.go, maintenance.go, recovery.go
│   │   ├── security_headers.go, session.go, throttle.go
│   │   ├── trim_strings.go, trusted_proxies.go
│   ├── request/
│   │   └── form_request.go # FormRequest interface + Validate()
│   ├── response/
│   │   ├── json.go         # JSON(), Ok(), Created(), Error()
│   │   ├── redirect.go     # Redirect() — thin wrapper
│   │   └── stream.go       # Stream(), Download()
│   ├── resources/          # API resource transformers (TBD)
│   └── client/             # HTTP client (TBD)
│
├── database/
│   ├── connection.go       # Connection{DB *sql.DB, Dialect dialect.Dialect}
│   ├── manager.go          # Manager (multi-connection, default conn)
│   ├── dialect/
│   │   ├── dialect.go      # Dialect interface + SelectQuery/WhereClause/JoinClause structs
│   │   └── sqlite.go       # SQLiteDialect (CompileSelect/Insert/Update/Delete)
│   ├── query/
│   │   └── builder.go      # Query Builder (291 lines) — all WHERE types, JOINs, aggregates
│   ├── orm/
│   │   ├── model.go        # Scope, Observer, Cast interfaces + DispatchModelEvent + global registries
│   │   ├── query.go        # ModelQuery[T] — generic CRUD (Find/Insert/Update/Delete/Save/First/Get)
│   │   └── relation.go     # HasOne/HasMany/BelongsTo/BelongsToMany types + eagerLoad (STUB)
│   ├── migration/          # Migrator (runs Up/Down, tracks in migrations table)
│   ├── pagination/         # Paginator struct (basic)
│   ├── schema/             # Schema builder (TBD)
│   ├── seeder/             # Database seeder
│   ├── factory/            # Model factory (interface only)
│   └── redis/              # Redis connection
│
├── validation/
│   └── validator.go        # 14 rules: required, email, min, max, between, size, numeric, string, in, regex, confirmed, unique, exists
│
├── session/
│   ├── store.go            # Store interface (Read/Write/Destroy)
│   ├── manager.go          # Manager (Get/Put/Flash/Reflash/Keep/Old/Regenerate/FlashInput)
│   ├── driver_file.go      # File session driver
│   ├── driver_cookie.go    # Cookie session driver
│   └── driver_redis.go     # Redis session driver
│
├── auth/
│   ├── guard.go            # Guard interface (Check/Guest/User/ID/Attempt/Login/Logout)
│   ├── provider.go         # UserProvider interface
│   ├── manager.go          # AuthManager + SessionGuard + Authenticatable interface
│   ├── access/             # Gate (Define/Allows/Denies)
│   ├── authorization/      # Authorization policies
│   ├── sanctum/
│   │   └── token.go        # TokenManager (CreateToken + Bearer middleware)
│   └── fortify/
│       └── fortify.go      # Headless auth routes (register/login/forgot/reset/verify/2fa) — STUB bodies
│
├── cache/
│   ├── store.go            # Store interface (Get/Put/Increment/Decrement/Forever/Forget/Flush)
│   ├── driver_memory.go    # In-memory cache with TTL
│   ├── driver_redis.go     # Redis cache driver
│   └── rate_limiter.go     # Token bucket rate limiter
│
├── queue/
│   ├── job.go              # Job interface (Handle/Failed)
│   ├── manager.go          # Queue Manager (Connection returns Driver)
│   ├── driver_sync.go      # SyncDriver (channel-based)
│   ├── driver_database.go  # Database queue driver
│   ├── driver_redis.go     # Redis queue driver
│   └── worker.go           # Worker (polls Pop(), processes jobs)
│
├── events/
│   ├── dispatcher.go       # Dispatcher (Listen/ListenAny/Dispatch)
│   └── manager.go          # Event Manager
│
├── mail/
│   ├── mailer.go           # Mailer struct + LogDriver + SmtpDriver (SMTP with attachments/CC/BCC)
│   └── message.go          # Message, Mailable interface, BaseMailable builder, Attachment struct
│
├── notifications/
│   ├── manager.go          # Notification Manager
│   └── channel.go          # Channel interface
│
├── broadcasting/
│   ├── manager.go          # Broadcasting Manager
│   └── channel.go          # LogBroadcaster + Public/Private/Presence channels
│
├── view/
│   ├── compiler.go         # Goblade compiler (regex-based directive → Go template transpilation)
│   ├── engine.go           # Engine (Make(name, data) → compile + parse + execute template)
│   └── composer.go         # View Composers registry
│
├── logging/
│   ├── logger.go           # Setup(Config) → slog.Logger with daily rotation + stack channels + JSON
│   └── service_provider.go
│
├── storage/
│   ├── filesystem.go       # Filesystem interface (Get/Put/Delete/Exists/URL)
│   └── storage.go          # Storage Manager + LocalDriver (working) + S3Driver (stub)
│
├── console/
│   ├── kernel.go           # Console Kernel (cobra.Command root)
│   └── schedule.go         # Task scheduling
│
├── cmd/artisan/            # Artisan CLI command files
│   ├── make_commands.go    # make:controller, make:model, make:migration, make:middleware
│   ├── make_auth.go        # Auth scaffolding
│   ├── cache_commands.go   # cache:clear, cache:forget
│   ├── config_commands.go  # config:cache, config:clear
│   ├── queue_commands.go   # queue:work, queue:retry, queue:failed
│   ├── schema_commands.go  # schema:dump
│   ├── schedule_run.go     # schedule:run
│   ├── maintenance_commands.go  # down, up
│   └── vendor_commands.go  # vendor:publish
│
├── support/
│   ├── str/str.go          # Camel, Studly, Snake, Kebab, Random, Limit, Contains
│   ├── arr/                # Array helpers (TBD)
│   ├── collection/collection.go  # Collection[T] with Map, Filter, Reduce, Each, Chunked
│   ├── fakes/              # Test fakes
│   ├── health/             # Health checks (TBD)
│   ├── httpclient/         # HTTP client wrapper (TBD)
│   ├── pennant/            # Feature flags (TBD)
│   ├── pipeline/           # Pipeline (TBD)
│   └── process/            # Process execution (TBD)
│
├── encryption/             # AES-256-GCM encryption
├── hashing/                # Bcrypt + Argon2id
├── cookie/                 # Cookie helpers
├── localization/           # Translator
├── testing/                # TestCase wrapping httptest
│   ├── testcase.go         # TestCase (Get/Post/AssertStatus/AssertJson)
│   └── db.go               # Database test helpers
└── tests/benchmarks/       # Performance benchmarks
```

---

## 3. Container (IoC)

**File**: `container/container.go` (279 lines)

### Key Types
```go
type Container struct {
    bindings            map[reflect.Type]*binding
    contextualBindings  map[reflect.Type]map[reflect.Type]*binding
    instances           map[reflect.Type]any
    aliases             map[string]reflect.Type
    frozen              bool
    resolutionStack     []reflect.Type
}

type binding struct {
    factory   any        // func(...) T or func(...) (T, error)
    singleton bool
}
```

### API
```go
c := container.New()
c.Bind((*MyInterface)(nil), func() *MyImpl { return &MyImpl{} })
c.Singleton((*DB)(nil), func() *DB { return connectDB() })
c.Instance((*Config)(nil), configInstance)
c.Freeze()  // No more bindings after this

// Generic resolution:
impl, err := container.Make[MyInterface](c)

// Contextual binding:
c.When((*PhotoController)(nil)).
    Needs((*Filesystem)(nil)).
    Give(func() *S3Filesystem { return NewS3() })

// Struct injection via `inject` tag:
type MyController struct {
    DB     *orm.DB    `inject:""`
    Cache  cache.Store `inject:""`
}
```

### How Factory Resolution Works
1. Factory function params are resolved via `c.Resolve(paramType)` recursively
2. If factory needs `*Container` itself, it's injected directly
3. Singletons are cached in `instances` after first resolution
4. Contextual bindings are checked when a parent type is in the `resolutionStack`

### IMPORTANT FOR IMPLEMENTORS
- Interface types are registered as `(*Interface)(nil)` — the container extracts the `reflect.Type` from the pointer
- Factory must be a `reflect.Func` — panics on non-func
- After `Freeze()`, `Bind`/`Singleton`/`Instance` return `ErrFrozen`
- `build()` (auto-wiring) only injects fields tagged with `inject:""`

---

## 4. Foundation / Application

**File**: `foundation/application.go` (61 lines)

```go
type Application struct {
    *container.Container    // Embeds the IoC container
    basePath  string
    booted    bool
    providers []ServiceProvider
}

// ServiceProvider interface (in provider.go):
type ServiceProvider interface {
    Register(app *Application)
    Boot(app *Application)
}
```

### Bootstrap Flow
```go
app := foundation.NewApplication(".")
app.RegisterProvider(&DatabaseProvider{})
app.RegisterProvider(&AuthProvider{})
app.Boot()  // Calls Boot() on all providers, then Freeze() on container
```

### IMPORTANT FOR IMPLEMENTORS
- `Application` embeds `*Container`, so `app.Bind(...)`, `app.Make(...)` work directly
- `registerBaseBindings()` registers `*Application` and `*Container` as instances
- `Boot()` is idempotent — calling twice does nothing
- `Freeze()` is called inside `Boot()` — no bindings after boot

---

## 5. Config

**File**: `config/repository.go` (99 lines)

```go
type Repository struct {
    values map[string]string   // flat key-value store
}

func NewRepository(basePath string) *Repository  // loads .env + OS env
func (r *Repository) Get(key string, defaultValue ...string) string
func (r *Repository) GetBool(key string, defaultValue ...bool) bool
func (r *Repository) GetInt(key string, defaultValue ...int) int
func (r *Repository) Set(key string, value string)
```

### How It Works
1. If `CachedValues` (from `cached.go`) is non-empty, uses that (for `config:cache`)
2. Otherwise loads `.env` via `godotenv.Load(basePath + "/.env")`
3. Then loads all `os.Environ()` into the map
4. `Get()` checks map first, then `os.Getenv()`, then default

### IMPORTANT FOR IMPLEMENTORS
- Keys are flat strings (e.g., `"DB_CONNECTION"`, `"APP_KEY"`) — **NOT** nested dot-notation
- When adding nested config like `Get("database.connections.mysql.host")`, you'd need to parse `.` in the key and look up a nested structure. Currently NOT supported.
- The `config:cache` command should serialize `values` into `cached.go` as a Go map literal

---

## 6. Routing

**File**: `routing/router.go` (253 lines)

### Key Types
```go
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Router struct {
    trees        map[string]*node       // method → radix tree
    namedRoutes  map[string]*Route
    middlewares  []func(http.Handler) http.Handler
    groupPrefix  string
}

type Route struct {
    Method      string
    Path        string
    Handler     HandlerFunc
    Middlewares []func(http.Handler) http.Handler
    Name        string
}
```

### API
```go
r := routing.NewRouter()

r.Get("/users", handler)
r.Post("/users", handler)
r.Put("/users/{id}", handler)
r.Patch("/users/{id}", handler)
r.Delete("/users/{id}", handler)
r.Options("/users", handler)
r.Head("/users", handler)

route := r.Get("/users/{id}", handler)
r.SetName(route, "users.show")

r.Group("/api", func(r *routing.Router) {
    r.Get("/posts", postsHandler)
})

r.Use(myMiddleware)  // Adds to subsequent routes

// Resource routes:
r.Resource("photos", &PhotoController{})    // 7 routes: index/create/store/show/edit/update/destroy
r.ApiResource("photos", &PhotoController{}) // 5 routes: index/store/show/update/destroy
```

### How Radix Tree Works
- Each HTTP method has its own tree (`trees["GET"]`, `trees["POST"]`, etc.)
- Path segments are split by `/` and inserted as tree nodes
- Wildcard segments `{id}` are stored under a special `{}` key
- Wildcard's `paramName` stores the param name (e.g., `"id"`)
- On match, params are put into `context.WithValue(req.Context(), ParamsKey, params)`

### How to Read Route Params in Handlers
```go
func showUser(w http.ResponseWriter, r *http.Request) {
    params := r.Context().Value(routing.ParamsKey).(map[string]string)
    id := params["id"]
}
```

### Route Model Binding (dispatcher.go)
- `Dispatcher.Wrap(handler)` reflects on handler params
- If a param is `*SomeStruct`, it looks for `{somestruct}` in route params
- Queries DB with `WHERE id = ?` using the query builder
- Hydrates the struct using reflection (same approach as ORM)
- **Limitation**: Only matches by lowercase struct name → param name

### Signed URLs (signed.go)
```go
gen := routing.NewURLGenerator(router, appKey)
url, _ := gen.SignedRoute("verify.email", params, time.Now().Add(1*time.Hour))
valid := gen.HasValidSignature(request)
```

### IMPORTANT FOR IMPLEMENTORS
- `groupPrefix` is mutable during `Group()` callback — restored after via oldPrefix pattern
- Middleware is snapshotted at route registration time (not dynamic)
- Route matching is prefix-based — trailing slashes are trimmed (except `/`)
- There is NO `Any()` or `Match()` method yet — would need to call `AddRoute` for each method
- `namedRoutes` is populated via `SetName()` — used by `URLGenerator` and could be used for `route:list`

---

## 7. HTTP Kernel & Middleware

**File**: `http/kernel.go` (41 lines)

```go
type Kernel struct {
    app         *foundation.Application
    router      *routing.Router
    middlewares []func(http.Handler) http.Handler
}

func (k *Kernel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var handler http.Handler = k.router
    for i := len(k.middlewares) - 1; i >= 0; i-- {
        handler = k.middlewares[i](handler)
    }
    handler.ServeHTTP(w, r)
}
```

### Middleware Signature
All middleware follows the standard Go pattern:
```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before
        next.ServeHTTP(w, r)
        // after
    })
}
```

### Existing Middleware (12 files)
| File | What It Does |
|------|-------------|
| `auth.go` | Checks `Guard.Check()`, returns 401 if not authenticated |
| `cors.go` | Sets CORS headers (origin, methods, headers, credentials) |
| `csrf.go` | Validates `_token` field or `X-CSRF-TOKEN` header |
| `encrypt_cookies.go` | Encrypts/decrypts cookie values |
| `error_handler.go` | Catches panics and custom HttpException, returns JSON/HTML error |
| `maintenance.go` | Returns 503 if maintenance mode file exists |
| `recovery.go` | Recovers from panics, logs stack trace |
| `security_headers.go` | Sets X-Frame-Options, X-Content-Type-Options, etc. |
| `session.go` | Loads session at start, saves at end |
| `throttle.go` | Rate limits using cache.RateLimiter |
| `trim_strings.go` | Trims whitespace from all request form values |
| `trusted_proxies.go` | Sets X-Forwarded-For from trusted proxy IPs |

### IMPORTANT FOR IMPLEMENTORS
- Middleware is applied in REVERSE order (last added wraps outermost)
- Global middleware on Kernel runs for ALL routes
- Route-specific middleware is set at `AddRoute()` time from `r.middlewares` slice
- There is NO middleware group system (`web`, `api`) yet — would be a map[string][]Middleware on Router
- There is NO middleware parameter parsing (`throttle:60,1`) yet

---

## 8. HTTP Request & Response

### Response Helpers (`http/response/`)
```go
response.JSON(w, 200, data)
response.Ok(w, data)
response.Created(w, data)
response.Error(w, 400, "Bad request")
response.Redirect(w, r, "/dashboard", 302)
response.Stream(w, func(w http.ResponseWriter) { w.Write(chunk) })
response.Download(w, r, "/path/to/file", "report.pdf")
```

### FormRequest (`http/request/form_request.go`)
```go
type FormRequest interface {
    Authorize() bool
    Rules() map[string][]string
}

// Usage:
errs, authorized := request.Validate(r, &CreateUserRequest{})
```
- Calls `r.ParseForm()` and converts to `map[string]any`
- Passes to `validation.NewValidator(data, rules).Validate()`

### IMPORTANT FOR IMPLEMENTORS
- There is NO Request wrapper — handlers receive raw `*http.Request`
- Input is accessed via `r.FormValue()`, `r.URL.Query()`, `json.NewDecoder(r.Body)` — no convenience methods
- The `Redirect()` function is a 1-line wrapper around `http.Redirect` — no fluent API, no `with()` flash data, no `route()` redirect

---

## 9. Database Layer

### Connection (`database/connection.go`)
```go
type Connection struct {
    DB      *sql.DB
    Dialect dialect.Dialect
}
```

### Manager (`database/manager.go`)
```go
type Manager struct {
    connections map[string]*Connection
    defaultConn string
    config      *config.Repository
}

m.AddConnection("sqlite", conn)
conn, _ := m.Connection("sqlite")
conn, _ := m.DB()  // default connection
m.Close()
```

### Dialect Interface (`database/dialect/dialect.go`)
```go
type Dialect interface {
    QuoteIdentifier(name string) string
    Placeholder(index int) string
    CompileSelect(query SelectQuery) (string, []any)
    CompileInsert(table string, columns []string, values [][]any) (string, []any)
    CompileUpdate(table string, values map[string]any, wheres []WhereClause) (string, []any)
    CompileDelete(table string, wheres []WhereClause) (string, []any)
}
```

### SelectQuery struct
```go
type SelectQuery struct {
    Table     string
    Columns   []string
    Wheres    []WhereClause
    OrderBys  []OrderByClause
    Joins     []JoinClause
    Aggregate *AggregateClause
    Limit     *int
    Offset    *int
}
```

### WhereClause struct
```go
type WhereClause struct {
    Type     string  // "Basic", "In", "Null", "NotNull", "Between", "Raw"
    Column   string
    Operator string
    Value    any
    Values   []any   // For IN, BETWEEN
    Boolean  string  // "AND" / "OR"
    RawSQL   string  // For raw queries
    RawArgs  []any
}
```

### Query Builder API (`database/query/builder.go`)
```go
b := query.NewBuilder(conn, dialect)
b.Table("users")
b.Select("name", "email")
b.Where("age", ">", 18)
b.OrWhere("role", "=", "admin")
b.WhereIn("status", []any{"active", "pending"})
b.WhereNull("deleted_at")
b.WhereNotNull("email_verified_at")
b.WhereBetween("created_at", []any{start, end})
b.OrderBy("name", "ASC")
b.Limit(10)
b.Offset(20)
b.Join("posts", "users.id", "=", "posts.user_id")
b.LeftJoin("profiles", "users.id", "=", "profiles.user_id")
b.RightJoin(...)
b.CrossJoin("settings")
b.When(onlyActive, func(b *Builder) { b.Where("active", "=", true) })

rows, _ := b.Get()                         // SELECT
res, _ := b.Insert(map[string]any{...})    // INSERT
res, _ := b.Update(map[string]any{...})    // UPDATE (uses current WHEREs)
res, _ := b.Delete()                       // DELETE (uses current WHEREs)
count, _ := b.Count("")                    // SELECT COUNT(*)
max, _ := b.Max("price")
min, _ := b.Min("price")
avg, _ := b.Avg("price")
sum, _ := b.Sum("amount")
sql, args := b.ToSQL()                     // Get raw SQL without executing
```

### Query Logging
```go
query.Listen(func(event query.QueryEvent) {
    log.Printf("SQL: %s | Duration: %v", event.SQL, event.Duration)
})
```

### IMPORTANT FOR IMPLEMENTORS
- Builder is NOT immutable — each method mutates internal `query` struct
- Only one dialect exists: `SQLiteDialect`. When adding MySQL/Postgres:
  - MySQL: `Placeholder` returns `?`, `QuoteIdentifier` uses backticks
  - Postgres: `Placeholder` returns `$1`, `$2`, etc., `QuoteIdentifier` uses double quotes
- The `WhereClause.Type` field controls SQL compilation in `compileWheres()`:
  - `"Basic"` → `column op ?`
  - `"In"` → `column IN (?, ?, ...)`
  - `"Null"` → `column IS NULL`
  - `"NotNull"` → `column IS NOT NULL`
  - `"Between"` → `column BETWEEN ? AND ?`
  - Add `"Raw"` type for raw WHERE expressions using `RawSQL`/`RawArgs` fields (already in struct, not yet compiled)
- **To add GROUP BY**: Add `GroupBys []string` to `SelectQuery`, compile after WHERE
- **To add HAVING**: Add `Havings []WhereClause` to `SelectQuery`, compile after GROUP BY
- **To add DISTINCT**: Add `Distinct bool` to `SelectQuery`, compile as `SELECT DISTINCT`
- **To add transactions**: Add methods on Builder or Connection:
  ```go
  func (c *Connection) Transaction(fn func(tx *sql.Tx) error) error {
      tx, _ := c.DB.Begin()
      err := fn(tx)
      if err != nil { tx.Rollback(); return err }
      return tx.Commit()
  }
  ```

---

## 10. ORM (Goquent)

### Model Interface
```go
type Model interface {
    TableName() string
}
```

### Struct Tag Convention
```go
type User struct {
    ID        int       `db:"id"    gow:"primaryKey,autoIncrement"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at" gow:"autoCreateTime"`
    UpdatedAt time.Time `db:"updated_at" gow:"autoUpdateTime"`
    DeletedAt *time.Time `db:"deleted_at" gow:"softDelete"`
}
```

### Tags recognized by ORM:
- `db:"column_name"` — maps struct field to DB column
- `gow:"primaryKey"` — marks the primary key field
- `gow:"autoIncrement"` — skipped during INSERT
- `gow:"autoCreateTime"` — set to `time.Now()` on INSERT
- `gow:"autoUpdateTime"` — set to `time.Now()` on INSERT and UPDATE
- `gow:"softDelete"` — recognized but **NOT filtered in queries** (bug)
- `gow:"hasMany"` / `gow:"belongsTo"` — relationship tags (used by relation.go)

### ModelQuery[T] API
```go
q := orm.NewQuery[User](db)
user, _ := q.Find(1)                              // WHERE id = 1 LIMIT 1
user, _ := q.Where("email", "=", "x@x.com").First()
users, _ := q.Where("active", "=", true).Get()
err := q.Insert(&user)                            // Sets auto timestamps, gets LastInsertId
err := q.Update(&user)                            // Finds PK, updates non-PK fields
err := q.Delete(&user)                            // Finds PK, deletes
err := q.Save(&user)                              // Insert if PK=0, Update otherwise
users, _ := q.With("Posts").Get()                  // Eager load (BROKEN — returns nil)
```

### How Insert Works Internally
1. Reflects on `*T`, iterates fields
2. Skips fields with `autoIncrement`
3. Sets `time.Now()` for `autoCreateTime`/`autoUpdateTime`
4. Builds `map[string]any` of column→value
5. Calls `builder.Insert(values)`
6. Gets `LastInsertId()` and sets `model.ID`

### How Update Works Internally
1. Reflects to find `primaryKey` or `id` field → gets PK value
2. Skips PK field from SET clause
3. Sets `time.Now()` for `autoUpdateTime`
4. Calls `builder.Where(pkColumn, "=", pkValue).Update(values)`

### How Save Works
1. Checks if PK is zero value (int=0, string="") → Insert
2. Otherwise → Update
3. **Bug**: The type assertion `v.(int64)` will panic if PK is `int` not `int64`

### How hydrateModel Works
1. Gets column names from `rows.Columns()`
2. Builds a `fieldMap[dbTag] = fieldIndex` from struct tags
3. Falls back to lowercase field name if no `db:` tag
4. Creates `scanArgs` pointing to struct field addresses
5. Dummy receiver for unmatched columns

### Relationship Types (declared but not functional)
```go
type HasMany[T any] struct { Models []T }
type BelongsTo[T any] struct { Model *T }
type HasOne[T any] struct { Model *T }
type BelongsToMany[T any] struct { Models []T }
```

### Scopes & Observers (declared, never called)
```go
type Scope interface { Apply(builder *Builder) *Builder }
type Observer interface {
    Creating(model any) bool
    Created(model any)
    Updating(model any) bool
    Updated(model any)
    Deleting(model any) bool
    Deleted(model any)
}
func AddGlobalScope(model string, scope Scope)
func Observe(model string, observer Observer)
func DispatchModelEvent(model any, event string) bool
```

### IMPORTANT FOR IMPLEMENTORS

**To implement working eager loading** (relation.go `loadRelation`):
1. Determine the target model type from the relation field (e.g., `Posts []Post` → `Post`)
2. Determine the foreign key: `belongsTo` → `UserID` on parent; `hasMany` → `user_id` on child table
3. Collect all parent IDs: `ids := collectFieldValues(models, "ID")`
4. Execute: `SELECT * FROM posts WHERE user_id IN (?, ?, ...)`
5. Hydrate results into `[]Post` using `hydrateModel`
6. Map results back to parents by foreign key value
7. Set the relation field on each parent model via reflection

**To implement soft deletes**:
1. In `ModelQuery.Get()` / `First()` / `Find()`, check if model struct has a `softDelete` tagged field
2. If yes, automatically add `WhereNull("deleted_at")` to the builder
3. Add `WithTrashed()` method that sets a flag to skip adding the clause
4. Add `OnlyTrashed()` method that adds `WhereNotNull("deleted_at")`
5. Add `Restore()` that sets `deleted_at = NULL`
6. Modify `Delete()` to set `deleted_at = time.Now()` instead of actual DELETE when soft delete is present
7. Add `ForceDelete()` for actual deletion

**To wire observers into CRUD**:
```go
func (q *ModelQuery[T]) Insert(model *T) error {
    if !orm.DispatchModelEvent(model, "creating") { return ErrCancelled }
    // ... existing insert logic ...
    orm.DispatchModelEvent(model, "created")
    return nil
}
```

**To apply global scopes**:
```go
func NewQuery[T any](db *DB) *ModelQuery[T] {
    // ... existing code ...
    modelName := reflect.TypeOf((*T)(nil)).Elem().Name()
    if scopes, ok := globalScopes[modelName]; ok {
        for _, scope := range scopes {
            scope.Apply(mq.builder)
        }
    }
    return mq
}
```

---

## 11. Validation

**File**: `validation/validator.go` (180 lines)

### API
```go
v := validation.NewValidator(data, map[string][]string{
    "email": {"required", "email"},
    "name":  {"required", "min:3", "max:255"},
    "age":   {"numeric", "between:18,99"},
    "role":  {"required", "in:admin,user,moderator"},
    "password": {"required", "min:8", "confirmed"},
    "username": {"required", "unique:users,username"},
})
v = v.WithDB(sqlDB)  // needed for unique/exists rules
errs := v.Validate()  // map[string][]error
```

### Rule Parsing
- Rules are strings like `"min:3"` → split on `:` → ruleName=`"min"`, ruleParams=`["3"]`
- Multiple params separated by `,` → `"between:18,99"` → params=`["18","99"]`

### Implemented Rules (14)
`required`, `email` (proper regex), `min`, `max`, `between`, `size`, `numeric`, `string`, `in`, `regex`, `confirmed`, `unique:table,column`, `exists:table,column`

### How min/max/size/between Work
- If value is numeric, compares the number value itself
- If value is string, compares `len(strVal)` (string length)
- This dual behavior matches Laravel

### IMPORTANT FOR IMPLEMENTORS
- **To add `nullable`**: If present and value is nil/empty, skip all other rules for that field
- **To add `sometimes`**: Only validate if the field exists in data
- **To add `bail`**: Stop running remaining rules for a field after first failure
- **To add `required_if:field,value`**: Check if another field equals value, then enforce required
- **To add `date`**: Parse with `time.Parse()` and validate
- **To add array validation (`items.*.name`)**: Split field key on `.`, iterate into nested data, apply rules to each
- Error messages are currently hardcoded English. For localization, replace with `translator.Get("validation.required", field)`

---

## 12. Session

**File**: `session/manager.go` (171 lines)

### Store Interface
```go
type Store interface {
    Read(id string) (map[string]any, error)
    Write(id string, data map[string]any) error
    Destroy(id string) error
}
```

### Manager API
```go
m := session.NewManager(store, sessionID) // generates ID if empty
m.Start()                        // loads from store, ages flash data
m.Get("key")                     // returns any
m.Put("key", "value")
m.Flash("status", "Saved!")      // available next request only
m.Reflash()                      // keep all flash data for another request
m.Keep("status", "error")        // keep specific flash keys
m.FlashInput(map[string]any{...}) // flash form input as "_old_input"
m.Old("email", "")               // retrieve old input from previous request
m.Regenerate()                    // new session ID, destroy old, save
m.Save()                          // writes to store, clears old flash data
m.ID()                            // current session ID
```

### How Flash Data Works
- `data["_flash"]` stores `map[string][]string{"old": [...], "new": [...]}`
- On `Start()`: `ageFlashData()` moves `new` → `old`, empties `new`
- On `Flash(key, value)`: puts value in data AND adds key to `new` array
- On `Save()`: `clearOldFlashData()` deletes all keys in `old` from data, clears `old`
- Result: flashed data survives exactly ONE more request

### Drivers
- **File**: Serializes to JSON files on disk
- **Cookie**: Stores data in encrypted cookie
- **Redis**: Uses Redis SET/GET with session ID as key

### IMPORTANT FOR IMPLEMENTORS
- Session middleware (`http/middleware/session.go`) calls `Start()` at request begin, `Save()` at end
- There is NO database driver yet — would need a `sessions` table with `id`, `payload`, `last_activity`
- There is NO encryption on file/cookie drivers currently
- `Regenerate()` destroys old session in store, generates new ID, saves immediately

---

## 13. Auth & Authorization

### Guard Interface
```go
type Guard interface {
    Check() bool
    Guest() bool
    User() any
    ID() string
    Attempt(credentials map[string]any) bool
    Login(user any)
    Logout()
}
```

### SessionGuard (working implementation)
- Stores user ID in session as `login_{guardName}_id`
- `User()`: gets ID from session, calls `provider.RetrieveByID(id)`
- `Attempt()`: calls `provider.RetrieveByCredentials()` + `ValidateCredentials()`
- `Login()`: requires user to implement `Authenticatable` interface
- `Logout()`: clears user from session, calls `Regenerate()` (session fixation protection)

### Authenticatable Interface
```go
type Authenticatable interface {
    GetAuthIdentifier() string   // returns user's ID as string
    GetAuthPassword() string     // returns hashed password
}
```

### UserProvider Interface
```go
type UserProvider interface {
    RetrieveByID(identifier string) any
    RetrieveByToken(identifier string, token string) any
    UpdateRememberToken(user any, token string)
    RetrieveByCredentials(credentials map[string]any) any
    ValidateCredentials(user any, credentials map[string]any) bool
}
```
**No concrete UserProvider implementation exists** — users must implement this themselves.

### AuthManager
```go
m := auth.NewManager()
m.AddGuard("web", sessionGuard)
m.AddGuard("api", sanctumGuard)
guard := m.Guard("web")
```

### Fortify (auth/fortify/fortify.go)
Routes registered: `/api/register`, `/api/login`, `/api/forgot-password`, `/api/reset-password`, `/api/two-factor/enable`, `/api/email/verify/{id}/{hash}`
**All handler bodies are stubs** — they return success messages without actual logic.

### Fortify.handleLogin — What's Implemented
- Decodes JSON body (email, password, remember)
- Calls `authManager.Guard("web").Attempt(credentials, remember)`
- **Bug**: `Guard.Attempt` takes `map[string]any` — doesn't have a `remember` param. The call `guard.Attempt(creds, creds.Remember)` passes 2 args but interface defines 1.

### Sanctum (auth/sanctum/token.go)
- `CreateToken()`: generates 40-char random token, SHA256 hashes it, stores in `personal_access_tokens` table
- Middleware: extracts `Bearer` token from Authorization header, hashes it, checks DB
- **Not implemented**: token expiration check, abilities/scopes, token pruning

### Gate (auth/access/gate.go)
- `Define(name, callback)` — registers an authorization callback
- `Allows(name, user, args...)` — runs the callback
- `Denies(name, user, args...)` — inverse of Allows

### IMPORTANT FOR IMPLEMENTORS
- `Guard.Attempt()` interface takes `map[string]any` — password comparison must happen in `UserProvider.ValidateCredentials()`
- The `Authenticatable` interface is required for `SessionGuard.Login()` — it type-asserts the user
- Remember token flow: `UserProvider.UpdateRememberToken()` exists but is never called — `Login()` doesn't handle remember me
- **To implement password reset**: Create a `password_resets` table, generate token, hash it, send via mail, verify on submit

---

## 14. Cache

### Store Interface
```go
type Store interface {
    Get(key string) (any, bool)
    Put(key string, value any, ttl time.Duration)
    Increment(key string, value int) int
    Decrement(key string, value int) int
    Forever(key string, value any)
    Forget(key string)
    Flush()
}
```

### Drivers
- **Memory**: `sync.Map` + TTL tracking map. Thread-safe.
- **Redis**: Uses redis client for GET/SET/INCR/DECR/DEL/FLUSHDB

### Rate Limiter
```go
limiter := cache.NewRateLimiter(store)
limiter.TooManyAttempts(key, maxAttempts) bool
limiter.Hit(key, decayMinutes)
limiter.Clear(key)
```

---

## 15. Queue

### Job Interface
```go
type Job interface {
    Handle() error
    Failed(err error)
}
```

### Driver Interface (implicit from manager.go)
```go
type Driver interface {
    Push(job Job) error
    Pop() (Job, error)
}
```

### Drivers
- **Sync**: Go channel-based, immediate processing
- **Database**: Serializes jobs to `jobs` table
- **Redis**: Uses Redis LIST (LPUSH/RPOP)

### Worker Loop
```go
// For Sync driver: range over channel
// For others: infinite loop calling Pop(), with no backoff (TODO: add sleep/backoff)
```

### IMPORTANT FOR IMPLEMENTORS
- There is NO delayed jobs, job chaining, batching, failed jobs table, retries, max attempts, or backoff
- The worker loop for non-Sync drivers has NO sleep — will spin CPU 100%. Add `time.Sleep(1*time.Second)` on empty Pop().

---

## 16. Events

### Dispatcher
```go
d := events.NewDispatcher()
d.Listen(UserCreated{}, func(event any) { sendWelcomeEmail(event.(UserCreated)) })
d.ListenAny(func(event any) { logEvent(event) })
d.Dispatch(UserCreated{User: user})
```

### Dispatch mechanism
- Type-based: uses `reflect.TypeOf(event)` as map key
- Wildcard listeners receive ALL events
- All synchronous — no queued listeners

### IMPORTANT FOR IMPLEMENTORS
- **To add queued listeners**: Wrap the listener func into a `Job` and push via `queue.Manager`
- **To add stoppable events**: Return `bool` from listener, break on `false`

---

## 17. Mail

### Architecture
```go
mailer := mail.NewMailer(&mail.SmtpDriver{Host: "smtp.example.com", Port: "587", ...})
mailer.Send(myMailable)
```

### Mailable Interface
```go
type Mailable interface { Build() *Message }
```

### BaseMailable Builder
```go
m := mail.NewMailable().
    From("noreply@example.com").
    To("user@example.com").
    Cc("manager@example.com").
    Bcc("audit@example.com").
    Subject("Welcome!").
    HTML("<h1>Welcome</h1>").
    Attach("report.pdf", fileBytes)
msg := m.Build()
```

### SmtpDriver
- Uses `net/smtp.PlainAuth` + `smtp.SendMail`
- Supports CC, BCC, attachments (multipart/mixed)
- **Bug**: Attachments are written as raw bytes, not base64-encoded (comment says "for simplicity")

---

## 18. View Engine (Goblade)

### Compiler Directives Supported
```
@if(cond)/@elseif(cond)/@else/@endif
@foreach(range)/@endforeach
@for(range)/@endfor
@switch(val)/@case(val)/@default/@endswitch
@yield('name') → {{ block "name" . }}{{ end }}
@section('name')/@endsection → {{ define "name" }}...{{ end }}
@include('partial') → {{ template "partial" . }}
@extends('layout') → {{/* extends "layout" */}} (comment — resolved by Engine)
@csrf → hidden input
@can(ability, resource)/@endcan
@cannot(ability, resource)/@endcannot
@checked(condition)
@class(...)
@once/@endonce → comments (NOT FUNCTIONAL)
@while(cond)/@endwhile → comments (NOT FUNCTIONAL)
<x-component /> → {{ template "components.component" . }}
```

### Engine (view/engine.go)
- `Engine.Make(name, data)` → resolves file path → compile → parse → execute
- Looks for `.goblade` (preferred), `.blade.php` (compat), `.html` extensions
- Detects `@extends` from compiled output via regex, resolves layout file
- Parses BOTH child + layout files via `template.ParseFiles(child, layout)`
- **Layout resolution IS implemented** — the engine regex-matches `{{/* extends "..." */}}` and loads the layout

### IMPORTANT FOR IMPLEMENTORS
- Engine needs `regexp` import in `engine.go` — currently uses it but verify it compiles
- Go templates use `{{ .FieldName }}` for data access, not `$variable`
- Custom template funcs (for `@can`, `@class`, `@checked`) need to be registered via `template.FuncMap` before parsing
- **For `$loop` variable**: Would need a custom `range` implementation that passes loop metadata (index, first, last, count) via a wrapper struct in the funcmap

---

## 19. Logging

**File**: `logging/logger.go` (117 lines)

### Features Implemented
- Daily rotation via `dailyWriter` (creates `gow-YYYY-MM-DD.log`)
- Stack channels via `io.MultiWriter`
- JSON output via `slog.NewJSONHandler`
- Request ID injection via `WithContext(ctx, logger)`
- Configurable levels: debug, info, warn, error

### Config
```go
type Config struct {
    Level    string    // "debug", "info", "warn", "error"
    Channel  string   // "single", "daily", "stack"
    Path     string   // log directory path
    Channels []string // used if Channel == "stack"
}
```

---

## 20. Storage

### Driver Interface
```go
type Driver interface {
    Put(path string, contents io.Reader) error
    Get(path string) (io.ReadCloser, error)
    Delete(path string) error
    Exists(path string) bool
    URL(path string) string
}
```

### LocalDriver — Fully Working
- Uses `os.Create`, `os.Open`, `os.Remove`, `os.Stat`
- Auto-creates directories via `os.MkdirAll`

### S3Driver — All Stubs (returns nil)

---

## 21. Console / CLI

### Entry Point (`artisan.go`)
```go
app := foundation.NewApplication(".")
kernel := console.NewKernel(app)
kernel.RegisterCommand(myCommand)
app.Boot()
kernel.Run()
```

### Existing Commands
- `make:controller`, `make:model`, `make:migration`, `make:middleware` (in make_commands.go)
- `make:policy`, `make:gate`, `make:guard` (in make_auth.go)
- `cache:clear`, `cache:forget` (in cache_commands.go)
- `config:cache`, `config:clear` (in config_commands.go)
- `queue:work`, `queue:retry`, `queue:failed` (in queue_commands.go)
- `schema:dump` (in schema_commands.go)
- `schedule:run` (in schedule_run.go)
- `down`, `up` (in maintenance_commands.go)
- `vendor:publish` (in vendor_commands.go)

### IMPORTANT FOR IMPLEMENTORS
- All commands use `cobra.Command` struct
- `make:*` commands should generate Go source files using `text/template`
- Generated files should follow the project structure in Section 2
- Template output should include proper imports and struct boilerplate

---

## 22. Support Packages

### Str (`support/str/str.go`)
```go
str.Camel("hello_world")   // "helloWorld"
str.Studly("hello_world")  // "HelloWorld"
str.Snake("helloWorld")     // "hello_world"
str.Kebab("helloWorld")     // "hello-world"
str.Random(40)              // random hex string
str.Limit("long text", 10, "...") // "long text..."
str.Contains("hello", "ell")      // true
```

### Collection (`support/collection/collection.go`)
```go
c := collection.Collect([]int{1, 2, 3, 4, 5})
c.Map(func(i int) int { return i * 2 })
c.Filter(func(i int) bool { return i > 2 })
c.Reduce(func(acc any, i int) any { return acc.(int) + i }, 0)
c.Each(func(i int) { fmt.Println(i) })
c.Chunked(2)  // []*Collection[int]
c.All()        // []int
```
**Bug**: `Chunk()` method returns nil — use `Chunked()` instead

---

## 23. Testing

### TestCase (`testing/testcase.go`)
```go
tc := testing.NewTestCase(t, router)
tc.Get("/api/users").AssertStatus(200).AssertJson(expected)
tc.Post("/api/users").AssertStatus(201)
```

---

## 24. Implementation Rules

### When Adding a New Feature:

1. **Package Location**: Follow existing structure — each feature is its own package under `gow/`
2. **Interface First**: Define the interface in the package, then implement drivers
3. **Tags**: Use `db:""` for DB columns, `gow:""` for ORM behavior, `inject:""` for DI
4. **Error Handling**: Return `error` — never panic (except in bootstrap/fatal)
5. **Thread Safety**: Use `sync.RWMutex` for any shared state (see Container, Router, Session)
6. **Context**: Use `context.Context` for values that need to flow (params, auth, request ID)
7. **Middleware Signature**: Always `func(http.Handler) http.Handler`
8. **Handler Signature**: Always `func(http.ResponseWriter, *http.Request)` (the `HandlerFunc` type)
9. **Naming**: Use Go conventions — exported = public API, unexported = internal
10. **Testing**: Use `testing` stdlib + `testify/assert`. Use SQLite `:memory:` for DB tests.

### When Adding a New Query Builder Feature:
1. Add the field to `dialect.SelectQuery` struct
2. Add the builder method on `query.Builder`
3. Add compilation logic in `sqlite.go`'s `CompileSelect()` (or relevant method)
4. When adding MySQL/Postgres, implement the same in their dialect files

### When Adding a New ORM Feature:
1. Add any needed struct tags to the tag parsing in `Insert()`/`Update()` etc.
2. Use `reflect` to check for tags — cache the metadata if performance matters
3. Wire into `DispatchModelEvent()` for observer integration

### When Adding a New Middleware:
1. Create file in `http/middleware/`
2. Follow signature: `func MyMiddleware(deps) func(http.Handler) http.Handler`
3. Register in the HTTP Kernel or Router

### When Adding a New Artisan Command:
1. Create file in `cmd/artisan/`
2. Define `cobra.Command` with `Use`, `Short`, `Run`
3. Register in `artisan.go` or via a registration function

### When Adding a New Validation Rule:
1. Add a `case "rulename":` in `applyRule()` switch statement
2. Parse params from `ruleParams` if the rule takes arguments
3. Return `errors.New("The " + field + " ...")` with Laravel-style message

### When Adding a New Cache/Session/Queue Driver:
1. Implement the interface (`cache.Store`, `session.Store`, queue's `Push/Pop`)
2. Create file named `driver_{name}.go` in the package
3. Register in the Manager's constructor or via a factory method

---

## Quick Reference: Key Type Relationships

```
foundation.Application
  └── embeds container.Container
  └── has []ServiceProvider
  └── creates http.Kernel
      └── has routing.Router
          └── has map[string]*node (radix trees per method)
          └── has map[string]*Route (named routes)
      └── has []Middleware

routing.Router
  └── creates routing.Route (per registered endpoint)
  └── uses routing.HandlerFunc = func(http.ResponseWriter, *http.Request)
  └── stores params in context via routing.ParamsKey

database.Manager
  └── has map[string]*database.Connection
      └── has *sql.DB
      └── has dialect.Dialect

query.Builder
  └── uses *sql.DB from Connection
  └── uses dialect.Dialect for SQL compilation
  └── builds dialect.SelectQuery struct

orm.DB
  └── has *sql.DB (from Connection)
  └── has *query.Builder

orm.ModelQuery[T]
  └── wraps *query.Builder
  └── wraps *orm.DB
  └── uses reflect to read struct tags

auth.AuthManager
  └── has map[string]Guard
      └── SessionGuard
          └── uses UserProvider
          └── uses session.Manager
      └── or Sanctum TokenManager

session.Manager
  └── uses session.Store (File/Cookie/Redis)
  └── manages flash data lifecycle

events.Dispatcher
  └── has map[reflect.Type][]Listener
  └── has []Listener (wildcards)

mail.Mailer
  └── uses mail.Driver (Log/Smtp)
  └── sends mail.Message built by mail.Mailable
```
