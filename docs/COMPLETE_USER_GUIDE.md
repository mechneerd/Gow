# GoW Framework — Complete User Guide

**Version**: 1.0 (Feature Parity Release — May 23, 2026)  
**Target**: Laravel-lite 1.2

This is the single complete documentation file for the GoW framework.  
It contains everything a new user needs after installing or cloning the project.

---

## Table of Contents

1. Installation
2. Getting Started
3. Building an API
4. Building a Full Web Application
5. Building a Simple Website
6. Routing Guide
7. ORM (Eloquent-Style) Guide
8. Authentication Guide
9. Artisan CLI Guide
10. Views & Templating Guide
11. Validation Guide
12. Testing Guide
13. Installation Scripts (for Publishing)
14. Final Notes

---

## 1. Installation

### Recommended One-Line Install

**macOS / Linux**
```bash
curl -sSfL https://raw.githubusercontent.com/yourusername/gow/main/install.sh | sh
```

**Windows (PowerShell)**
```powershell
iwr -useb https://raw.githubusercontent.com/yourusername/gow/main/install.ps1 | iex
```

After installation, verify:
```bash
gow --version
```

### Manual Installation

```bash
go install github.com/yourusername/gow/cmd/gow@latest
```

---

## 2. Getting Started

### Create a New Project

GoW supports three main project types:

| Command                        | Best For                        | Includes                              |
|--------------------------------|---------------------------------|---------------------------------------|
| `gow new myapp`                | Full web applications           | Blade views, Auth, ORM, Routing       |
| `gow new myapp --api`          | REST / JSON APIs                | Sanctum, Resources, Form Requests     |
| `gow new mysite --minimal`     | Simple websites / landing pages | Minimal routing + views               |

### Run the Application

```bash
cd myapp
go run main.go
# or
gow serve
```

Default address: **http://localhost:8080**

---

## 3. Building an API

### Recommended Command
```bash
gow new myapi --api
cd myapi
gow serve
```

### Typical API Development Flow

1. Create Model + Migration
   ```bash
   gow make:model Post --migration
   gow migrate
   ```

2. Create Form Request
   ```bash
   gow make:request StorePostRequest
   ```

3. Create API Controller
   ```bash
   gow make:controller Api/PostController --api
   ```

4. Define Routes in `routes/api.go`
   ```go
   router.Group("/api", func(r *routing.Router) {
       r.Middleware("auth:sanctum")
       r.Resource("posts", controllers.PostController{})
   })
   ```

5. Use Eloquent ORM + Resources

You now have a production-ready JSON API with validation, authentication, and relationships.

---

## 4. Building a Full Web Application

### Recommended Command
```bash
gow new myblog
cd myblog
gow serve
```

### Typical Web App Development Flow

1. Create Model + Migration
2. Create Controller
3. Add routes in `routes/web.go`
4. Create Blade views in `resources/views/`
5. Use built-in authentication (login/register already work)

Features included:
- Blade templating with components
- Session authentication
- Form validation with old input
- Eloquent ORM

---

## 5. Building a Simple Website

### Recommended Command
```bash
gow new mysite --minimal
cd mysite
gow serve
```

Use basic routing and views. Ideal for landing pages and marketing sites.

---

## 6. Routing Guide

GoW provides a fast, expressive routing system.

### Basic Routes
```go
router.Get("/", handler)
router.Post("/users", handler)
router.Put("/users/{id}", handler)
router.Delete("/users/{id}", handler)
```

### Route Parameters
```go
router.Get("/users/{id}", func(c *http.Context) {
    id := c.Param("id")
})
```

### Route Groups & Middleware
```go
router.Group("/admin", func(r *routing.Router) {
    r.Middleware("auth")
    r.Get("/dashboard", ...)
})
```

### Named Routes & URL Generation
```go
router.Get("/posts/{id}", handler).Name("posts.show")
url := router.Route("posts.show", map[string]string{"id": "42"})
```

### Resource Routes
```go
router.Resource("posts", PostController{})
```

### Route Model Binding
```go
router.Bind("user", resolver)
user := routing.Model[models.User](c, "user")
```

### Signed URLs
```go
url := routing.TemporarySignedRoute("verify-email", 30*time.Minute, map[string]string{"id": "123"})
```

### Content Negotiation
```go
if c.WantsJson() {
    return c.JSON(data)
}
return c.View("posts.index", data)
```

---

## 7. ORM (Eloquent-Style) Guide

One of GoW's strongest areas.

### Defining a Model
```go
type Post struct {
    orm.Model
    Title   string
    Content string
    UserID  uint
}

func (Post) TableName() string { return "posts" }
```

### Basic Querying
```go
user := models.User{}.Find(1)
users := models.User{}.Where("active", true).Get()
```

### Mass Assignment Protection
```go
func (u User) Fillable() []string { return []string{"name", "email"} }
```

### Soft Deletes
```go
user.Delete()
models.User{}.WithTrashed().Get()
user.Restore()
```

### Relationships
- HasOne / HasMany
- BelongsTo
- BelongsToMany (with Attach/Detach/Sync/Toggle)
- Polymorphic (MorphOne, MorphMany, MorphTo)
- HasOneThrough / HasManyThrough

