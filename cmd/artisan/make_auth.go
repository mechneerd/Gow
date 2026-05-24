package artisan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var MakeAuthCmd = &cobra.Command{
	Use:   "make:auth",
	Short: "Scaffold basic login and registration views and routes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scaffolding Authentication...")

		// Generate AuthController
		createFile("app/Http/Controllers/AuthController.go", authControllerStub)

		// Generate Form Requests
		createFile("app/Http/Requests/LoginRequest.go", loginRequestStub)
		createFile("app/Http/Requests/RegisterRequest.go", registerRequestStub)

		// Generate User Model
		createFile("app/Models/User.go", userModelStub)

		// Generate Views
		createFile("resources/views/auth/login.gohtml", loginViewStub)
		createFile("resources/views/auth/register.gohtml", registerViewStub)

		fmt.Println("Auth scaffolding complete. Please manually register your routes in routes/web.go.")
	},
}

func createFile(path, content string) {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("File %s already exists. Skipping.\n", path)
		return
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating %s: %s\n", path, err)
		return
	}
	fmt.Printf("Created: %s\n", path)
}

const authControllerStub = `package controllers

import (
	"net/http"
)

type AuthController struct {}

func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	// Render resources/views/auth/login.gohtml
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// Validate LoginRequest, Attempt auth, redirect
}

func (c *AuthController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	// Render resources/views/auth/register.gohtml
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	// Validate RegisterRequest, Hash password, create User, log in, redirect
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Logout user, redirect
}
`

const loginRequestStub = `package requests

import "gow/http/request"

type LoginRequest struct {}

func (r *LoginRequest) Authorize() bool {
	return true
}

func (r *LoginRequest) Rules() map[string][]string {
	return map[string][]string{
		"email":    {"required", "email"},
		"password": {"required"},
	}
}
`

const registerRequestStub = `package requests

import "gow/http/request"

type RegisterRequest struct {}

func (r *RegisterRequest) Authorize() bool {
	return true
}

func (r *RegisterRequest) Rules() map[string][]string {
	return map[string][]string{
		"name":     {"required"},
		"email":    {"required", "email"},
		"password": {"required"}, // Add min:8, confirmed etc. later
	}
}
`

const userModelStub = `package models

type User struct {
	ID       string ` + "`db:\"id\"`" + `
	Name     string ` + "`db:\"name\"`" + `
	Email    string ` + "`db:\"email\"`" + `
	Password string ` + "`db:\"password\"`" + `
}

func (u *User) GetAuthIdentifier() string {
	return u.ID
}

func (u *User) GetAuthPassword() string {
	return u.Password
}
`

const loginViewStub = `<!-- resources/views/auth/login.gohtml -->
<!DOCTYPE html>
<html>
<head><title>Login</title></head>
<body>
    <h2>Login</h2>
    <form method="POST" action="/login">
        @csrf
        <div>
            <label>Email:</label>
            <input type="email" name="email" required>
        </div>
        <div>
            <label>Password:</label>
            <input type="password" name="password" required>
        </div>
        <button type="submit">Login</button>
    </form>
</body>
</html>
`

const registerViewStub = `<!-- resources/views/auth/register.gohtml -->
<!DOCTYPE html>
<html>
<head><title>Register</title></head>
<body>
    <h2>Register</h2>
    <form method="POST" action="/register">
        @csrf
        <div>
            <label>Name:</label>
            <input type="text" name="name" required>
        </div>
        <div>
            <label>Email:</label>
            <input type="email" name="email" required>
        </div>
        <div>
            <label>Password:</label>
            <input type="password" name="password" required>
        </div>
        <button type="submit">Register</button>
    </form>
</body>
</html>
`
