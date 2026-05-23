> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Routing

> **Status**: ✅ Implemented


GoW provides an expressive and highly performant routing engine, utilizing a trie-based structure capable of parameter extraction and middleware grouping.

## Basic Routing

You will define all of your application routes in the `routes/web.go` or `routes/api.go` files.

```go
import "gow/http/router"

router.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to GoW!"))
})

// Route parameters
router.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := router.Param(r, "id")
    w.Write([]byte("User ID: " + id))
})
```

## Route Groups & Middleware

You can group routes to apply middleware or prefixes to multiple routes simultaneously.

```go
api := router.Group("/api")
api.Use(middleware.Throttle(60, 1)) // 60 requests per 1 minute

api.Get("/users", UserController.Index)
api.Post("/users", UserController.Store)
```

## Signed & Temporary Signed URLs

Generate URLs that cannot be tampered with:

```go
url, _ := router.SignedRoute("user.show", map[string]string{"id": "1"}, time.Now().Add(30*time.Minute))
```

Protect routes with middleware:

```go
router.Get("/verify", handler).Middleware(middleware.ValidateSignature(urlGenerator))
```

## Advanced Rate Limiting

```go
// Simple
router.Use(middleware.Throttle(60, 1))

// Named limiter (per user)
router.Use(middleware.ThrottleNamed("login", 5, 1))
```

Supports per-user throttling via `X-User-ID` header.

## Form Request Validation (Full Features)

```go
type LoginRequest struct {
    request.FormRequest
}

func (r *LoginRequest) Rules() map[string][]string {
    return map[string][]string{
        "email":    {"required", "email"},
        "password": {"required", "min:8"},
    }
}

func (r *LoginRequest) Messages() map[string]string { ... }
func (r *LoginRequest) PrepareForValidation(data map[string]any) map[string]any { ... }
```

Use in controllers:

```go
errs, ok := request.Validate(req, &LoginRequest{})
```

## Resource Routes

GoW provides automatic resource routing to quickly assign CRUD routes to a controller.

```go
// Generates standard routes: index, create, store, show, edit, update, destroy
router.Resource("photos", PhotoController{})

// For API-only controllers (excludes create and edit)
router.ApiResource("photos", PhotoController{})
```

## Route Macros

You can define custom Route Macros to extend the routing capabilities dynamically.

```go
router.Macro("customAction", func(r *router.Router, path string) {
    r.Get(path, CustomHandler)
})

router.customAction("/custom")
```

## Signed URLs

GoW can generate cryptographically signed URLs to prevent tampering, which is especially useful for actions like email verification.

```go
urlGenerator := routing.NewURLGenerator(router, appKey)

// Generate a signed URL valid for 30 minutes
signedURL, err := urlGenerator.SignedRoute("unsubscribe", map[string]string{"user": "1"}, time.Now().Add(30 * time.Minute))

// Verify a signed URL in a controller
isValid := urlGenerator.HasValidSignature(r)
```
