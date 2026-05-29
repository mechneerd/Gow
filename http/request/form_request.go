package request

import (
	"encoding/json"
	"github.com/mechneerd/gow/validation"
	"net/http"
	"strings"
)

// FormRequest is the base interface for form validation requests.
// Only Authorize() and Rules() are required. Other methods are optional
// and checked via type assertions.
type FormRequest interface {
	Authorize() bool
	Rules() map[string][]string
}

// Optional interfaces for extended FormRequest features.
// These are checked via type assertions in Validate().

// FormRequestMessages provides custom validation messages.
type FormRequestMessages interface {
	Messages() map[string]string
}

// FormRequestAttributes provides custom attribute names for errors.
type FormRequestAttributes interface {
	Attributes() map[string]string
}

// FormRequestPrepare allows modifying data before validation.
type FormRequestPrepare interface {
	PrepareForValidation(data map[string]any) map[string]any
}

// FormRequestPassed is called after successful validation.
type FormRequestPassed interface {
	PassedValidation()
}

// FormRequestFailed is called after failed validation.
type FormRequestFailed interface {
	FailedValidation(errors map[string][]error)
}

// Validate executes authorization and validation against the HTTP request.
// Supports both form-urlencoded and JSON request bodies.
func Validate(r *http.Request, form FormRequest) (map[string][]error, bool) {
	if !form.Authorize() {
		return nil, false
	}

	data := make(map[string]any)

	// Parse based on Content-Type
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		// Parse JSON body
		if r.Body != nil {
			var jsonData map[string]any
			if err := json.NewDecoder(r.Body).Decode(&jsonData); err == nil {
				data = jsonData
			}
		}
	} else {
		// Parse form data
		r.ParseForm()
		for k, v := range r.Form {
			if len(v) > 0 {
				data[k] = v[0]
			}
		}
	}

	// Prepare for validation hook
	if preparer, ok := form.(FormRequestPrepare); ok {
		data = preparer.PrepareForValidation(data)
	}

	validator := validation.NewValidator(data, form.Rules())

	// Custom messages (if the form implements it)
	if msgProvider, ok := form.(FormRequestMessages); ok {
		_ = msgProvider.Messages()
	}

	errs := validator.Validate()

	if len(errs) > 0 {
		if failHook, ok := form.(FormRequestFailed); ok {
			failHook.FailedValidation(errs)
		}
		return errs, true
	}

	if passHook, ok := form.(FormRequestPassed); ok {
		passHook.PassedValidation()
	}

	return nil, true
}

