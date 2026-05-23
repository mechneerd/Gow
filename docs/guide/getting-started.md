# Getting Started with GoW Framework (Complete Beginner Guide)

**For People With Zero Experience**

This guide is written for someone who has **never used Go**, **never built a web app**, and has **no previous experience with GoW**.

We will go from **downloading Go** all the way to running a **full web application** with:
- Backend in Go
- Database
- Authentication
- Modern frontend using **React** or **Vue**

---

## Important: Read This First

- This guide is **very detailed** on purpose. Every single step is explained.
- Copy and paste commands exactly as shown.
- If something doesn't work, read the "Troubleshooting" section under each step.
- We recommend using **Windows** with **PowerShell**, **macOS Terminal**, or **Linux Terminal**.

---

## Step 0: What is GoW?

GoW is a **web framework** for the Go programming language. It is heavily inspired by Laravel (a very popular PHP framework).

With GoW you can build:
- Websites
- Web applications (like SaaS products)
- APIs
- Admin panels

It helps you write less code and follow best practices.

---

## Step 1: Install Go (The Programming Language)

GoW is written in Go, so you must install Go first.

### 1.1 For Windows Users

1. Open your web browser and go to: https://go.dev/dl
2. Download the file that says **"Windows"** and ends with `.msi` (example: `go1.26.3.windows-amd64.msi`)
3. Double-click the downloaded `.msi` file
4. Click **Next** → **Next** → **Install**
5. When installation finishes, click **Finish**

**Verify installation:**

1. Press `Windows + S`, type `PowerShell`, and open it.
2. Type this command and press Enter:

```powershell
go version
```

You should see something like:
```
go version go1.26.3 windows/amd64
```

If you see an error, restart your computer and try again.

### 1.2 For macOS Users

1. Go to https://go.dev/dl
2. Download the `.pkg` file for macOS (example: `go1.26.3.darwin-amd64.pkg`)
3. Double-click the file and follow the installer.
4. Open **Terminal** (press `Command + Space`, type Terminal).

Run this command:

```bash
go version
```

You should see the version number.

### 1.3 For Linux Users (Ubuntu / Debian example)

Open Terminal and run:

```bash
wget https://go.dev/dl/go1.26.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.3.linux-amd64.tar.gz
```

Add Go to your PATH by editing your shell profile:

```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

Verify:

```bash
go version
```

---

## Step 2: Install Required Tools

### 2.1 Install Visual Studio Code (Recommended Editor)

1. Go to https://code.visualstudio.com
2. Download and install VS Code for your operating system.
3. Open VS Code.

**Recommended Extensions** (install these in VS Code):
- Go (by Go Team at Google)
- Go Nightly (optional)

### 2.2 Install Git (Optional but Recommended)

- Windows: https://git-scm.com/download/win
- macOS: Usually already installed. Check with `git --version`
- Linux: `sudo apt install git`

---

## Step 3: Create Your First GoW Project (Step by Step)

We will create a folder called `my-first-gow-app`.

### 3.1 Open Terminal / PowerShell

**Windows**: Press `Windows + S`, search for "PowerShell", open it.

**macOS**: Press `Command + Space`, search "Terminal".

**Linux**: Open your terminal.

### 3.2 Create the Project Folder

Run these commands one by one:

```bash
mkdir my-first-gow-app
cd my-first-gow-app
```

**What just happened?**
- `mkdir` = make directory (create folder)
- `cd` = change directory (go inside the folder)

### 3.3 Initialize Go Module

Run this command:

```bash
go mod init my-first-gow-app
```

This creates a file called `go.mod`. It tells Go that this folder is a Go project.

### 3.4 Add GoW Framework

Now install GoW:

```bash
go get gow@latest
```

This may take 30–60 seconds. You will see many lines downloading.

After it finishes, you should see a `go.sum` file created.

---

## Step 4: Create the Recommended Folder Structure

GoW works best with a specific folder structure (similar to Laravel).

Run these commands one by one:

```bash
mkdir -p app/http/controllers
mkdir -p app/models
mkdir -p app/http/middleware
mkdir -p routes
mkdir -p resources/views
mkdir -p config
mkdir -p bootstrap
mkdir -p database/migrations
mkdir -p storage/logs
mkdir -p public
```

**Explanation of folders:**

| Folder                    | Purpose |
|---------------------------|-------|
| `app/http/controllers`    | Your controllers (logic for pages) |
| `app/models`              | Your database models |
| `routes`                  | All your route definitions |
| `resources/views`         | Your HTML templates (Goblade) |
| `bootstrap`               | Application startup code |
| `public`                  | CSS, JS, images that browser can access |

---

## Step 5: Create Your First Running Application

### 5.1 Create `main.go`

In your project root (`my-first-gow-app`), create a new file called `main.go`.

**Using VS Code:**
1. Open VS Code
2. Click `File` → `Open Folder`
3. Select the `my-first-gow-app` folder
4. Create new file `main.go`

Paste this code:

```go
package main

import (
	"fmt"
	"net/http"

	"gow/http/router"
)

