# GoW Quick Start — One Page (Printable)

**Goal**: Get a running GoW app in under 15 minutes.

---

## 1. Install Go

Download from https://go.dev/dl

Verify in terminal:

```bash
go version
```

---

## 2. Create Project

```bash
mkdir my-gow-app && cd my-gow-app
go mod init my-gow-app
go install github.com/mechneerd/gow/cmd/gow@latest
mkdir -p routes bootstrap resources/views
```

---

## 3. Create Files

### `main.go`

```go
package main

import (
	"fmt"
	"net/http"
	"gow/http/router"
	"my-gow-app/routes"
	"my-gow-app/bootstrap"
)

func main() {
	bootstrap.NewApplication()
	r := router.New()
	routes.RegisterWebRoutes(r)
	fmt.Println("Running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
```

### `bootstrap/app.go`

```go
package bootstrap
import "gow/foundation"
func NewApplication() *foundation.Application {
	app := foundation.NewApplication(".")
	app.Boot()
	return app
}
```

### `routes/web.go`

```go
package routes
import "gow/http/router"
func RegisterWebRoutes(r *router.Router) {
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello from GoW!"))
	})
}
```

---

## 4. Run

```bash
go run main.go
```

Open browser → **http://localhost:8080**

---

## 5. Optional: Add a View

Create `resources/views/home.blade.go`:

```html
<h1>Hello {{ .Name }}!</h1>
```

Update route:

```go
import "gow/view"

r.Get("/", func(w, req) {
	view.Make("home", map[string]any{"Name": "World"}).Render(w)
})
```

---

## 6. React / Vue (Fast)

```bash
npm create vite@latest resources/js -- --template react   # or vue
cd resources/js && npm install @inertiajs/inertia @inertiajs/inertia-react
```

Return from GoW:

```go
import "gow/http/inertia"
return inertia.Render(w, "Home", map[string]any{"message": "Hi"})
```

---

## Next Steps

- Full Guide: `getting-started.md`
- ORM: `guide/orm.md`
- Auth: `guide/authentication.md`
- Socialite: `guide/socialite.md`

---

**That's it. You now have a running GoW application.**

---

**Need more details?**  
Full guides:  
- [Getting Started (Detailed)](getting-started.md)  
- [Quick Start (Normal)](quick-start.md)  
- [Upgrade Guide](UPGRADE.md)
