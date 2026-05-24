# GoW — Laravel for Go

**GoW** is a modern, Laravel-inspired web framework for Go.  
It delivers the **developer experience** and **productivity** of Laravel while staying 100% native Go.

> **Current Target**: Laravel-lite 1.2 — the 80% of Laravel that 80% of developers use every day, done right.

---

## Features

- ✅ Beautiful routing with groups, resources, model binding, signed URLs
- ✅ Powerful Eloquent-style ORM (relations, scopes, soft deletes, casting, upsert, locking, chunking, etc.)
- ✅ Blade-like templating engine with components, slots, and directives
- ✅ Production-ready authentication (Session, Sanctum, Socialite, 2FA, Lockout, Email Verification, Gate)
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
curl -sSfL https://raw.githubusercontent.com/yourusername/gow/main/install.sh | sh
```

**Windows (PowerShell)**
```powershell
iwr -useb https://raw.githubusercontent.com/yourusername/gow/main/install.ps1 | iex
```

After installation:
```bash
gow --version
```

### Manual Installation

```bash
go install github.com/yourusername/gow/cmd/gow@latest
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

| Command                              | Best For                          | Status     |
|--------------------------------------|-----------------------------------|------------|
| `gow new myapp --yes`                | Full web app (with auth)          | Ready      |
| `gow new myapp --api --yes`          | REST / JSON APIs                  | Ready      |
| `gow new mysite --minimal --yes`     | Simple websites / landing pages   | Ready      |
| `minimal-api` kit                    | Ultra-light API only              | Ready      |
| `full`, `admin-panel`, etc.          | More advanced starters            | Planned    |

You can also run `gow new myapp` without flags to use the **interactive wizard**.

Use `--skeleton` to load custom or local templates (including the new planned kits).

---

## Why GoW?

- Excellent generic container + service provider system
- World-class ORM (one of the most complete in Go)
- Production-ready Authentication stack (Session + Sanctum + Socialite + 2FA)
- Very powerful `gow new` scaffolding + rich Artisan CLI
- Strong testing utilities
- Native WebSocket broadcasting
- Clean, modern Go codebase
- Laravel-like developer experience in Go

---

## Contributing

We welcome contributions! Please read [CONTRIBUTING.md](.github/CONTRIBUTING.md) (will be added before release).

---

## License

MIT License © 2026 GoW Contributors

---

**GoW is now genuinely production-viable** for building real-world applications (SaaS, admin panels, APIs, full-stack apps).

Start building today with `gow new myblog --yes` (or run `gow new myblog` for the interactive wizard)
