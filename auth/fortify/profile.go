package fortify

import (
	"encoding/json"
	"net/http"

	"github.com/mechneerd/gow/auth"
)

// ProfileController handles user profile updates.
type ProfileController struct {
	authManager *auth.Manager
}

// NewProfileController creates a new profile controller.
func NewProfileController(authManager *auth.Manager) *ProfileController {
	return &ProfileController{authManager: authManager}
}

// ProfileRequest represents a profile update request.
type ProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// RegisterProfileRoutes registers profile management routes.
func (pc *ProfileController) RegisterProfileRoutes(router interface {
	Post(path string, handler func(http.ResponseWriter, *http.Request) error)
	Get(path string, handler func(http.ResponseWriter, *http.Request) error)
	Put(path string, handler func(http.ResponseWriter, *http.Request) error)
	Delete(path string, handler func(http.ResponseWriter, *http.Request) error)
}) {
	router.Put("/api/user/profile", pc.handleUpdateProfile)
	router.Put("/api/user/password", pc.handleUpdatePassword)
	router.Delete("/api/user", pc.handleDeleteAccount)
}

// handleUpdateProfile handles profile information updates.
func (pc *ProfileController) handleUpdateProfile(w http.ResponseWriter, r *http.Request) error {
	var req ProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pc.jsonResponse(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid request body",
		})
		return nil
	}

	// Validate
	if req.Name == "" {
		pc.jsonResponse(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Validation failed",
			"errors": map[string][]string{
				"name": {"The name field is required"},
			},
		})
		return nil
	}

	// Get user from context
	userID := auth.UserID(r)
	if userID == "" {
		pc.jsonResponse(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthenticated",
		})
		return nil
	}

	// Update user (stub - in production would update database)
	pc.jsonResponse(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":    userID,
			"name":  req.Name,
			"email": req.Email,
		},
	})
	return nil
}

// handleUpdatePassword handles password updates.
func (pc *ProfileController) handleUpdatePassword(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		CurrentPassword string `json:"current_password"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirmation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pc.jsonResponse(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid request body",
		})
		return nil
	}

	// Validate
	if req.CurrentPassword == "" || req.Password == "" {
		pc.jsonResponse(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Validation failed",
			"errors": map[string][]string{
				"current_password": {"The current password field is required"},
				"password":         {"The password field is required"},
			},
		})
		return nil
	}

	if req.Password != req.PasswordConfirm {
		pc.jsonResponse(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Validation failed",
			"errors": map[string][]string{
				"password": {"The password confirmation does not match"},
			},
		})
		return nil
	}

	// Get user from context
	userID := auth.UserID(r)
	if userID == "" {
		pc.jsonResponse(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthenticated",
		})
		return nil
	}

	// Update password (stub - in production would update database)
	pc.jsonResponse(w, http.StatusOK, map[string]any{
		"message": "Password updated successfully",
	})
	return nil
}

// handleDeleteAccount handles account deletion.
func (pc *ProfileController) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// Get user from context
	userID := auth.UserID(r)
	if userID == "" {
		pc.jsonResponse(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthenticated",
		})
		return nil
	}

	// Delete user (stub - in production would delete from database)
	pc.jsonResponse(w, http.StatusOK, map[string]any{
		"message": "Account deleted successfully",
	})
	return nil
}

func (pc *ProfileController) jsonResponse(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
