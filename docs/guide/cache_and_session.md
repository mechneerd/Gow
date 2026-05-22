> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Cache & Session

> **Status**: ✅ Implemented


GoW provides an expressive, unified API for various caching and session backends.

## Cache

The Cache manager supports multiple stores. By default, GoW is configured to use the `memory` store, but you can easily swap this to `redis` or `database` in your `config/cache.go` file.

### Basic Usage

```go
// Store an item in the cache for 10 minutes
cache.Put("key", "value", 10 * time.Minute)

// Retrieve an item
value, err := cache.Get("key")

// Retrieve an item or store a default value if it doesn't exist
value, err := cache.Remember("users", 60 * time.Minute, func() any {
    return orm.Table("users").Get()
})
```

## Session

HTTP driven applications are stateless. Sessions provide a way to store information about the user across multiple requests.

GoW supports several session drivers, including `cookie`, `database`, and `redis`.

### Basic Usage

```go
// Store data in the session
r.Session().Put("key", "value")

// Retrieve data from the session
value := r.Session().Get("key")

// Flash data to the session (only available during the subsequent request)
r.Session().Flash("status", "Profile updated!")
```
