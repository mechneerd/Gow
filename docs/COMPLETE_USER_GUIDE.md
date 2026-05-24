# GoW Framework — Complete User Guide

**Version**: 1.0 (Feature Parity Release — Updated May 24, 2026)  
**Target**: Laravel-lite 1.2

This is the single complete documentation file for the GoW framework.  
It contains everything a new user needs after installing or cloning the project.

> **Note**: The `gow new` command and scaffolding system received major improvements in May 2026 (local skeletons, `--yes`, better testing, and multiple new planned starter kits).

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

The `gow new` command is the primary way to bootstrap new GoW applications. It is designed to feel like Laravel's `laravel new` but with modern Go tooling.

#### Default Behavior (Recommended for Most Users)

When you run `gow new` without any flags, it launches an **interactive wizard**:

```bash
gow new myblog
```

The wizard will ask you:
- Which starter kit you want
- Which database driver to use
- Custom module path (optional)

#### Non-Interactive Mode (`--yes`)

For CI, scripting, or when you want fast defaults, use `--yes`:

```bash
gow new myblog --yes
```

This creates a project using the recommended defaults:
- **Starter kit**: Web + Auth (full-stack with authentication)
- **Database**: SQLite

#### Available Starter Kits

| Flag / Kit        | Example                                      | Description                                              | Status     |
|-------------------|----------------------------------------------|----------------------------------------------------------|------------|
| `--auth`          | `gow new app --auth --yes`                   | Full web app + Auth + RBAC + Superadmin seeder           | Ready      |
| `--api`           | `gow new api --api --yes`                    | REST API with Sanctum                                    | Ready      |
| `--minimal`       | `gow new site --minimal --yes`               | Lightweight project (no auth)                            | Ready      |
| `minimal-api`     | via `--skeleton`                             | Ultra-light API only (with basic auth + user table)      | Ready      |
| `web-auth` / `full` | via `--skeleton`                           | Complete Auth + full RBAC + 9 migrations + Superadmin seeder | Ready    |
| `admin-panel`, `inertia-*`, etc. | via `--skeleton`                    | More advanced starters                                   | Planned    |

#### Useful Flags

| Flag              | Description                                      | Example |
|-------------------|--------------------------------------------------|--------|
| `--yes`           | Skip interactive prompts                         | `gow new app --yes` |
| `--db`            | Choose database driver (`sqlite`, `mysql`, `postgres`) | `gow new app --db=postgres` |
| `--module`        | Set custom Go module path                        | `gow new app --module=github.com/user/myapp` |
| `--force`         | Overwrite existing directory                     | `gow new app --force` |
| `--no-git`        | Skip automatic `git init`                        | `gow new app --no-git` |
| `--skeleton`      | Use custom skeleton (remote URL **or local path**) | `gow new app --skeleton=...` |

#### Using Custom Skeletons

You can use any custom skeleton (remote or local) with the `--skeleton` flag:

```bash
# Remote custom skeleton
gow new myapp --skeleton=https://github.com/yourcompany/gow-skeleton.git --yes

# Local skeleton (very useful during development)
gow new myapp --skeleton=/path/to/your/custom-skeleton --yes
```

**Note**: On Windows you can use `D:\path\to\your\custom-skeleton` (forward or backslashes both work).
```

#### Common Real-World Examples

```bash
# Fast full-stack project
gow new myblog --yes

# API with PostgreSQL + custom module
gow new myapi --api --yes --db=postgres --module=github.com/company/myapi

# Force create + skip git
gow new existing-folder --yes --force --no-git
```

**Tip**: Running `gow new myapp` without any flags will launch the interactive wizard. This is the most beginner-friendly experience.

After creating a project, run `gow list` to see every available Artisan command.

### Run the Application

After creating a project:

```bash
cd myblog
gow serve
# or
go run main.go
```

Default address: **http://localhost:8080**

---

## 3. Building an API

### Recommended Command
```bash
# Interactive
gow new myapi --api

# Or non-interactive
gow new myapi --api --yes --db=postgres
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

5. Use Eloquent ORM + API Resources + Sanctum

You now have a production-ready JSON API with validation, authentication, and relationships.

---

## 4. Building a Full Web Application

### Recommended Command
```bash
# Interactive (recommended)
gow new myblog

# Fast non-interactive with auth
gow new myblog --auth --yes
```

### Typical Web App Development Flow

1. Create Model + Migration
2. Create Controller
3. Add routes in `routes/web.go`
4. Create Blade views in `resources/views/`
5. Use built-in authentication (login/register already work out of the box when using `--auth` or the default)

The default / `--auth` starter includes:
- Blade-like templating with components and layouts
- Session-based authentication
- Form validation with old input
- CSRF protection
- Eloquent ORM + relationships

---

## 5. Building a Simple Website

### Recommended Command
```bash
gow new mysite --minimal
# or
gow new mysite --minimal --yes
```

Use basic routing + views. Perfect for:
- Landing pages
- Marketing sites
- Static content websites

This starter has the lightest footprint.

---

## 6. Routing Guide

GoW uses a fast, expressive router (similar to Laravel + chi).

### Key Features
- Clean `router.Get/Post/Put/Delete` methods
- Route groups + middleware
- Named routes + URL generation
- Route model binding
- Content negotiation (`WantsJson()` / `WantsHtml()`)

```go
router.Get("/posts/{id}", handler).Name("posts.show")
router.Group("/api", func(r *routing.Router) {
    r.Middleware("auth:sanctum")
    r.Get("/users", handler)
})
```

