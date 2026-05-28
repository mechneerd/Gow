# GoW Framework Documentation

> The Laravel experience in pure Go.

---

## Table of Contents

### Getting Started
1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Directory Structure](#directory-structure)

### Core Concepts
4. [Application Lifecycle](#application-lifecycle)
5. [Service Providers](#service-providers)
6. [IoC Container](#ioc-container)

### Routing & Controllers
7. [Routing](#routing)
8. [Route Groups](#route-groups)
9. [Route Model Binding](#route-model-binding)
10. [Controllers](#controllers)
11. [Resource Controllers](#resource-controllers)
12. [Middleware](#middleware)

### Views & Templates
13. [GoBlade Templates](#goblade-templates)
14. [Components & Slots](#components--slots)
15. [Layouts](#layouts)

### Database
16. [Query Builder](#query-builder)
17. [ORM (Goquent)](#orm-goquent)
18. [Migrations](#migrations)
19. [Seeders](#seeders)
20. [Schema Builder](#schema-builder)

### Authentication
21. [Authentication](#authentication)
22. [Authorization & Gates](#authorization--gates)
23. [RBAC (Roles & Permissions)](#rbac)

### Frontend
24. [Livewire (Reactive UI)](#livewire)

### Services
25. [Cache](#cache)
26. [Session](#session)
27. [Queue](#queue)
28. [Events](#events)
29. [Mail](#mail)
30. [Notifications](#notifications)
31. [Broadcasting](#broadcasting)
32. [Storage](#storage)
33. [Logging](#logging)
34. [Encryption](#encryption)
35. [Localization](#localization)

### Testing
36. [Testing](#testing)

### CLI
37. [Artisan CLI](#artisan-cli)
38. [Scheduling](#scheduling)

---

# Getting Started

## Installation

### Quick Install

**macOS / Linux**
```bash
curl -sSfL https://raw.githubusercontent.com/mechneerd/gow/main/install.sh | sh
```

**Windows (PowerShell)**
```powershell
iwr -useb https://raw.githubusercontent.com/mechneerd/gow/main/install.ps1 | iex
```

### Create a New Project

```bash
# Interactive wizard
gow new myblog

# Non-interactive with defaults
gow new myblog --yes

# REST API with PostgreSQL
gow new myapi --api --yes --db=postgres

# Web app with authentication
gow new myapp --auth --yes

# Minimal website
gow new mysite --minimal --yes
```

### Run the Application

```bash
cd myblog
gow serve
# or
go run main.go
```

Visit `http://localhost:8080`.

---

## Configuration

GoW loads configuration from `.env` files and environment variables.

### Environment Variables

```env
APP_NAME=MyApp
APP_ENV=local
APP_PORT=8080
APP_KEY=base64:your-secret-key

DB_CONNECTION=sqlite
DB_DATABASE=database.sqlite
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=

CACHE_DRIVER=memory
SESSION_DRIVER=file
QUEUE_DRIVER=sync
```

### Accessing Configuration

```go
import "github.com/mechneerd/gow/config"

cfg := config.NewRepository(".")

// Get string value
name := cfg.Get("APP_NAME", "MyApp")

// Get typed values
port := cfg.GetInt("APP_PORT", 8080)
debug := cfg.GetBool("APP_DEBUG", false)
```

---

## Directory Structure

```
myapp/
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   └── Auth/
│   │   └── Middleware/
│   ├── Livewire/
│   └── Models/
├── bootstrap/
│   └── app.go              # Application bootstrap
├── config/
│   ├── app.go
│   ├── auth.go
│   └── database.go
├── database/
│   ├── migrations/
│   ├── seeders/
│   └── factory/
├── public/
│   └── js/
├── resources/
│   └── views/
│       ├── welcome.goblade
│       └── components/
├── routes/
│   ├── web.go
│   └── api.go
├── storage/
├── .env
├── go.mod
├── go.sum
└── main.go
```

---

# Core Concepts

## Application Lifecycle

```go
// main.go
package main

import (
    "myapp/routes"
    "myapp/bootstrap"
)

func main() {
    // Register routes
    routes.RegisterWebRoutes()
    routes.RegisterAPIRoutes()

    // Bootstrap and serve
    app := bootstrap.NewApplication()
    app.Serve()
}
```

### Bootstrap

```go
// bootstrap/app.go
package bootstrap

import (
    "fmt"
    "log"
    "net/http"
    "myapp/config"
)

type Application struct {
    Config config.AppConfig
}

func NewApplication() *Application {
    cfg := config.Load()
    return &Application{Config: cfg}
}

func (a *Application) Serve() {
    addr := ":" + a.Config.Port
    fmt.Printf("%s is running on http://localhost%s\n", a.Config.AppName, addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}
```

---

## Service Providers

Service providers are the foundation of the GoW application. They register services into the container.

```go
package myprovider

import "github.com/mechneerd/gow/foundation"

type ServiceProvider struct {
    foundation.BaseServiceProvider
}

func (p *ServiceProvider) Register(app *foundation.Application) {
    // Register services
    app.Singleton((*MyService)(nil), func() *MyService {
        return &MyService{}
    })
}

func (p *ServiceProvider) Boot(app *foundation.Application) {
    // Called after all providers are registered
}
```

### Registering Providers

```go
// bootstrap/app.go
func NewApplication() *foundation.Application {
    app := foundation.NewApplication(".")
    app.RegisterProvider(&config.ServiceProvider{})
    app.RegisterProvider(&logging.ServiceProvider{})
    app.Boot()
    return app
}
```

---

## IoC Container

The container manages dependency injection.

### Basic Usage

```go
import "github.com/mechneerd/gow/container"

c := container.New()

// Bind an interface to implementation
c.Bind((*Logger)(nil), func() Logger {
    return &FileLogger{}
})

// Singleton binding
c.Singleton((*Database)(nil), func() *Database {
    return connectDB()
})

// Register an existing instance
c.Instance((*Config)(nil), configInstance)

// Resolve with generics
logger, err := container.Make[Logger](c)
```

### Contextual Bindings

```go
c.When((*PhotoController)(nil)).
    Needs((*Filesystem)(nil)).
    Give(func() *Filesystem {
        return NewS3Filesystem()
    })
```

### Struct Injection

```go
type UserController struct {
    DB     *orm.DB    `inject:""`
    Cache  cache.Store `inject:""`
}
```

---

# Routing & Controllers

## Routing

```go
import "github.com/mechneerd/gow/routing"

router := routing.NewRouter()

router.Get("/users", listUsers)
router.Post("/users", createUser)
router.Put("/users/{id}", updateUser)
router.Delete("/users/{id}", deleteUser)
```

### Handler Signature

```go
func myHandler(w http.ResponseWriter, r *http.Request) error {
    // Your logic here
    return nil
}
```

### Named Routes

```go
route := router.Get("/dashboard", dashboardHandler)
router.SetName(route, "dashboard")
```

### Reading Route Params

```go
func showUser(w http.ResponseWriter, r *http.Request) error {
    params := r.Context().Value(routing.ParamsKey).(map[string]string)
    id := params["id"]
    // ...
    return nil
}
```

---

## Route Groups

```go
router.Group("/api", func(r *routing.Router) {
    r.Get("/users", apiListUsers)
    r.Post("/users", apiCreateUser)
})
```

### Middleware on Groups

```go
router.Group("/admin", func(r *routing.Router) {
    r.Use(authMiddleware)
    r.Use(roleMiddleware("admin"))

    r.Get("/dashboard", adminDashboard)
})
```

### Middleware Aliases

```go
router.Alias("auth", authMiddleware)
router.Alias("throttle", throttleMiddleware)

router.Group("/api", func(r *routing.Router) {
    r.Middleware("auth", "throttle:60,1")
    r.Get("/profile", getProfile)
})
```

### Middleware Groups

```go
router.GroupMiddleware("web", sessionMiddleware, csrfMiddleware, trimMiddleware)

router.Group("/", func(r *routing.Router) {
    r.Middleware("web")
    r.Get("/", homePage)
})
```

---

## Route Model Binding

```go
// Automatic model resolution by ID
router.Get("/users/{user}", func(w http.ResponseWriter, r *http.Request) error {
    user, ok := routing.Model[models.User](r, "user")
    if !ok {
        return response.Error(w, 404, "User not found")
    }
    return response.JSON(w, 200, user)
})
```

### Custom Bindings

```go
router.Bind("user", func(id string) (any, error) {
    return userRepo.FindByID(id)
})
```

---

## Controllers

```go
package controllers

import (
    "net/http"
    "github.com/mechneerd/gow/http/response"
)

type UserController struct{}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) error {
    users := getAllUsers()
    return response.JSON(w, 200, users)
}

func (c *UserController) Store(w http.ResponseWriter, r *http.Request) error {
    // Create user
    return response.Created(w, user)
}
```

---

## Resource Controllers

```go
router.Resource("posts", &PostController{})
// Creates: index, create, store, show, edit, update, destroy

router.ApiResource("photos", &PhotoController{})
// Creates: index, store, show, update, destroy (no create/edit)
```

### Controller Macros

```go
router.WithMacro("only", func(r *routing.Router, controller *PostController, actions []string) {
    // Register only specified actions
})
```

---

## Middleware

### Creating Middleware

```go
package middleware

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Before
        if isRateLimited(r) {
            http.Error(w, "Too Many Requests", 429)
            return
        }
        next.ServeHTTP(w, r)
        // After
    })
}
```

### Available Middleware

| Middleware | Description |
|-----------|-------------|
| `auth` | Checks authentication status |
| `cors` | Cross-Origin Resource Sharing headers |
| `csrf` | CSRF token validation |
| `encrypt_cookies` | Encrypts/decrypts cookie values |
| `error_handler` | Catches panics, returns JSON/HTML errors |
| `maintenance` | Returns 503 in maintenance mode |
| `recovery` | Recovers from panics, logs stack traces |
| `security_headers` | Sets X-Frame-Options, X-Content-Type-Options, etc. |
| `session` | Loads/saves session data |
| `throttle` | Rate limiting using cache |
| `trim_strings` | Trims whitespace from form values |
| `trusted_proxies` | Sets X-Forwarded-For from trusted proxies |

---

# Views & Templates

## GoBlade Templates

GoW uses GoBlade, a Blade-like template engine for Go.

### Template Syntax

```html
{{-- Echo --}}
{{ $name }}
{!! $rawHtml !!}  {{-- Unescaped --}}

{{-- If/Else --}}
@if($user->active)
    <p>Active</p>
@elseif($user->pending)
    <p>Pending</p>
@else
    <p>Inactive</p>
@endif

{{-- Loops --}}
@foreach($users as $user)
    <p>{{ $user->name }} ({{ $loop->index }})</p>
@endforeach

@for($i = 0; $i < 10; $i++)
    <p>{{ $i }}</p>
@endfor

{{-- Switch --}}
@switch($status)
    @case('active')
        <p>Active</p>
    @break
    @default
        <p>Unknown</p>
@endswitch

{{-- CSRF --}}
<form method="POST">
    @csrf
    <input type="text" name="email">
</form>

{{-- Authorization --}}
@auth
    <p>Welcome, {{ auth()->name }}</p>
@endauth

@guest
    <p>Please log in</p>
@endguest

@can('update', $post)
    <a href="/posts/{{ $post->id }}/edit">Edit</a>
@endcan
```

### Rendering Views

```go
import "github.com/mechneerd/gow/view"

engine := view.NewEngine([]string{"resources/views"}, "storage/framework/views")

html, err := engine.Make("welcome", map[string]any{
    "AppName": "MyApp",
    "Year":    "2026",
})

w.Header().Set("Content-Type", "text/html; charset=utf-8")
w.Write([]byte(html))
```

---

## Components & Slots

```html
{{-- resources/views/components/alert.goblade --}}
<div class="alert alert-{{ $type }}">
    {{ $slot }}
</div>

{{-- Usage --}}
<x-alert type="success">
    <p>Operation completed successfully!</p>
</x-alert>
```

### Component Attributes

```html
{{-- resources/views/components/button.goblade --}}
<button class="btn btn-{{ $variant ?? 'primary' }}" {{ $attributes }}>
    {{ $slot }}
</button>

{{-- Usage --}}
<x-button variant="danger" onclick="confirm()">Delete</x-button>
```

---

## Layouts

### Defining a Layout

```html
{{-- resources/views/layouts/app.goblade --}}
<!DOCTYPE html>
<html>
<head>
    <title>@yield('title', 'MyApp')</title>
</head>
<body>
    @yield('content')
</body>
</html>
```

### Using a Layout

```html
@extends('layouts.app')

@section('title', 'Dashboard')

@section('content')
    <h1>Dashboard</h1>
    <p>Welcome back!</p>
@endsection
```

### Yielding Content

```html
@yield('content')
@yield('content', 'Default content')
```

### Section & Show

```html
@section('sidebar')
    <nav>Sidebar content</nav>
@show

@section('content')
    <p>Main content</p>
@endsection
```

---

# Database

## Query Builder

```go
import "github.com/mechneerd/gow/database/query"

b := query.NewBuilder(conn, dialect)

// Select
rows, err := b.Table("users").
    Select("name", "email").
    Where("age", ">", 18).
    OrWhere("role", "=", "admin").
    WhereIn("status", []any{"active", "pending"}).
    WhereNull("deleted_at").
    WhereBetween("created_at", []any{start, end}).
    OrderBy("name", "ASC").
    Limit(10).
    Offset(20).
    Get()

// Joins
rows, err := b.Table("users").
    Join("posts", "users.id", "=", "posts.user_id").
    LeftJoin("profiles", "users.id", "=", "profiles.user_id").
    Where("users.active", "=", true).
    Get()

// Aggregates
count, _ := b.Table("users").Count("")
maxAge, _ := b.Table("users").Max("age")
avgPrice, _ := b.Table("products").Avg("price")
totalSales, _ := b.Table("orders").Sum("amount")

// Insert
result, _ := b.Table("users").Insert(map[string]any{
    "name":  "John",
    "email": "john@example.com",
})

// Update
result, _ := b.Table("users").
    Where("id", "=", 1).
    Update(map[string]any{"name": "Jane"})

// Delete
result, _ := b.Table("users").
    Where("id", "=", 1).
    Delete()

// Conditional Clauses
b.Table("users").
    When(isAdmin, func(b *query.Builder) {
        b.Where("role", "=", "admin")
    }).
    Get()

// Raw Queries
rows, _ := b.Table("users").
    SelectRaw("COUNT(*) as total, role").
    GroupBy("role").
    Having("total", ">", 5).
    Get()
```

### Query Events

```go
query.Listen(func(event query.QueryEvent) {
    log.Printf("SQL: %s | Duration: %v", event.SQL, event.Duration)
})
```

---

## ORM (Goquent)

### Defining Models

```go
package models

import (
    "time"
    "github.com/mechneerd/gow/database/orm"
)

type User struct {
    orm.Model
    Name            string     `db:"name"`
    Email           string     `db:"email"`
    Password        string     `db:"password"`
    EmailVerifiedAt *time.Time `db:"email_verified_at"`
    RememberToken   *string    `db:"remember_token"`
    DeletedAt       *time.Time `db:"deleted_at" gow:"softDelete"`
}

func (User) TableName() string {
    return "users"
}
```

### Struct Tags

| Tag | Description |
|-----|-------------|
| `db:"column_name"` | Maps field to DB column |
| `gow:"primaryKey"` | Marks the primary key |
| `gow:"autoIncrement"` | Auto-increment field |
| `gow:"autoCreateTime"` | Set to `time.Now()` on INSERT |
| `gow:"autoUpdateTime"` | Set to `time.Now()` on INSERT and UPDATE |
| `gow:"softDelete"` | Enables soft deletes |
| `gow:"table=users"` | Custom table name |

### CRUD Operations

```go
q := orm.NewQuery[User](db)

// Find by ID
user, _ := q.Find(1)

// First with conditions
user, _ := q.Where("email", "=", "john@example.com").First()

// Get all
users, _ := q.Where("active", "=", true).Get()

// Insert
user := &User{Name: "John", Email: "john@example.com"}
err := q.Insert(user)

// Update
user.Name = "Jane"
err := q.Update(user)

// Save (Insert if new, Update if exists)
err := q.Save(user)

// Delete
err := q.Delete(user)
```

### Soft Deletes

```go
// Automatically excludes soft-deleted records
users, _ := q.Where("active", "=", true).Get()

// Include soft-deleted
users, _ := q.WithTrashed().Get()

// Only soft-deleted
users, _ := q.OnlyTrashed().Get()

// Restore
err := q.Restore(user)

// Force delete
err := q.ForceDelete(user)
```

### Eager Loading

```go
// Load relationships
users, _ := q.With("Posts", "Profile").Get()

// Nested eager loading
users, _ := q.With("Posts.Comments").Get()

// Conditional eager loading
users, _ := q.With("Posts", func(q *orm.ModelQuery[models.Post]) {
    q.Where("published", "=", true)
}).Get()
```

### Relationships

```go
// HasMany
type User struct {
    orm.Model
    Posts []Post `gow:"hasMany"`
}

// BelongsTo
type Post struct {
    orm.Model
    UserID uint   `db:"user_id"`
    User   *User  `gow:"belongsTo"`
}

// HasOne
type User struct {
    orm.Model
    Profile *Profile `gow:"hasOne"`
}

// BelongsToMany
type User struct {
    orm.Model
    Roles []Role `gow:"belongsToMany"`
}

// Pivot table operations
q.Attach(user.ID, role.ID)      // Attach
q.Detach(user.ID, role.ID)     // Detach
q.Sync(user.ID, []uint{1,2,3}) // Sync
q.Toggle(user.ID, role.ID)     // Toggle
```

### Lifecycle Hooks

```go
type User struct {
    orm.Model
    // ...
}

func (u *User) BeforeCreate() error {
    // Hash password before insert
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

func (u *User) AfterCreate() error {
    // Send welcome email
    return mail.SendWelcomeEmail(u)
}
```

### Available Hooks
- `BeforeCreate`, `AfterCreate`
- `BeforeUpdate`, `AfterUpdate`
- `BeforeSave`, `AfterSave`
- `BeforeDelete`, `AfterDelete`
- `BeforeFind`, `AfterFind`

### Global Scopes

```go
type ActiveScope struct{}

func (s ActiveScope) Apply(builder *query.Builder) *query.Builder {
    return builder.Where("active", "=", true)
}

func init() {
    orm.AddGlobalScope("users", ActiveScope{})
}
```

### Model Observers

```go
type UserObserver struct{}

func (o *UserObserver) Creating(model any) bool {
    // Return false to cancel creation
    return true
}

func (o *UserObserver) Created(model any) {
    // After creation
}

func init() {
    orm.Observe("users", &UserObserver{})
}
```

### Attribute Casting

```go
type User struct {
    orm.Model
    Settings map[string]any `db:"settings" gow:"json"`
    IsActive bool           `db:"is_active" gow:"bool"`
    Balance  float64        `db:"balance" gow:"float"`
    Metadata []string       `db:"metadata" gow:"array"`
}
```

### Pagination

```go
// Simple pagination
page, _ := q.Paginate(1, 15)
// page.Data, page.Total, page.LastPage, page.PerPage, page.CurrentPage

// Cursor pagination (for infinite scroll)
page, _ := q.CursorPaginate(cursor, 15)
// page.Data, page.NextCursor, page.HasMore

// Length-aware pagination
page, _ := q.LengthAwarePaginate(1, 15)
// Includes links: page.FirstPage, page.LastPage, page.NextPage, page.PrevPage
```

---

## Migrations

### Creating Migrations

```bash
gow make:migration create_users_table
```

### Migration File

```go
package migrations

import (
    "github.com/mechneerd/gow/database/migration"
    "github.com/mechneerd/gow/database/schema"
)

type CreateUsersTable struct{}

func (CreateUsersTable) Up(m *schema.Builder) error {
    return m.Create("users", func(table *schema.Blueprint) {
        table.ID()
        table.String("name", 255)
        table.String("email", 255).Unique()
        table.String("password", 255)
        table.Boolean("is_active").Default(true)
        table.Timestamps()
    })
}

func (CreateUsersTable) Down(m *schema.Builder) error {
    return m.Drop("users")
}

func init() {
    migration.Register("2026_05_24_000001_create_users_table", CreateUsersTable{})
}
```

### Running Migrations

```bash
gow migrate              # Run all pending migrations
gow migrate:rollback     # Rollback last migration
gow migrate:fresh        # Drop all tables and re-run
gow migrate:refresh      # Rollback and re-run
gow migrate:status       # Show migration status
gow migrate:run <name>   # Run specific migration
```

### Schema Builder

```go
// Create table
m.Create("posts", func(table *schema.Blueprint) {
    table.ID()
    table.String("title", 255)
    table.Text("content").Nullable()
    table.Integer("user_id")
    table.Boolean("published").Default(false)
    table.Timestamps()
    table.SoftDeletes()
})

// Modify table
m.Table("users", func(table *schema.Blueprint) {
    table.String("phone", 20).Nullable()
})

// Drop table
m.Drop("old_table")

// Drop if exists
m.DropIfExists("temp_table")
```

### Column Types

| Method | SQL Type |
|--------|----------|
| `table.ID()` | INTEGER PRIMARY KEY AUTOINCREMENT |
| `table.String("name", 255)` | VARCHAR(255) |
| `table.Text("content")` | TEXT |
| `table.Integer("count")` | INTEGER |
| `table.BigInteger("big_id")` | BIGINT |
| `table.Float("price")` | REAL |
| `table.Double("precise")` | REAL |
| `table.Boolean("active")` | BOOLEAN |
| `table.Date("birthday")` | DATE |
| `table.Timestamp("created_at")` | TIMESTAMP |
| `table.Timestamps()` | created_at, updated_at |
| `table.SoftDeletes()` | deleted_at |
| `table.Binary("data")` | BLOB |
| `table.JSON("metadata")` | TEXT (JSON string) |
| `table.Decimal("amount", 8, 2)` | DECIMAL(8,2) |
| `table.Uuid("uuid")` | VARCHAR(36) |

### Column Modifiers

| Modifier | Description |
|----------|-------------|
| `.Nullable()` | Allow NULL values |
| `.Default(value)` | Set default value |
| `.Unique()` | Add unique constraint |
| `.Primary()` | Mark as primary key |
| `.Unsigned()` | Unsigned integer (MySQL) |
| `.Index()` | Add index |
| `.Comment("text")` | Column comment |

---

## Seeders

```bash
gow make:seeder RoleSeeder
```

```go
package seeders

import (
    "database/sql"
    "github.com/mechneerd/gow/database/seeder"
)

type RoleSeeder struct{}

func (s *RoleSeeder) Run(db *sql.DB) error {
    _, err := db.Exec(`INSERT INTO roles (name) VALUES (?)`, "admin")
    return err
}

func init() {
    seeder.Register("RoleSeeder", &RoleSeeder{})
}
```

```bash
gow db:seed              # Run all seeders
gow db:seed RoleSeeder   # Run specific seeder
```

---

# Authentication

## Authentication

### Session Guard

```go
import "github.com/mechneerd/gow/auth"

// Create guard
guard := auth.NewSessionGuard("web", userProvider, sessionManager)

// Attempt login
success := guard.Attempt(map[string]any{
    "email":    "john@example.com",
    "password": "secret",
})

// Check authentication
if guard.Check() {
    user := guard.User()
}

// Login
guard.Login(user)

// Logout
guard.Logout()
```

### Auth Manager

```go
manager := auth.NewManager()
manager.AddGuard("web", sessionGuard)
manager.AddGuard("api", sanctumGuard)

guard := manager.Guard("web")
```

---

## Authorization & Gates

### Defining Gates

```go
import "github.com/mechneerd/gow/auth/access"

gate := access.NewGate()

// Simple gate
gate.Define("edit-post", func(user any, post any) bool {
    u := user.(*models.User)
    p := post.(*models.Post)
    return u.ID == p.UserID
})

// Using in handlers
if gate.Allows(user, "edit-post", post) {
    // Allow
}
```

### Authorization in Templates

```go
// In route handler
data["Gate"] = gate
```

```html
@can('edit-post', $post)
    <a href="/posts/{{ $post->id }}/edit">Edit</a>
@endcan

@cannot('delete-post', $post)
    <p>You cannot delete this post</p>
@endcannot
```

### Policies

```go
type PostPolicy struct{}

func (p *PostPolicy) Update(user any, post any) bool {
    return user.(*models.User).ID == post.(*models.Post).UserID
}

func (p *PostPolicy) Delete(user any, post any) bool {
    return user.(*models.User).ID == post.(*models.Post).UserID
}
```

---

## RBAC

### Assigning Roles

```go
import "github.com/mechneerd/gow/auth/rbac"

// Assign role
user.AssignRole("admin")

// Check role
if user.HasRole("admin") {
    // Admin logic
}

// Remove role
user.RemoveRole("admin")

// Sync roles
user.SyncRoles([]string{"admin", "editor"})
```

### Role Middleware

```go
router.Group("/admin", func(r *routing.Router) {
    r.Use(middleware.RoleMiddleware("admin"))
    r.Get("/dashboard", adminDashboard)
})
```

### Permission Middleware

```go
router.Group("/posts", func(r *routing.Router) {
    r.Use(middleware.PermissionMiddleware("posts.manage"))
    r.Get("/", listPosts)
    r.Post("/", createPost)
})
```

---

# Frontend

## Livewire

GoW Livewire provides reactive UI components that update without page reloads.

### Setting Up

```html
<!-- In your layout -->
<script src="/js/livewire.js"></script>
```

### Counter Component

```go
// routes/web.go
var counterValue int

http.HandleFunc("/livewire/counter", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if r.Method == http.MethodPost {
        var req struct {
            Count int `json:"count"`
            Delta int `json:"delta"`
        }
        json.NewDecoder(r.Body).Decode(&req)
        counterValue = req.Count + req.Delta
    }

    json.NewEncoder(w).Encode(map[string]any{
        "count": counterValue,
        "html":  fmt.Sprintf(`<div class="counter">%d</div>`, counterValue),
    })
})
```

```html
<!-- In your view -->
<div id="counter">
    <span id="count">0</span>
    <button onclick="updateCounter(1)">+1</button>
    <button onclick="updateCounter(-1)">-1</button>
</div>

<script>
function updateCounter(delta) {
    fetch('/livewire/counter', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ count: parseInt(document.getElementById('count').textContent), delta })
    })
    .then(res => res.json())
    .then(data => {
        document.getElementById('count').textContent = data.count;
    });
}
</script>
```

---

# Services

## Cache

```go
import "github.com/mechneerd/gow/cache"

// Create cache store
store := cache.NewMemoryStore()

// Basic operations
store.Put("key", "value", 5*time.Minute)
value, ok := store.Get("key")
store.Forget("key")
store.Flush()

// Increment/Decrement
store.Increment("counter", 1)
store.Decrement("counter", 1)

// Forever
store.Forever("key", "value")

// Rate limiting
limiter := cache.NewRateLimiter(store)
if limiter.TooManyAttempts("user:123", 60) {
    http.Error(w, "Too Many Requests", 429)
}
limiter.Hit("user:123", 1) // 1 minute decay
```

### Cache Drivers

| Driver | Description |
|--------|-------------|
| `memory` | In-memory with TTL (default) |
| `redis` | Redis-backed cache |

---

## Session

```go
import "github.com/mechneerd/gow/session"

// Create session
store := session.NewFileStore()
mgr := session.NewManager(store, "")

// Start session
mgr.Start()

// Get/Put
value := mgr.Get("key")
mgr.Put("key", "value")

// Flash data
mgr.Flash("status", "Saved!")
status := mgr.Old("status", "")
mgr.Reflash()     // Keep all flash for another request
mgr.Keep("error") // Keep specific key

// Flash input for form validation
mgr.FlashInput(map[string]any{
    "email": r.FormValue("email"),
})

// Regenerate session ID
mgr.Regenerate()

// Save
mgr.Save()
```

### Session Drivers

| Driver | Description |
|--------|-------------|
| `file` | File-based sessions |
| `cookie` | Cookie-based sessions |
| `redis` | Redis-backed sessions |

---

## Queue

```go
import "github.com/mechneerd/gow/queue"

// Define a job
type SendEmailJob struct {
    To      string
    Subject string
    Body    string
}

func (j *SendEmailJob) Handle() error {
    // Send email logic
    return nil
}

func (j *SendEmailJob) Failed(err error) {
    log.Printf("Job failed: %v", err)
}

// Dispatch job
q := queue.NewManager(queue.NewSyncDriver())
q.Dispatch(&SendEmailJob{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Hello!",
})
```

### Queue Drivers

| Driver | Description |
|--------|-------------|
| `sync` | Synchronous (default) |
| `memory` | In-memory channel |
| `database` | Database-backed queue |
| `redis` | Redis-backed queue |

### Worker

```go
import "github.com/mechneerd/gow/queue"

mgr := queue.NewManager(queue.NewSyncDriver())
worker := queue.NewWorker(mgr, 5) // 5 concurrent workers
worker.Start()
```

---

## Events

```go
import "github.com/mechneerd/gow/events"

// Define event
type UserCreated struct {
    User *models.User
}

// Create dispatcher
dispatcher := events.NewDispatcher()

// Listen for events
dispatcher.Listen(UserCreated{}, func(event any) {
    e := event.(UserCreated)
    sendWelcomeEmail(e.User)
})

// Wildcard listener (receives all events)
dispatcher.ListenAny(func(event any) {
    log.Printf("Event dispatched: %T", event)
})

// Dispatch event
dispatcher.Dispatch(UserCreated{User: user})
```

---

## Mail

```go
import "github.com/mechneerd/gow/mail"

// Create mailer
driver := &mail.SmtpDriver{
    Host:     "smtp.example.com",
    Port:     "587",
    Username: "user@example.com",
    Password: "secret",
}
mailer := mail.NewMailer(driver)

// Build message
msg := mail.NewMailable().
    From("noreply@example.com").
    To("user@example.com").
    Subject("Welcome!").
    HTML("<h1>Welcome to our app!</h1>").
    Attach("report.pdf", fileBytes)

// Send
mailer.Send(msg)
```

### Mailable Interface

```go
type WelcomeMail struct {
    User *models.User
}

func (m *WelcomeMail) Build() *mail.Message {
    return mail.NewMailable().
        From("noreply@example.com").
        To(m.User.Email).
        Subject("Welcome!").
        HTML("<h1>Welcome, " + m.User.Name + "!</h1>")
}
```

---

## Notifications

```go
import "github.com/mechneerd/gow/notifications"

// Create notification manager
mgr := notifications.NewManager()

// Register channels
mgr.RegisterChannel("mail", notifications.NewMailChannel(mailer))
mgr.RegisterChannel("database", notifications.NewDatabaseChannel(db))

// Send notification
notification := &UserCreatedNotification{User: user}
mgr.Send(notification, user, "mail", "database")
```

---

## Broadcasting

```go
import "github.com/mechneerd/gow/broadcasting"

// Create hub
hub := broadcasting.NewHub()
go hub.Run()

// Create broadcaster
broadcaster := broadcasting.NewWebSocketDriver(hub)

// Broadcast to channels
broadcaster.Broadcast([]string{"notifications"}, "NewMessage", map[string]any{
    "message": "Hello!",
    "user":    "John",
})
```

### WebSocket Endpoint

```go
router.Get("/ws", func(w http.ResponseWriter, r *http.Request) error {
    broadcasting.ServeWs(hub, w, r)
    return nil
})
```

---

## Storage

```go
import "github.com/mechneerd/gow/storage"

// Create storage
fs := storage.NewLocalDriver("storage/app")

// Basic operations
fs.Put("file.txt", contents)
reader, _ := fs.Get("file.txt")
fs.Delete("file.txt")
exists := fs.Exists("file.txt")
url := fs.URL("file.txt")
```

---

## Logging

```go
import "github.com/mechneerd/gow/logging"

// Setup logger
cfg := logging.Config{
    Level:   "debug",
    Channel: "daily",
    Path:    "storage/logs",
}
logger := logging.Setup(cfg)

// Usage
logger.Info("Server started", "port", 8080)
logger.Error("Database error", "err", err)
```

---

## Encryption

```go
import "github.com/mechneerd/gow/encryption"

// Create encrypter
encrypter, _ := encryption.NewEncrypter([]byte("secret-key-32-bytes-long!!!!!"))

// Encrypt
encrypted, _ := encrypter.Encrypt("sensitive data")

// Decrypt
decrypted, _ := encrypter.Decrypt(encrypted)
```

---

## Localization

```go
import "github.com/mechneerd/gow/localization"

// In templates
// {{ __ "welcome.message" }}
// {{ trans "messages.greeting" }}
```

---

# Testing

## Testing

### Creating Tests

```bash
gow make:test UserControllerTest
```

### TestCase

```go
package tests

import (
    "testing"
    "github.com/mechneerd/gow/testing"
)

func TestGetUsers(t *testing.T) {
    tc := gowtesting.NewTestCase(t, router)

    tc.Get("/api/users").
        AssertStatus(200).
        AssertJson(map[string]any{
            "data": []any{},
        })
}

func TestCreateUser(t *testing.T) {
    tc := gowtesting.NewTestCase(t, router)

    tc.Post("/api/users", map[string]any{
        "name":  "John",
        "email": "john@example.com",
    }).
        AssertStatus(201).
        AssertJsonStructure([]string{"id", "name", "email"})
}

func TestUnauthorizedAccess(t *testing.T) {
    tc := gowtesting.NewTestCase(t, router)

    tc.Get("/api/admin").
        AssertStatus(401)
}
```

### Acting As

```go
func TestAdminDashboard(t *testing.T) {
    tc := gowtesting.NewTestCase(t, router)

    tc.ActingAs(adminUser).
        Get("/admin/dashboard").
        AssertStatus(200)
}
```

### Database Assertions

```go
func TestUserCreation(t *testing.T) {
    tc := gowtesting.NewTestCase(t, router)

    tc.Post("/api/users", userData).
        AssertStatus(201).
        AssertDatabaseHas("users", map[string]any{
            "email": "john@example.com",
        })
}
```

### Available Assertions

| Assertion | Description |
|-----------|-------------|
| `AssertStatus(code)` | HTTP status code |
| `AssertJson(data)` | JSON response matches |
| `AssertJsonStructure(keys)` | JSON has specified keys |
| `AssertHeader(key, value)` | Response header |
| `AssertSee(text)` | Response contains text |
| `AssertDontSee(text)` | Response doesn't contain text |
| `AssertDatabaseHas(table, data)` | Database has record |
| `AssertDatabaseMissing(table, data)` | Database missing record |
| `AssertDatabaseCount(table, count)` | Database record count |

---

# CLI

## Artisan CLI

```bash
# Project management
gow new myapp                    # Create new project
gow serve                        # Start dev server
gow about                        # Framework info

# Database
gow migrate                      # Run migrations
gow migrate:rollback             # Rollback migrations
gow migrate:fresh                # Drop & re-run
gow migrate:status               # Show status
gow db:seed                      # Run seeders

# Generators
gow make:controller PostController    # Create controller
gow make:model Post                   # Create model
gow make:migration create_posts       # Create migration
gow make:middleware Auth              # Create middleware
gow make:job SendEmailJob             # Create job
gow make:seeder RoleSeeder            # Create seeder
gow make:request StorePostRequest     # Create form request
gow make:command MyCommand            # Create custom command

# Cache
gow cache:clear                   # Clear all cache
gow cache:forget key              # Forget cache key

# Config
gow config:cache                  # Cache config
gow config:clear                  # Clear config cache

# Queue
gow queue:work                    # Start queue worker
gow queue:retry id                # Retry failed job
gow queue:failed                  # List failed jobs
gow queue:flush                   # Flush failed jobs

# Other
gow key:generate                  # Generate APP_KEY
gow env                           # Show environment
gow list                          # List all commands
```

### Custom Commands

```bash
gow make:command SendNewsletter
```

```go
package commands

import (
    "fmt"
    "github.com/spf13/cobra"
)

var SendNewsletterCmd = &cobra.Command{
    Use:   "newsletter:send",
    Short: "Send newsletter to all users",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Sending newsletter...")
        // Your logic here
    },
}
```

### Registering Commands

```go
// In artisan.go or your main.go
kernel.RegisterCommand(SendNewsletterCmd)
```

---

## Scheduling

```go
import "github.com/mechneerd/gow/console"

scheduler := console.NewScheduler()

// Run every minute
scheduler.EveryMinute("cache:clear")

// Run hourly
scheduler.Hourly("newsletter:send")

// Run on cron schedule
scheduler.Cron("0 9 * * * *", "report:generate")

// Without overlapping
scheduler.WithoutOverlapping("data:sync", 10*time.Minute)
```

### Running Scheduler

```bash
gow schedule:run
```

---

# Appendices

## String Helpers

```go
import "github.com/mechneerd/gow/support/str"

str.Camel("hello_world")        // "helloWorld"
str.Studly("hello_world")       // "HelloWorld"
str.Snake("helloWorld")         // "hello_world"
str.Kebab("helloWorld")         // "hello-world"
str.Random(40)                   // random hex string
str.Limit("long text", 10, "...") // "long text..."
str.Contains("hello", "ell")     // true
```

## Collection

```go
import "github.com/mechneerd/gow/support/collection"

c := collection.Collect([]int{1, 2, 3, 4, 5})

c.Map(func(i int) int { return i * 2 })      // [2, 4, 6, 8, 10]
c.Filter(func(i int) bool { return i > 2 })  // [3, 4, 5]
c.Each(func(i int) { fmt.Println(i) })
c.Chunked(2)                                  // [[1,2], [3,4], [5]]
c.All()                                       // []int{1, 2, 3, 4, 5}
```

---

**GoW Framework** — The Laravel experience in pure Go.
