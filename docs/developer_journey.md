# The Developer Journey

Here is how a developer will experience and use the **GoW Framework** from start to finish to build a web application:

### 1. Installation & Scaffolding
The developer's journey begins in their terminal. They don't need to manually create directories or copy boilerplate code. Instead, they run:
```bash
npx gow new my-awesome-app
```
This CLI tool instantly scaffolds a brand-new GoW project. It generates the `app`, `config`, `routes`, and `database` directories. It sets up the `main.go` entry point and creates a fresh `.env` file with a secure `APP_KEY`. The developer simply `cd my-awesome-app` and they are ready to go.

### 2. Configuration
Before writing code, the developer opens the `.env` file. Here, they configure their database credentials (e.g., SQLite, MySQL, or PostgreSQL), set up their cache/queue driver (switching from `memory` to `redis`), and add SMTP credentials for sending emails. Because GoW uses standard `.env` files, it feels instantly familiar.

### 3. Creating Models & Migrations
The developer wants to build a blog, so they need a `Post` model. Instead of creating files manually, they use the Artisan CLI:
```bash
go run artisan make:model Post -m
```
This single command generates two things:
- A `Post` struct in `app/models/post.go` embedded with Goquent ORM tags.
- A database migration file in `database/migrations/`. 

The developer opens the migration file, uses the fluent schema builder to add a `title` and `body` column, and runs `go run artisan migrate`. The database table is instantly created.

### 4. Routing & Controllers
Next, the developer needs to create endpoints. They open `routes/web.go` and use the expressive routing API to map URLs to controller actions:
```go
router.Get("/posts", PostController.Index)
router.Post("/posts", PostController.Store)
```
To quickly generate that controller, they run `go run artisan make:controller PostController`. Inside the controller's `Index` method, they use the **Goquent ORM** to fetch data effortlessly:
```go
posts := orm.Table("posts").With("Author").Get()
```
They didn't have to write raw SQL or worry about N+1 query problems because Goquent handles the eager loading automatically.

### 5. Views & UI (Goblade)
To show the posts to the user, the developer uses the **Goblade** template engine. In the controller, they return a view:
```go
view.Make(w, "posts.index", map[string]any{"posts": posts})
```
Inside `resources/views/posts/index.gohtml`, they write beautiful, logic-injected HTML using Blade directives:
```html
@extends('layouts.app')

@section('content')
    @foreach(posts)
        <div @class(['highlight': .IsFeatured])>
            <h1>{{ .Title }}</h1>
        </div>
    @endforeach
@endsection
```
The Goblade engine compiles this down to pure, highly-performant Go templates behind the scenes.

### 6. Authentication & Security
The developer realizes they need users to log in before creating posts. Because GoW has "Batteries Included," they don't have to build auth from scratch. 
They register the **Fortify** backend in their router, which instantly exposes `/api/login` and `/api/register` endpoints. They add the `auth` middleware to their post-creation routes:
```go
router.Post("/posts", PostController.Store).Use(middleware.Auth)
```
They also define a **Policy** to ensure users can only edit their *own* posts, and use the `@can('edit', post)` directive directly in their Goblade views to hide the edit button from unauthorized users.

### 7. Background Jobs (Queues)
When a user publishes a post, the app needs to send a notification email to all subscribers. Sending emails during the HTTP request is slow, so the developer dispatches a background job using GoW's Queue system:
```go
queue.Push("SendSubscriberEmails", postData)
```
In a separate terminal tab, they run `go run artisan queue:work`. This background worker seamlessly picks up the job and dispatches the emails using GoW's built-in Mail and Notification managers, ensuring the web request remains lightning fast.

### 8. Testing
Before deploying, the developer writes tests using GoW’s fluent testing wrapper. They don't have to spin up a live web server or mock databases manually. They just write:
```go
tc := testing.NewTestCase(t, router)
tc.Post("/posts", map[string]string{"title": "My Post"}).AssertStatus(201)
tc.AssertDatabaseHas("posts", map[string]any{"title": "My Post"})
```
The tests execute incredibly fast using Go's native `go test` command.

### 9. Deployment
Finally, it’s time to deploy. Because this is a Go application, the entire framework, views, and business logic compile down to a **single, statically linked binary executable**. 
The developer runs `go build -o server cmd/app/main.go`. They drop that single binary onto a Linux server (or into a minimal Docker container), set their production `.env` variables, and start the server. 

Because it's powered by Go, it uses a fraction of the RAM of a PHP/Laravel or Node.js application, handles thousands of concurrent requests with ease, and provides a breathtakingly fast experience for the end-user.
