# GoW Framework — Gaps, Weaknesses & Technical Debt Analysis

**Date**: May 24, 2026  
**Version**: Post Feature Parity + Artisan & Scaffolding Upgrade + Remediation (May 24)  
**Status**: Active remediation phase — majority of high-priority gaps have been closed or strongly mitigated. Document now reflects current reality.

This document provides an honest assessment of the framework's gaps and weaknesses. Historical problems are retained for context, while the top sections and Summary Table reflect the post-remediation state as of May 24, 2026.

---

## Production Readiness Status (May 24, 2026)

**Overall Assessment**: Substantial progress has been made. The framework has crossed from "promising prototype" into "usable for real projects with some caveats".

**Major Areas Now Usable in Generated Projects**:
- Full migration workflow (including rollback by step, status, fresh, refresh)
- Working RBAC with pragmatic global DB wiring (`rbac.SetDefaultDB`)
- Strong and mostly functional Artisan CLI
- Production-oriented starter kits (`web-auth`, `full`) with real migrations + working Super Admin seeder
- Basic RBAC middleware available out of the box
- Automatic seeder discovery via `gow db:seed`
- Auth kits now ship with injected RBAC examples + protected route patterns in `bootstrap/app.go`

**Current Recommended Path for New Projects**:
```bash
gow new myapp --auth --yes
cd myapp
gow key:generate
gow migrate
gow db:seed
# RBAC examples + wiring instructions are already in bootstrap/app.go
gow serve
```

**Remaining Areas Needing Attention**:
- Directory structure consistency (actively standardized in generators; more work remains)
- Some non-critical Artisan commands still lightweight
- End-to-end testing coverage
- Deeper integration of services in generated `bootstrap/app.go` (examples now present)

As of the end of this remediation session, GoW is considered **production-viable for many real-world applications**, especially internal tools, APIs, and small-to-medium SaaS products, provided the team is comfortable with some manual wiring for advanced use cases.

---

## 1. Executive Summary

GoW has made impressive progress, particularly in the ORM and the Artisan CLI during May 2026. However, it currently exhibits a classic **"feature-complete on paper"** problem.

**Update (May 24, 2026 - Remediation Phase)**: Active remediation in progress. We are systematically closing gaps to reach production viability.

**Recently Addressed (this session)**
- Migration commands (`migrate`, `migrate:rollback --step=N`, `migrate:run`, `migrate:status`, `migrate:fresh`, `migrate:refresh`) are now functional.
- Basic RBAC implemented (`HasRole`, `HasPermission`, `AssignRole`) + global DB helper + middleware.
- `make:model --migration` now generates real migration files.
- Skeleton RoleSeeder improved (now works with nil DB + direct SQL).
- Migrator helper improved (Postgres support + better errors).
- Added pragmatic global DB wiring for RBAC (`auth/rbac/db.go` + `SetDefaultDB`).
- Directory naming consistency fixes in generators (`database/seeders/`, `database/migrations/`).
- `db:seed` now performs automatic seeder discovery (scans database/seeders/ and lists *Seeder.go files).
- Generated auth kits now receive injected RBAC middleware examples + protected route patterns + `SetDefaultDB` guidance directly into `bootstrap/app.go` via scaffolding post-processor.

**Remaining High-Priority Gaps (being worked on)**
- Full automatic registration of auth middleware in bootstrap (partially improved — examples now auto-injected for auth kits)
- Remaining non-critical placeholder commands (cache, queue, etc.)
- Broader directory naming standardization (more generators + skeletons still need alignment)

**Major Gaps Considered Closed or Strongly Mitigated (as of this session)**
- Artisan migration system (fully functional with step support, status, fresh, refresh, run)
- Core RBAC functionality + practical global DB wiring (`rbac.SetDefaultDB`)
- `make:model --migration` now generates real files
- Skeleton seeder reliability + easy one-line DB setup
- Bootstrap examples + comments for auth + RBAC in generated projects (now auto-injected via scaffold post-processor)
- Basic RBAC middleware available
- Directory naming standardization started in generators (seeders + migrations paths normalized; ongoing)
- Automatic seeder discovery in `db:seed` command

More work is ongoing.

---

## 1b. Executive Summary (Original — Pre-Remediation Context)

GoW has made impressive progress, particularly in the ORM and the Artisan CLI during May 2026. However, it currently exhibits a classic **"feature-complete on paper"** problem.

Many features are declared as implemented in documentation and marketing materials, but the actual implementation is either:
- Placeholder code
- Partially wired
- Non-functional in real usage scenarios

The gap between **claimed capability** and **working reality** is the single largest weakness of the project right now.

---

