# Skeleton Patches - GoW Starter Kits (All Kits)

This folder contains the exact files you need to add to the `gow-skeleton` repository so that **all** starter kits (`minimal`, `api`, `web`, `auth`) include:

- Beautiful modern landing page (`welcome.goblade`)
- Livewire.js client
- Example reactive Livewire component (Counter)
- Full no-page-reload interactivity on the landing page

---

## How to Apply These Patches

For **each** kit folder below, copy the contents into the corresponding template in the skeleton repo:

### 1. Minimal Kit
**Skeleton path:** `templates/minimal/`

Copy from here:
- `skeleton-patches/minimal/public/js/livewire.js`  
  → `templates/minimal/public/js/livewire.js`

- `skeleton-patches/minimal/resources/views/welcome.goblade`  
  → `templates/minimal/resources/views/welcome.goblade`

- `skeleton-patches/minimal/app/Livewire/Counter.go`  
  → `templates/minimal/app/Livewire/Counter.go`

### 2. API Kit
**Skeleton path:** `templates/api/`

Same files as above → into `templates/api/...`

### 3. Web Kit (default)
**Skeleton path:** `templates/web/`

Same files → into `templates/web/...`

### 4. Auth Kit (recommended)
**Skeleton path:** `templates/web-auth/`

Same files → into `templates/web-auth/...`

---

## 5. Route Registration (Critical Step)

You must also wire Livewire in the generated project's routes.

### Recommended place: `routes/web.go` (or main router setup)

Add this:

```go
import (
    "gow/http/livewire"
    "gow/routing"
)

func RegisterWebRoutes(router *routing.Router) {
    // ... your existing routes ...

    // Enable GoW Livewire (for reactive UI without reload)
    livewire.RegisterRoutes(router)
}
```

Call `RegisterWebRoutes(router)` during application bootstrap (usually in `bootstrap/app.go` or `main.go`).

---

## 6. Optional: Add Demo Route for Landing Page

The landing page expects `/livewire/counter` to exist for the demo.

`livewire.RegisterRoutes(router)` already includes a basic demo endpoint.

---

## After Applying

Once you push these changes to the skeleton repo, running:

```bash
gow new myapp --minimal --yes
# or
gow new myapp --api --yes
# or
gow new myapp --auth --yes
```

Will automatically give users the beautiful landing page with a working Livewire counter that updates **without any page reload**.

---

## Notes

- All templates now use `.goblade` (Go-native, no PHP naming).
- The landing page is the same across all kits for consistency.
- You can later customize per-kit (e.g. API kit can have a more API-focused landing page).

This completes the "include Livewire + beautiful landing page for all starter packs" task.
