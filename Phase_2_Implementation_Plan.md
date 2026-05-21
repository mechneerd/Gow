# Phase 2 Implementation Plan (Database & ORM Layer)

This document outlines the detailed implementation of Phase 2 (Goquent) within the GoW framework.

## 1. Database Connection Manager [P2.1]
- **`database/dialect/dialect.go`**: Abstract interface for handling SQL compilation differences across database engines (quoting, placeholders, pagination).
- **`database/dialect/sqlite.go`**: Implementation of the SQLite dialect.
- **`database/connection.go`**: Wrapper for `*sql.DB` containing its resolved dialect.
- **`database/manager.go`**: Connection manager to retrieve and cache connections based on configuration.

## 2. Query Builder [P2.2]
- **`database/query/builder.go`**: Fluent query builder featuring `Select()`, `Where()`, `OrWhere()`, `OrderBy()`, `Limit()`, and `Offset()`.
- Implements `ToSQL()` by passing its internal representation to the assigned dialect for compilation.

## 3. Schema Builder & Blueprint [P2.3, P2.4]
- **`database/schema/blueprint.go`**: Go-native DSL for defining table schemas (`ID()`, `String()`, `Timestamps()`, `SoftDeletes()`).
- **`database/schema/builder.go`**: Schema execution runner (`Create`, `Drop`, `DropIfExists`).
- **`database/migration/migrator.go`**: Runner logic to track executed migrations in a `migrations` table and rollback executed batches.

## 4. Base ORM (Goquent) & Relationships [P2.5, P2.6, P2.7]
- **`database/orm/query.go`**: Model hydration mechanism relying on reflection. Maps database columns to struct fields using `db` tags.
- Identifies primary keys and auto-increment properties via `gow:"primaryKey,autoIncrement"`.
- Automatically injects `created_at` times via `gow:"autoCreateTime"`.
- Features `First()` and `Get()` returning typed `*T` or `[]*T`.

## 5. Seeders & Factories [P2.8]
- **`database/factory/factory.go`**: Base abstraction for generating mock struct data.
- **`database/seeder/seeder.go`**: Seeder runner to register and execute multiple seeders during development.

## 6. Pagination [P2.9]
- **`database/pagination/paginator.go`**: Paginator generic struct designed to wrap query results and return paginated metadata (Total, PerPage, CurrentPage, LastPage).
