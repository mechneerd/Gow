# Middleware

Middleware provide a convenient mechanism for inspecting and filtering HTTP requests entering your application. For example, GoW includes a middleware that verifies the user of your application is authenticated.

## Defining Middleware

To create a new middleware, define a function that takes an `http.Handler` and returns an `http.Handler`.

```go
func EnsureTokenIsValid(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Global Middleware

If you want a middleware to run during every HTTP request to your application, list the middleware in your `app/http/kernel.go` global middleware stack.

## Route Middleware

If you would like to assign middleware to specific routes, you can invoke the `Use` method on a route or route group:

```go
router.Get("/profile", ProfileHandler).Use(EnsureTokenIsValid)

// Or on a group
api := router.Group("/api")
api.Use(EnsureTokenIsValid)
```

## Included Middleware

GoW includes several middleware out of the box:

- **TrimStrings**: Automatically trims whitespace from all incoming string fields in the request body.
- **TrustedProxies**: Safely inspects the `X-Forwarded-For` and `X-Forwarded-Proto` headers when your application sits behind a load balancer. Configure trusted proxy IPs in your `config/app.go` file.
