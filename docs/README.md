# GoW Documentation

> ✅ Implemented · 🚧 In Progress · 📋 Planned

Welcome to the official documentation for the **GoW Framework** — a Laravel-inspired web framework for Go.

> **Standout Feature**: GoW ships with one of the most complete Artisan-style CLIs in the Go ecosystem (`gow` / `artisan`). See the dedicated section below.

## Guide — ✅ Implemented Features

**Start Here if you're new**:

- **[Quick Start](guide/quick-start.md)** — Fastest way to get a running app (10–20 minutes)
- **[Quick Start - One Page Printable](guide/quick-start-one-page.md)** — Ultra-short version (fits on one printed page)
- **[Getting Started (Zero to Advanced)](guide/getting-started.md)** — Very detailed, step-by-step guide from zero experience to React/Vue apps

These features are fully implemented and tested.

| Section | Status |
|---|---|
| [Getting Started (Beginner → Advanced)](guide/getting-started.md) | ✅ New |
| [Installation](guide/installation.md) | ✅ |
| [Routing](guide/routing.md) | ✅ |
| [Middleware](guide/middleware.md) | ✅ |
| [Database & Tooling](guide/database.md) | ✅ |
| [ORM (Goquent)](guide/orm.md) | ✅ |
| [Cache & Session](guide/cache_and_session.md) | ✅ |
| [Views (Goblade)](guide/views.md) | ✅ |
| [Authentication](guide/authentication.md) | ✅ |
| [Events](guide/events.md) | ✅ |
| [Queues](guide/queues.md) | ✅ |
| [Utilities](guide/utilities.md) | 🚧 |
| [Socialite (OAuth)](guide/socialite.md) | ✅ New |
| [Testing](guide/testing.md) | ✅ New |
| [CLI & Artisan](https://github.com/mechneerd/gow) | ✅ Very Strong |

## Artisan CLI (gow / artisan)

GoW includes a first-class Artisan-style command-line interface built on Cobra. It is one of the most complete and production-ready parts of the framework.

### Highlights

- Laravel-like experience (`gow make:*`, `gow migrate:*`, etc.)
- Excellent code generators that produce real, usable code
- Full migration system with rollback and status support
- Easy to extend with `make:command`

### Key Commands (as of May 24, 2026)

| Category              | Command                              | Description                                      | Example |
|-----------------------|--------------------------------------|--------------------------------------------------|---------|
| **Core**              | `gow serve`                          | Start the development server                     | `gow serve` |
| **Core**              | `gow key:generate`                   | Generate secure APP_KEY                          | `gow key:generate` |
| **Core**              | `gow about`                          | Show framework + environment info                | `gow about` |
| **Core**              | `gow env`                            | Display current environment                      | `gow env` |
| **Database**          | `gow migrate`                        | Run pending migrations                           | `gow migrate` |
| **Database**          | `gow migrate:rollback [--step=N]`    | Rollback migrations                              | `gow migrate:rollback --step=1` |
| **Database**          | `gow migrate:status`                 | Show migration status                            | `gow migrate:status` |
| **Database**          | `gow db:seed`                        | Run database seeders                             | `gow db:seed` |
| **Generators**        | `gow make:model`                     | Create an ORM model (`--migration`)              | `gow make:model Post --migration` |
| **Generators**        | `gow make:migration`                 | Create a new migration file                      | `gow make:migration create_posts_table` |
| **Generators**        | `gow make:controller`                | Create a controller (`--resource`, `--api`)      | `gow make:controller PostController --resource` |
| **Generators**        | `gow make:request`                   | Create a FormRequest (validation)                | `gow make:request StorePostRequest` |
| **Generators**        | `gow make:seeder`                    | Create a database seeder                         | `gow make:seeder RoleSeeder` |
| **Generators**        | `gow make:command`                   | Create a custom Artisan command                  | `gow make:command ClearCache` |
| **Utilities**         | `gow route:list`                     | List all routes (`--path`, `--method`)           | `gow route:list --method=GET --path=posts` |
| **Utilities**         | `gow list`                           | List all available Artisan commands              | `gow list` |

The Artisan CLI is designed to be used every day during development. It is intentionally one of the most polished parts of GoW.

It is also deeply integrated into `gow new` — every new project you scaffold comes with helpful next-step instructions powered by the same command system.

**Pro tip**: After creating a new project with `gow new`, run `gow list` to see every available command.

### Why GoW's Artisan is Special

Unlike many Go CLIs that feel like afterthoughts, GoW's Artisan is a first-class citizen. It was designed from day one to match (and in some areas exceed) the developer experience of Laravel's Artisan. The generators produce idiomatic, modern Go code that integrates cleanly with the framework's ORM, validation, and routing systems.

### Recent Artisan & Scaffolding Improvements (May 24, 2026)

During this release cycle the following CLI and scaffolding enhancements were completed:

- **New commands added**:
  - `gow migrate:rollback --step=N`
  - `gow migrate:run <migration>`
  - `gow migrate:status`
  - `gow key:generate`
  - `gow about`
  - `gow env`
  - `gow list`
  - `gow make:seeder`
  - `gow make:request`
  - `gow make:command` (meta generator)

- **Existing generators significantly improved**:
  - `make:model --migration`
  - `make:controller --api`
  - `route:list` now supports `--path` and `--method` filters + shows middleware

- **Scaffolding system**:
  - Official `gow-skeleton` now contains real, production-ready files for `web-auth` and `full` kits
  - Includes 9 database migrations (users, roles, permissions, pivots, sessions, verification/reset tokens)
  - `RoleSeeder` creates default roles + a ready-to-use Super Admin (`superadmin` / `12345678`)

- **Developer experience**:
  - `gow new` post-install instructions now recommend running `gow list`, `gow migrate`, `gow db:seed`, and `gow key:generate`
  - All major generators now produce modern, framework-aligned code

These changes make GoW one of the most "Laravel-like" experiences available in the Go ecosystem while remaining fully idiomatic.

### Common Artisan Workflows

```bash
# Full web + auth + RBAC project (recommended)
gow new myapp --auth --yes
cd myapp
gow key:generate
gow migrate
gow db:seed          # Creates superadmin user (superadmin / 12345678)

# API project with migrations
gow new myapi --api --yes --db=postgres
cd myapi
gow make:migration create_posts_table
gow make:model Post
gow migrate:status

# Daily development loop
gow serve
gow make:model Comment --migration
gow make:controller CommentController --api
gow make:request StoreCommentRequest
gow route:list --path=comments
```

### Artisan + Scaffolding Integration

The CLI is deeply integrated with `gow new`:

- When you run `gow new myapp --auth`, it pulls from the official `gow-skeleton` repository.
- The generated project includes pre-built migrations + seeders (including the Super Admin user with `superadmin / 12345678`).
- After scaffolding, the post-install output gives you the exact Artisan commands to run next (`gow migrate`, `gow db:seed`, etc.).

This combination (powerful generators + production-ready starter kits with full RBAC) is one of GoW’s biggest advantages over typical Go frameworks.

### Extending the CLI

You can create your own commands with:
```bash
gow make:command MyCustomCommand
```

Then register them in your console kernel (in `artisan.go` or your console kernel).

This level of CLI tooling is rare in Go frameworks and is a major productivity advantage.

The CLI is designed to be Laravel-like while staying idiomatic Go. Many generators produce production-ready stubs that integrate with the ORM, validation, and RBAC systems.

---

## Roadmap — In Progress & Planned

Most major features are now **implemented**. The roadmap is kept for future ecosystem work.

| Section | Status |
|---|---|
| [Testing](roadmap/testing.md) | ✅ Largely Complete |
| [File Storage](roadmap/storage.md) | ✅ Complete |
| [Authorization](roadmap/authorization.md) | ✅ Complete (Gates + Policies) |
| [Mail & Notifications](roadmap/mail_and_notifications.md) | ✅ Good (Markdown + Channels + Slack) |
| [Broadcasting](roadmap/broadcasting.md) | ✅ Complete (Native WebSocket) |
| Socialite / OAuth | ✅ New Guide |
| Advanced ORM Features | ✅ Documented in `guide/orm.md` |

---

*Note for Contributors: This documentation is maintained in Markdown to ensure minimal dependencies and maximal contribution ease. If you wish to propose edits, feel free to open a PR directly modifying these Markdown files!*

---

## Latest Release

**May 23, 2026 — Feature Parity Release**

> **May 24, 2026 Update**: Major Artisan CLI and scaffolding work completed. The command-line experience is now one of the strongest aspects of GoW (see Artisan section below). Official starter kits (`web-auth`, `full`) now ship with complete RBAC migrations + production-ready seeders.

This is the largest release in GoW history. Most remaining critical gaps have been closed.

**Release Materials** (all new or heavily updated on 2026-05-23):

- Full Release Notes: [2026-05-23.md](whats_new/2026-05-23.md)
- GitHub Releases Version: [2026-05-23-github-release.md](whats_new/2026-05-23-github-release.md)
- Short Announcement Summary: [2026-05-23-release-summary.md](whats_new/2026-05-23-release-summary.md)
- **Upgrade / Migration Guide**: [UPGRADE.md](UPGRADE.md) ← Recommended starting point for existing projects
- Announcement Templates: [2026-05-23-announcement-template.md](whats_new/2026-05-23-announcement-template.md)
- Blog Post Draft: [2026-05-23-blog-post-draft.md](whats_new/2026-05-23-blog-post-draft.md)
- Twitter Thread: [2026-05-23-twitter-thread.md](whats_new/2026-05-23-twitter-thread.md)
- Release Email Template: [2026-05-23-release-email-template.md](whats_new/2026-05-23-release-email-template.md)
- One-Page Migration Checklist: [2026-05-23-migration-checklist.md](whats_new/2026-05-23-migration-checklist.md)

See also the current project scope: [LARAVEL_LITE_1.2.md](../LARAVEL_LITE_1.2.md) and [CHANGELOG.md](../CHANGELOG.md).