### Eager Loading
```go
users := models.User{}.With("Posts", "Roles").Get()
```

### Advanced Features
- Attribute Casting
- Accessors & Mutators
- Local Scopes (auto-discovery)
- Upsert
- Pessimistic Locking (`LockForUpdate`, `SharedLock`)
- Chunk / Lazy loading
- Pagination (all types)
- Transactions

---

## 8. Authentication Guide

Production-ready authentication system.

### Session Authentication (Web)
```go
auth.Attempt(email, password)
auth.Login(user, remember)
auth.Logout()
```

### Sanctum (API Tokens)
```go
token := user.CreateToken("api-token", []string{"read", "write"})
```

### Socialite (OAuth)
- Google
- GitHub
- Extensible provider system

### Advanced Security
- Two-Factor Authentication (TOTP)
- Account Lockout
- Remember Me
- Email Verification
- Password Reset

### Authorization
- Gate (`auth.Gate().Define`, `Allows`, `Denies`)
- Policies
- Blade directives: `@auth`, `@can`, `@cannot`

---

## 9. Artisan CLI Guide

Powerful command line tool (`gow`).

### Project Commands
- `gow new myapp`
- `gow new myapp --api`
- `gow new mysite --minimal`
- `gow serve`

### Generator Commands
- `make:controller`
- `make:model`
- `make:migration`
- `make:request`
- `make:resource`
- `make:mail`, `make:job`, `make:event`, `make:listener`, `make:policy`, `make:test`, etc.

### Database & Maintenance
- `migrate`, `migrate:fresh`, `migrate:refresh`
- `db:seed`
- `route:list`
- `cache:clear`, `config:clear`, `view:clear`
- `queue:work`, `queue:retry`
- `tinker` (interactive REPL)

---

## 10. Views & Templating Guide

Blade-inspired templating engine.

### Basic Syntax
```blade
{{ $name }}                    <!-- Escaped -->
{!! $html !!}                  <!-- Raw -->
```

### Conditionals & Loops
```blade
@if ... @elseif ... @else ... @endif
@foreach ... @endforeach
```

### $loop Variable
- `$loop.Index`, `$loop.Iteration`, `$loop.First`, `$loop.Last`, `$loop.Even`, `$loop.Odd`, `$loop.Remaining`

### Layouts & Sections
```blade
@extends('layouts.app')
@section('content') ... @endsection
```

### Components & Slots
```blade
<x-alert type="success">
    <x-slot name="title">Success</x-slot>
    {{ $slot }}
</x-alert>
```

### Auth Directives
```blade
@auth ... @endauth
@guest ... @endguest
@can('update-post', $post) ... @endcan
```

---

## 11. Validation Guide

### Form Requests (Recommended)
```go
type StorePostRequest struct {
    http.FormRequest
}

func (r *StorePostRequest) Rules() map[string]string { ... }
func (r *StorePostRequest) Messages() map[string]string { ... }
func (r *StorePostRequest) Authorize() bool { ... }
```

### Available Rules
`required`, `nullable`, `email`, `url`, `uuid`, `min`, `max`, `between`, `date`, `before`, `after`, `required_if`, `same`, `different`, `gt`, `gte`, `lt`, `lte`, `json`, `boolean`, `integer`, `array`, `alpha`, `alpha_num`, and more.

### Usage in Controllers
```go
func (c *PostController) Store(req *requests.StorePostRequest) {
    data := req.Validated()
}
```

---

## 12. Testing Guide

### Base Test Class
```go
suite := testing.NewTestCase(t)
suite.RefreshDatabase()
```

### Authentication in Tests
```go
suite.ActingAs(user)
```

### HTTP Assertions
```go
suite.Get("/dashboard").AssertOk()
suite.Post("/users", data).AssertCreated()
suite.Get("/posts/1").AssertJsonStructure([]string{"id", "title"})
```

### File Uploads
```go
suite.Post("/upload").Upload("file", "test.png", content)
```

### Database Assertions
- `AssertDatabaseHas`
- `AssertDatabaseMissing`
- `AssertDatabaseCount`
- `AssertDatabaseExactly`

### Fakes
- MailFake
- QueueFake
- EventFake

---

## 13. Installation Scripts (for Publishing)

### install.sh (macOS / Linux)
```bash
#!/usr/bin/env bash
set -e
echo "Installing GoW CLI..."
# ... (full script as previously provided)
```

### install.ps1 (Windows)
```powershell
Write-Host "Installing GoW CLI for Windows..."
# ... (full script as previously provided)
```

---

## 14. Final Notes

- GoW reached **genuine production viability** on May 23, 2026 (Feature Parity Release).
- Core pillars (Foundation, HTTP/Routing, Database/ORM) + Authentication are complete and stable.
- Start with `docs/guide/getting-started.md` or this `COMPLETE_USER_GUIDE.md`.
- Use `gow list` to discover all available commands.
- Read individual guides for deep dives.

---

**End of GoW Complete User Guide**

Thank you for using GoW!
```

