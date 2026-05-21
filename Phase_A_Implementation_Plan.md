# Phase A: Core Usability Implementation Plan

## Goal
Implement **Phase A** of the GoW framework roadmap to elevate the core from a skeletal prototype into a usable, developer-friendly foundation. This phase focuses on the Query Builder, ORM capabilities, comprehensive Validation, Blade templating features, and core utility helpers (Collections, Strings).

## Proposed Changes

We will divide Phase A into four logical tracks of execution:

---

### Track 1: Database & ORM Mastery (COMPLETED)

#### [MODIFY] `database/query/builder.go` & `database/dialect/*.go`
- **Feature**: Complete Query Builder.
- Add `Join()`, `LeftJoin()`, `RightJoin()`, `CrossJoin()` methods.
- Add `WhereIn()`, `WhereNull()`, `WhereNotNull()`, `WhereBetween()`.
- Add aggregate methods: `Count()`, `Max()`, `Min()`, `Avg()`, `Sum()`.
- Implement `When(condition bool, callback func(*Builder))` for conditional clauses.
- Update `SelectQuery` struct and Dialects to compile JOINs and aggregates into SQL.

#### [MODIFY] `database/orm/model.go` & `database/orm/relation.go`
- **Feature**: Complete ORM Relationships & CRUD.
- Add `Find()`, `Save()`, `Update()`, `Delete()` directly to the base ORM layer.
- Build relationship structs (`HasMany`, `BelongsTo`, `HasOne`, `BelongsToMany`).
- Implement `Load(relations ...string)` for N+1 safe eager loading using `WhereIn`.
- Implement soft deletes (`deleted_at` filtering).

---

### Track 2: Core Web (Routing, Validation, Session) (COMPLETED)

#### [MODIFY] `routing/router.go`
- **Feature**: HTTP Verb Helpers.
- Add fluent `.Put()`, `.Patch()`, `.Delete()`, `.Options()`, and `.Head()` methods mirroring Laravel's router.

#### [MODIFY] `validation/validator.go`
- **Feature**: Finish Validation.
- Implement rules: `email`, `min`, `max`, `between`, `size`, `numeric`, `string`, `in`, `regex`, `confirmed`.
- Add database-aware rules: `unique:table,column`, `exists:table,column`.

#### [MODIFY] `session/manager.go` & `session/store.go`
- **Feature**: Session Enhancements.
- Add `Flash(key, value)` for temporary data.
- Add `Keep()` and `Reflash()` methods.
- Add `FlashInput(req)` and `Old(key)` helpers for repopulating forms.

---

### Track 3: The Front-End (Goblade & Binding) (COMPLETED)

#### [MODIFY] `view/goblade/compiler.go`
- **Feature**: Advanced Control Structures & Layouts.
- Implement `@extends('layout')`, `@section('name')`, `@yield('name')`.
- Implement `@for`, `@while`, `@switch`, `@case`, `@break`.
- Enhance the regex engine to support nested structures.

#### [MODIFY] `routing/dispatcher.go` (or a new binding package)
- **Feature**: Route Model Binding.
- When resolving controller parameters via reflection, if a parameter is an ORM struct and matches a route parameter (e.g., `{user}` -> `User`), automatically query the DB to inject the model.

---

### Track 4: DX (Helpers, Collections, Artisan) (COMPLETED)

#### [NEW] `support/collection/collection.go`
- **Feature**: Generic Collections.
- Build `type Collection[T any] struct { items []T }`.
- Implement `Map()`, `Filter()`, `Reduce()`, `Each()`, `Pluck()`, `Chunk()`.

#### [NEW] `support/str/str.go` & `support/arr/arr.go`
- **Feature**: Utility Helpers.
- `str`: `Camel()`, `Snake()`, `Studly()`, `Kebab()`, `Random()`, `Limit()`, `Contains()`.
- `arr`: `IsAssoc()`, `Flatten()`, `Only()`, `Except()`, `Wrap()`.

#### [MODIFY] `cmd/artisan/main.go` & `cmd/artisan/make/`
- **Feature**: Artisan Generators.
- Implement `make:controller`, `make:model`, `make:migration`, `make:middleware`, `make:command` with stub template generation.

#### [MODIFY] `foundation/container.go`
- **Feature**: Advanced Container Bindings.
- Add context binding (inject different implementations based on the dependent class).
- Add Facades pattern foundation (global accessor to container bindings).

## Verification Plan

### Automated Tests
- Write unit tests for new `str` and `arr` helpers.
- Write unit tests for `collection` map/filter/reduce logic.
- Expand query builder tests to verify `JOIN` and aggregate SQL compilation.

### Manual Verification
- Generate a new CRUD controller/model via `artisan make:*`.
- Use Route Model Binding to fetch a record.
- Render an `@extends` layout view with flashed session data.
