package http

import (
	"encoding/json"
	"net/http"

	"github.com/mechneerd/gow/validation"
)

// FormRequest represents a form request with validation
type FormRequest interface {
	Rules() map[string][]string
	Messages() map[string]string
	Authorize(r *http.Request) bool
}

// ValidateRequest validates an incoming request against form request rules
func ValidateRequest(r *http.Request, request FormRequest) (map[string]any, error) {
	// Parse request body
	var data map[string]any
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return nil, err
		}
	} else {
		// Parse query parameters
		data = make(map[string]any)
		for key, values := range r.URL.Query() {
			if len(values) == 1 {
				data[key] = values[0]
			} else {
				data[key] = values
			}
		}
	}

	// Get validation rules
	rules := request.Rules()
	_ = request.Messages() // Messages available for future use

	// Run validation
	v := validation.NewValidator(data, rules)
	errs := v.Validate()

	if len(errs) > 0 {
		validationErrors := make(map[string][]string)
		for field, errors := range errs {
			errStrs := make([]string, len(errors))
			for i, err := range errors {
				errStrs[i] = err.Error()
			}
			validationErrors[field] = errStrs
		}
		return nil, &ValidationError{Errors: validationErrors}
	}

	return data, nil
}

// ValidationError represents validation errors
type ValidationError struct {
	Errors map[string][]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

// AuthorizeFormRequest checks if the request is authorized
func AuthorizeFormRequest(r *http.Request, request FormRequest) bool {
	return request.Authorize(r)
}

// HandleValidationErrors sends a validation error response
func HandleValidationErrors(w http.ResponseWriter, err error) {
	if validationErr, ok := err.(*ValidationError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "The given data was invalid.",
			"errors":  validationErr.Errors,
		})
		return
	}
	http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
}
