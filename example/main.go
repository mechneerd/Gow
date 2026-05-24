package main

import (
	"fmt"
	"github.com/mechneerd/gow/foundation"
	gowhttp "github.com/mechneerd/gow/http"
	"github.com/mechneerd/gow/routing"
	"net/http"
)

func main() {
	// 1. Create Application
	app := foundation.NewApplication(".")
	
	// 2. Setup Router
	router := routing.NewRouter()
	
	// 3. Register Routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("Welcome to GoW Framework!"))
		return nil
	})

	router.Group("/api", func(r *routing.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				next.ServeHTTP(w, req)
			})
		})
		
		r.Get("/status", func(w http.ResponseWriter, req *http.Request) error {
			w.Write([]byte(`{"status": "ok"}`))
			return nil
		})
	})

	// 4. Create HTTP Kernel
	kernel := gowhttp.NewKernel(app, router)
	
	// 5. Add Global Middleware
	kernel.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Printf("[%s] %s\n", req.Method, req.URL.Path)
			next.ServeHTTP(w, req)
		})
	})

	// 6. Boot App
	app.Boot()

	// 7. Start Server
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", kernel); err != nil {
		panic(err)
	}
}

