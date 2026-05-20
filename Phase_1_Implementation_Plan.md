# GoW Framework - Phase 1 Implementation Plan

This plan covers the implementation of Phase 1 of the GoW framework as outlined in the `GoW_SPEC.md` and `CHANGELOG.md` documents. 

Phase 1 focuses on building the core foundation of the framework, including the directory structure, Service Container (IoC), and Application Bootstrap/Kernel.

## User Review Required

Please review the proposed approach for the Service Container and Routing. We will be using generics (since we are targeting Go 1.24+) for type-safe bindings where possible, combined with reflection for constructor dependency injection.

## Open Questions

1. Do you want to build a custom radix-tree router from scratch, or adapt an existing high-performance router like `httprouter` or `chi` to support named routes and route caching?
2. For the Service Container, should we prioritize thread safety (using mutexes for bindings/instances) from the start?

## Proposed Changes

### 1. Project Scaffolding & Directory Structure [P1.1]

We will initialize the Go module and create the core directory structure in a new directory or the current one.
We need to clarify if the framework itself is being built as a library (e.g. `github.com/gow/core`) or if we are building a scaffolding tool (`gow new`) that generates a project which *uses* the framework.
For the purpose of development, we will build the core framework library first.

#### [NEW] `go.mod`
Initialize the `gow` module.

#### [NEW] `gow.go`
Core framework types and interfaces.

### 2. Service Container [P1.2]

The IoC container is the heart of the framework.

#### [NEW] `container/container.go`
Implementation of the IoC container.
- Support for `Bind`, `Singleton`, `Instance`.
- `Make` and `MakeWith` for resolving dependencies.
- Reflection logic to analyze constructor parameters and recursively inject dependencies.

### 3. Application Bootstrap & Kernel [P1.3]

#### [NEW] `foundation/application.go`
The `Application` struct which embeds the `Container` and manages the application lifecycle (Bootstrap, Register, Boot).

#### [NEW] `http/kernel.go`
The HTTP Kernel which receives a request, passes it through the middleware pipeline, and dispatches it to the router.
- `Handle(req *http.Request) *http.Response`
- Implements `http.Handler`

### 4. HTTP Router & Middleware Pipeline [P1.4, P1.5]

#### [NEW] `routing/router.go`
Route registration, groups, and parameters.

#### [NEW] `routing/pipeline.go`
Middleware pipeline execution logic.

## Verification Plan

### Automated Tests
- We will write table-driven unit tests for the `container` to ensure bindings and singleton instances are resolved correctly, especially edge cases like circular dependencies (if we choose to detect them) and complex nested dependencies.
- Tests for the HTTP kernel to ensure request flow and middleware execution order.

### Manual Verification
- We will create a small `main.go` entrypoint to boot the framework, register a simple service, and serve a hello-world route to verify end-to-end functionality of Phase 1 components.
