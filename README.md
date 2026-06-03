# GoW — Laravel for Go

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/mechneerd/gow)](https://github.com/mechneerd/gow/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/mechneerd/gow)](https://goreportcard.com/report/github.com/mechneerd/gow)

**GoW** is a modern, Laravel-inspired web framework for Go.  
It delivers the **developer experience** and **productivity** of Laravel while staying 100% native Go.

> **Current Target**: Laravel-lite 1.2 — the 80% of Laravel that 80% of developers use every day, done right.

---

## Features

> **Note**: GoW is production-viable for many real-world apps but some advanced features (full RBAC wiring, queue/mail drivers, etc.) may require additional application-level code.

- ✅ Beautiful routing with groups, resources, model binding, signed URLs
- ✅ Powerful Eloquent-style ORM (relations, scopes, soft deletes, casting, upsert, locking, chunking, etc.)
- ✅ Blade-like templating engine with components, slots, and directives
- ✅ Solid authentication foundation (Session, Fortify, Socialite, 2FA, Email Verification) + RBAC (manual wiring for advanced cases)
- ✅ Form Request validation
- ✅ Queue, Events, Broadcasting (WebSocket), Cache, Session
- ✅ Mail + Notifications (with Markdown support)
- ✅ Full Artisan-style CLI (`gow` command) with dozens of generators
- ✅ Excellent testing tools (`actingAs`, fakes, database assertions)

---

## Installation

### Quick Install (Recommended)

**macOS / Linux**
```bash
curl -sSfL https://raw.githubusercontent.com/mechneerd/gow/main/install.sh | sh
```

**Windows (PowerShell)**
```powershell
iwr -useb https://raw.githubusercontent.com/mechneerd/gow/main/install.ps1 | iex
```

After installation:
```bash
gow --version
```

### Manual Installation

```bash
go install github.com/mechneerd/gow/cmd/gow@latest
```

---

## Quick Start

### Create a New Project

`gow new` now includes an **interactive wizard** by default:

```bash
gow new myblog
```

Or use `--yes` for fast, non-interactive creation:

```bash
# Recommended full-stack app (with auth)
gow new myblog --yes

# REST API with PostgreSQL
gow new myapi --api --yes --db=postgres

# Minimal website
gow new mysite --minimal --yes
```

**Supported flags** (most useful ones):
- `--yes` — Skip prompts and use defaults
- `--auth` / `--api` / `--minimal` — Choose starter kit
- `--db=sqlite|mysql|postgres` — Choose database
- `--module=github.com/user/myapp` — Custom module path

### Run the Application

```bash
cd myblog
gow serve
# or
go run main.go
```

---

## Documentation

**Complete documentation is available in one file:**

- [Complete User Guide](docs/COMPLETE_USER_GUIDE.md) — Everything you need in a single document (includes full `gow new` reference)

Individual guides are also available in `docs/guide/`:
- Getting Started
- Routing
- ORM
- Authentication
- Artisan CLI
- Views
- Validation
- Testing

---

## Project Types

| Command                              | Best For                                              | Status     |
|--------------------------------------|-------------------------------------------------------|------------|
| `gow new myapp --yes`                | Full web app + Auth + RBAC                            | Ready      |
| `gow new myapp --api --yes`          | REST / JSON APIs                                      | Ready      |
| `gow new mysite --minimal --yes`     | Simple websites / landing pages                       | Ready      |
| `minimal-api` kit (via --skeleton)   | Ultra-light API only (with basic auth + user table)   | Ready      |
| `web-auth` / `full` (via --skeleton) | Complete Auth + RBAC + Superadmin seeder              | Ready      |
| `admin-panel`, `inertia-*`, etc.     | More advanced starters                                | Planned    |

You can also run `gow new myapp` without flags to use the **interactive wizard**.

Use `--skeleton` to load custom or local templates. The official `gow-skeleton` repo now ships **production-ready kits** (`web-auth`, `full`) with:
- 9 database migrations (users, roles, permissions, pivots, sessions, verification tokens, etc.)
- Full RBAC + Super Admin seeder (`superadmin` / `12345678`)
- Real bootstrap, config, routes, and models

The Artisan CLI (`gow`) has received major upgrades and is now one of GoW’s strongest features.

---

## Why GoW?

- Excellent generic container + service provider system
- World-class ORM (one of the most complete in Go)
- Production-ready Authentication stack (Session + Sanctum + Socialite + 2FA)
- Extremely powerful `gow new` scaffolding + one of the richest Artisan CLIs in Go (dozens of make:* and migrate:* commands)
- Strong testing utilities
- Native WebSocket broadcasting
- Clean, modern Go codebase
- Laravel-like developer experience in Go

---

## Code Examples

### Routing

```go
import "github.com/mechneerd/gow/routing"

router := routing.NewRouter()

// Basic routes
router.Get("/users", ListUsers)
router.Post("/users", CreateUser)
router.Get("/users/{id}", GetUser)

// Route groups with middleware
router.Group("/api", func(r *routing.Router) {
    r.Get("/profile", GetProfile)
    r.Put("/profile", UpdateProfile)
}, authMiddleware)

// Resource routes (RESTful)
router.Resource("/posts", PostController{})

// Named routes & URL generation
router.Get("/dashboard", Dashboard).Name("dashboard")
url := router.Route("dashboard", map[string]string{"id": "1"})
```

### ORM (Eloquent-style)

```go
import "github.com/mechneerd/gow/database/orm"

type User struct {
    orm.Model
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
}

// Find
user, err := orm.Find[User](db, 1)

// Query with scopes
users, err := orm.NewQuery[User](db).
    Where("active", true).
    With("Posts").
    OrderBy("created_at", "desc").
    Paginate(1, 15)

// Create
user = &User{Name: "John", Email: "john@example.com"}
err = orm.Create(db, user)

// Relationships
type Post struct {
    orm.Model
    Title    string
    AuthorID uint
    Author   User `gorm:"foreignKey:AuthorID"`
}

// Eager load
posts, _ := orm.NewQuery[Post](db).With("Author").Get()
```

### Authentication

```go
import (
    "github.com/mechneerd/gow/auth"
    "github.com/mechneerd/gow/auth/sanctum"
)

// Session-based auth
manager := auth.NewManager()
guard := auth.NewSessionGuard("web", provider, sessionManager)
guard.Login(user)

// Sanctum API tokens
token, _ := sanctum.CreateToken(user, []string{"server:update"})
if token.Can("server:update") {
    // Authorized
}

// Middleware
router.Get("/admin", AdminDashboard, auth.Authenticate(manager, "web"))
```

### Blade-like Views

```blade
{{-- resources/views/users/show.html --}}
@extends("layouts.app")

@section("content")
    <h1>{{ user.Name }}</h1>

    @if(len(user.Posts) > 0)
        @foreach(user.Posts as post)
            <article>
                <h2>{{ post.Title }}</h2>
                <p>{!! post.Body !!}</p>
            </article>
        @endforeach
    @else
        <p>No posts yet.</p>
    @endif
@endsection
```

### Validation

```go
import "github.com/mechneerd/gow/validation"

rules := map[string][]string{
    "name":  {"required", "string", "min:2", "max:255"},
    "email": {"required", "email", "unique:users,email"},
    "age":   {"required", "integer", "gte:18"},
}

errors := validation.Validate(data, rules)
if errors.HasErrors() {
    return response.JSON(422, errors.Errors())
}
```

### Queue Jobs

```go
import "github.com/mechneerd/gow/queue"

type SendEmailJob struct {
    queue.Job
    To      string
    Subject string
}

func (j *SendEmailJob) Handle() error {
    // Send the email
    return mail.Send(j.To, j.Subject, j.Body)
}

func (j *SendEmailJob) RetryAfter() int { return 60 }
func (j *SendEmailJob) MaxRetries() int { return 3 }

// Dispatch
queue.Dispatch(&SendEmailJob{To: "user@example.com", Subject: "Welcome!"})

// Dispatch with delay
queue.DispatchLater(&SendEmailJob{To: "user@example.com"}, 5*time.Minute)
```

### Testing

```go
import (
    "testing"
    "github.com/mechneerd/gow/testing/assert"
)

func TestCreateUser(t *testing.T) {
    // Acting as authenticated user
    assert.ActingAs(t, user, func() {
        resp := http.Get("/api/profile")
        assert.Equal(t, 200, resp.StatusCode)
    })

    // Database assertions
    assert.DatabaseHas(t, "users", map[string]any{"email": "test@example.com"})
    assert.DatabaseCount(t, "users", 1)

    // Mail assertions
    assert.MailSent(t, 1)
    assert.MailSentTo(t, "user@example.com")
}
```

### Artisan CLI

```bash
# Scaffolding
gow make:model User --migration
gow make:controller UserController --resource
gow make:request StoreUserRequest
gow make:seeder RoleSeeder

# Database
gow migrate
gow migrate:fresh --seed
gow db:seed

# Development
gow serve
gow route:list
gow about

# Queue
gow queue:work
gow queue:failed
```

---

## Contributing

We welcome contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

MIT License © 2026 GoW Contributors

---

**GoW is now genuinely production-viable** for building real-world applications (SaaS, admin panels, APIs, full-stack apps).

Start building today with `gow new myblog --yes` (or run `gow new myblog` for the interactive wizard)
