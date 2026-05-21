package main

import (
	"fmt"
	"gow/http/router"
	"log"
	"net/http"
)

// This is the entry point for the canonical Phase 7 GoW Demo Blog.

func main() {
	// Initialize Router
	r := router.New()

	// Register Routes
	RegisterRoutes(r)

	fmt.Println("🚀 GoW Blog Demo starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func RegisterRoutes(r *router.Router) {
	// API Routes (Content Negotiation Example)
	api := r.Group("/api")
	api.Get("/posts", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id": 1, "title": "Welcome to GoW", "author": "Kimi"}]`))
	})

	// Web Routes
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
		<html>
			<head><title>GoW Blog</title></head>
			<body style="font-family: sans-serif; max-width: 800px; margin: 40px auto;">
				<h1>Welcome to the GoW Blog</h1>
				<p>This is the canonical example application demonstrating GoW's capabilities.</p>
				<hr>
				<h2>Latest Posts</h2>
				<article>
					<h3><a href="/posts/1">Welcome to GoW</a></h3>
					<p><em>By Kimi</em></p>
					<p>GoW is a modern, full-stack Go web framework...</p>
				</article>
			</body>
		</html>
		`
		w.Write([]byte(html))
	})

	r.Get("/posts/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := router.Param(req, "id")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf("<h1>Post %s Details</h1><p>Full post content goes here.</p><a href='/'>&larr; Back</a>", id)))
	})
}
