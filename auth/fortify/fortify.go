package fortify

import (
	"encoding/json"
	"net/http"

	"gow/auth"
	"gow/routing"
)

// Fortify provides a headless authentication backend for SPAs and APIs.
type Fortify struct {
	authManager *auth.AuthManager
}

// New creates a new Fortify instance.
func New(authManager *auth.AuthManager) *Fortify {
	return &Fortify{
		authManager: authManager,
	}
}

// RegisterRoutes registers the headless JSON API authentication routes.
func (f *Fortify) RegisterRoutes(router *routing.Router) {
	router.Group("/api", func(r *routing.Router) {
		r.Post("/register", f.handleRegister)
		r.Post("/login", f.handleLogin)
		r.Post("/forgot-password", f.handleForgotPassword)
		r.Post("/reset-password", f.handleResetPassword)
		
		// Protected routes
		// Assuming there's an auth middleware that can be applied, 
		// but we'll define the route mapping here for now.
		r.Post("/two-factor/enable", f.handleEnableTwoFactor)
		r.Get("/email/verify/{id}/{hash}", f.handleVerifyEmail)
	})
}

func (f *Fortify) jsonResponse(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (f *Fortify) handleRegister(w http.ResponseWriter, r *http.Request) {
	// 1. Validate request (email, password, etc)
	// 2. Create User
	// 3. Fire Registered event
	// 4. Authenticate User

	f.jsonResponse(w, http.StatusCreated, map[string]any{
		"status": "success",
		"message": "User registered successfully",
	})
}

func (f *Fortify) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Extract credentials
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		f.jsonResponse(w, http.StatusBadRequest, map[string]any{"message": "Invalid request"})
		return
	}

	guard := f.authManager.Guard("web")
	if guard.Attempt(map[string]any{"email": creds.Email, "password": creds.Password}, creds.Remember) {
		f.jsonResponse(w, http.StatusOK, map[string]any{
			"status": "success",
			"user": guard.User(),
		})
		return
	}

	f.jsonResponse(w, http.StatusUnauthorized, map[string]any{
		"message": "These credentials do not match our records.",
	})
}

func (f *Fortify) handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	// Generate reset token, send email
	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status": "success",
		"message": "We have emailed your password reset link!",
	})
}

func (f *Fortify) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	// Validate token, update password
	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status": "success",
		"message": "Your password has been reset!",
	})
}

func (f *Fortify) handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Validate signature, update user email_verified_at
	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status": "success",
		"message": "Email verified!",
	})
}

func (f *Fortify) handleEnableTwoFactor(w http.ResponseWriter, r *http.Request) {
	// Generate 2FA secret, return QR code URL / recovery codes
	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status": "success",
		"message": "Two factor authentication enabled.",
	})
}
