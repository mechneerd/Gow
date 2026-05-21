# Phase 3 Implementation Plan (Web Layer)

This document outlines the detailed implementation of Phase 3 within the GoW framework.

## 1. Template Engine (Goblade) [P3.1]
- **`view/compiler.go`**: Implements a lightweight scanner/parser transpiler that translates Blade-like directives (`@if`, `@foreach`, `@yield`, `@section`, `@csrf`) into native Go `html/template` syntax.
- Pre-compiles templates and caches the output in `storage/framework/views/` using a SHA-256 hash of the content to ensure zero overhead in production.
- **`view/engine.go`**: The view engine wrapper that resolves view names (e.g. `auth.login`), triggers compilation, and executes the cached `html/template` binary.

## 2. Session Management & Cookies [P3.3, P3.4]
- **`session/store.go`**: Abstract interface for Session drivers (`Read`, `Write`, `Destroy`).
- **`session/driver_file.go`**: The default local development driver. Stores session payloads securely as JSON files inside `storage/framework/sessions/{id}`.
- **`session/driver_cookie.go`**: A stub driver ready for production cookie-based payload storage.
- **`session/manager.go`**: Handles retrieving, flashing, and regenerating session IDs.
- **`http/middleware/session.go`**: The `StartSession` middleware that extracts the `gow_session` cookie, initializes the `SessionManager`, injects it into the request `context.Context`, and saves it automatically upon response.

## 3. CSRF Protection [P3.7]
- **`http/middleware/csrf.go`**: The `VerifyCsrfToken` middleware designed to prevent Cross-Site Request Forgery.
- Automatically generates a token and stores it in the session.
- Validates all mutating requests (`POST`, `PUT`, `DELETE`) by checking the `X-CSRF-TOKEN` header or the `_token` form value against the session state. Aborts with `419` on mismatch.

## 4. File Storage / Filesystem [P3.5]
- **`storage/filesystem.go`**: A Flysystem-inspired abstraction layer for managing file storage.
- Includes `LocalFilesystem` driver to provide unified `Put()`, `Get()`, `Exists()`, and `Delete()` functionality over the local disk.

## 5. HTTP Client [P3.10]
- **`http/client/client.go`**: A fluent, chainable wrapper around Go's native `net/http` package.
- Supports adding headers (`WithHeader`), executing requests (`Get`, `Post`), and easily mapping response bodies to JSON structs (`res.JSON(&target)`).

## 6. View Composers [P3.2]
- **`view/composer.go`**: Implemented `ComposerRegistry` to bind data globally (`*`) or to specific views prior to rendering.

## 7. Cache System [P3.4]
- **`cache/store.go` & `cache/driver_memory.go`**: Interface defining `Get`, `Put`, `Increment`, `Forever`, `Forget`. Includes a thread-safe in-memory driver implementation (`sync.RWMutex`).

## 8. Validation & Form Requests [P3.5, P3.6]
- **`validation/validator.go`**: Rule-based execution engine handling validations like `required` and `email`.
- **`http/request/form_request.go`**: A `Validate` wrapper that handles form authorization and runs the `Validator` engine against request payloads automatically.

## 9. Localization (i18n) [P3.8]
- **`localization/translator.go`**: A JSON-based language file loader that exposes a `Translate` helper (`__('key')`), supporting dynamic parameter replacement.

## 10. Cookie Management [P3.9]
- **`cookie/manager.go`**: Provides a standalone manager utilizing AES-256-GCM and HMAC-SHA256 to securely encrypt and sign outgoing cookies.
- **`http/middleware/encrypt_cookies.go`**: Middleware that intercepts and decrypts incoming cookies automatically, except those explicitly listed in an exemption list.
