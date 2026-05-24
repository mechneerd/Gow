# GoW Framework — Bug Fixes & Remediation Log

**Date Range**: May 2026 (Active Remediation Phase)  
**Maintained by**: Kilo (AI-assisted)  
**Purpose**: Track concrete bug fixes, build issues resolved, and code quality improvements made during the post-feature-parity remediation.

This file serves as the authoritative record of **specific bugs that were fixed**, separate from the broader gap analysis in `GAPS_AND_WEAKNESSES.md`.

---

## 2026-05-24 — Artisan Package Full Build Cleanup + Core Service Provider Fixes

### Critical Build Blockers Resolved

**Problem**: `go build ./cmd/artisan/...` was completely broken due to multiple layered issues.

**Root Causes Fixed**:

1. **Missing Database Drivers (Module Errors)**
   - `cmd/artisan/migrator_helper.go` blank-imported `github.com/go-sql-driver/mysql`, `github.com/lib/pq`, and `github.com/mattn/go-sqlite3`.
   - These were never declared in `go.mod`.
   - **Fix**: Added all three drivers via `go get` + `go mod tidy`.
   - Files changed: `go.mod`, `go.sum`

2. **Incorrect Container Resolution Pattern (Reflection Type Errors)**
   - Multiple service providers were calling `app.Resolve(...)` incorrectly:
     - `logging/service_provider.go:17`: `app.Resolve(config.Repository{})` (value instead of `reflect.Type`)
     - `notifications/service_provider.go:18`: `app.Resolve((*interface{})(nil))` (invalid type)
     - `notifications/service_provider.go:24`: Similar misuse for `*mail.Mailer`
     - `cmd/artisan/tinker.go`: Same broken pattern for `*orm.DB`
   - **Fix**: Migrated all calls to the proper generic helper:
     ```go
     container.Make[config.Repository](app.Container)
     container.Make[*mail.Mailer](app.Container)
     ```
   - Also fixed downstream usage (e.g., `Setup(Config{...})` in logging).

3. **Missing Methods on `*orm.DB`**
   - `notifications/channel_database.go:43` and `testing/testing.go` called non-existent methods:
     - `c.db.RawDB()`
     - `c.db.Dialect()`
   - **Fix**:
     - Added `dialect` field to `orm.DB` struct.
     - Implemented `RawDB() query.QueryExecer` and `Dialect() dialect.Dialect` methods.
   - File: `database/orm/query.go`

4. **Duplicate Command Declarations**
   - Many Artisan commands were defined in multiple files (e.g., `CacheClearCmd`, `ConfigClearCmd`, `QueueRetryCmd`, `MakeJobCmd` appeared in both dedicated `*_commands.go` files and `cli_more.go` / `make_more.go`).
   - Caused redeclaration compile errors.
   - **Fix**: Removed duplicate definitions from `cli_more.go` and `make_more.go`. Kept authoritative versions in their proper files.

5. **Syntax Rot & Leftover Garbage Code**
   - Previous edits left malformed code blocks (stray backticks, incomplete if-statements, code outside functions) in:
     - `make_commands.go`
     - `config_commands.go`
   - **Fix**: Cleaned up all syntax errors and removed dead/garbage code.

6. **Initialization Cycle**
   - `list_command.go:10`: `ListCmd` referenced itself inside its own `Run` function, causing a compile-time initialization cycle.
   - **Fix**: Removed `ListCmd` from its own command list.

7. **Type Mismatches**
   - `route_list.go:51`: `strings.Join` called on `[]func(http.Handler) http.Handler` instead of `[]string`.
   - `cmd/gow/main.go` (serve command): Default route handler did not match `routing.HandlerFunc` signature (missing `error` return).
   - **Fixes**:
     - Converted middleware list to type names before joining.
     - Updated handler to `func(...) error { ... return nil }`.

8. **Missing Imports**
   - `tinker.go` was missing `reflect`, `strings`, and `github.com/spf13/cobra`.
   - `schema_commands.go` had imports stripped too aggressively.
   - **Fix**: Restored required imports.

### Additional Improvements (Same Day)

- `gow serve` command fully implemented (real HTTP server using `http.Kernel` with graceful shutdown).
- Continued directory naming standardization in generators (`app/Mail/`, `app/Events/`, `app/Listeners/`, `app/Policies/`, `app/Http/Resources/`, `app/Notifications/`, etc.).
- `db:seed` command now performs real seeder discovery.
- Multiple non-critical Artisan commands received real behavior and better output.

### Result

- `go build ./cmd/artisan/...` → **Clean success**
- `go build ./cmd/gow` → **Clean success**

These changes eliminated the largest remaining technical blocker preventing real-world usage of the Artisan CLI and generated projects.

---

## Future Entries

All subsequent concrete bug fixes during remediation or normal development should be appended here with:
- Date
- Clear problem description
- Root cause
- Files changed
- How it was fixed

This keeps `GAPS_AND_WEAKNESSES.md` focused on strategic gaps while this file tracks tactical fixes.

---

**End of current remediation fixes as of 2026-05-24**
