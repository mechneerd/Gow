# GoW Framework

![GoW Banner](https://via.placeholder.com/1200x300?text=GoW+Framework)

**GoW** is a modern, full-stack Go web framework that perfectly marries the uncompromised performance of Go with the breathtaking developer experience of Laravel.

## 🚀 Why GoW?

The Go ecosystem has historically prioritized minimal routing libraries and "bring-your-own-architecture" philosophies. GoW challenges this by providing a cohesive, batteries-included framework that helps you ship faster without sacrificing Go's legendary performance.

- **Developer Experience First**: Familiar, expressive APIs inspired by Laravel.
- **Goquent ORM**: A fluent query builder and Active Record implementation.
- **Goblade Templates**: Server-side rendering made beautiful with structural directives.
- **Artisan CLI**: Generate controllers, models, migrations, and serve your app instantly.
- **Batteries Included**: Routing, Middleware, Validation, Localization, Queues, Mail, and more.

## 📦 Quick Start

Ensure you have Go 1.24+ installed.

```bash
# Initialize a new GoW application
npx gow new my-app
cd my-app

# Start the development server
go run cmd/app/main.go
```

## 📖 Documentation

The full documentation is located in the [/docs](/docs) directory.

- [Installation & Setup](/docs/installation.md)
- [Routing](/docs/routing.md)
- [Goquent ORM](/docs/orm.md)

## 🧪 Testing

GoW ships with a fluent testing wrapper integrating `httptest` and `testify/assert`.

```go
func TestHealthCheck(t *testing.T) {
    tc := testing.NewTestCase(t, app.Router)
    
    tc.Get("/api/health").
       AssertStatus(200).
       AssertJson(map[string]any{"status": "up"})
}
```

## 🏛️ Project Status

GoW is currently heavily in development. See `CHANGELOG.md` for phase tracking.

## 📜 License

The GoW framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
