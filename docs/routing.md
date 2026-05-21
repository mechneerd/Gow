# Routing

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