## 2. Historical Critical Gaps (Mostly Addressed — May 2026 Remediation)

> **Note**: The detailed evidence below describes the state *before* the active remediation session. Many of these issues have been mitigated or closed (see "Major Gaps Considered Closed" above). The sections are retained for historical accuracy.

### 2.1 Artisan Commands Were Mostly Placeholders (Pre-Remediation)

**Severity**: Critical  
**Impact**: The CLI (one of the project's biggest selling points) is largely non-functional for core operations.

#### Implementation Details

**Files involved**:
- `cmd/artisan/migrate_commands.go`
- `cmd/artisan/db_seed_command.go`
- `cmd/artisan/list_command.go` (partially)
- `artisan.go` (registration point)

**Code Evidence** (from `cmd/artisan/migrate_commands.go`):

```go
var MigrateRollbackCmd = &cobra.Command{
    Use:   "migrate:rollback",
    Short: "Rollback the last database migration(s)",
    Run: func(cmd *cobra.Command, args []string) {
        steps, _ := cmd.Flags().GetInt("step")
        fmt.Printf("Rolling back last %d migration(s)...\n", steps)

        // if migrator != nil { migrator.RollbackSteps(steps) }
        fmt.Println("Rollback completed.")
    },
}
```

Similar pattern exists in:
- `MigrateRunCmd`
- `MigrateStatusCmd`
- `DbSeedCmd`

**Root Cause**:
The commands were registered in `artisan.go`:
```go
kernel.RegisterCommand(artisan.MigrateRollbackCmd)
kernel.RegisterCommand(artisan.MigrateRunCmd)
kernel.RegisterCommand(artisan.MigrateStatusCmd)
kernel.RegisterCommand(artisan.DbSeedCmd)
```

However, there is no mechanism to inject a real `*migration.Migrator` instance when the `artisan` binary runs standalone.

**Related Real Implementation That Is Unused**:
In `database/migration/migrator.go`, the following powerful methods exist but are never called from the CLI:

```go
func (m *Migrator) RollbackSteps(steps int) error { ... }
func (m *Migrator) RollbackMigration(name string) error { ... }
func (m *Migrator) MigrateOne(name string) error { ... }
func (m *Migrator) Status() error { ... }
```

These methods are well implemented but dead code from the CLI's perspective.

---

### 2.2 RBAC Implementation Was Incomplete (Stubs Only — Pre-Remediation)

**Severity**: Critical  
**Impact**: The "Full Auth + RBAC" marketing claim in starter kits is misleading.

#### Implementation Details

**Files in skeleton**:
- `D:/gow-skeleton/templates/web-auth/app/Models/Role.go`
- `D:/gow-skeleton/templates/web-auth/app/Models/User.go`
- `D:/gow-skeleton/templates/web-auth/database/seeders/RoleSeeder.go`
- `D:/gow-skeleton/templates/full/app/Models/Role.go` (same problem)

**Code Evidence** (`Role.go` in web-auth kit):

```go
func (r *Role) AttachPermission(perm Permission) {
    // Placeholder for many-to-many logic
}

func (r *Role) AttachUser(user *User) {
    // Placeholder for many-to-many logic
}
```

**In the seeder** (`RoleSeeder.go`):

```go
superAdminRole.AttachPermission(p)   // Does nothing
superAdminRole.AttachUser(&adminUser) // Does nothing
```

**Framework side** (`auth/rbac/has_roles.go`):

```go
func (h *HasRoles) HasRole(roleName string) bool {
    // Real implementation will query the role_user pivot
    // This is a placeholder that always returns false until full wiring
    return false
}
```

**Current State**: RBAC tables are created correctly (good migration work), but the entire permission/role assignment and checking layer is non-functional.

---

### 2.3 Generated Projects Were Not Production-Ready (Pre-Remediation)

**Severity**: Critical

After running `gow new myapp --auth --yes`, the resulting project still requires significant manual work to become functional, despite claims in `README.md` and `docs/README.md`.

---

## 3. Major Weaknesses

### 3.1 Inconsistent Project Structure and Naming (In Progress)

**Severity**: High (improving)

**Evidence** (original):
- `app/Models/` (capital M)
- `app/Http/Controllers/` (mixed)
- `app/console/commands/` (lowercase in some generators)
- `database/seeders/` vs `app/Models/`

**Progress**: Generators now consistently use `database/seeders/` and `database/migrations/` (lowercase). Broader alignment across `app/`, `Http/`, and skeleton templates is ongoing.

**Impact**: Still breaks expectations in some areas but actively being standardized.

---

### 3.2 Artisan Binary vs Application Runtime Disconnect

**Implementation Location**:
- `artisan.go` (standalone binary)
- `console/kernel.go`
- `cmd/gow/main.go` (for `gow` commands)

The `artisan` binary uses its own minimal bootstrap and does not share the same service container or database connection as a real application started via `gow serve` or `go run main.go`.

This is the root cause of most CLI non-functionality.

---

### 3.3 Heavy Reliance on Reflection

Found extensively in:
- `container/container.go`
- `database/orm/model.go` and relation loading
- View engine and some routing features

While functional, this creates performance overhead and makes debugging difficult.

---

## 4. CLI & Scaffolding Specific Gaps (Historical + Current Status)

### 4.1 `migrate:status` — **RESOLVED**
The `Migrator.Status()` method is now wired and functional via `gow migrate:status`.

### 4.2 `make:model --migration` — **RESOLVED**
Now generates real migration files.

### 4.3 Bootstrap & Auth Integration in Generated Kits — **SIGNIFICANTLY IMPROVED**
- RBAC examples + `SetDefaultDB` guidance + protected route patterns are now automatically injected into `bootstrap/app.go` for auth kits.
- Remaining: Full automatic middleware registration (still requires some manual uncommenting in some cases).

---

## 5. Testing & Quality Gaps

- Very few integration tests for the full `gow new → migrate → seed` flow.
- CLI commands have almost no unit tests.
- The powerful new migrator methods (`RollbackSteps`, `Status`, etc.) have no tests.
- End-to-end scaffolding tests are minimal (`cmd/gow/scaffold/scaffold_test.go` only covers basic replacement).

---

## 6. Documentation vs Reality Gap (Improving)

Significant progress has been made during the May 2026 remediation. The gap has narrowed considerably.

**Remaining Risk**:
- Some older docs and marketing claims still overstate maturity.
- Full end-to-end "zero manual wiring" experience for advanced RBAC is not yet 100% automatic.

**Status**: Much improved. The top-level Production Readiness Status in this document is now the most accurate view.

---

## 7. Summary Table (Post-Remediation — May 24, 2026)

| ID | Gap/Weakness                              | Severity      | Primary Location(s)                  | Implementation Status                          |
|----|-------------------------------------------|---------------|--------------------------------------|------------------------------------------------|
| G1 | Artisan commands were placeholders        | Critical      | `cmd/artisan/*_commands.go`          | **Mitigated** — Core migration + seed commands now functional |
| G2 | RBAC logic was stubs                      | Critical      | Skeleton + `auth/rbac/`              | **Mitigated** — Core + global DB wiring + middleware + examples injected |
| G3 | CLI not wired to real migrator            | Critical      | `artisan.go`, `migrator.go`          | **Mitigated** — Migrator is wired in practice for most commands |
| G4 | Inconsistent directory casing             | High          | Generators + skeleton                | **In Progress** — seeders/migrations normalized; broader cleanup ongoing |
| G5 | Weak testing of CLI & scaffolding         | High          | `cmd/gow/scaffold/`, `cmd/artisan/`  | Still limited (real gap)                       |
| G6 | Documentation overstates maturity         | Medium        | Older docs                           | **Improving** — This document now leads with honest status |
| G7 | `make:* --migration` incomplete           | Medium-High   | `make_commands.go`                   | **Resolved** — Real files generated            |
| G8 | Heavy reflection in core paths            | Medium        | `container/`, `database/orm/`        | Still present (architectural trade-off)        |
| G9 | Generated projects require heavy manual work | High       | All skeleton kits                    | **Significantly Reduced** — Auth kits now include injected RBAC guidance |

---

## 8. Recommended Next Steps (Current — Post-Remediation)

**Completed / Strongly Mitigated**:
- Real wiring for migration + seeding commands
- Core RBAC + practical global DB helper + middleware + auto-injected examples
- `make:* --migration` now generates real files
- Directory naming fixes started in generators

**Remaining Priorities** (ranked):
1. Broader directory naming standardization (generators + all skeletons + `app/` structure)
2. End-to-end integration tests for the full `gow new --auth → migrate → db:seed` flow
3. Reduce remaining lightweight/placeholder commands (cache, queue, schedule, etc.)
4. Optional: Deeper automatic middleware registration in generated `bootstrap/app.go` (beyond commented examples)
5. Audit and update older marketing/docs claims for accuracy

The framework is now considered production-viable for most real-world use cases with the caveats noted in the Production Readiness Status section above.

---

**End of Document**

This document has been actively maintained and cleaned up during the May 24, 2026 remediation session to accurately reflect both historical issues and current post-fix reality.

**Maintained by**: Kilo (AI-assisted analysis)  
**Last Updated**: 2026-05-24 (Remediation updates applied)
