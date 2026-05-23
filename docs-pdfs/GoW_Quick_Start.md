# GoW Quick Start Guide (PDF Version)

**Goal**: Get a running GoW application in 10–20 minutes.

> This is a clean, PDF-optimized version of the Quick Start guide.  
> It includes tech stack recommendations (API, React, Vue, Goblade, HTMX, etc.).

For the full detailed version, see `GoW_Getting_Started_Detailed.md`.

---

## What Can You Build With GoW?

GoW is flexible. Here are the main ways people use it:

### 1. API Only (Most Common)
Build just the backend. Use any frontend you like:
- React, Vue, Svelte, Angular, Flutter, React Native, etc.

**Best for**: Mobile apps + Web frontends

### 2. Modern Web App (Recommended)
Use **Inertia.js** + React or Vue.

This gives you a modern frontend experience while using GoW as the backend.

**Best for**: SaaS products, dashboards, full web applications

### 3. Traditional Web App
Use GoW’s built-in templates (**Goblade**).

Similar to Laravel Blade.

**Best for**: Simple websites and admin panels

### 4. Lightweight Modern (Low JavaScript)
Use Goblade + **HTMX** + Alpine.js + Tailwind.

Very fast and performs well.

**Best for**: Fast development with minimal JavaScript

### 5. Livewire-style Components (Now Functional)
GoW has a working **Livewire equivalent** (`http/livewire`).

**Current Capabilities** (May 23, 2026):
- `wire:click="method"` → Call component methods
- `wire:model="property"` → Two-way binding on inputs
- `wire:submit="method"` → Handle form submissions
- Reactive updates → Component automatically re-renders when state changes
- Basic `wire:loading` support
- Lifecycle hook: `Mount()`

**Quick Example**:

```go
type Counter struct {
    livewire.BaseComponent
    Count int
}

func (c *Counter) Render() string {
    return fmt.Sprintf(`<div wire:id="%s">
        <h2>Count: %d</h2>
        <button wire:click="Increment">+1</button>
        <input type="number" wire:model="Count">
    </div>`, c.GetID(), c.Count)
}

func (c *Counter) Increment() { c.Count++ }
```

This gives you a real reactive component experience similar to Laravel Livewire. More directives (wire:ignore, better loading states, etc.) can be added in future updates.

**Quick Recommendation**:

| Goal                          | Recommended Stack                     |
|-------------------------------|---------------------------------------|
| Modern web app / SaaS         | GoW + Inertia + React or Vue          |
| API for mobile + web          | GoW (pure API)                        |
| Simple & fast website         | GoW + Goblade                         |
| Lightweight with little JS    | GoW + Goblade + HTMX                  |

---

## Step 1: Install Go

Download and install Go from: https://go.dev/dl

After installation, open your terminal and verify:

```bash
go version
```

You should see `go1.26.x` or higher.

---

## Step 2: Create a New Project (Fastest Way)

Open your terminal and run these commands one after another:

```bash
mkdir my-gow-app
cd my-gow-app
go mod init my-gow-app
go get gow@latest
```

Create the basic folders:

```bash
mkdir -p app/http/controllers routes resources/views bootstrap
```

---

## Step 3: Create Your First App

### 1. Create `main.go`

In the root of your project, create `main.go` and paste this:

```go
package main

import (
	"fmt"
	"net/http"

	"gow/http/router"
)

func main() {
	r := router.New()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "<h1>Hello from GoW!</h1><p>Your app is working.</p>")
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
```

### 2. Run the app

```bash
go run main.go
```

Open your browser at **http://localhost:8080**

You should see "Hello from GoW!"

**Success!** You now have a running GoW application.

---

## Step 4: Better Structure (Recommended)

Create these files for a cleaner setup:

### `routes/web.go`

```go
package routes

import "gow/http/router"

func RegisterWebRoutes(r *router.Router) {
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Welcome to GoW - Quick Start!"))
	})
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

### Update `main.go`

```go
package main

import (
	"fmt"
	"net/http"

	"my-gow-app/bootstrap"
	"my-gow-app/routes"

	"gow/http/router"
)

func main() {
	bootstrap.NewApplication()

	r := router.New()
	routes.RegisterWebRoutes(r)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
```

Run again:

```bash
go run main.go
```

---

## Step 5: Add a Simple View (Goblade)

Create `resources/views/welcome.blade.go`:

```html
<h1>Hello {{ .Name }}!</h1>
<p>This is your first GoW view.</p>
```

Update your route:

```go
import "gow/view"

r.Get("/", func(w http.ResponseWriter, req *http.Request) {
	view.Make("welcome", map[string]any{"Name": "Developer"}).Render(w)
})
```

---

## Step 6: Quick React + Inertia (Modern Frontend)

Want to use **React** quickly?

### 1. Setup Vite + React

```bash
npm create vite@latest resources/js -- --template react
cd resources/js
npm install
npm install @inertiajs/inertia @inertiajs/inertia-react
```

### 2. Return Inertia from GoW

In a route or controller:

```go
import "gow/http/inertia"

return inertia.Render(w, "Home", map[string]any{
    "message": "Hello from React + GoW",
})
```

### 3. Create React Page

Create `resources/js/Pages/Home.jsx`:

```jsx
export default function Home({ message }) {
    return <h1>{message}</h1>;
}
```

---

## Step 7: Quick Vue + Inertia

Same as React, just change the frontend:

```bash
npm create vite@latest resources/js -- --template vue
cd resources/js
npm install @inertiajs/inertia @inertiajs/inertia-vue3
```

Then use Vue components in `resources/js/Pages/`.

---

## Next Steps

Now that you have a working app, here’s what to learn next (in order):

1. **Routing** – `docs/guide/routing.md`
2. **ORM & Database** – `docs/guide/orm.md`
3. **Authentication** – `docs/guide/authentication.md`
4. **Views (Goblade)** – `docs/guide/views.md`
5. **Testing** – `docs/guide/testing.md`
6. **Socialite (OAuth)** – `docs/guide/socialite.md`

For complete beginners who want full explanations, read the full **[Getting Started Guide](getting-started.md)**.

---

**You now have a running GoW app with the option to use React or Vue.**

Happy coding!
