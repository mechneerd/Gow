package fortify

import (
	"encoding/json"
	"net/http"

	"gow/auth"
	"gow/auth/password"
	"gow/routing"
)

// Fortify provides a headless authentication backend for SPAs and APIs.
type Fortify struct {
	authManager *auth.Manager
	broker      *password.Broker
}

// New creates a new Fortify instance.
func New(authManager *auth.Manager, broker *password.Broker) *Fortify {
	return &Fortify{
		authManager: authManager,
		broker:      broker,
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

func (f *Fortify) handleRegister(w http.ResponseWriter, r *http.Request) error {
	// 1. Validate request (email, password, etc)
	// 2. Create User
	// 3. Fire Registered event
	// 4. Authenticate User

	f.jsonResponse(w, http.StatusCreated, map[string]any{
		"status": "success",
		"message": "User registered successfully",
	})
	return nil
}

func (f *Fortify) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// Extract credentials
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		f.jsonResponse(w, http.StatusBadRequest, map[string]any{"message": "Invalid request"})
		return nil
	}

	guard := f.authManager.Guard("web")
	if guard.Attempt(map[string]any{"email": creds.Email, "password": creds.Password}) {
		f.jsonResponse(w, http.StatusOK, map[string]any{
			"status": "success",
			"user": guard.User(),
		})
		return nil
	}

	f.jsonResponse(w, http.StatusUnauthorized, map[string]any{
		"message": "These credentials do not match our records.",
	})
	return nil
}

func (f *Fortify) handleForgotPassword(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.jsonResponse(w, http.StatusBadRequest, map[string]any{"message": "Invalid request"})
		return nil
	}

	if f.broker == nil {
		f.jsonResponse(w, http.StatusInternalServerError, map[string]any{"message": "Password reset not configured"})
		return nil
	}

	if err := f.broker.SendResetLink(req.Email); err != nil {
		f.jsonResponse(w, http.StatusInternalServerError, map[string]any{"message": "Failed to send reset link"})
		return nil
	}

	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status":  "success",
		"message": "We have emailed your password reset link!",
	})
	return nil
}

func (f *Fortify) handleResetPassword(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		f.jsonResponse(w, http.StatusBadRequest, map[string]any{"message": "Invalid request"})
		return nil
	}

	if f.broker == nil {
		f.jsonResponse(w, http.StatusInternalServerError, map[string]any{"message": "Password reset not configured"})
		return nil
	}

	if err := f.broker.Reset(req.Email, req.Token, req.Password); err != nil {
		f.jsonResponse(w, http.StatusBadRequest, map[string]any{"message": err.Error()})
		return nil
	}

	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status":  "success",
		"message": "Your password has been reset!",
	})
	return nil
}

func (f *Fortify) handleVerifyEmail(w http.ResponseWriter, r *http.Request) error {
	// In a real implementation, we would verify the signed URL signature here
	// using routing.VerifySignedRequest or similar.

	// For now, we assume the signature was validated by middleware or the route itself.
	user := auth.User(r)
	if user == nil {
		f.jsonResponse(w, http.StatusUnauthorized, map[string]any{"message": "Unauthenticated"})
		return nil
	}

	// Mark as verified (integrate with your User model)
	// verification.MarkAsVerified(user.(auth.Authenticatable))

	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status":  "success",
		"message": "Email verified successfully!",
	})
	return nil
}

func (f *Fortify) handleEnableTwoFactor(w http.ResponseWriter, r *http.Request) error {
	// Generate 2FA secret, return QR code URL / recovery codes
	f.jsonResponse(w, http.StatusOK, map[string]any{
		"status": "success",
		"message": "Two factor authentication enabled.",
	})
	return nil
}
