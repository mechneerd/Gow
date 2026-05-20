
# GOW FRAMEWORK SPECIFICATION
## A Laravel-Inspired Web Development Framework for Go

**Version:** 0.1.0-alpha  
**Date:** 2026-05-20  
**Status:** Specification / Planning Phase  
**Language:** Go 1.24+  
**Architecture:** MVC-inspired with Service Container, Middleware Pipeline, and Modular Design

---

## TABLE OF CONTENTS
1. [Directory Structure](#1-directory-structure)
2. [Core Architecture](#2-core-architecture)
3. [Routing System](#3-routing-system)
4. [Middleware System](#4-middleware-system)
5. [Controller Layer](#5-controller-layer)
6. [ORM / Database Layer (Goquent)](#6-orm--database-layer-goquent)
7. [Migration System](#7-migration-system)
8. [Seeding & Factories](#8-seeding--factories)
9. [View / Template Engine (Blade-like)](#9-view--template-engine-blade-like)
10. [Service Container & IoC](#10-service-container--ioc)
11. [Service Providers](#11-service-providers)
12. [Configuration Management](#12-configuration-management)
13. [Authentication & Authorization](#13-authentication--authorization)
14. [Validation System](#14-validation-system)
15. [Session Management](#15-session-management)
16. [Cache System](#16-cache-system)
17. [Queue & Job System](#17-queue--job-system)
18. [Event & Listener System](#18-event--listener-system)
19. [Mail System](#19-mail-system)
20. [Notification System](#20-notification-system)
21. [File Storage](#21-file-storage)
22. [Broadcasting / WebSockets](#22-broadcasting--websockets)
23. [Pagination](#23-pagination)
24. [Localization (i18n)](#24-localization-i18n)
25. [Task Scheduling](#25-task-scheduling)
26. [Testing Framework](#26-testing-framework)
27. [CLI / Artisan Equivalent](#27-cli--artisan-equivalent)
28. [Logging System](#28-logging-system)
29. [Security Features](#29-security-features)
30. [Error Handling & Exceptions](#30-error-handling--exceptions)
31. [Package Development](#31-package-development)
32. [Implementation Roadmap](#32-implementation-roadmap)
33. [Go-Specific Implementation Notes](#33-go-specific-implementation-notes)
34. [Agent Continuation Guidelines](#34-agent-continuation-guidelines)
35. [API Resources / Transformers](#35-api-resources--transformers)
36. [Collections & Lazy Collections](#36-collections--lazy-collections)
37. [Helpers, Stringable & Fluent Utilities](#37-helpers-stringable--fluent-utilities)
38. [Signed URLs](#38-signed-urls)
39. [Route Model Binding](#39-route-model-binding)
40. [Process / CLI Command Execution](#40-process--cli-command-execution)
41. [Schema Dumping](#41-schema-dumping)
42. [Model Pruning](#42-model-pruning)
43. [Health Checks](#43-health-checks)
44. [Feature Flags](#44-feature-flags)
45. [Advanced Blade / Template Features](#45-advanced-blade--template-features)
46. [Request & Response Macros / Old Input / Flashing](#46-request--response-macros--old-input--flashing)
47. [Query Logging & Database Profiling](#47-query-logging--database-profiling)
48. [Package Auto-Discovery & Publishing](#48-package-auto-discovery--publishing)
49. [Testing Utilities](#49-testing-utilities)
50. [Modern Laravel Parity](#50-modern-laravel-parity)

---

## 1. DIRECTORY STRUCTURE

```
gow/
├── app/                          # Application logic
│   ├── console/                  # CLI commands (Artisan equivalent)
│   │   └── commands/             # Custom commands
│   ├── http/                     # HTTP layer
│   │   ├── controllers/          # HTTP controllers
│   │   ├── middleware/           # HTTP middleware
│   │   └── requests/             # Form request validation classes
│   ├── models/                   # Eloquent/ORM models
│   ├── providers/                # Service providers
│   ├── services/                 # Business logic services
│   ├── repositories/             # Data access layer (optional)
│   ├── jobs/                     # Queue jobs
│   ├── listeners/                # Event listeners
│   ├── mail/                     # Mailable classes
│   ├── notifications/            # Notification classes
│   ├── policies/                 # Authorization policies
│   ├── rules/                    # Custom validation rules
│   └── facades/                  # Facade accessors (if implemented)
├── bootstrap/                    # Framework bootstrapping
│   ├── app.go                    # Application instance creation
│   └── cache/                    # Framework cache (compiled templates, etc.)
├── config/                       # Configuration files
│   ├── app.go                    # App config
│   ├── database.go               # Database config
│   ├── cache.go                  # Cache config
│   ├── queue.go                  # Queue config
│   ├── mail.go                   # Mail config
│   ├── filesystems.go            # Storage config
│   ├── session.go                # Session config
│   ├── auth.go                   # Auth config
│   ├── logging.go                # Logging config
│   ├── services.go               # Third-party services
│   └── cors.go                   # CORS config
├── database/                     # Database-related files
│   ├── migrations/               # Migration files
│   ├── seeders/                  # Database seeders
│   ├── factories/                # Model factories
│   └── sqlite/                   # SQLite databases (dev)
├── public/                       # Web server root
│   ├── index.go                  # Front controller (compiled binary or handler)
│   ├── css/                      # Compiled CSS
│   ├── js/                       # Compiled JS
│   ├── images/                   # Public images
│   └── .htaccess / nginx.conf    # Server configs
├── resources/                    # Raw assets & views
│   ├── views/                    # Template files (.gohtml / .blade equivalent)
│   │   ├── layouts/              # Layout templates
│   │   ├── components/           # Reusable components
│   │   └── partials/             # Partial templates
│   ├── lang/                     # Language files (i18n)
│   │   ├── en/                   # English translations
│   │   └── es/                   # Spanish translations
│   ├── css/                      # Source CSS (Tailwind/Sass)
│   ├── js/                       # Source JS
│   └── raw/                      # Other raw assets
├── routes/                       # Route definitions
│   ├── web.go                    # Web routes
│   ├── api.go                    # API routes
│   ├── console.go                # Console routes/commands
│   └── channels.go               # Broadcasting channels
├── storage/                      # Writable storage
│   ├── app/                      # App-generated files
│   │   └── public/               # User-uploaded public files
│   ├── framework/                # Framework internal
│   │   ├── cache/                # Cached data
│   │   ├── sessions/             # Session files
│   │   ├── views/                # Compiled templates
│   │   └── testing/              # Test snapshots
│   └── logs/                     # Application logs
├── tests/                        # Test suites
│   ├── feature/                  # Feature/Integration tests
│   ├── unit/                     # Unit tests
│   ├── bootstrap.go              # Test bootstrap
│   └── testdata/                 # Test fixtures
├── vendor/                       # Dependencies (Go modules)
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── .env                          # Environment variables
├── .env.example                  # Environment template
├── gow.go                        # Framework core entry point
├── main.go                       # Application entry point
├── artisan.go                    # CLI entry point
├── README.md
└── CHANGELOG.md                  # Changes tracking
```

---

## 2. CORE ARCHITECTURE

### 2.1 Application Lifecycle
1. **Bootstrap**: Load environment, config, service providers
2. **Register**: Providers register bindings in the service container
3. **Boot**: Providers boot (run startup logic)
4. **Handle**: Kernel receives request, pushes through middleware, dispatches to router
5. **Terminate**: Response sent, terminable middleware runs, cleanup

### 2.2 Key Components
- **Application**: Central container holding all bindings, config, and providers
- **Kernel**: HTTP/Console kernel handling request lifecycle
- **Router**: Route registration, matching, and dispatch
- **Container**: Dependency injection container (constructor injection + interface resolution)
- **Pipeline**: Middleware pipeline for request/response transformation

### 2.3 Go Implementation Strategy
- Use `interface{}` or generics (Go 1.18+) for container bindings
- Use `http.Handler` and `http.HandlerFunc` as base HTTP abstractions
- Use `context.Context` for request-scoped data propagation
- Use Go modules for dependency management
- Use `embed` package for template/assets embedding in production

---

## 3. ROUTING SYSTEM

### 3.1 Features Required
- Route registration (GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD)
- Route groups with shared attributes (prefix, middleware, namespace)
- Named routes for URL generation
- Route parameters (required, optional, regex constraints)
- Route model binding (automatic model injection)
- Route caching for production (compile to map/tree)
- Fallback routes
- Redirect routes
- View routes (route directly returns view)
- Rate limiting per route/group
- Subdomain routing
- Route middleware assignment

### 3.2 How It Works in Laravel
- Routes defined in `routes/web.php` and `routes/api.php`
- Route facade accesses router instance
- Route matching uses Symfony's URL matcher or FastRoute
- Parameters extracted from URI and injected into controller actions
- Named routes allow `route('name', params)` URL generation

### 3.3 Go Implementation
- Use a radix tree or trie for O(n) route matching (similar to httprouter/julienschmidt)
- Route definitions as Go code in `routes/web.go`:
  ```go
  func WebRoutes(router *gow.Router) {
      router.Get("/", homeController.Index)
      router.Get("/users/{id}", userController.Show).Name("users.show")
      router.Group("/admin", func(r *gow.Router) {
          r.Use(middleware.Auth)
          r.Get("/dashboard", adminController.Dashboard)
      })
  }
  ```
- Route parameters parsed via regex or segment matching
- Named routes stored in a map for reverse URL generation
- Support route caching: compile routes to a serialized structure on first boot

---

## 4. MIDDLEWARE SYSTEM

### 4.1 Features Required
- Global middleware (applies to all routes)
- Route/group middleware
- Middleware parameters
- Terminable middleware (runs after response sent)
- Priority/sorting of middleware
- Built-in middleware: auth, guest, verified, throttle, cors, csrf, etc.

### 4.2 How It Works in Laravel
- PSR-15 inspired pipeline: request passes through layers of middleware
- Each middleware can modify request, short-circuit, or pass to next layer
- Terminable middleware implements `terminate(request, response)`
- Middleware assigned in `app/Http/Kernel.php` (Go equivalent: config or bootstrap)

### 4.3 Go Implementation
- Middleware signature: `func(http.Handler) http.Handler`
- Pipeline built using nested handlers:
  ```go
  func (m Middleware) Handle(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          // Before request
          next.ServeHTTP(w, r)
          // After request (for terminable, use goroutine or defer)
      })
  }
  ```
- Middleware stack applied in reverse order so first registered runs first
- Support for middleware groups (web, api) similar to Laravel's `$middlewareGroups`
- Configurable in `bootstrap/app.go` or `config/app.go`

---

## 5. CONTROLLER LAYER

### 5.1 Features Required
- Base controller with helper methods
- Resource controllers (index, create, store, show, edit, update, destroy)
- Invokable controllers (single action)
- Controller middleware (defined in controller constructor)
- Request injection into actions
- Response helpers (view, json, redirect, download, etc.)
- Form requests (validation + authorization per controller action)

### 5.2 How It Works in Laravel
- Controllers extend base controller
- Actions are public methods
- Dependency injection via constructor or method signature
- Form requests validate before controller action runs
- Resource routing auto-maps to standard CRUD methods

### 5.3 Go Implementation
- Controllers are structs with injected dependencies:
  ```go
  type UserController struct {
      users *services.UserService
      validator *validation.Validator
  }
  ```
- Actions are methods with signature compatible with HTTP handlers
- Use reflection or code generation for dependency injection into actions
- Resource controllers implement a standard interface:
  ```go
  type ResourceController interface {
      Index(w http.ResponseWriter, r *http.Request)
      Show(w http.ResponseWriter, r *http.Request)
      // ... etc
  }
  ```
- Form requests are structs with validation rules and authorize method

---

## 6. ORM / DATABASE LAYER (Goquent)

### 6.1 Features Required
- Query Builder (fluent interface)
- Model definitions with structs
- Relationships: HasOne, HasMany, BelongsTo, BelongsToMany, HasManyThrough, MorphOne, MorphMany, MorphTo, MorphToMany
- Eager loading (with `.With()`)
- Lazy loading protection (optional)
- Mass assignment protection (fillable/guarded)
- Mutators & accessors (get/set transformations)
- Casting (dates, JSON, encrypted, custom)
- Scopes (local and global)
- Soft deletes
- Timestamps (created_at, updated_at)
- Touching parent timestamps
- Collections (model collections with map/filter/reduce)
- Pagination integration
- Transactions
- Database connections (read/write splitting)
- Raw queries with model hydration
- Observers (creating, created, updating, updated, saving, saved, deleting, deleted, restoring, restored)

### 6.2 How It Works in Laravel
- Eloquent extends Model base class
- Uses ActiveRecord pattern
- Magic methods for dynamic property access (`__get`, `__set`)
- Query builder returns Builder instance for chaining
- Relationships defined as methods returning relation types
- Uses PDO for database abstraction

### 6.3 Go Implementation
- Use `database/sql` with driver abstraction or `sqlx`/`gorm` as reference
- Models are structs with tags:
  ```go
  type User struct {
      ID        uint      `db:"id" gow:"primaryKey;autoIncrement"`
      Name      string    `db:"name" gow:"fillable"`
      Email     string    `db:"email" gow:"fillable;unique"`
      Password  string    `db:"password" gow:"hidden"`
      CreatedAt time.Time `db:"created_at" gow:"autoCreateTime"`
      UpdatedAt time.Time `db:"updated_at" gow:"autoUpdateTime"`
      DeletedAt *time.Time `db:"deleted_at" gow:"softDelete"`
  }
  ```
- Query builder using method chaining (fluent API):
  ```go
  users, err := gow.DB().Table("users").
      Where("active", true).
      Where("age", ">", 18).
      OrderBy("name", "asc").
      With("posts", "comments").
      Paginate(15)
  ```
- Relationships via struct fields with tags or methods:
  ```go
  type User struct {
      // ... fields
      Posts []Post `gow:"hasMany:posts;foreignKey:user_id"`
  }
  ```
- Use reflection for struct scanning and field mapping
- Connection pooling via `sql.DB`
- Support multiple drivers: PostgreSQL, MySQL, SQLite, SQL Server

---

## 7. MIGRATION SYSTEM

### 7.1 Features Required
- Create/alter/drop tables
- Column types and modifiers
- Indexes (primary, unique, index, fulltext, spatial)
- Foreign keys with constraints
- Running migrations (up)
- Rolling back migrations (down)
- Migration status
- Fresh (drop all and remigrate)
- Refresh (rollback and remigrate)
- Seed after migrate
- Single file migrations with timestamps
- Schema builder (blueprint)

### 7.2 How It Works in Laravel
- Migration files contain `up()` and `down()` methods
- Uses Schema facade with Blueprint for table definitions
- Tracks ran migrations in `migrations` table
- Artisan commands: migrate, rollback, status, fresh, refresh

### 7.3 Go Implementation
- Migration files as Go files in `database/migrations/`:
  ```go
  // 2026_05_20_000001_create_users_table.go
  type CreateUsersTable struct{}

  func (m *CreateUsersTable) Up(schema *gow.Schema) {
      schema.Create("users", func(table *gow.Blueprint) {
          table.ID()
          table.String("name", 255)
          table.String("email", 255).Unique()
          table.String("password", 255)
          table.Timestamps()
          table.SoftDeletes()
      })
  }

  func (m *CreateUsersTable) Down(schema *gow.Schema) {
      schema.Drop("users")
  }
  ```
- Migration runner scans directory, sorts by timestamp, checks `migrations` table
- Schema builder generates DDL SQL per driver
- CLI command: `go run artisan.go migrate`

---

## 8. SEEDING & FACTORIES

### 8.1 Features Required
- Database seeders for inserting test/initial data
- Model factories for generating fake data
- Factory states (e.g., admin state for User)
- Factory relationships (create with related models)
- Seeding during tests
- Calling seeders from migrations

### 8.2 How It Works in Laravel
- Seeders extend Seeder class, define `run()` method
- Factories use Faker for fake data generation
- `factory(User::class)->create()` generates model
- States modify factory definitions

### 8.3 Go Implementation
- Seeders in `database/seeders/`:
  ```go
  type UserSeeder struct{}

  func (s *UserSeeder) Run(db *gow.DB) {
      factory.User().Count(50).Create()
  }
  ```
- Factories using `gofakeit` or `faker` library:
  ```go
  type UserFactory struct{}

  func (f *UserFactory) Definition() map[string]interface{} {
      return map[string]interface{}{
          "name": gofakeit.Name(),
          "email": gofakeit.Email(),
          "password": hash.Make("password"),
      }
  }
  ```
- Factory methods: `Make()`, `Create()`, `Count(n)`, `State()`, `Sequence()`

---

## 9. VIEW / TEMPLATE ENGINE (Blade-like)

### 9.1 Features Required
- Template inheritance (extends, sections, yields)
- Includes & components
- Layouts
- Control structures (if, unless, for, foreach, while, switch)
- Echo with auto-escaping (`{{ }}`) and raw (`{!! !!}`)
- Comments (`{{-- --}}`)
- Stacks for pushing content to layouts
- Components & slots (anonymous and class-based)
- Directives (@auth, @guest, @csrf, @method, @error, etc.)
- Custom directives registration
- Template compilation/caching
- Sharing data across all views
- View composers (bind data to views)
- Localization in views

### 9.2 How It Works in Laravel
- Blade templates compiled to plain PHP
- `@extends('layout')` defines parent template
- `@section('content')` / `@yield('content')` for sections
- Components are classes or anonymous templates
- Views resolved via dot notation (`view('users.index')`)

### 9.3 Go Implementation
- Use Go's `html/template` as base with custom parsing layer
- Create Blade-like syntax parser that compiles to Go templates:
  - `@extends("layouts.app")` → `{{template "layouts/app" .}}`
  - `@section("content")` → `{{define "content"}}`
  - `@yield("content")` → `{{template "content" .}}`
  - `{{ $var }}` → `{{html .Var}}` (auto-escape)
  - `{!! $var !!}` → `{{.Var}}` (raw HTML)
- Pre-compile templates on boot and cache in `storage/framework/views/`
- Component system:
  ```go
  // resources/views/components/alert.gohtml
  <div class="alert alert-{{ .Type }}">
      {{ .Slot }}
  </div>
  ```
- View composers registered in providers:
  ```go
  app.View().Composer("layouts.app", func(data ViewData) {
      data.Set("notifications", getNotifications())
  })
  ```

---

## 10. SERVICE CONTAINER & IoC

### 10.1 Features Required
- Binding interfaces to implementations
- Singleton bindings
- Scoped bindings (request-scoped)
- Instance binding (existing object)
- Primitive binding (with context)
- Tagging (bind multiple implementations to a tag)
- Resolving with automatic dependency injection
- Contextual binding (different implementations per consumer)
- Container events (resolving callbacks)
- Service provider integration

### 10.2 How It Works in Laravel
- Container is the core of the framework
- `bind()`, `singleton()`, `scoped()`, `instance()` methods
- Reflection-based resolution of constructor dependencies
- `app()->make(ClassName::class)` resolves with dependencies
- Service providers register bindings during bootstrap

### 10.3 Go Implementation
- Container using map[string]interface{} with type keys:
  ```go
  type Container struct {
      bindings map[reflect.Type]binding
      instances map[reflect.Type]interface{}
      aliases map[string]reflect.Type
  }
  ```
- Binding methods:
  ```go
  container.Bind((*Interface)(nil), &Implementation{})
  container.Singleton((*Interface)(nil), &Implementation{})
  container.Instance("config", configInstance)
  ```
- Resolution uses reflection to inspect constructor parameters and resolve recursively
- Support for factory functions:
  ```go
  container.Bind((*Cache)(nil), func(c *Container) *RedisCache {
      return &RedisCache{Conn: c.Make("redis")}
  })
  ```
- Go's lack of generics in older versions requires interface{} with type assertions; use Go 1.18+ generics for type safety where possible

---

## 11. SERVICE PROVIDERS

### 11.1 Features Required
- Register method (bind services to container)
- Boot method (run after all providers registered)
- Deferred providers (only load when needed)
- Provider discovery (auto-register)
- Package service providers

### 11.2 How It Works in Laravel
- Providers listed in `config/app.php`
- `register()` should only bind, no external services
- `boot()` can use other services, routes, views
- Deferred providers implement `DeferrableProvider` with `provides()` method

### 11.3 Go Implementation
- Provider interface:
  ```go
  type ServiceProvider interface {
      Register(app *Application)
      Boot()
      IsDeferred() bool
      Provides() []string  // binding names this provider offers
  }
  ```
- Base provider struct for easy extension:
  ```go
  type RouteServiceProvider struct {
      BaseProvider
  }

  func (p *RouteServiceProvider) Boot() {
      p.LoadRoutes()
      p.ConfigureRateLimiting()
  }
  ```
- Provider manifest for deferred loading: cache which provider provides which bindings

---

## 12. CONFIGURATION MANAGEMENT

### 12.1 Features Required
- Multi-environment config (dev, testing, production)
- .env file support
- Config caching for production (compile to single file)
- Nested config access (dot notation)
- Default values
- Type casting (string, int, bool, array)
- Environment-specific overrides
- Encryption of sensitive config values

### 12.2 How It Works in Laravel
- Config files in `config/` return arrays
- `config('app.name')` accesses nested values
- `.env` loaded by DotEnv library
- `php artisan config:cache` compiles all config to one file

### 12.3 Go Implementation
- Config as Go structs in `config/` files:
  ```go
  // config/app.go
  type AppConfig struct {
      Name string `env:"APP_NAME" default:"Gow"`
      Env  string `env:"APP_ENV" default:"production"`
      Debug bool `env:"APP_DEBUG" default:"false"`
      URL string `env:"APP_URL" default:"http://localhost"`
      Key string `env:"APP_KEY" required:"true"`
  }
  ```
- Use `joho/godotenv` for .env parsing
- Config loader uses reflection to populate structs from env vars
- Config cache: serialize loaded config to gob/json in `bootstrap/cache/config.gob`
- Access via:
  ```go
  config.Get("app.name")  // string
  config.GetBool("app.debug")
  config.GetInt("database.connections.mysql.port")
  ```

---

## 13. AUTHENTICATION & AUTHORIZATION

### 13.1 Authentication Features
- Multi-guard authentication (web, api)
- Session-based auth (cookies)
- Token-based auth (API tokens, Sanctum-like)
- JWT support
- Registration & login
- Password reset / forgot password
- Email verification
- Remember me functionality
- Auth middleware (auth, guest)
- Route protection
- User provider abstraction (database, LDAP, etc.)
- Password confirmation (for sensitive actions)
- Two-factor authentication (2FA)

### 13.2 Authorization Features
- Gates ( closures for authorization logic)
- Policies (class-based authorization per model)
- Roles & permissions (Spatie-like package integration)
- Blade directives (@can, @cannot, @canany)
- Middleware (can:update,post)
- Exception handling (403 Forbidden)

### 13.3 How It Works in Laravel
- Auth facade accesses guard instances
- Guards define how users are authenticated (session, token)
- User providers retrieve users from storage
- Gates defined in `AuthServiceProvider`
- Policies mapped to models in `AuthServiceProvider`
- `authorize()` method in controllers throws 403 if fails

### 13.4 Go Implementation
- Auth system:
  ```go
  type AuthManager struct {
      guards map[string]Guard
      providers map[string]UserProvider
  }

  type Guard interface {
      Check() bool
      Guest() bool
      User() Authenticatable
      ID() interface{}
      Validate(credentials map[string]string) bool
      Login(user Authenticatable, remember bool)
      Logout()
  }
  ```
- Session guard: uses encrypted cookies + session store
- Token guard: checks Authorization header for Bearer/API tokens
- Password hashing using bcrypt (golang.org/x/crypto/bcrypt)
- Policies:
  ```go
  type PostPolicy struct{}

  func (p *PostPolicy) Update(user *models.User, post *models.Post) bool {
      return user.ID == post.UserID
  }
  ```
- Gate closures:
  ```go
  app.Gate().Define("edit-settings", func(user *models.User) bool {
      return user.IsAdmin
  })
  ```
- Middleware: `middleware.Auth`, `middleware.Can("update", "post")`

---

## 14. VALIDATION SYSTEM

### 14.1 Features Required
- Rule-based validation (required, string, email, min, max, confirmed, unique, exists, etc.)
- Custom error messages
- Custom validation rules (class-based or closure)
- Conditional validation (sometimes, when)
- Array validation (nested rules)
- File validation (image, mimes, size, dimensions)
- Form request validation (auto-validates before controller)
- Validation error bags (named error bags)
- After validation hooks
- Translation of error messages
- Validator facade / helper

### 14.2 How It Works in Laravel
- `Validator::make($data, $rules, $messages)`
- Rules are pipe-separated strings or arrays
- Form requests validate in `authorize()` and `rules()` methods
- Errors stored in session for redirect back
- JSON responses return 422 with error details

### 14.3 Go Implementation
- Validation using struct tags or rule maps:
  ```go
  type CreateUserRequest struct {
      Name     string `validate:"required|string|min:3|max:255"`
      Email    string `validate:"required|email|unique:users,email"`
      Password string `validate:"required|string|min:8|confirmed"`
      Avatar   *multipart.FileHeader `validate:"nullable|image|max:2048"`
  }
  ```
- Validator interface:
  ```go
  type Validator interface {
      Validate(data map[string]interface{}, rules Rules) ValidationResult
      ValidateStruct(s interface{}) ValidationResult
  }
  ```
- Form request struct:
  ```go
  type Request interface {
      Authorize() bool
      Rules() map[string]string
      Messages() map[string]string
      Validate() error
  }
  ```
- Validation middleware checks if request implements `Request` and auto-validates
- Error response standardized:
  ```json
  {
      "message": "The given data was invalid.",
      "errors": {
          "email": ["The email field is required.", "The email must be a valid email address."]
      }
  }
  ```

---

## 15. SESSION MANAGEMENT

### 15.1 Features Required
- Multiple drivers: file, cookie, database, redis, memcached, dynamodb, array (testing)
- Session ID regeneration (security)
- Flash data (available only on next request)
- Session encryption
- CSRF token storage in session
- Session fixation protection
- Idle timeout / lifetime configuration
- Concurrent request handling (locking)
- Session start and middleware
- `Session` facade/helper

### 15.2 How It Works in Laravel
- Session started by StartSession middleware
- Session ID in encrypted cookie
- Data serialized and stored in chosen driver
- Flash data marked for deletion after read
- Regenerate ID on login for security

### 15.3 Go Implementation
- Session manager with drivers:
  ```go
  type SessionManager struct {
      driver Driver
      lifetime time.Duration
  }

  type Driver interface {
      Get(id string) (map[string]interface{}, error)
      Put(id string, data map[string]interface{}, lifetime int) error
      Destroy(id string) error
      GC(lifetime int) error
  }
  ```
- Cookie-based session ID with HttpOnly, Secure, SameSite flags
- Middleware starts session, loads data into request context
- Flash data: separate flash bag in session, cleared after read
- Redis driver using `go-redis`, database driver using `database/sql`
- Session interface:
  ```go
  type Session interface {
      Get(key string, defaultValue interface{}) interface{}
      Put(key string, value interface{})
      Forget(keys ...string)
      Flush()
      Flash(key string, value interface{})
      Reflash()
      Regenerate() string
      ID() string
      Token() string  // CSRF token
  }
  ```

---

## 16. CACHE SYSTEM

### 16.1 Features Required
- Drivers: file, database, redis, memcached, dynamodb, array, null (testing), APCu
- Cache tags (group cache keys for bulk operations)
- Atomic locks (for job queues, etc.)
- Cache forever / TTL
- Remember pattern (get or store)
- Cache events (hit, miss, delete)
- Increment/decrement
- Cache warming
- Cache flushing

### 16.2 How It Works in Laravel
- `Cache::get()`, `Cache::put()`, `Cache::remember()`
- Tags allow `Cache::tags(['people', 'artists'])->flush()`
- Locks for distributed systems
- File driver stores serialized data in files

### 16.3 Go Implementation
- Cache store interface:
  ```go
  type Store interface {
      Get(key string) (interface{}, bool)
      Many(keys []string) map[string]interface{}
      Put(key string, value interface{}, seconds int) bool
      PutMany(values map[string]interface{}, seconds int) bool
      Increment(key string, value int) int
      Decrement(key string, value int) int
      Forever(key string, value interface{}) bool
      Forget(key string) bool
      Flush() bool
      GetLock(key string) Lock
  }
  ```
- Redis implementation using `go-redis`
- File implementation using gob-encoded files in `storage/framework/cache/`
- Remember helper:
  ```go
  value := cache.Remember("users", 60, func() interface{} {
      return db.Table("users").Get()
  })
  ```
- Tags implemented via separate tag index (set of keys per tag)

---

## 17. QUEUE & JOB SYSTEM

### 17.1 Features Required
- Multiple drivers: sync (immediate), database, redis, sqs, beanstalkd, rabbitmq
- Job classes with handle method
- Job middleware
- Delayed dispatching
- Job chaining (sequential execution)
- Batch jobs (parallel with completion callback)
- Failed job handling (retry, max attempts, backoff)
- Job events (before, after, exception)
- Rate limiting for queues
- Worker CLI command
- Job serialization (store class + payload)
- Unique jobs (prevent duplicates)
- Job timeouts
- Job retries with exponential backoff

### 17.2 How It Works in Laravel
- Jobs implement `ShouldQueue` interface
- `dispatch()` sends to queue driver
- Workers run via `php artisan queue:work`
- Failed jobs stored in `failed_jobs` table
- Batches track completion via `job_batches` table

### 17.3 Go Implementation
- Job interface:
  ```go
  type Job interface {
      Handle() error
      Failed(err error)
      RetryUntil() time.Time
      Backoff() []int  // seconds between retries
      Timeout() int    // seconds
  }
  ```
- Queue drivers:
  ```go
  type Queue interface {
      Push(job Job, queue string) error
      Later(delay time.Duration, job Job, queue string) error
      Pop(queue string) (Job, error)
      Size(queue string) int
  }
  ```
- Database queue: table with `id`, `queue`, `payload`, `attempts`, `reserved_at`, `available_at`, `created_at`
- Worker command: `go run artisan.go queue:work --queue=default --sleep=3 --tries=3`
- Job payload serialized as JSON with type info and data
- Failed jobs table with error message and stack trace

---

## 18. EVENT & LISTENER SYSTEM

### 18.1 Features Required
- Event classes (plain structs)
- Listener classes with handle method
- Event auto-discovery
- Queued listeners
- Event subscribers (class with multiple handlers)
- Wildcard listeners
- Event broadcasting integration
- Stoppable events
- Event fake for testing

### 18.2 How It Works in Laravel
- `event(new OrderShipped($order))` dispatches event
- Listeners registered in `EventServiceProvider`
- `ShouldQueue` interface queues listener
- `broadcastOn()` method for broadcasting

### 18.3 Go Implementation
- Event dispatcher:
  ```go
  type Dispatcher struct {
      listeners map[reflect.Type][]Listener
      wildcards []WildcardListener
      queue Queue
  }

  type Listener interface {
      Handle(event interface{}) error
      ShouldQueue() bool
  }
  ```
- Event registration:
  ```go
  app.Events().Listen((*OrderShipped)(nil), &SendShipmentNotification{})
  app.Events().Subscribe(&UserEventSubscriber{})
  ```
- Dispatch:
  ```go
  app.Events().Dispatch(&OrderShipped{Order: order})
  ```
- Queued listeners dispatched via job queue
- Wildcard listeners: `app.Events().Listen("*", &LogAllEvents{})`

---

## 19. MAIL SYSTEM

### 19.1 Features Required
- Mailable classes (build method for content)
- Multiple drivers: smtp, sendmail, mailgun, postmark, ses, log (testing)
- Markdown mail support
- Attachments (files, raw data, from storage)
- CC, BCC, Reply-To
- HTML and plain text views
- Queueing mail
- Mail localization
- Preview in browser (dev mode)
- Templated mail (third-party templates)

### 19.2 How It Works in Laravel
- Mailable classes with `envelope()`, `content()`, `attachments()` methods
- Blade views for HTML/text content
- Markdown mailables convert markdown to HTML
- `Mail::to($user)->send(new OrderShipped($order))`

### 19.3 Go Implementation
- Mailable interface:
  ```go
  type Mailable interface {
      Envelope() Envelope
      Content() Content
      Attachments() []Attachment
  }

  type Envelope struct {
      From Address
      To []Address
      CC []Address
      BCC []Address
      ReplyTo []Address
      Subject string
      Tags []string
      Metadata map[string]string
  }
  ```
- Mailer interface with drivers:
  ```go
  type Mailer interface {
      Send(mailable Mailable) error
      Queue(mailable Mailable) error
      Later(delay time.Duration, mailable Mailable) error
  }
  ```
- Use `gomail.v2` or `email` package for SMTP
- Markdown support via `goldmark` parser
- Log driver for testing writes to `storage/logs/mail.log`

---

## 20. NOTIFICATION SYSTEM

### 20.1 Features Required
- Notifiable trait/interface (User can receive notifications)
- Notification classes with via() and channel methods
- Channels: mail, database, broadcast, sms, slack, discord
- On-demand notifications (to non-users)
- Queueing notifications
- Notification localization
- Mark as read/unread
- Database notification storage
- Notification events

### 20.2 How It Works in Laravel
- `User->notify(new InvoicePaid($invoice))`
- `via()` returns channels to use
- Each channel has `toMail()`, `toDatabase()`, etc.
- Database notifications stored in `notifications` table (polymorphic)

### 20.3 Go Implementation
- Notifiable interface:
  ```go
  type Notifiable interface {
      RouteNotificationFor(channel string) interface{}
      GetKey() interface{}
  }
  ```
- Notification interface:
  ```go
  type Notification interface {
      Via(notifiable Notifiable) []string
      ToMail(notifiable Notifiable) *MailMessage
      ToDatabase(notifiable Notifiable) map[string]interface{}
      ToBroadcast(notifiable Notifiable) *BroadcastMessage
      ShouldQueue() bool
  }
  ```
- Database channel stores JSON in `notifications` table:
  ```sql
  CREATE TABLE notifications (
      id UUID PRIMARY KEY,
      type VARCHAR(255),
      notifiable_type VARCHAR(255),
      notifiable_id BIGINT,
      data JSON,
      read_at TIMESTAMP,
      created_at TIMESTAMP
  );
  ```
- Notification manager dispatches to appropriate channels

---

## 21. FILE STORAGE

### 21.1 Features Required
- Multiple disks: local, public, s3, sftp, gcs, azure, ftp
- File upload handling
- File visibility (public/private)
- Temporary URLs (for private files)
- File streaming
- File metadata (size, mime, last modified)
- Directory creation/listing/deletion
- File copying/moving
- Storage facade / helper
- UploadedFile wrapper
- Image manipulation integration (optional)

### 21.2 How It Works in Laravel
- `Storage::disk('s3')->put('file.txt', 'contents')`
- Filesystems configured in `config/filesystems.php`
- `store()` method on UploadedFile handles moving from temp to disk
- Public disk linked to `storage/app/public`

### 21.3 Go Implementation
- Storage manager:
  ```go
  type StorageManager struct {
      disks map[string]Filesystem
      defaultDisk string
  }

  type Filesystem interface {
      Exists(path string) bool
      Get(path string) ([]byte, error)
      Put(path string, contents []byte) error
      Delete(paths ...string) bool
      Copy(from, to string) bool
      Move(from, to string) bool
      Size(path string) int64
      LastModified(path string) int64
      Files(directory string) []string
      AllFiles(directory string) []string
      Directories(directory string) []string
      MakeDirectory(path string) bool
      DeleteDirectory(path string) bool
      URL(path string) string
      TemporaryURL(path string, expiration time.Time) string
      SetVisibility(path string, visibility string) bool
  }
  ```
- Local driver using `os` package
- S3 driver using AWS SDK for Go v2
- Public files served via symlink or copied to `public/storage/`
- UploadedFile struct wrapping `multipart.FileHeader`:
  ```go
  type UploadedFile struct {
      *multipart.FileHeader
  }

  func (f *UploadedFile) Store(path string, disk string) string {
      // Move to storage disk
  }
  ```

---

## 22. BROADCASTING / WEBSOCKETS

### 22.1 Features Required
- Drivers: pusher, ably, redis (pub/sub), log (testing)
- Public channels
- Private channels (auth required)
- Presence channels (who's online)
- Channel authorization (channel routes)
- Event broadcasting (implements ShouldBroadcast)
- Broadcast queues
- Socket ID tracking
- Whisper (client events)
- Echo-like JS client integration

### 22.2 How It Works in Laravel
- Events implement `ShouldBroadcast` interface
- `broadcastOn()` returns channel array
- Private channels prefixed with `private-`, presence with `presence-`
- `routes/channels.php` defines authorization closures
- Pusher/Ably handle WebSocket delivery

### 22.3 Go Implementation
- Broadcaster interface:
  ```go
  type Broadcaster interface {
      Auth(request *http.Request) (map[string]interface{}, error)
      ValidAuthenticationResponse(request *http.Request, result map[string]interface{}) (map[string]interface{}, error)
      Broadcast(channels []string, event string, payload map[string]interface{}) error
  }
  ```
- Pusher driver using Pusher Go library
- Redis driver using Redis pub/sub for self-hosted WebSocket servers
- Channel authorization in `routes/channels.go`:
  ```go
  func ChannelRoutes(broadcaster *gow.Broadcaster) {
      broadcaster.Private("orders.{orderId}", func(user *models.User, orderId int) bool {
          return user.ID == orderId
      })
  }
  ```
- Presence channel data:
  ```go
  broadcaster.Presence("chat.{roomId}", func(user *models.User, roomId int) map[string]interface{} {
      return map[string]interface{}{
          "id": user.ID,
          "name": user.Name,
      }
  })
  ```

---

## 23. PAGINATION

### 23.1 Features Required
- Length-aware pagination (total count known)
- Simple pagination (next/prev only, no total)
- Cursor pagination (for large datasets, performance)
- API resource pagination (JSON response)
- Bootstrap / Tailwind / Custom views
- Pagination links generation
- Query string preservation
- Fragment appending (#section)
- Per-page configuration

### 23.2 How It Works in Laravel
- `Model::paginate(15)` returns LengthAwarePaginator
- `Model::simplePaginate(15)` returns Paginator
- `Model::cursorPaginate(15)` returns CursorPaginator
- Views render pagination links
- JSON responses include links/meta

### 23.3 Go Implementation
- Paginator structs:
  ```go
  type LengthAwarePaginator struct {
      Data []interface{}
      Total int64
      PerPage int
      CurrentPage int
      LastPage int
      Path string
      Query map[string]string
      Fragment string
  }

  func (p *LengthAwarePaginator) Links() []PageLink {
      // Generate prev, page numbers, next links
  }

  func (p *LengthAwarePaginator) ToJSON() PaginatorJSON {
      return PaginatorJSON{
          Data: p.Data,
          Meta: Meta{
              CurrentPage: p.CurrentPage,
              LastPage: p.LastPage,
              Total: p.Total,
              PerPage: p.PerPage,
          },
          Links: p.Links(),
      }
  }
  ```
- Query builder integration:
  ```go
  users, err := db.Table("users").Where("active", true).Paginate(15, request.Page())
  ```
- Cursor pagination uses encoded cursor (last ID + sort values):
  ```go
  users, err := db.Table("users").CursorPaginate(15, request.Cursor())
  ```

---

## 24. LOCALIZATION (i18n)

### 24.1 Features Required
- Language files (JSON and Go map formats)
- Fallback locale
- Pluralization (one, many, other)
- Replacement parameters (`:name`)
- Choice/plural forms
- Locale detection (header, cookie, session, URL parameter)
- Translation helpers (__(), trans(), trans_choice())
- JSON group translations
- Publishing language files

### 24.2 How It Works in Laravel
- Files in `resources/lang/{locale}/`
- `__('messages.welcome')` or `trans('messages.welcome')`
- JSON files for simple key-value: `{"Welcome": "Bienvenue"}`
- `trans_choice()` handles pluralization

### 24.3 Go Implementation
- Translator:
  ```go
  type Translator struct {
      loader Loader
      locale string
      fallback string
      loaded map[string]map[string]string
  }

  type Loader interface {
      Load(locale, group string) map[string]string
      LoadJSON(locale string) map[string]string
  }
  ```
- Translation files in `resources/lang/en/messages.json`:
  ```json
  {
      "welcome": "Welcome, :name!",
      "apples": "{0} No apples|{1} One apple|[2,*] :count apples"
  }
  ```
- Helper functions:
  ```go
  lang.Get("messages.welcome", map[string]string{"name": "John"})
  lang.Choice("messages.apples", 5)  // "5 apples"
  ```
- Pluralization using CLDR rules or simple one/other
- Locale middleware sets locale from Accept-Language header or URL prefix

---

## 25. TASK SCHEDULING

### 25.1 Features Required
- Cron-like scheduling (everyMinute, hourly, daily, weekly, monthly, etc.)
- Custom cron expressions
- Task callbacks (closures, commands, jobs)
- Overlap prevention (withoutOverlapping)
- On-one-server (for multi-node deployments)
- Background/foreground execution
- Email output
- Task hooks (before, after, onSuccess, onFailure)
- Maintenance mode handling
- Timezone support
- Mutex/locking for overlap prevention

### 25.2 How It Works in Laravel
- `app/Console/Kernel.php` `schedule()` method
- `php artisan schedule:run` triggered by system cron every minute
- `php artisan schedule:work` runs scheduler in foreground
- Uses cache locks for overlap prevention

### 25.3 Go Implementation
- Scheduler:
  ```go
  type Scheduler struct {
      events []Event
      mutex Mutex
  }

  type Event struct {
      expression string  // cron expression
      callback func()
      mutexName string
      withoutOverlapping bool
      onOneServer bool
      timezone *time.Location
  }
  ```
- Cron expression parser (use `robfig/cron` library)
- Schedule registration in `app/console/kernel.go`:
  ```go
  func (k *Kernel) Schedule(scheduler *gow.Scheduler) {
      scheduler.Call(func() {
          // Cleanup old records
      }).Daily()

      scheduler.Command("report:generate").Weekly().Mondays().At("13:00")

      scheduler.Job(&CleanupJob{}).EveryFiveMinutes().WithoutOverlapping()
  }
  ```
- System cron calls `go run artisan.go schedule:run` every minute
- Overlap prevention using file locks or Redis locks

---

## 26. TESTING FRAMEWORK

### 26.1 Features Required
- PHPUnit-like test runner (use Go's testing package)
- HTTP testing (GET, POST, etc. with assertions)
- Database transactions per test
- Database seeding for tests
- Mocking/faking: mail, queue, events, notifications, storage, cache, broadcasting
- Authentication helpers (actingAs)
- Form request testing
- Console command testing
- Browser testing integration (Laravel Dusk equivalent - optional)
- Factory states in tests
- Snapshot testing (optional)
- Parallel testing support

### 26.2 How It Works in Laravel
- `Tests\TestCase` base class
- `->get('/')->assertStatus(200)`
- Database migrations refreshed for each test class
- `Mail::fake()`, `Queue::fake()` intercept dispatches
- `actingAs($user)` sets authenticated user

### 26.3 Go Implementation
- Test case base:
  ```go
  type TestCase struct {
      *httptest.Server
      app *Application
      t *testing.T
  }

  func (tc *TestCase) Get(uri string) *TestResponse {
      // Make request, return response wrapper
  }
  ```
- Response assertions:
  ```go
  response := tc.Get("/")
  response.AssertStatus(200)
  response.AssertJson(map[string]interface{}{"status": "ok"})
  response.AssertViewIs("welcome")
  response.AssertRedirect("/login")
  ```
- Fakes:
  ```go
  func (tc *TestCase) Setup() {
      tc.FakeMail()
      tc.FakeQueue()
      tc.FakeEvents()
  }
  ```
- Database testing:
  ```go
  func TestUserCreation(t *testing.T) {
      tc := NewTestCase(t)
      tc.RefreshDatabase()  // run migrations
      tc.Seed(&UserSeeder{})

      // Test...
  }
  ```
- Use `testify` for assertions

---

## 27. CLI / ARTISAN EQUIVALENT

### 27.1 Features Required
- Command registration and discovery
- Built-in commands: serve, migrate, make:controller, make:model, make:migration, route:list, cache:clear, config:cache, queue:work, schedule:run, tinker (REPL), optimize, db:seed, key:generate
- Command arguments and options
- Input/output helpers (ask, confirm, choice, table, progress bar)
- Exit codes
- Command scheduling integration
- REPL / interactive shell (optional)
- Code generation scaffolding

### 27.2 How It Works in Laravel
- Artisan commands extend `Command` class
- Signature defines arguments/options
- `handle()` method contains logic
- `php artisan list` shows all commands
- Generators create boilerplate files

### 27.3 Go Implementation
- Use `spf13/cobra` or `urfave/cli` as base
- Command interface:
  ```go
  type Command interface {
      Name() string
      Description() string
      Handle(ctx *Context) error
  }
  ```
- Built-in commands structure:
  ```go
  // cmd/gow/main.go or artisan.go
  var rootCmd = &cobra.Command{Use: "artisan"}

  func init() {
      rootCmd.AddCommand(commands.ServeCmd)
      rootCmd.AddCommand(commands.MigrateCmd)
      rootCmd.AddCommand(commands.MakeControllerCmd)
      // ... etc
  }
  ```
- Code generators use text/template to create files from stubs:
  ```go
  // stubs/controller.stub
  package controllers

  type {{.Name}}Controller struct {
      BaseController
  }
  ```
- REPL using `chzyer/readline` with evaluation of Go expressions (limited compared to PHP)

---

## 28. LOGGING SYSTEM

### 28.1 Features Required
- Multiple channels: single, daily, slack, syslog, errorlog, custom
- Log levels: debug, info, notice, warning, error, critical, alert, emergency
- Contextual data (arrays/objects with message)
- Formatters (line, JSON)
- Processors (add extra data)
- Log rotation (daily, size-based)
- Slack/email alerts for critical errors
- Tap/customization
- Monolog-like channel stacking
- Request ID injection
- Async logging (optional)

### 28.2 How It Works in Laravel
- `Log::info('Message', ['context' => 'data'])`
- Channels defined in `config/logging.php`
- Stack channel combines multiple channels
- Daily channel rotates files by date

### 28.3 Go Implementation
- Use `uber-go/zap` or `sirupsen/logrus` as base, or standard `log/slog` (Go 1.21+)
- Logger interface:
  ```go
  type Logger interface {
      Debug(msg string, fields ...Field)
      Info(msg string, fields ...Field)
      Warn(msg string, fields ...Field)
      Error(msg string, fields ...Field)
      Fatal(msg string, fields ...Field)
      With(fields ...Field) Logger
      WithContext(ctx context.Context) Logger
  }
  ```
- Channel configuration:
  ```go
  // config/logging.go
  type LoggingConfig struct {
      Default string
      Channels map[string]Channel
  }

  type Channel struct {
      Driver string  // single, daily, slack, syslog
      Path string    // for file drivers
      Level string   // minimum level
      Days int       // retention for daily
      WebhookURL string // for slack
  }
  ```
- Request ID middleware adds `request_id` to all log entries in context
- Slack alert channel sends webhooks for critical+

---

## 29. SECURITY FEATURES

### 29.1 Required Security Features
1. **CSRF Protection**
   - Token generation per session
   - Token validation on state-changing requests (POST, PUT, PATCH, DELETE)
   - Exclude routes from CSRF (API routes)
   - Token refresh on login

2. **XSS Protection**
   - Auto-escaping in template engine (`{{ }}`)
   - Raw output explicit (`{!! !!}`)
   - Content Security Policy (CSP) headers middleware
   - HTMLPurifier integration for WYSIWYG content

3. **SQL Injection Prevention**
   - Parameterized queries exclusively in ORM
   - Query builder binds all values
   - Raw query parameter binding enforcement
   - No string concatenation in query generation

4. **Encryption**
   - AES-256-GCM for data at rest
   - App key management (APP_KEY in .env)
   - Encrypted cookies
   - Encrypted session data (optional)
   - Encrypted casts in models

5. **Hashing**
   - Bcrypt for passwords (cost factor 12+)
   - Argon2id support (optional)
   - Rehash on login if parameters changed
   - No plain text password storage anywhere

6. **CORS**
   - Configurable allowed origins, methods, headers
   - Preflight request handling
   - Credential support
   - Max age configuration

7. **Rate Limiting**
   - Per-route or global rate limits
   - By IP, user, or custom key
   - Multiple limit tiers (e.g., 60/min, 1000/hour)
   - Response headers (X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After)
   - Rate limit exceptions for specific users/routes

8. **Input Sanitization**
   - Trim strings middleware
   - Convert empty strings to null (optional)
   - Strip tags middleware for specific fields
   - File upload validation (extension, mime, size, content scanning)

9. **Security Headers**
   - X-Content-Type-Options: nosniff
   - X-Frame-Options: DENY/SAMEORIGIN
   - Strict-Transport-Security (HSTS)
   - X-XSS-Protection (legacy browsers)
   - Referrer-Policy
   - Permissions-Policy
   - Content-Security-Policy

10. **Authentication Security**
    - Brute force protection (throttle logins)
    - Password confirmation for sensitive actions
    - Session fixation protection (regenerate ID on login)
    - Concurrent session limit (optional)
    - Device tracking (optional)
    - Login alerts (optional)

11. **Authorization Security**
    - Gate/Policy authorization before any sensitive action
    - Mass assignment protection (fillable/guarded)
    - Route authorization middleware
    - API scope checking (OAuth2/Sanctum-like)

12. **Dependency Security**
    - Vulnerability scanning in Go modules (`govulncheck`)
    - Minimal dependency policy
    - Regular security updates
    - No cgo unless necessary (reduces attack surface)

13. **Error Handling Security**
    - No stack traces or sensitive data in production responses
    - Generic error messages to users, detailed logs internally
    - Exception filtering for sensitive fields

14. **Secret Management**
    - .env file outside web root
    - Encrypted environment variables (optional)
    - Secret rotation support
    - No secrets in code or logs

### 29.2 Implementation Notes
- Use `golang.org/x/crypto` for bcrypt, scrypt, argon2
- Use `crypto/aes` with GCM mode for encryption
- Use `net/http` security headers middleware
- Implement rate limiting via Redis or in-memory token bucket
- CSRF token in session + hidden form field + double-submit cookie for SPA/API hybrid

---

## 30. ERROR HANDLING & EXCEPTIONS

### 30.1 Features Required
- Exception/Error classes (HTTP exceptions, validation exceptions, etc.)
- Global exception handler
- Custom error pages (404, 500, 503)
- Error reporting (log, sentry, etc.)
- Renderable exceptions (custom response)
- Reportable exceptions (custom logging)
- Whoops-like detailed error pages (dev only)
- Stack trace filtering
- Ignored exceptions (don't report)
- Exception transformation (all to JSON in API context)
- Maintenance mode (503 with retry-after)

### 30.2 How It Works in Laravel
- `app/Exceptions/Handler.php` handles all exceptions
- `report()` logs, `render()` creates response
- `abort(404)` throws HttpException
- Maintenance mode via `php artisan down`

### 30.3 Go Implementation
- Go uses errors, not exceptions. Adapt pattern:
  ```go
  type HttpException struct {
      StatusCode int
      Message string
      Headers map[string]string
  }

  func (e *HttpException) Error() string {
      return e.Message
  }
  ```
- Exception handler:
  ```go
  type ExceptionHandler struct {
      dontReport []reflect.Type
      reporters []Reporter
      renderers map[reflect.Type]Renderer
  }

  func (h *ExceptionHandler) Handle(err error, w http.ResponseWriter, r *http.Request) {
      // Check if renderable
      // Log if reportable
      // Return appropriate response
  }
  ```
- Error pages:
  - Development: detailed stack trace using `runtime/debug.Stack()`
  - Production: generic error view, log detailed to file
- Maintenance mode middleware checks `storage/framework/down` file
- Panic recovery middleware catches unexpected panics and converts to 500

---

## 31. PACKAGE DEVELOPMENT

### 31.1 Features Required
- Package service providers
- Package discovery (auto-register providers)
- Vendor publishing (config, views, migrations, assets)
- Package routing
- Package views/namespaces
- Package configuration merging
- Package migrations
- Package assets
- Package commands

### 31.2 How It Works in Laravel
- `vendor:publish` copies package resources to app
- `PackageServiceProvider` registers package resources
- `loadViewsFrom()`, `loadMigrationsFrom()`, `mergeConfigFrom()`
- Auto-discovery via `composer.json` extra

### 31.3 Go Implementation
- Package provider:
  ```go
  type PackageServiceProvider struct {
      BaseProvider
      packageRoot string
  }

  func (p *PackageServiceProvider) Boot() {
      p.LoadViewsFrom(p.packageRoot+"/resources/views", "package-name")
      p.LoadMigrationsFrom(p.packageRoot+"/database/migrations")
      p.MergeConfigFrom(p.packageRoot+"/config/package.go", "package")
      p.Publishes(map[string]string{
          p.packageRoot+"/config/package.go": config.Path("package.go"),
      }, "package-config")
  }
  ```
- Publish command copies files from package to app directories
- View namespaces: `view("package::viewname")`
- Config namespaces: `config.Get("package.key")`

---

## 32. IMPLEMENTATION ROADMAP

### Phase 1: Foundation (Core Framework)
- [ ] Project scaffolding and directory structure
- [ ] Service Container with basic binding/resolution
- [ ] Application bootstrap and kernel
- [ ] HTTP Router with groups, parameters, named routes
- [ ] Middleware pipeline
- [ ] Base Controller and Response helpers
- [ ] Configuration system with .env support
- [ ] Service Provider system
- [ ] Error handling and exception rendering
- [ ] Logging system
- [ ] CLI framework (Artisan equivalent)
- [ ] Helpers, Stringable & Fluent Utilities (Str, Arr, Number, Date, Benchmark)
- [ ] Request & Response Macros foundation

### Phase 2: Database & ORM
- [ ] Database connection manager (multi-driver)
- [ ] Query Builder (fluent API)
- [ ] Schema Builder & Blueprint
- [ ] Migration system with runner
- [ ] Model base with timestamps, soft deletes
- [ ] Basic relationships (HasOne, HasMany, BelongsTo)
- [ ] Eager loading (With)
- [ ] Seeders and Factories
- [ ] Pagination
- [ ] Collections & Lazy Collections
- [ ] Query Logging & Database Profiling
- [ ] Schema Dumping
- [ ] Model Pruning

### Phase 3: Web Layer
- [ ] Template engine (Blade-like compilation)
- [ ] View composers and creators
- [ ] Session management (multi-driver)
- [ ] Cache system (multi-driver)
- [ ] Validation system
- [ ] Form requests
- [ ] File Storage / Uploads
- [ ] Localization (i18n)
- [ ] API Resources / Transformers
- [ ] Advanced Blade / Template Directives (loop vars, @once, @fragment, @aware, @inject)
- [ ] Signed URLs
- [ ] Route Model Binding (implicit, explicit, custom fields, scoping)
- [ ] Old Input, Flashing, and Request casting helpers

### Phase 4: Auth & Security
- [ ] Authentication system (multi-guard)
- [ ] Password hashing and reset
- [ ] Authorization (Gates & Policies)
- [ ] CSRF protection
- [ ] Encryption system
- [ ] Rate limiting
- [ ] CORS handling
- [ ] Security headers middleware
- [ ] Feature Flags (Pennant-like)
- [ ] Health Checks

### Phase 5: Advanced Features
- [ ] Queue system (multi-driver)
- [ ] Job classes and workers
- [ ] Event & Listener system
- [ ] Mail system (multi-driver)
- [ ] Notification system
- [ ] Task Scheduling
- [ ] Broadcasting / WebSockets
- [ ] Cache tagging and atomic locks
- [ ] Process / CLI Command Execution
- [ ] Modern Laravel Parity (Reverb WebSocket server, Pulse monitoring, Prompts CLI UX, Context propagation)

### Phase 6: Testing & Tooling
- [ ] Test case base with HTTP assertions
- [ ] Fakes (Mail, Queue, Event, Notification)
- [ ] Database testing helpers
- [ ] Code generation commands (make:controller, make:model, etc.)
- [ ] Route caching
- [ ] Config caching
- [ ] View caching
- [ ] Package development support
- [ ] Testing Utilities (Database assertions, Time manipulation, Auth assertions)
- [ ] Package Auto-Discovery & Publishing

### Phase 7: Optimization & Polish
- [ ] Performance benchmarking
- [ ] Route compilation optimization
- [ ] Template compilation optimization
- [ ] Connection pooling optimization
- [ ] Memory leak testing
- [ ] Documentation
- [ ] Example application
- [ ] Package ecosystem seeding

---

## 33. GO-SPECIFIC IMPLEMENTATION NOTES

### 33.1 Language Considerations
- **No Magic Methods**: Go lacks PHP's `__get`, `__set`, `__call`. Use explicit methods, code generation, or reflection sparingly.
- **Static Typing**: Embrace interfaces and struct embedding over dynamic features.
- **Error Handling**: Use Go's explicit error returns. Wrap framework errors with context.
- **Generics**: Use Go 1.18+ generics for container resolution, collections, and type-safe helpers.
- **Concurrency**: Leverage goroutines for queue workers, scheduled tasks, and terminable middleware.
- **Context**: Use `context.Context` extensively for request-scoped values, cancellation, and timeouts.

### 33.2 Dependency Management
- Minimize external dependencies
- Preferred libraries:
  - Router: `julienschmidt/httprouter` or custom radix tree
  - ORM: Custom built on `database/sql` (GORM as reference but avoid dependency if possible)
  - Migrations: Custom (Goose or Migrate as reference)
  - Templates: `html/template` with custom parser layer
  - Redis: `go-redis/redis`
  - Testing: `stretchr/testify`
  - CLI: `spf13/cobra` or `urfave/cli`
  - Env: `joho/godotenv`
  - Validation: Custom with `go-playground/validator` as reference
  - Logging: `uber-go/zap` or standard `log/slog`
  - Cryptography: `golang.org/x/crypto`
  - Mail: `go-mail/mail` or `jordan-wright/email`
  - Queue: Custom with Redis driver
  - Cache: Custom with Redis/File driver

### 33.3 Performance Targets
- Router: < 1μs per request matching
- Template rendering: < 5ms for 1000 iterations
- Database query: Minimal overhead over raw `database/sql`
- Memory: < 50MB base footprint per worker
- Concurrency: Handle 10k+ concurrent connections

### 33.4 Go Idioms to Follow
- Accept interfaces, return structs
- Use struct tags for metadata (ORM, validation, config)
- Prefer composition over inheritance (embed structs)
- Use `sync.RWMutex` for concurrent map access in container/cache
- Use `context.WithTimeout` for DB queries and external calls
- Use `defer` for resource cleanup
- Use `go test` natively, wrap with helpers rather than replacing

### 33.5 Code Generation Strategy
Since Go is statically typed and lacks PHP's runtime flexibility, use code generation for:
- Route caching (compile routes to Go map)
- Config caching (compile config structs)
- View compilation (Blade-like to Go templates)
- ORM model helpers (generate query methods)
- Migration runners (generate migration registry)
- Service provider manifests (deferred loading)

Use `go generate` directives and custom CLI commands for this.

---

## 34. AGENT CONTINUATION GUIDELINES

### 34.1 For Any AI Agent Continuing This Project

**MANDATORY READING**: Read this entire specification before writing any code.

**Architecture Principles**:
1. Follow the directory structure in Section 1 exactly unless there's a compelling Go-specific reason to deviate.
2. Prefer interfaces over concrete types for all core framework components.
3. Keep the framework agnostic of specific drivers (DB, cache, queue, mail, storage).
4. Implement features in the order of the Implementation Roadmap (Section 32).
5. Never sacrifice Go idioms to mimic Laravel exactly. Adapt patterns to fit Go's strengths.

**Code Standards**:
- Use Go 1.24+ features where beneficial (generics, slog, etc.)
- All public APIs must be documented with Go doc comments
- Use `gofmt` / `goimports` formatting
- Error messages must be clear and actionable
- Configuration must always support .env overrides
- Never hardcode paths; use configurable base paths

**Testing Requirements**:
- Every core component must have unit tests
- Every feature must have at least one integration test
- Use table-driven tests where possible
- Mock external services (DB, mail, HTTP) in unit tests
- Maintain >80% code coverage for framework core

**Documentation Requirements**:
- Update this specification when adding/modifying features
- Add implementation notes in code for non-obvious decisions
- Maintain CHANGELOG.md for all changes (see separate file)
- Document breaking changes immediately

**Security Checklist for Every Feature**:
- [ ] Input validation and sanitization
- [ ] Output escaping where applicable
- [ ] Authentication/authorization hooks where applicable
- [ ] Rate limiting considerations
- [ ] Audit logging capability
- [ ] No secrets in code or default configs

**Before Implementing Any Feature**:
1. Check if it exists in this spec. If not, add it first.
2. Write the interface definitions
3. Write tests that fail (TDD approach)
4. Implement the minimum viable version
5. Add to CHANGELOG.md
6. Update this spec's Phase checklist

**Feature Parity Notes**:
- Not every Laravel feature needs 1:1 parity. Some PHP-specific features (Facades with magic methods, runtime trait behavior) should be adapted or omitted.
- Focus on the developer experience: fluent APIs, clear conventions, helpful error messages.
- Performance is a feature. Profile before optimizing, but design for concurrency from the start.

**When in Doubt**:
- Check how Laravel does it for the conceptual approach
- Check how Go frameworks (Gin, Echo, Beego, Revel, Buffalo) do it for the implementation approach
- Choose the approach that best serves Go developers while maintaining Laravel's elegance

**Files to Maintain**:
- `GOW_SPEC.md` (this file) - Update with implementation details as built
- `CHANGELOG.md` - Log all changes, decisions, and known issues
- `ROADMAP.md` (create when starting) - Track Phase progress
- `ARCHITECTURE.md` (create when starting) - Document core architectural decisions and patterns

---

## 35. API RESOURCES / TRANSFORMERS
*Added by Kimi K2.6 — Framework: GoW*

#### 35.1 Features Required
- Resource classes for transforming models to JSON
- Resource collections with pagination metadata
- Conditional attributes (`when`, `whenLoaded`, `whenPivotLoaded`)
- Relationship embedding (`with`, `mergeWhen`)
- Resource wrapping (`data` key wrapping)
- Pagination integration in resource responses
- Nested resources (resources within resources)
- Transformation helpers (`only`, `except`, `merge`)
- API Resource route auto-generation (optional)

#### 35.2 How It Works in Laravel
- `php artisan make:resource UserResource`
- `return new UserResource($user)` in controller
- `UserResource::collection($users)` for collections
- `return new UserCollection($users)` for custom collection classes
- Conditional fields: `$this->when(Auth::user()->isAdmin(), ...)`
- Relationships: `$this->whenLoaded('posts')`

#### 35.3 Go Implementation
- Resource interface:
  ```go
  type Resource interface {
      ToJSON() map[string]interface{}
      With(key string, value interface{}) Resource
      When(condition bool, callback func() map[string]interface{}) map[string]interface{}
      WhenLoaded(relationship string, callback func() interface{}) interface{}
  }
  ```
- Base resource struct:
  ```go
  type BaseResource struct {
      Model interface{}
  }

  func (r *BaseResource) ToJSON() map[string]interface{} {
      // Default reflection-based serialization
  }
  ```
- Usage:
  ```go
  func (c *UserController) Show(w http.ResponseWriter, r *http.Request) {
      user := c.users.Find(1)
      response.Json(w, resources.NewUser(user).ToJSON(), 200)
  }
  ```
- Collection with pagination:
  ```go
  users, _ := db.Table("users").Paginate(15)
  response.Json(w, resources.NewUserCollection(users).ToJSON(), 200)
  ```
- Conditional attributes:
  ```go
  func (r *UserResource) ToJSON() map[string]interface{} {
      return map[string]interface{}{
          "id":   r.Model.ID,
          "name": r.Model.Name,
          "admin_notes": r.When(auth.User().IsAdmin(), func() interface{} {
              return r.Model.AdminNotes
          }),
          "posts": r.WhenLoaded("posts", func() interface{} {
              return resources.NewPostCollection(r.Model.Posts)
          }),
      }
  }
  ```

---

## 36. COLLECTIONS & LAZY COLLECTIONS
*Added by Kimi K2.6 — Framework: GoW*

#### 36.1 Features Required
- Fluent collection API (`map`, `filter`, `reduce`, `each`, `first`, `last`, `pluck`, `sum`, `avg`, `chunk`, etc.)
- Immutable vs mutable operation modes
- Higher-order messages
- Collection macros
- Lazy collections for memory-efficient large datasets (generator-based)
- Eloquent Collection extending base Collection with model-specific methods
- Partition, `groupBy`, `keyBy`, `zip`, `combine`, `flatten`, `dot`, `undot`
- Debugging (`dd`, `dump`)
- Serialization (`toJSON`, `toArray`)

#### 36.2 How It Works in Laravel
- `collect($items)` creates a Collection
- Method chaining: `collect($users)->where('active', true)->sortBy('name')`
- `LazyCollection::make(function () { yield $item; })` for generators
- Eloquent relationships return `Eloquent\Collection`

#### 36.3 Go Implementation
- Collection using generics:
  ```go
  type Collection[T any] struct {
      items []T
  }

  func Collect[T any](items []T) *Collection[T] {
      return &Collection[T]{items: items}
  }

  func (c *Collection[T]) Filter(fn func(T) bool) *Collection[T] { /* ... */ }
  func (c *Collection[T]) Map(fn func(T) T) *Collection[T] { /* ... */ }
  func (c *Collection[T]) First() T { /* ... */ }
  func (c *Collection[T]) Pluck(key string) *Collection[any] { /* reflection-based */ }
  func (c *Collection[T]) Reduce(fn func(any, T) any, initial any) any { /* ... */ }
  func (c *Collection[T]) Chunk(size int) *Collection[[]T] { /* ... */ }
  func (c *Collection[T]) GroupBy(key string) map[string]*Collection[T] { /* ... */ }
  ```
- Lazy collection:
  ```go
  type LazyCollection[T any] struct {
      generator func(yield func(T) bool)
  }

  func Lazy[T any](gen func(yield func(T) bool)) *LazyCollection[T] {
      return &LazyCollection[T]{generator: gen}
  }

  func (lc *LazyCollection[T]) Filter(fn func(T) bool) *LazyCollection[T] {
      return Lazy[T](func(yield func(T) bool) {
          lc.generator(func(item T) bool {
              if fn(item) { return yield(item) }
              return true
          })
      })
  }
  ```
- Eloquent Collection:
  ```go
  type EloquentCollection[T Model] struct {
      Collection[T]
  }
  func (ec *EloquentCollection[T]) Load(relations ...string) error { /* eager load */ }
  func (ec *EloquentCollection[T]) Fresh() *EloquentCollection[T] { /* reload from DB */ }
  ```
- Macros:
  ```go
  collection.Macro("toUpper", func(c *Collection[string]) *Collection[string] {
      return c.Map(func(s string) string { return strings.ToUpper(s) })
  })
  ```

---

## 37. HELPERS, STRINGABLE & FLUENT UTILITIES
*Added by Kimi K2.6 — Framework: GoW*

#### 37.1 Features Required
- **String helpers**: `Camel`, `Snake`, `Kebab`, `Studly`, `Slug`, `Limit`, `Random`, `Uuid`, `Ulid`, `Headline`, `Squish`, `Transliterate`, `Password`, `Wrap`, `Reverse`, `InlineMarkdown`, `Markdown`
- **Array helpers**: `First`, `Last`, `Random`, `Map`, `MapWithKeys`, `Prepend`, `Pull`, `Set` (dot notation), `Sort`, `Undot`, `Where`, `Wrap`, `Query` (http query string)
- **Number helpers**: `Currency`, `Format`, `Percentage`, `Spell`, `Ordinal`, `ForHumans`, `FileSize`, `Abbreviate`
- **Date helpers**: `Now`, `Parse`, `Create`, fluent manipulation (`AddDays`, `SubMonths`, `StartOfDay`, `EndOfMonth`, `IsWeekend`, `DiffForHumans`)
- **Stringable fluent interface**: `Str::of("hello")->upper()->append(" WORLD")`
- **Benchmark helper**: measure execution time, `DD()` for debug-die
- **Sleep helper**: expressive pause (`Sleep.For(5).Minutes()`)
- **General helpers**: `value()`, `tap()`, `with()`, `blank()`, `filled()`, `optional()`, `retry()`, `throw_if()`, `throw_unless()`

#### 37.2 How It Works in Laravel
- Static methods on `Illuminate\Support\Str`, `Arr`, `Number`, `Date`
- `Str::of()` returns `Stringable` for fluent chaining
- `Number::currency(1000, 'USD')` → `"$1,000.00"`
- `Date::now()->addDays(5)` fluent Carbon-like interface
- Global helper functions available

#### 37.3 Go Implementation
- String package:
  ```go
  package str

  func Camel(value string) string { /* ... */ }
  func Snake(value string) string { /* ... */ }
  func Slug(value string, separator string) string { /* ... */ }
  func Limit(value string, limit int, end string) string { /* ... */ }
  func Random(length int) string { /* ... */ }
  func Uuid() string { /* ... */ }
  func Ulid() string { /* ... */ }
  func Headline(value string) string { /* ... */ }
  func Squish(value string) string { /* ... */ }
  func Transliterate(value string) string { /* ... */ }
  func Password(length int, letters, numbers, symbols, spaces bool) string { /* ... */ }

  type Stringable struct { value string }
  func Of(value string) *Stringable { return &Stringable{value: value} }
  func (s *Stringable) Upper() *Stringable { s.value = strings.ToUpper(s.value); return s }
  func (s *Stringable) Append(v string) *Stringable { s.value += v; return s }
  func (s *Stringable) Value() string { return s.value }
  ```
- Array package:
  ```go
  package arr

  func First[K comparable, V any](m map[K]V, cb func(V, K) bool) (V, bool) { /* ... */ }
  func Last[K comparable, V any](m map[K]V) (V, bool) { /* ... */ }
  func Map[T, U any](items []T, fn func(T) U) []U { /* ... */ }
  func MapWithKeys[T any, K comparable, V any](items []T, fn func(T) (K, V)) map[K]V { /* ... */ }
  func Pull[K comparable, V any](m map[K]V, key K, def V) (V, map[K]V) { /* ... */ }
  func Set[V any](m map[string]V, key string, value V) map[string]V { /* dot notation */ }
  func Sort[T any](items []T, fn func(T, T) bool) []T { /* ... */ }
  func Undot[V any](flat map[string]V) map[string]interface{} { /* ... */ }
  func Where[T any](items []T, fn func(T) bool) []T { /* ... */ }
  func Wrap[T any](value T) []T { return []T{value} }
  func Query(m map[string]string) string { /* url.Values.Encode() */ }
  ```
- Number package:
  ```go
  package number

  func Currency(value float64, currency, locale string) string { /* ... */ }
  func Format(value float64, precision int) string { /* ... */ }
  func Percentage(value float64, precision int) string { /* ... */ }
  func Spell(value int, locale string) string { /* ... */ }
  func Ordinal(value int, locale string) string { /* ... */ }
  func ForHumans(value int, units string) string { /* ... */ }
  func FileSize(bytes int64, precision int) string { /* ... */ }
  func Abbreviate(value int) string { /* ... */ }
  ```
- Date package:
  ```go
  package date

  type Date struct { time.Time }
  func Now() *Date { return &Date{Time: time.Now()} }
  func Parse(layout, value string) (*Date, error) { /* ... */ }
  func Create(year, month, day, hour, min, sec int) *Date { /* ... */ }
  func (d *Date) AddDays(days int) *Date { d.Time = d.Time.AddDate(0, 0, days); return d }
  func (d *Date) SubMonths(months int) *Date { d.Time = d.Time.AddDate(0, -months, 0); return d }
  func (d *Date) StartOfDay() *Date { /* truncate time */ return d }
  func (d *Date) EndOfMonth() *Date { /* ... */ return d }
  func (d *Date) IsWeekend() bool { /* ... */ }
  func (d *Date) DiffForHumans(other *Date) string { /* ... */ }
  ```
- Benchmark:
  ```go
  package benchmark
  func Measure(fn func()) time.Duration { start := time.Now(); fn(); return time.Since(start) }
  func DD(fn func()) { dur := Measure(fn); fmt.Printf("%v\\n", dur); os.Exit(0) }
  ```
- Sleep:
  ```go
  package sleep
  type Sleeper struct{ d time.Duration }
  func For(d time.Duration) *Sleeper { return &Sleeper{d: d} }
  func (s *Sleeper) Seconds() { time.Sleep(s.d * time.Second) }
  func (s *Sleeper) Milliseconds() { time.Sleep(s.d * time.Millisecond) }
  func Until(t time.Time) { time.Sleep(time.Until(t)) }
  ```

---

## 38. SIGNED URLS
*Added by Kimi K2.6 — Framework: GoW*

#### 38.1 Features Required
- Generate signed URLs with HMAC signature
- Generate temporary signed URLs with expiration timestamp
- Validate signature and expiration
- Custom signature algorithms (HMAC-SHA256)
- Exclude specific query parameters from signature
- Signed route middleware (auto-validate)
- Relative vs absolute signed URLs
- Key rotation considerations

#### 38.2 How It Works in Laravel
- `URL::signedRoute('unsubscribe', ['user' => 1])`
- `URL::temporarySignedRoute('unsubscribe', now()->addMinutes(30), ['user' => 1])`
- `URL::hasValidSignature($request)`
- `SignedUrlMiddleware` validates automatically

#### 38.3 Go Implementation
- URL signer:
  ```go
  type URLSigner struct {
      key []byte
  }

  func (s *URLSigner) SignedRoute(name string, params map[string]string, absolute bool) string {
      url := router.Generate(name, params)
      sig := s.Sign(url)
      return url + "?signature=" + sig
  }

  func (s *URLSigner) TemporarySignedRoute(name string, expiration time.Time, params map[string]string, absolute bool) string {
      url := router.Generate(name, params)
      expires := strconv.FormatInt(expiration.Unix(), 10)
      sig := s.Sign(url + "?expires=" + expires)
      return url + "?expires=" + expires + "&signature=" + sig
  }

  func (s *URLSigner) Sign(url string) string {
      mac := hmac.New(sha256.New, s.key)
      mac.Write([]byte(url))
      return base64.URLEncoding.EncodeToString(mac.Sum(nil))
  }

  func (s *URLSigner) HasValidSignature(r *http.Request) bool {
      // Verify HMAC and expiration timestamp if present
  }
  ```
- Middleware:
  ```go
  func SignedURLMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          signer := app.Make("url.signer").(*URLSigner)
          if !signer.HasValidSignature(r) {
              abort(403)
              return
          }
          next.ServeHTTP(w, r)
      })
  }
  ```

---

## 39. ROUTE MODEL BINDING
*Added by Kimi K2.6 — Framework: GoW*

#### 39.1 Features Required
- Implicit model binding (auto-resolve from route parameter by ID)
- Explicit model binding (custom resolution logic)
- Custom binding fields (e.g., bind by `slug` instead of `id`)
- Fallback behavior (404 if not found)
- Scoping for nested resources (child must belong to parent)
- Soft delete inclusion (`withTrashed`)
- Multiple bindings per route
- Resource route auto-registration (`Route::resource`)
- API resource routes (exclude create/edit)
- Shallow nesting

#### 39.2 How It Works in Laravel
- `Route::get('/posts/{post}', ...)` auto-resolves Post model
- `Route::bind('post', fn($value) => Post::where('slug', $value)->firstOrFail())`
- `Route::model('user', User::class)` custom model
- `Route::scopeBindings()` for nested resources
- `Route::get('/posts/{post:slug}', ...)` custom field
- `Route::resource('photos', PhotoController::class)` generates 7 routes
- `Route::apiResource('photos', PhotoController::class)` generates 5 routes

#### 39.3 Go Implementation
- Binder interface:
  ```go
  type ModelBinder interface {
      Bind(value string, model interface{}) error
      BindWithField(field string, value string, model interface{}) error
  }
  ```
- Implicit binding:
  ```go
  router.Get("/users/{user}", userController.Show).
      Bind("user", (*models.User)(nil)) // auto-resolve by ID
  ```
- Custom binding:
  ```go
  router.Bind("post", func(value string) (interface{}, error) {
      return models.PostQuery().Where("slug", value).First()
  })
  ```
- Scoping:
  ```go
  router.Group("/users/{user}/posts", func(r *Router) {
      r.ScopeBindings()
      r.Get("/{post}", postController.Show)
  })
  ```
- Custom field:
  ```go
  router.Get("/posts/{post:slug}", postController.Show)
  ```
- Resource routes:
  ```go
  router.Resource("users", userController) // index, create, store, show, edit, update, destroy
  router.ApiResource("users", userController) // index, store, show, update, destroy
  router.Resource("users.posts", postController) // nested
  router.ShallowResource("users.posts", postController) // shallow nesting
  ```
- Middleware:
  ```go
  func BindModelsMiddleware(next http.Handler) http.Handler {
      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          // Resolve route parameters to models, attach to request context
          next.ServeHTTP(w, r)
      })
  }
  ```

---

## 40. PROCESS / CLI COMMAND EXECUTION
*Added by Kimi K2.6 — Framework: GoW*

#### 40.1 Features Required
- Run external CLI commands with fluent API
- Process pools / concurrency
- Input/output piping
- Timeout handling
- Environment variable injection
- Working directory specification
- Fake / testing mode (assert commands were run)
- Retrieve stdout, stderr, exit code
- Asynchronous execution
- Signal handling (SIGTERM, SIGKILL)

#### 40.2 How It Works in Laravel
- `Process::run('ls -la')`
- `Process::timeout(60)->run('long-running-command')`
- `Process::path('/custom/dir')->env(['KEY' => 'value'])->run('command')`
- `Process::concurrently([fn() => Process::run('cmd1'), fn() => Process::run('cmd2')])`
- `Process::fake()` for testing

#### 40.3 Go Implementation
- Process runner:
  ```go
  type Process struct {
      command string
      args    []string
      timeout time.Duration
      dir     string
      env     []string
      input   string
  }

  func NewProcess(command string, args ...string) *Process {
      return &Process{command: command, args: args}
  }

  func (p *Process) Timeout(d time.Duration) *Process { p.timeout = d; return p }
  func (p *Process) Path(dir string) *Process { p.dir = dir; return p }
  func (p *Process) Env(env map[string]string) *Process {
      for k, v := range env { p.env = append(p.env, k+"="+v) }
      return p
  }
  func (p *Process) Input(input string) *Process { p.input = input; return p }

  func (p *Process) Run() (*ProcessResult, error) {
      ctx := context.Background()
      if p.timeout > 0 {
          var cancel context.CancelFunc
          ctx, cancel = context.WithTimeout(ctx, p.timeout)
          defer cancel()
      }
      cmd := exec.CommandContext(ctx, p.command, p.args...)
      cmd.Dir = p.dir
      cmd.Env = append(os.Environ(), p.env...)
      if p.input != "" { cmd.Stdin = strings.NewReader(p.input) }
      output, err := cmd.CombinedOutput()
      return &ProcessResult{
          Output:   string(output),
          ExitCode: cmd.ProcessState.ExitCode(),
          Error:    err,
      }, err
  }

  type ProcessResult struct {
      Output   string
      ExitCode int
      Error    error
  }

  func (r *ProcessResult) Successful() bool { return r.ExitCode == 0 }
  func (r *ProcessResult) Failed() bool     { return r.ExitCode != 0 }
  func (r *ProcessResult) Throw() error {
      if r.Failed() {
          return fmt.Errorf("process failed [%d]: %s", r.ExitCode, r.Output)
      }
      return nil
  }
  ```
- Fake for testing:
  ```go
  type FakeProcess struct {
      commands []string
  }
  func (f *FakeProcess) Run(command string, args ...string) (*ProcessResult, error) {
      f.commands = append(f.commands, command+" "+strings.Join(args, " "))
      return &ProcessResult{Output: "", ExitCode: 0}, nil
  }
  func (f *FakeProcess) AssertRan(command string) {
      // assertion logic
  }
  ```

---

## 41. SCHEMA DUMPING
*Added by Kimi K2.6 — Framework: GoW*

#### 41.1 Features Required
- Dump current database schema to SQL file
- Load schema from SQL file for fast restore
- Schema pruning (remove old migration files after dump)
- Multi-driver SQL dump support
- Schema as baseline for fresh installs
- `schema:dump` and `schema:load` commands
- `--prune` flag to delete migration files after dumping

#### 41.2 How It Works in Laravel
- `php artisan schema:dump` dumps to `database/schema/{connection}-schema.dump`
- `php artisan schema:dump --prune` dumps and deletes migrations
- Fresh installs load schema instead of running all migrations
- Supports MySQL, PostgreSQL, SQLite

#### 41.3 Go Implementation
- Schema dumper interface:
  ```go
  type SchemaDumper interface {
      Dump() (string, error)
      Load(sql string) error
      Prune(migrationsPath string) error
  }
  ```
- Driver-specific implementations:
  ```go
  type MySQLDumper struct{ db *sql.DB }
  type PostgresDumper struct{ db *sql.DB }
  type SQLiteDumper struct{ db *sql.DB }
  ```
- CLI commands:
  ```go
  // artisan schema:dump [--prune]
  func SchemaDumpCommand(ctx *Context) error {
      dumper := app.Make("schema.dumper").(SchemaDumper)
      schema, err := dumper.Dump()
      if err != nil { return err }
      os.WriteFile("database/schema/gow-schema.sql", []byte(schema), 0644)
      if ctx.Option("prune") {
          dumper.Prune("database/migrations/")
      }
      return nil
  }
  ```

---

## 42. MODEL PRUNING
*Added by Kimi K2.6 — Framework: GoW*

#### 42.1 Features Required
- Prune models matching custom criteria (age, status, etc.)
- Scheduled pruning command
- Prunable interface on models
- Configurable prune frequency
- Chunked pruning for large datasets
- Events before/after pruning
- `model:prune` command with `--model` option

#### 42.2 How It Works in Laravel
- Models implement `Prunable` with `prunable()` method returning query builder
- `php artisan model:prune` runs pruning
- Scheduled in `App\Console\Kernel::schedule()`
- Default trait prunes soft-deleted records older than N days

#### 42.3 Go Implementation
- Prunable interface:
  ```go
  type Prunable interface {
      Prunable() *goquent.Builder
  }
  ```
- Model implementation:
  ```go
  type Session struct {
      BaseModel
      // ...
  }
  func (s *Session) Prunable() *goquent.Builder {
      return s.Query().Where("last_activity", "<", time.Now().AddDate(0, 0, -30))
  }
  ```
- Pruner service:
  ```go
  type Pruner struct {
      models []Prunable
  }
  func (p *Pruner) Prune() error {
      for _, model := range p.models {
          query := model.Prunable()
          for {
              ids, err := query.Limit(1000).Pluck("id")
              if len(ids) == 0 { break }
              query.GetConnection().Table("sessions").WhereIn("id", ids).Delete()
          }
      }
      return nil
  }
  ```
- Command: `artisan model:prune [--model=Session]`

---

## 43. HEALTH CHECKS
*Added by Kimi K2.6 — Framework: GoW*

#### 43.1 Features Required
- Health check endpoint (`/up` or configurable)
- Multiple indicators: database, cache, redis, disk space, external services
- Custom health checks registration
- Results: passing, warning, failing
- JSON response with component status details
- Exclusion from logs/throttling
- Scheduling integration for periodic checks
- Notification on failure

#### 43.2 How It Works in Laravel
- Laravel 11+ `Route::health()` or `Health` facade
- Checks: DatabaseCheck, CacheCheck, RedisCheck, StorageCheck
- Returns 200 if healthy, 503 if unhealthy
- JSON with detailed component status

#### 43.3 Go Implementation
- Health check interface:
  ```go
  type HealthCheck interface {
      Name() string
      Check() (HealthStatus, error)
  }
  type HealthStatus struct {
      Status   string // "ok", "warning", "error"
      Message  string
      Metadata map[string]interface{}
  }
  ```
- Built-in checks:
  ```go
  type DatabaseCheck struct{ db *sql.DB }
  func (c *DatabaseCheck) Check() (HealthStatus, error) {
      err := c.db.Ping()
      if err != nil { return HealthStatus{Status: "error", Message: err.Error()}, err }
      return HealthStatus{Status: "ok"}, nil
  }
  type CacheCheck struct{ cache Cache }
  type RedisCheck struct{ redis *redis.Client }
  type DiskSpaceCheck struct{ path string; threshold float64 }
  ```
- Handler:
  ```go
  func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
      checks := app.Health().RunChecks()
      status, code := "ok", 200
      for _, check := range checks {
          if check.Status == "error" { status = "error"; code = 503 }
          else if check.Status == "warning" && status == "ok" { status = "warning" }
      }
      response.Json(w, map[string]interface{}{"status": status, "checks": checks}, code)
  }
  ```
- Route: `router.Get("/up", healthCheckHandler).Name("health")`

---

## 44. FEATURE FLAGS
*Added by Kimi K2.6 — Framework: GoW*

#### 44.1 Features Required
- Define feature flags in code or config
- Check features for users/teams/globally
- Multiple drivers: array, database, redis
- Rich feature values (not just boolean)
- Scope resolution (user-scoped, global)
- A/B testing support (percentage rollout)
- Feature caching
- Blade directives (`@feature`, `@endfeature`)
- Middleware for feature-gated routes
- Artisan commands to activate/deactivate

#### 44.2 How It Works in Laravel
- `Feature::active('new-api')`
- `Feature::activate('new-api', 'user-123')`
- `Feature::for($user)->active('beta-feature')`
- Database stores features in `features` table
- Blade: `@feature('new-api') ... @endfeature`

#### 44.3 Go Implementation
- Feature manager:
  ```go
  type FeatureManager struct {
      driver   FeatureDriver
      resolver ScopeResolver
  }
  type FeatureDriver interface {
      Get(feature, scope string) (interface{}, error)
      Set(feature, scope string, value interface{}) error
      Delete(feature, scope string) error
  }
  type ScopeResolver interface {
      Resolve(r *http.Request) string
  }
  ```
- Usage:
  ```go
  if app.Features().Active("new-dashboard") { /* show new dashboard */ }
  if app.Features().For(user).Active("beta-feature") { /* beta access */ }
  value := app.Features().Value("theme", user) // "dark", "light", etc.
  ```
- Blade directives:
  ```go
  // @feature('new-api') ... @elsefeature ... @endfeature
  ```
- Middleware:
  ```go
  router.Get("/beta", betaController.Index).
      Middleware(middleware.Feature("beta-access"))
  ```
- Database:
  ```sql
  CREATE TABLE features (
      id BIGINT PRIMARY KEY, name VARCHAR(255), scope VARCHAR(255),
      value JSON, created_at TIMESTAMP, updated_at TIMESTAMP
  );
  ```

---

## 45. ADVANCED BLADE / TEMPLATE FEATURES
*Added by Kimi K2.6 — Framework: GoW*

#### 45.1 Features Required
- Loop variables (`$loop->index`, `$loop->first`, `$loop->last`, `$loop->even`, `$loop->odd`, `$loop->depth`, `$loop->parent`)
- `@once` directive (render once per rendering cycle)
- `@fragment` / `@endfragment` (partial template rendering for Turbo/Hotwire)
- `@aware` directive (pass data to child components)
- `@inject` directive (service injection in templates)
- `@env` and `@production` directives
- `@checked`, `@selected`, `@readonly`, `@disabled` form helpers
- `@class` directive (conditional classes array)
- `@style` directive (conditional styles)
- `@pushIf` and `@prepend` directives
- `@method` support for all HTTP verbs

#### 45.2 How It Works in Laravel
- `$loop` available in `@foreach`
- `@once` ensures content renders once even in loops
- `@fragment('content')` allows partial rendering
- `@aware(['color' => 'gray'])` in components
- `@inject('metrics', App\Services\MetricsService)`
- `@class(['p-4', 'font-bold' => $isActive])`

#### 45.3 Go Implementation
- Loop context:
  ```go
  type LoopContext struct {
      Index     int // 0-based
      Iteration int // 1-based
      Remaining int
      Count     int
      First     bool
      Last      bool
      Even      bool
      Odd       bool
      Depth     int
      Parent    *LoopContext
  }
  ```
  In template compilation, inject `.Loop` variable in `range` loops.
- Directive compilation:
  ```go
  // @once → {{ if $.OnceStore.Add "unique-id" }}...{{ end }}
  // @fragment("name") → {{ define "fragment-name" }}...{{ end }}
  // @aware(["color" => "gray"]) → set defaults in component scope
  // @inject("metrics", "services.MetricsService") → resolve from container
  // @env("local") → {{ if eq .Env "local" }}...{{ end }}
  // @production → {{ if eq .Env "production" }}...{{ end }}
  // @class(map) → {{ template "class-helper" .ClassMap }}
  // @checked(condition) → {{ if .Condition }}checked{{ end }}
  // @selected(condition) → {{ if .Condition }}selected{{ end }}
  // @readonly(condition) → {{ if .Condition }}readonly{{ end }}
  // @disabled(condition) → {{ if .Condition }}disabled{{ end }}
  // @style(map) → {{ template "style-helper" .StyleMap }}
  // @pushIf(condition, stack) → conditional push
  ```

---

## 46. REQUEST & RESPONSE MACROS / OLD INPUT / FLASHING
*Added by Kimi K2.6 — Framework: GoW*

#### 46.1 Features Required
- Request macros (add custom methods to request)
- Response macros (add custom methods to response)
- Old input retrieval (`old('email')`, `old('email', 'default')`)
- Flash input to session (`flash()`, `flashOnly()`, `flashExcept()`)
- Reflash session data (`reflash()`, `keep()`)
- Request boolean/integer/string/date input casting
- Request `has()`, `filled()`, `missing()`, `whenHas()`, `whenFilled()`
- Request `merge()`, `replace()`
- Request `is()`, `routeIs()`, `fullUrlIs()`
- Request `cookie()`, `hasCookie()`
- Response `withHeaders()`, `withCookie()`, `withException()`
- Response `jsonp()`, `streamDownload()`, `streamedJson()`
- Input sanitization middleware (`TrimStrings`, `ConvertEmptyStringsToNull`)
- Form Request `PrepareForValidation()` and `PassedValidation()` hooks

#### 46.2 How It Works in Laravel
- `Request::macro('foo', fn() => ...)`
- `Response::macro('foo', fn() => ...)`
- `$request->old('key')` gets flashed input
- `$request->flash()`, `$request->flashOnly(['email'])`
- `$request->boolean('active')`, `$request->integer('age')`, `$request->date('birthday')`
- `$request->has('name')`, `$request->filled('name')`
- `$request->merge(['extra' => 'data'])`

#### 46.3 Go Implementation
- Request wrapper:
  ```go
  type Request struct {
      *http.Request
      session Session
      macros  map[string]func(...interface{}) interface{}
  }
  func (r *Request) Macro(name string, fn func(...interface{}) interface{}) {
      r.macros[name] = fn
  }
  func (r *Request) Old(key string, defaultValue ...string) string {
      if r.session == nil { return "" }
      old := r.session.Get("_old_input", map[string]string{}).(map[string]string)
      if val, ok := old[key]; ok { return val }
      if len(defaultValue) > 0 { return defaultValue[0] }
      return ""
  }
  func (r *Request) Flash() { r.session.FlashInput(r.FormValue) }
  func (r *Request) FlashOnly(keys ...string) {
      input := make(map[string]string)
      for _, key := range keys { input[key] = r.FormValue(key) }
      r.session.FlashInput(input)
  }
  func (r *Request) Boolean(key string) bool {
      val := r.FormValue(key)
      return val == "1" || val == "true" || val == "on" || val == "yes"
  }
  func (r *Request) Integer(key string) int {
      val, _ := strconv.Atoi(r.FormValue(key))
      return val
  }
  func (r *Request) Date(key, layout string) (time.Time, error) {
      return time.Parse(layout, r.FormValue(key))
  }
  func (r *Request) Has(keys ...string) bool {
      for _, key := range keys { if r.FormValue(key) == "" { return false } }
      return true
  }
  func (r *Request) Filled(keys ...string) bool {
      for _, key := range keys { if strings.TrimSpace(r.FormValue(key)) == "" { return false } }
      return true
  }
  func (r *Request) Missing(keys ...string) bool { return !r.Has(keys...) }
  func (r *Request) Merge(data map[string]string) { /* add to request context */ }
  func (r *Request) Is(patterns ...string) bool {
      for _, p := range patterns {
          if matched, _ := filepath.Match(p, r.URL.Path); matched { return true }
      }
      return false
  }
  func (r *Request) RouteIs(names ...string) bool {
      routeName := r.Context().Value("route.name").(string)
      for _, n := range names { if n == routeName { return true } }
      return false
  }
  ```
- Response macros:
  ```go
  type Response struct {
      http.ResponseWriter
      macros map[string]func(...interface{}) interface{}
  }
  func (r *Response) WithCookie(c *http.Cookie) *Response { http.SetCookie(r, c); return r }
  func (r *Response) WithHeaders(h map[string]string) *Response {
      for k, v := range h { r.Header().Set(k, v) }
      return r
  }
  func (r *Response) Jsonp(callback string, data interface{}) { /* wrap JSON in callback */ }
  ```
- Form Request hooks:
  ```go
  type FormRequest interface {
      Authorize() bool
      Rules() map[string]string
      Messages() map[string]string
      PrepareForValidation() error
      PassedValidation() error
      Validate() error
  }
  ```

---

## 47. QUERY LOGGING & DATABASE PROFILING
*Added by Kimi K2.6 — Framework: GoW*

#### 47.1 Features Required
- Query listener / logging (`DB::listen` equivalent)
- Query duration tracking
- Query count per request
- Slow query detection and logging
- Query `explain` support
- Raw expression support (`DB::raw`)
- Database transaction monitoring
- Query log retrieval and clearing
- Duplicate query detection (N+1 warning)
- Query binding logging

#### 47.2 How It Works in Laravel
- `DB::listen(function ($query) { Log::info($query->sql); })`
- `DB::getQueryLog()` returns array of queries
- `DB::enableQueryLog()` / `DB::disableQueryLog()`
- `DB::raw('count(*) as total')` for raw expressions
- `explain()` on query builder

#### 47.3 Go Implementation
- Query logger:
  ```go
  type QueryLog struct {
      SQL        string
      Bindings   []interface{}
      Time       float64 // milliseconds
      Connection string
  }
  type QueryLogger struct {
      logs          []QueryLog
      listeners     []func(QueryLog)
      enabled       bool
      slowThreshold float64
  }
  func (ql *QueryLogger) Listen(callback func(QueryLog)) {
      ql.listeners = append(ql.listeners, callback)
  }
  func (ql *QueryLogger) Log(sql string, bindings []interface{}, duration time.Duration, connection string) {
      if !ql.enabled { return }
      log := QueryLog{SQL: sql, Bindings: bindings, Time: float64(duration.Microseconds()) / 1000.0, Connection: connection}
      ql.logs = append(ql.logs, log)
      for _, l := range ql.listeners { l(log) }
      if ql.slowThreshold > 0 && log.Time > ql.slowThreshold {
          logger.Warn("Slow query detected", zap.Any("query", log))
      }
  }
  func (ql *QueryLogger) GetQueryLog() []QueryLog { return ql.logs }
  func (ql *QueryLogger) FlushQueryLog() { ql.logs = nil }
  ```
- Raw expressions:
  ```go
  type RawExpression struct {
      Value    string
      Bindings []interface{}
  }
  func Raw(value string, bindings ...interface{}) *RawExpression {
      return &RawExpression{Value: value, Bindings: bindings}
  }
  ```
- Explain:
  ```go
  func (b *Builder) Explain() ([]map[string]interface{}, error) {
      return b.GetConnection().Query("EXPLAIN "+b.ToSQL(), b.GetBindings())
  }
  ```
- N+1 detector:
  ```go
  type NPlusOneDetector struct {
      queryCounts map[string]int
  }
  func (d *NPlusOneDetector) Check(log QueryLog) {
      // Detect repeated similar queries
  }
  ```

---

## 48. PACKAGE AUTO-DISCOVERY & PUBLISHING
*Added by Kimi K2.6 — Framework: GoW*

#### 48.1 Features Required
- Auto-discovery of service providers from dependencies
- Publishing configs, views, migrations, assets, routes, translations
- Vendor publish command with tags
- Package manifest generation
- Overwrite protection when publishing
- Selective publishing by provider or tag
- Package asset compilation hooks
- Package route, view, translation namespace registration
- Package config merging with app config

#### 48.2 How It Works in Laravel
- `vendor:publish` copies package resources to app
- `--provider` and `--tag` options for selective publishing
- Auto-discovery reads `composer.json` extra
- `PackageServiceProvider` has `publishes()`, `mergeConfigFrom()`, `loadViewsFrom()`, `loadMigrationsFrom()`, `loadTranslationsFrom()`, `loadRoutesFrom()`

#### 48.3 Go Implementation
- Package manifest:
  ```go
  type PackageManifest struct {
      packages     []Package
      manifestPath string
  }
  type Package struct {
      Name      string
      Providers []string
      Aliases   map[string]string
  }
  func (m *PackageManifest) Build() error {
      // Scan go.mod or gow.packages config for packages with provider files
      // Cache to bootstrap/cache/packages.gob
  }
  ```
- Publishing:
  ```go
  type Publisher struct {
      files map[string]string
      tags  map[string][]string
  }
  func (p *Publisher) Publishes(paths map[string]string, tag string) {
      for src, dst := range paths {
          p.files[src] = dst
          p.tags[tag] = append(p.tags[tag], src)
      }
  }
  func (p *Publisher) Publish(tag string, force bool) error {
      // Copy files from package to app directories
  }
  ```
- Command: `artisan vendor:publish [--provider=...] [--tag=...] [--force]`

---

## 49. TESTING UTILITIES
*Added by Kimi K2.6 — Framework: GoW*

#### 49.1 Features Required
- Database assertions: `assertDatabaseHas`, `assertDatabaseMissing`, `assertSoftDeleted`, `assertNotSoftDeleted`, `assertDatabaseCount`
- Time manipulation: `travel()`, `travelTo()`, `freezeTime()`, `unfreezeTime()`
- Exception testing: `assertThrows`, `expectException`
- View assertions: `assertViewHas`, `assertViewMissing`, `assertViewIs`
- Authentication: `actingAs`, `assertAuthenticated`, `assertGuest`, `assertAuthenticatedAs`
- HTTP assertions: `assertRedirectToRoute`, `assertValid`, `assertInvalid`
- Mail/Queue/Event/Notification assertions
- Storage assertions: `assertExists`, `assertMissing`
- Artisan command testing: `assertSuccessful`, `assertFailed`, `expectsOutput`, `expectsConfirmation`
- Browser testing foundation (Dusk-like): page objects, element interaction, screenshots (optional/future)

#### 49.2 How It Works in Laravel
- `assertDatabaseHas('users', ['email' => 'foo@bar.com'])`
- `travel(5)->days()` modifies `Carbon::now()`
- `actingAs($user)` sets authenticated user
- `Mail::assertSent(OrderShipped::class)`
- `Queue::assertPushed(ProcessPodcast::class)`

#### 49.3 Go Implementation
- Database assertions:
  ```go
  func (tc *TestCase) AssertDatabaseHas(table string, conditions map[string]interface{}) {
      var count int
      query := tc.DB().Table(table)
      for col, val := range conditions { query = query.Where(col, val) }
      query.Count(&count)
      if count == 0 {
          tc.t.Errorf("Failed asserting that %s has row matching %v", table, conditions)
      }
  }
  func (tc *TestCase) AssertDatabaseMissing(table string, conditions map[string]interface{}) { /* inverse */ }
  func (tc *TestCase) AssertSoftDeleted(model interface{}) { /* check deleted_at */ }
  func (tc *TestCase) AssertDatabaseCount(table string, count int) { /* assert row count */ }
  ```
- Time manipulation:
  ```go
  type TestTime struct{}
  func Travel(d time.Duration) *TestTime { /* set global time offset */ return &TestTime{} }
  func (tt *TestTime) Days(n int) *TestTime { /* ... */ return tt }
  func (tt *TestTime) Freeze() { /* freeze */ }
  func (tt *TestTime) Return() { /* restore */ }
  ```
- Authentication:
  ```go
  func (tc *TestCase) ActingAs(user Authenticatable) *TestCase {
      tc.app.Auth().Login(user)
      return tc
  }
  func (tc *TestCase) AssertAuthenticated() {
      if !tc.app.Auth().Check() { tc.t.Error("Expected authenticated user") }
  }
  ```

---

## 50. MODERN LARAVEL PARITY
*Added by Kimi K2.6 — Framework: GoW*

#### 50.1 Features Required
- **Reverb-like WebSocket Server**: First-party WebSocket server
  - Channel authorization, presence channels
  - Horizontal scaling via Redis pub/sub
  - Binary: `gow reverb:start`
- **Pulse-like Monitoring**: Application performance monitoring
  - Slow request tracking, exception frequency
  - Queue throughput, cache hit/miss rates
  - Database slow query collection
  - User activity dashboard data
- **Prompts-like CLI UX**: Rich terminal interactions
  - Interactive select/search, multi-select
  - Confirmations, progress bars with ETA
  - Spinners for long tasks, autocomplete
- **Context-like Propagation**: Cross-layer context (Laravel 11+)
  - Store data in request context that flows to jobs
  - Hidden context (not serialized to logs)
  - Context hydration in queued jobs

#### 50.2 How It Works in Laravel
- Reverb: `php artisan reverb:start`, scales with Redis
- Pulse: records to `pulse_*` tables, dashboard at `/pulse`
- Prompts: `select()`, `multiselect()`, `confirm()`, `search()`, `spin()`
- Context: `Context::addHidden('key', 'value')`, flows to jobs and logs

#### 50.3 Go Implementation
- Reverb server:
  ```go
  type ReverbServer struct {
      hub   *Hub
      redis *redis.Client
  }
  func (s *ReverbServer) Start(addr string) error {
      // WebSocket upgrade, Redis pub/sub for multi-node
  }
  func (s *ReverbServer) Broadcast(channel, event string, payload interface{}) error {
      // Publish to Redis or local hub
  }
  ```
- Pulse recorder:
  ```go
  type Pulse struct {
      ingestChannel chan PulseEntry
      storage       PulseStorage
  }
  type PulseEntry struct {
      Type      string // "request", "exception", "query", "queue", "cache"
      Key       string
      Value     float64
      Timestamp time.Time
  }
  func (p *Pulse) Record(entry PulseEntry) { p.ingestChannel <- entry }
  func (p *Pulse) Values(entryType string, lookback time.Duration) []PulseEntry {
      return p.storage.Read(entryType, lookback)
  }
  ```
- Prompts (CLI):
  ```go
  package prompts
  func Select(label string, options []string, defaultValue string) (string, error) { /* ... */ }
  func MultiSelect(label string, options []string) ([]string, error) { /* ... */ }
  func Confirm(label string, defaultValue bool) (bool, error) { /* ... */ }
  func Search(label string, searchFn func(string) []string) (string, error) { /* ... */ }
  func Spin(label string, fn func() error) error { /* spinner while fn runs */ }
  func Progress(label string, total int) *ProgressBar { /* ... */ }
  ```
- Context:
  ```go
  type Context struct {
      data   map[string]interface{}
      hidden map[string]interface{}
  }
  func FromContext(ctx context.Context) *Context { /* ... */ }
  func (c *Context) Add(key string, value interface{}) { /* ... */ }
  func (c *Context) AddHidden(key string, value interface{}) { /* ... */ }
  func (c *Context) Get(key string) interface{} { /* ... */ }
  func (c *Context) Hydrate(job *Job) { /* copy context to job payload */ }
  ```

---

## APPENDIX A: LARAVEL FEATURE CHECKLIST

This checklist maps Laravel features to this spec for tracking parity:

### Core
- [x] Application / Container
- [x] Service Providers
- [x] Facades (adapted - no magic methods)
- [x] Configuration
- [x] Environment / .env
- [x] Helpers (custom Go helpers)

### HTTP Layer
- [x] Routing
- [x] Middleware
- [x] Controllers
- [x] Requests (Form Requests)
- [x] Responses
- [x] URL Generation
- [x] Session
- [x] Validation
- [x] Views / Blade
- [x] CSRF Protection

### Database
- [x] Query Builder
- [x] Migrations
- [x] Seeding
- [x] Eloquent ORM (Goquent)
- [x] Pagination
- [x] Redis

### Auth
- [x] Authentication
- [x] Authorization (Gates/Policies)
- [x] Email Verification
- [x] Encryption
- [x] Hashing
- [x] Password Reset

### Advanced
- [x] Artisan Console
- [x] Broadcasting
- [x] Cache
- [x] Collections
- [x] Events
- [x] File Storage
- [x] Helpers
- [x] Mail
- [x] Notifications
- [x] Package Development
- [x] Queues
- [x] Rate Limiting
- [x] Strings / Pluralization (adapted)
- [x] Task Scheduling
- [x] Testing

### Security (Modern Framework Requirements)
- [x] CSRF
- [x] XSS Prevention
- [x] SQL Injection Prevention
- [x] Encryption at Rest
- [x] Secure Hashing
- [x] CORS
- [x] Rate Limiting
- [x] Input Sanitization
- [x] Security Headers
- [x] Brute Force Protection
- [x] Session Security
- [x] Mass Assignment Protection
- [x] Dependency Vulnerability Management
- [x] Secure Error Handling
- [x] Secret Management

### Additional Features (Added 2026-05-20 by Kimi K2.6)
- [x] API Resources / Transformers
- [x] Collections & Lazy Collections
- [x] Helpers (Str, Arr, Number, Date, Benchmark)
- [x] Stringable / Fluent Strings
- [x] Signed URLs
- [x] Route Model Binding (Implicit & Explicit)
- [x] Resource Route Registration
- [x] Process / CLI Execution
- [x] Schema Dumping
- [x] Model Pruning
- [x] Health Checks
- [x] Feature Flags (Pennant-like)
- [x] Advanced Blade Directives (@once, @fragment, @aware, @inject, @class, @style)
- [x] Request/Response Macros & Old Input
- [x] Query Logging & Profiling
- [x] Package Auto-Discovery
- [x] Testing Utilities (assertDatabaseHas, Time travel, actingAs)
- [x] WebSocket Server (Reverb-like)
- [x] Application Monitoring (Pulse-like)
- [x] CLI Prompts (Select, Confirm, Search, Spin)
- [x] Context Propagation (Laravel 11+)

---

## APPENDIX B: NAMING CONVENTIONS

| Laravel Concept | Go Equivalent | File/Package |
|----------------|---------------|--------------|
| Eloquent | Goquent | `gow/orm` |
| Blade | Goblade / Templates | `gow/view` |
| Artisan | Artisan / gow-cli | `cmd/artisan` |
| Facade | Accessor / Helper | `gow/facade` (optional) |
| Middleware | Middleware | `app/http/middleware` |
| Controller | Controller | `app/http/controllers` |
| Model | Model | `app/models` |
| Migration | Migration | `database/migrations` |
| Seeder | Seeder | `database/seeders` |
| Factory | Factory | `database/factories` |
| Policy | Policy | `app/policies` |
| Gate | Gate | `gow/auth` |
| Job | Job | `app/jobs` |
| Event | Event | `app/events` |
| Listener | Listener | `app/listeners` |
| Mailable | Mailable | `app/mail` |
| Notification | Notification | `app/notifications` |
| Request | Request / FormRequest | `app/http/requests` |
| Resource | Resource / Transformer | `app/http/resources` |
| Rule | Rule | `app/rules` |
| Provider | Provider | `app/providers` |
| Channel | Channel | `routes/channels` |
| API Resource | Resource / Transformer | `app/http/resources` |
| Collection | Collection | `gow/support/collection` |
| Lazy Collection | LazyCollection | `gow/support/lazy` |
| Stringable | Stringable | `gow/support/str` |
| Process | Process | `gow/support/process` |
| Feature Flag | Feature | `gow/feature` |
| Health Check | HealthCheck | `gow/health` |
| Pulse | Pulse | `gow/pulse` |
| Reverb | Reverb | `gow/reverb` |
| Prompts | Prompts | `gow/prompts` |
| Context | Context | `gow/context` |
| Schema Dump | SchemaDumper | `gow/database/schema` |
| Model Pruner | Pruner | `gow/database/pruning` |

---

## APPENDIX C: KNOWN ISSUES & CONSIDERATIONS

| ID | Issue | Severity | Phase | Notes |
|----|-------|----------|-------|-------|
| KI-001 | Go lacks runtime class loading | MEDIUM | 1 | Use code generation and interfaces instead |
| KI-002 | No generics in interfaces before 1.18 | LOW | 1 | Target 1.24+, use generics freely |
| KI-003 | Error handling is explicit, not exceptions | MEDIUM | 1 | Adapt Laravel's exception pattern to Go errors |
| KI-004 | Template compilation complexity | MEDIUM | 3 | Blade-like syntax requires custom parser |
| KI-005 | ORM eager loading performance | MEDIUM | 2 | Must optimize batch queries |
| KI-006 | Queue job serialization in Go | MEDIUM | 5 | Need type registry for deserialization |
| KI-007 | Session locking in high concurrency | LOW | 3 | Consider read/write locks |
| KI-008 | Rate limiting distributed across nodes | LOW | 4 | Redis required for multi-node |
| KI-009 | Broadcasting WebSocket server | MEDIUM | 5 | May need separate WS server or use Pusher |
| KI-010 | Testing fake implementations maintenance | LOW | 6 | Fakes must stay in sync with real drivers |
| KI-011 | Go template loop variables require custom context injection | MEDIUM | 3 | Blade-like $loop variable needs custom parsing |
| KI-012 | Process fake requires command interception | LOW | 5 | Testing external process calls |
| KI-013 | Feature flag serialization in distributed systems | MEDIUM | 4 | Redis/database consistency |
| KI-014 | Reverb WebSocket server horizontal scaling | MEDIUM | 5 | Redis pub/sub complexity |
| KI-015 | Lazy collection memory efficiency vs channel overhead | LOW | 2 | Benchmark generator patterns |
| KI-016 | Signed URL key rotation without invalidating existing URLs | MEDIUM | 4 | Cryptographic design |

---

*End of Specification*  
*Last Updated: 2026-05-20*  
*Status: Phase 0 - Specification Complete*  
*Framework: GoW (Go Web)*
'''