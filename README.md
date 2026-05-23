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

```bash
# Full-stack web application (recommended)
gow new myapp

# API-only project
gow new myapi --api

# Minimal website / landing page
gow new mysite --minimal
```

### Run the Application

```bash
cd myapp
go run main.go
# or
gow serve
```

---

## Documentation

**Complete documentation is available in one file:**

- [Complete User Guide](docs/COMPLETE_USER_GUIDE.md) — Everything you need in a single document

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

| Command                        | Best For                     |
|--------------------------------|------------------------------|
| `gow new myapp`                | Full web applications        |
| `gow new myapp --api`          | REST / JSON APIs             |
| `gow new mysite --minimal`     | Simple websites / landing pages |

---

## Why GoW?

- Excellent generic container + service provider system
- World-class ORM (one of the most complete in Go)
- Production-ready Authentication stack
- Very rich Artisan CLI and developer experience
- Strong testing utilities
- Native WebSocket broadcasting
- Clean, modern Go codebase

---

## Contributing

We welcome contributions! Please read [CONTRIBUTING.md](.github/CONTRIBUTING.md) (will be added before release).

---

## License

MIT License © 2026 GoW Contributors

---

**GoW is now genuinely production-viable** for building real-world applications (SaaS, admin panels, APIs, full-stack apps).

Start building today with `gow new myapp`
