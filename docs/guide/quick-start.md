# Quick Start Guide — GoW Framework

**Goal**: Get a working GoW application running in **10–20 minutes**, even if you have little experience.

This is the **fast-track** version. For full explanations, see the [Getting Started Guide](getting-started.md).

> **Need it even shorter?**  
> See the [One-Page Printable Version](quick-start-one-page.md) — fits on a single printed page.

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
