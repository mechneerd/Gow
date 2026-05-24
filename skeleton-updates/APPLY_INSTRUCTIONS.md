# APPLY THESE CHANGES TO THE SKELETON REPO

This folder (`skeleton-updates`) contains the complete implementation for adding:

- Beautiful modern landing page (`welcome.goblade`)
- GoW Livewire support (reactive UI without page reload)
- Example Livewire component

To **all four starter kits** in the gow-skeleton repository.

---

## Step-by-Step Instructions

### 1. Clone the skeleton repo (if you haven't already)

```bash
git clone https://github.com/mechneerd/gow-skeleton.git
cd gow-skeleton
```

### 2. Copy the files from this folder

Copy the entire contents of `templates/` from this `skeleton-updates` folder into your local clone of the skeleton repo.

You can do it with:

**PowerShell (Windows):**
```powershell
Copy-Item -Path "D:\Go framework\skeleton-updates\templates\*" `
          -Destination "path\to\your\gow-skeleton\templates" `
          -Recurse -Force
```

**Git Bash / Linux / macOS:**
```bash
cp -r "D:/Go framework/skeleton-updates/templates/"* /path/to/your/gow-skeleton/templates/
```

This will update:
- `templates/minimal/`
- `templates/api/`
- `templates/web/`
- `templates/web-auth/`

### 3. Add Livewire Route Registration (Important!)

You must wire Livewire in the generated projects.

#### Best place: Inside the skeleton's bootstrap or route setup.

Most kits have a file like:
- `templates/<kit>/routes/web.go`
- or `templates/<kit>/bootstrap/app.go`

Add the following import and call:

```go
import (
    "gow/http/livewire"
    "gow/routing"
)

func RegisterWebRoutes(router *routing.Router) {
    // ... your existing routes here ...

    // Enable GoW Livewire (reactive components without page reload)
    livewire.RegisterRoutes(router)
}
```

Make sure `RegisterWebRoutes(router)` is called when the app boots.

### 4. Commit and Push

```bash
cd /path/to/your/gow-skeleton
git add .
git commit -m "feat: add beautiful landing page + GoW Livewire to all starter kits"
git push origin main
```

---

## What Was Added to Every Kit

| File | Description |
|------|-------------|
| `public/js/livewire.js` | GoW Livewire JavaScript client |
| `resources/views/welcome.goblade` | Beautiful landing page with working Livewire demo |
| `app/Livewire/Counter.go` | Example reactive component (Counter) |

The landing page uses **GoBlade** (`.goblade` extension) for a pure Go-native feel.

---

## After This Change

When anyone runs:

```bash
gow new myapp --minimal --yes
gow new myapp --api --yes
gow new myapp --yes
gow new myapp --auth --yes
```

They will get:
- A professional, modern landing page at `/`
- A working Livewire counter that updates **without any page reload**

---

## Notes

- This uses the improved Livewire system from the main GoW framework (as of May 24, 2026).
- No PHP naming remains — everything uses `.goblade`.
- The demo counter is intentionally simple so users can immediately see Livewire working.

You can customize per-kit later (e.g. make the API kit landing page more API-focused).

---

**Ready to apply.**  
Just copy the `templates/` folder as described above, add the route registration, and push.