func main() {
	// Create a new router
	r := router.New()

	// Define a route for the homepage
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "<h1>Hello from GoW!</h1><p>Welcome to your first application.</p>")
	})

	// Define another route
	r.Get("/about", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "This is the about page.")
	})

	// Start the server on port 8080
	fmt.Println("Server is running at http://localhost:8080")
	fmt.Println("Press CTRL+C to stop the server")
	http.ListenAndServe(":8080", r)
}
```

### 5.2 Run the Application

Go back to your terminal (make sure you are inside the `my-first-gow-app` folder) and run:

```bash
go run main.go
```

You should see:

```
Server is running at http://localhost:8080
Press CTRL+C to stop the server
```

Now open your browser and go to:

→ **http://localhost:8080**

You should see a webpage with "Hello from GoW!"

Also try: **http://localhost:8080/about**

**Congratulations!** You have a running GoW web application.

To stop the server, press `CTRL + C` in the terminal.

---

## Step 6: Understanding What You Just Did (Important)

- `router.New()` creates a router that handles web requests.
- `r.Get("/", ...)` means "when someone visits the homepage, run this code".
- `http.ListenAndServe(":8080", r)` starts a web server on port 8080.

This is the foundation of every GoW application.

---

## Step 7: Using Better Project Structure (Recommended)

The `main.go` above is fine for learning, but real projects use a better structure.

We will now improve it.

### 7.1 Create `bootstrap/app.go`

Create the file `bootstrap/app.go` and paste:

```go
package bootstrap

import (
	"gow/foundation"
)

func NewApplication() *foundation.Application {
	app := foundation.NewApplication(".")
	app.Boot()
	return app
}
```

### 7.2 Create `routes/web.go`

Create `routes/web.go`:

```go
package routes

import "gow/http/router"

func RegisterWebRoutes(r *router.Router) {
	r.Get("/", func(w, req) {
		w.Write([]byte("Welcome to the improved GoW app!"))
	})
}
```

### 7.3 Update `main.go`

Replace your `main.go` with this cleaner version:

```go
package main

import (
	"fmt"
	"net/http"

	"my-first-gow-app/bootstrap"
	"my-first-gow-app/routes"

	"gow/http/router"
)

func main() {
	app := bootstrap.NewApplication()

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

## Step 8: Using Goblade Templates (Views)

Instead of writing HTML inside Go code, we use templates.

### 8.1 Create a View File

Create `resources/views/home.blade.go`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>My GoW App</title>
</head>
<body>
    <h1>Hello {{ .Name }}!</h1>
    <p>Welcome to GoW.</p>
</body>
</html>
```

### 8.2 Render the View from a Route

Update your route in `routes/web.go`:

```go
import (
	"gow/view"
	"gow/http/router"
)

func RegisterWebRoutes(r *router.Router) {
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		data := map[string]any{
			"Name": "John",
		}
		view.Make("home", data).Render(w)
	})
}
```

Now when you visit the homepage, it will show a proper HTML page.

---

## Step 9: Adding a Database (SQLite - Easiest for Beginners)

### 9.1 Create `.env` file

In your project root, create a file named `.env`:

```env
APP_NAME=MyGoWApp
APP_ENV=local
DB_CONNECTION=sqlite
DB_DATABASE=database.sqlite
```

### 9.2 Create the Database File

Run this in terminal:

```bash
touch database.sqlite
```

(For Windows PowerShell use: `New-Item database.sqlite`)

### 9.3 Basic Database Connection

We will use the built-in database manager in later steps.

For now, you can manually create tables using SQLite browser or SQL commands.

---

## Step 10: Moving to React or Vue (Advanced Frontend)

This is the part most beginners want.

### 10.1 React + Inertia Setup (Detailed)

#### Step 1: Create a frontend folder

```bash
mkdir -p resources/js/Pages
```

#### Step 2: Initialize Vite + React

```bash
npm create vite@latest resources/js -- --template react
cd resources/js
npm install
```

#### Step 3: Install Inertia

```bash
npm install @inertiajs/inertia @inertiajs/inertia-react
```

#### Step 4: In GoW, return Inertia response

In a controller or route:

```go
import "gow/http/inertia"

return inertia.Render(w, "Dashboard", map[string]any{
    "user": currentUser,
})
```

#### Step 5: Create React Component

Create `resources/js/Pages/Dashboard.jsx`:

```jsx
import React from 'react';

export default function Dashboard({ user }) {
    return (
        <div>
            <h1>Welcome, {user.name}</h1>
        </div>
    );
}
```

This is the modern way to build web apps with GoW + React.

---

## Step 11: Recommended Learning Path

After finishing this guide, follow this order:

1. Learn more about **Routing** (`docs/guide/routing.md`)
2. Deep dive into **ORM** (`docs/guide/orm.md`)
3. Add **Authentication** using Fortify
4. Choose one frontend:
   - Goblade (fastest)
   - React + Inertia
   - Vue + Inertia
5. Learn **Queues** and **Broadcasting**
6. Write **Tests**
7. Deploy your application

---

## Final Words

You have now gone from **zero** to having a running GoW application with the ability to use modern frontend frameworks.

Take your time. The most important thing is to type the commands yourself and understand what each step does.

If you get stuck at any point, read the error message carefully — Go usually gives good errors.

> **Want a shorter version?**  
> Check the [Quick Start - One Page Printable](quick-start-one-page.md) if you just want the fastest path.

---

**Next Recommended Action:**

Open the file `docs/guide/orm.md` and start learning how to work with the database properly.

Welcome to the GoW community!
