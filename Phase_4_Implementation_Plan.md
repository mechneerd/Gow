# Phase 4 Implementation Plan (Security & Authentication)

This document outlines the detailed implementation of Phase 4 within the GoW framework.

## 1. Password Hashing [P4.2]
- **`hashing/hasher.go`**: Defines the `Hasher` interface (`Make`, `Check`, `NeedsRehash`) for abstracting password storage.
- **`hashing/bcrypt.go`**: Implements the default Bcrypt driver utilizing `golang.org/x/crypto/bcrypt` with a configurable cost (default: 12).
- **`hashing/argon2.go`**: Implements Argon2id using OWASP recommended defaults (`m=64MB, t=3, p=2`). Features `NeedsRehash` to allow hot-upgrades on login when migrating drivers or changing cost factors.

## 2. Authentication System [P4.1]
- **`auth/manager.go`**: Resolves guards dynamically. Contains `SessionGuard` logic which manages state using standard HTTP sessions (`gow/session`).
- **`auth/guard.go`**: Abstraction interface defining `Check()`, `Attempt()`, `Login()`, `Logout()`.
- **`auth/provider.go`**: Abstraction interface defining `RetrieveByID` and `RetrieveByCredentials` to decouple the Auth system from Goquent (the ORM).
- **`http/middleware/auth.go`**: Intercepts requests and enforces authentication, returning `401 Unauthorized` for APIs or a `/login` redirect for Web requests.

## 3. Authorization (Gates & Policies) [P4.3]
- **`auth/access/gate.go`**: Implements Laravel-style Gates. Provides a `Gate.Define()` method to register closures globally, and `Gate.Allows()` / `Gate.Denies()` helper functions to execute them securely across the app.

## 4. Rate Limiting [P4.4]
- **`cache/rate_limiter.go`**: Introduces a Cache-backed `RateLimiter` utilizing `Increment` and `Forget` functions to manage attempt counts per key (e.g. per IP address or user ID).
- **`http/middleware/throttle.go`**: Wraps the `RateLimiter` to protect routes, automatically injecting `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `Retry-After` HTTP headers on every request.

## 5. Artisan Auth Scaffold [P4.5]
- **`cmd/artisan/make_auth.go`**: Exposes the `artisan make:auth` command. 
- Scaffolds a complete backend login/registration flow including `AuthController`, `LoginRequest`, `RegisterRequest`, `User` Model, and simple framework-agnostic `Goblade` templates (`login.gohtml` / `register.gohtml`), deliberately avoiding NPM/Tailwind complexity.

## 6. Encryption System [P4.6]
- **`encryption/encrypter.go`**: Utilizes AES-256-GCM to provide secure, randomized encryption and decryption of values. Returns base64 encoded strings ready for safe storage.

## 7. CORS & Security Middleware [P4.7, P4.8]
- **`http/middleware/cors.go`**: Implements Cross-Origin Resource Sharing handling, injecting `Access-Control-Allow-Origin`, handling preflight `OPTIONS` requests, and supporting configurable allowed methods/headers.
- **`http/middleware/security_headers.go`**: Automatically injects industry standard security headers (`X-Frame-Options`, `X-Content-Type-Options`, `X-XSS-Protection`, `Strict-Transport-Security`, `Content-Security-Policy`) into responses.
