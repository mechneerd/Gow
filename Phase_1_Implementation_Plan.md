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

### 5. Base Controller & Response Helpers [P1.6]

#### [NEW] `http/response/json.go` & `http/response/redirect.go`
- JSON, Ok, Created, and Error helpers for standardized responses.
- Redirect helper for HTTP redirects.

### 6. Configuration System [P1.7]

#### [NEW] `config/repository.go` & `config/service_provider.go`
- `Repository` to manage configuration values.
- Integration with `joho/godotenv` to load `.env` files.
- Configuration Service Provider to bind the repository to the container.

### 7. Service Provider System [P1.8]

#### [NEW] `foundation/provider.go`
- `ServiceProvider` interface with `Register` and `Boot` lifecycle methods.
- Application updated to manage and boot a list of providers.

### 8. Error Handling & Exception Rendering [P1.9]

#### [NEW] `http/exception.go` & `http/middleware/recovery.go`
- `HttpException` type and `Abort(status, message)` helper to emulate Laravel's abort functionality via panic.
- Panic recovery middleware to catch exceptions and render structured JSON HTTP errors.

### 9. Logging System [P1.10]

#### [NEW] `logging/logger.go` & `logging/service_provider.go`
- Centralized logger using Go's native `log/slog` with JSON formatting.
- Reads `LOG_LEVEL` from configuration.
- Provider to bind the logger into the IoC container.

### 10. CLI Framework (Artisan) [P1.11]

#### [NEW] `console/kernel.go` & `artisan.go`
- Console Kernel utilizing `spf13/cobra` for command registration.
- `artisan.go` root entrypoint for executing CLI commands.

## Verification Plan

### Automated Tests
- We will write table-driven unit tests for the `container` to ensure bindings and singleton instances are resolved correctly, especially edge cases like circular dependencies (if we choose to detect them) and complex nested dependencies.
- Tests for the HTTP kernel to ensure request flow and middleware execution order.

### Manual Verification
- We will create a small `main.go` entrypoint to boot the framework, register a simple service, and serve a hello-world route to verify end-to-end functionality of Phase 1 components.