See full details in `docs/guide/routing.md`.

---

## 7. ORM (Eloquent-Style) Guide

Powerful ActiveRecord-style ORM with full relationship support.

### Highlights
- `orm.Model` base + automatic timestamps/soft deletes
- Eager loading, scopes, mass assignment protection
- All relationship types (HasOne/HasMany/BelongsTo/Many-to-Many/Polymorphic)
- Transactions, chunking, pessimistic locking, upsert

```go
type Post struct {
    orm.Model
    Title  string
    UserID uint
}

posts := models.Post{}.With("User", "Comments").Where("published", true).Get()
```

See full details in `docs/guide/eloquent.md` and `docs/guide/relationships.md`.

---

## 8. Authentication Guide

Production-ready auth out of the box.

### Included
- Session auth (login/register/logout) for web
- Sanctum API tokens (SPA + mobile)
- Socialite (Google, GitHub, etc.)
- 2FA, email verification, password reset, rate limiting

```go
auth.Attempt(email, password)
token := user.CreateToken("mobile", []string{"*"})
```

Blade directives: `@auth`, `@guest`, `@can`

Full guide: `docs/guide/authentication.md` + `docs/guide/sanctum.md`

---

## 9. Artisan CLI Guide

GoW ships with one of the most complete and useful Artisan-style CLIs in the Go ecosystem (`gow` / `artisan`), built on Cobra.

It is intentionally designed to feel like Laravel's Artisan while staying fully idiomatic Go. Many generators produce production-ready code that integrates with the ORM, validation, and RBAC systems.

### Major CLI Improvements (May 2026)

The CLI received significant upgrades:
- Full migration workflow (`migrate`, `migrate:rollback --step=N`, `migrate:run`, `migrate:status`)
- Excellent code generators (`make:model --migration`, `make:controller --api`, `make:request`, `make:seeder`, `make:command`)
- Convenience commands (`key:generate`, `about`, `env`, `list`)
- Powerful route inspection with filters

### `gow new` — Project Scaffolding (Most Important Command)

`gow new` is the main way to create new GoW projects. It supports both **interactive** and **non-interactive** modes.

#### Basic Usage

```bash
# Interactive wizard (recommended for first-time users)
gow new myblog

# Non-interactive with defaults (Web + Auth + SQLite)
gow new myblog --yes
```

After scaffolding, run `gow list` to see every available command.

### Generator Commands
- `make:controller` (`--resource`, `--api`)
- `make:model` (`--migration`)
- `make:migration`
- `make:request`
- `make:seeder`
- `make:command` (create your own Artisan commands)
- `make:mail`, `make:job`, `make:event`, `make:listener`, `make:policy`, `make:test`, etc.

### Database & Maintenance
- `migrate`, `migrate:fresh`, `migrate:refresh`
- `migrate:rollback` (supports `--step=N`)
- `migrate:run <migration_name>`
- `migrate:status`
- `db:seed`
- `route:list` (with `--path` and `--method` filters)
- `cache:clear`, `config:clear`, `view:clear`
- `queue:work`, `queue:retry`
- `tinker` (interactive REPL)
- `key:generate`
- `about`
- `env`

---

## 10. Views & Templating Guide

Blade-inspired templating engine with components and layouts.

### Core Syntax
```blade
{{ $name }}                    {{-- escaped --}}
{!! $html !!}                  {{-- raw --}}

@if ... @elseif ... @endif
@foreach ($items as $item) ... @endforeach

@auth ... @endauth
@can('update', $post) ... @endcan
```

Layouts via `@extends('layouts.app')` + `@section` / `@yield`.

Reusable components: `<x-alert type="success">...</x-alert>`

See `docs/guide/views.md` and `docs/guide/blade-components.md`.

---

## 11. Validation Guide

Strong form validation with Form Requests (recommended pattern).

```go
type StorePostRequest struct {
    http.FormRequest
}

func (r *StorePostRequest) Rules() map[string]string {
    return map[string]string{
        "title": "required|min:3|max:255",
        "email": "required|email",
    }
}
```

Usage:
```go
func (c *PostController) Store(req *requests.StorePostRequest) {
    data := req.Validated()
}
```

See `docs/guide/validation.md` for all rules and custom messages.

---

## 12. Testing Guide

PHPUnit-style testing built on `testing` package.

```go
suite := testing.NewTestCase(t)
suite.RefreshDatabase()

suite.ActingAs(user)
suite.Post("/posts", data).AssertCreated()
suite.Get("/posts/1").AssertJsonStructure([]string{"id", "title"})
```

Fakes: MailFake, QueueFake, EventFake.

Database assertions: `AssertDatabaseHas`, `AssertDatabaseCount`, etc.

Full guide: `docs/guide/testing.md`

---

## 13. Installation Scripts (for Publishing)

Ready-to-use installers for end users:

- `install.sh` (macOS/Linux)
- `install.ps1` (Windows PowerShell)

They handle `go install`, PATH updates, and verification.

See the scripts in the repo root when preparing a release.

---

## 14. Final Notes

GoW reached genuine production viability in the May 2026 Feature Parity Release.

Core pillars (Foundation, HTTP/Routing, Database/ORM, Authentication) are complete and stable.

### Next Steps
- Start with `gow new myblog --yes`
- Explore `docs/guide/` for deep dives
- Use `gow list` to discover all commands
- Join the community and contribute!

**GoW is now genuinely production-viable** for building real-world applications.

Happy building! 🚀
