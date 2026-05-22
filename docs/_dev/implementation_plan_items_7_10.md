# GoW Framework — Items 7–10 Implementation Plan

> **Session**: Top-10 Roadmap Final Sprint (Items 7–10)  
> **Date**: 2026-05-22  
> **Status**: All decisions locked before implementation began.

---

## Goal

Tackle the remaining top-priority tasks (Items 7 through 10) from `SUGGESTIONS.md` to complete the framework's core roadmap. This improves ORM performance, provides a native concurrency queue option, brings documentation up to date, and establishes a solid developer experience via a dedicated CLI tool.

---

## Decision Log

All decisions were locked before any code was written:

| Area | Decision |
|---|---|
| CLI binary name | `gow` via Cobra |
| Queue driver name | `memory` |
| Queue `Pop()` | Blocking; `TryPop()` non-blocking returning nil |
| Queue `Len()` | `len(d.jobs)` — useful for health checks |
| ORM cache key | `reflect.Type` (not string name) |
| Docs badges | Emoji — ✅ 🚧 📋 |
| `gow new` module ref | Remote module path (`github.com/yourname/gow`) — works anywhere |

---

## Proposed Changes

### 1. ORM Metadata Caching (Item 7)

Currently, the ORM evaluates `reflect.Type.NumField()` and parses struct tags (`db`, `gow`) on every query execution (`hydrateModel`, `Insert`, `Update`, etc.). We introduced a global metadata cache using `sync.Map`.

#### [NEW] `database/orm/metadata.go`
- `FieldMeta` struct: `Name`, `Column`, `Index`, `IsPrimary`, `IsAuto`
- `ModelMetadata` struct: `TableName`, `PrimaryKey`, `SoftDeleteCol`, `Fields []FieldMeta`
- `getMetadata(reflect.Type) *ModelMetadata` — parses the struct once, caches forever
- Skips unexported fields via `field.PkgPath != ""`

#### [MODIFY] `database/orm/query.go`
- `NewQuery()` → reads `meta.TableName` and `meta.SoftDeleteCol` from cache
- `Insert()` → loops `meta.Fields` instead of `typ.NumField()`
- `Update()` → same, reads `meta.PrimaryKey` for WHERE clause
- `Delete()` → same
- `Restore()` → same
- `Save()` → same
- `hydrateModel()` → builds `fieldMap` from `meta.Fields`, eliminates all per-call reflection

#### [MODIFY] `database/orm/relation.go`
- `loadRelation()` → uses `getMetadata(targetType)` for table name + primary key instead of `getTableName()` + manual loops

#### [DELETE] `database/orm/softdelete.go`
- The `softDeleteCache sync.Map` and `getSoftDeleteColumn()` function are fully subsumed by `getMetadata()`. File removed.

---

### 2. Goroutine-Based Queue Driver (Item 8)

Added a new queue driver leveraging Go's native concurrency using buffered channels.

#### [NEW] `queue/driver_memory.go`

```go
type MemoryDriver struct {
    jobs chan Job
}

func (d *MemoryDriver) Push(job Job) error      // send to buffered chan
func (d *MemoryDriver) Pop() (Job, error)        // blocking <-d.jobs
func (d *MemoryDriver) TryPop() Job             // non-blocking select
func (d *MemoryDriver) Len() int                // len(d.jobs)
```

#### [NEW] `queue/driver_memory_test.go`
- `TestMemoryDriverConcurrentPushPop` — 500 concurrent push goroutines + 500 concurrent pop goroutines
- `TestMemoryDriverTryPop` — empty queue returns nil, non-empty returns job

#### [MODIFY] `queue/manager.go`
- `NewManager()` auto-registers `memory` driver with 10,000-job buffer
- Also fixed pre-existing `unused import` errors in `driver_database.go` and `manager.go`

---

### 3. Reconcile Docs (Item 9)

Aligned all documentation with actual implementation status.

#### [MODIFY] `docs/guide/*.md`
- Added badge legend at top of every file: `> ✅ Implemented · 🚧 In Progress · 📋 Planned`
- All `guide/` docs marked `✅ Implemented`: ORM, Routing, Views, Auth, Middleware, Events, Cache & Session, Database
- `utilities.md` marked `🚧 In Progress` (feature flags/health checks partially stubbed)

#### [MODIFY] `docs/roadmap/*.md`
- `testing.md` → `🚧 In Progress` (core suite exists, HTTP test helpers pending)
- `storage.md` → `🚧 In Progress` (local disk implemented, S3 stubbed)
- `authorization.md` → `📋 Planned`
- `broadcasting.md` → `📋 Planned`
- `mail_and_notifications.md` → `📋 Planned`

#### [NEW] `docs/guide/queues.md`
- Queues **promoted from roadmap stub to guide** — completely rewritten to reflect actual `MemoryDriver` API: `Push`, `Pop`, `TryPop`, `Len`, custom driver registration

#### [MODIFY] `docs/README.md`
- Rewritten as a proper status-table index replacing the old unordered list

---

### 4. `gow` CLI Binary + Scaffolding (Item 10)

Built the official `gow` CLI tool using Cobra.

#### [NEW] `cmd/gow/main.go`

**Commands**:
| Command | Behaviour |
|---|---|
| `gow version` | Prints `GoW Framework version 1.0.0` |
| `gow new <name>` | Scaffolds full project structure (see below) |
| `gow serve` | Prints clear "not yet implemented" message |

**`gow new <name>` scaffold**:
```
myapp/
├── go.mod                          ← requires github.com/yourname/gow
├── cmd/app/main.go
├── routes/web.go
├── app/http/controllers/
│   └── home_controller.go
├── resources/views/welcome.html
├── storage/.gitkeep
└── .env.example
```

---

## Verification Plan

### Automated Tests
- `go test ./database/orm/` — ensure metadata caching does not break any existing ORM tests
- `go test ./queue/` — verify concurrent push/pop, TryPop, Len

### Manual Verification
- `go build -o gow.exe ./cmd/gow`
- `./gow.exe version`
- `./gow.exe new testapp`
- `./gow.exe serve`
