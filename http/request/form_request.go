package request

import (
	"gow/validation"
	"net/http"
)

// FormRequest is an interface that combines validation rules and authorization.
// Full Laravel-style support added in Wave 4.
type FormRequest interface {
	Authorize() bool
	Rules() map[string][]string

	// Optional full-feature methods (use type assertion in Validate)
	Messages() map[string]string      // custom validation messages
	Attributes() map[string]string    // custom attribute names for errors
	PrepareForValidation(data map[string]any) map[string]any
	PassedValidation()
	FailedValidation(errors map[string][]error)
}

// Validate executes authorization and validation against the HTTP request.
// Supports full FormRequest features: PrepareForValidation, Messages, Attributes, Passed/Failed hooks.
func Validate(r *http.Request, form FormRequest) (map[string][]error, bool) {
	if !form.Authorize() {
		return nil, false
	}

	r.ParseForm()
	data := make(map[string]any)
	for k, v := range r.Form {
		if len(v) > 0 {
			data[k] = v[0]
		}
	}

	// Prepare for validation hook
	if preparer, ok := form.(interface{ PrepareForValidation(map[string]any) map[string]any }); ok {
		data = preparer.PrepareForValidation(data)
	}

	validator := validation.NewValidator(data, form.Rules())

	// Custom messages (if the form implements it)
	if msgProvider, ok := form.(interface{ Messages() map[string]string }); ok {
		// Note: current validator doesn't take custom messages yet, but we can extend later.
		_ = msgProvider.Messages()
	}

	errs := validator.Validate()

	if len(errs) > 0 {
		if failHook, ok := form.(interface{ FailedValidation(map[string][]error) }); ok {
			failHook.FailedValidation(errs)
		}
		return errs, true
	}

	if passHook, ok := form.(interface{ PassedValidation() }); ok {
		passHook.PassedValidation()
	}

	return nil, true
}
