# Wave 2 Implementation Plan

> **Date**: 2026-05-22
> **Scope**: All 7 unimplemented items from SUGGESTIONS.md Wave 2

## Items

| # | Item | Priority | Status |
|---|------|----------|--------|
| 1 | Graceful Shutdown | 🔴 High | 📋 Planned |
| 2 | Connection Pool Config | 🟡 Medium | 📋 Planned |
| 3 | MySQL & PostgreSQL Dialects | 🔴 High | 📋 Planned |
| 4 | Plugin / Service Provider Auto-Discovery | 🟡 Medium | 📋 Planned |
| 5 | Health Check `/up` Endpoint | 🟢 Low | 📋 Planned |
| 6 | Native WebSocket Broadcasting Integration | 🟡 Medium | 📋 Planned |
| 7 | Type-Safe `Map[T,U]` Collection Helper | 🟢 Low | 📋 Planned |

---

## 1. Graceful Shutdown (🔴 High)

### MODIFY `http/kernel.go`

Add a `Serve(addr string) error` method that:
- Creates `*http.Server` with the kernel as handler
- Starts `ListenAndServe` in a goroutine
- Listens for `SIGINT`/`SIGTERM` via `signal.Notify`
- Calls `srv.Shutdown(ctx)` with a configurable timeout (default 30s)
- Exposes `OnShutdown(fn func())` for registering cleanup hooks (close DB, flush logs, drain queue)

```go
func (k *Kernel) Serve(addr string) error { ... }
func (k *Kernel) OnShutdown(fn func()) { ... }
```

---

## 2. Connection Pool Config (🟡 Medium)

### NEW `database/pool.go`

Create a `PoolConfig` struct that reads from the config repository:

```go
type PoolConfig struct {
    MaxOpenConns    int           // DB_MAX_OPEN_CONNS, default 25
    MaxIdleConns    int           // DB_MAX_IDLE_CONNS, default 5
    ConnMaxLifetime time.Duration // DB_CONN_MAX_LIFETIME, default 5m
    ConnMaxIdleTime time.Duration // DB_CONN_MAX_IDLE_TIME, default 3m
}
```

Add `ConfigurePool(db *sql.DB, cfg PoolConfig)` function.

### MODIFY `database/manager.go`

Update `AddConnection` to apply pool config from the config repository automatically.

---

## 3. MySQL & PostgreSQL Dialects (🔴 High)

### NEW `database/dialect/mysql.go`

MySQL dialect implementing the `Dialect` interface:
- `Placeholder`: always `?`
- `QuoteIdentifier`: backtick quoting `` `col` ``
- All Compile methods with MySQL-specific syntax

### NEW `database/dialect/postgres.go`

PostgreSQL dialect implementing the `Dialect` interface:
- `Placeholder`: positional `$1`, `$2`, etc.
- `QuoteIdentifier`: double-quote `"col"`
- All Compile methods with Postgres-specific syntax

### NEW `database/dialect/mysql_test.go` and `database/dialect/postgres_test.go`

SQL compilation tests for each dialect.

---

## 4. Plugin / Service Provider Auto-Discovery (🟡 Medium)

### MODIFY `foundation/provider.go`

Extend with a `PublishableProvider` interface:

```go
type PublishableProvider interface {
    ServiceProvider
    Publishes() map[string]string
}
```

### NEW `foundation/discovery.go`

Auto-discovery system:
- `ProviderRegistry` that manages provider lists
- `AutoDiscover(app *Application)` function
- `PublishAssets(provider PublishableProvider, basePath string)` function

### MODIFY `foundation/application.go`

Add `DiscoverProviders()` and `PublishProvider()` methods.

---

## 5. Health Check `/up` Endpoint (🟢 Low)

### MODIFY `support/health/health.go`

Add `time` and `version` fields to the default response.

### NEW `support/health/checkers.go`

Built-in checkers: `DatabaseChecker` (pings `*sql.DB`), `CacheChecker`.

---

## 6. Native WebSocket Broadcasting Integration (🟡 Medium)

### NEW `broadcasting/service_provider.go`

Service provider that creates Hub, registers WebSocketDriver, wires `/ws` route.

### MODIFY `broadcasting/server.go`

Add `Stop()` method on `Hub` for graceful shutdown.

---

## 7. Type-Safe `Map[T,U]` Collection Helper (🟢 Low)

### MODIFY `support/collection/collection.go`

Add standalone generic functions:

```go
func Map[T, U any](c *Collection[T], fn func(T) U) *Collection[U]
func FlatMap[T, U any](c *Collection[T], fn func(T) []U) *Collection[U]
func Reduce[T, U any](c *Collection[T], fn func(U, T) U, initial U) U
```

---

## Verification

```bash
go build ./...
go test ./...
```
