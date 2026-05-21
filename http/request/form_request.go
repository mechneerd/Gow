package request

import (
	"gow/validation"
	"net/http"
)

// FormRequest is an interface that combines validation rules and authorization.
type FormRequest interface {
	Authorize() bool
	Rules() map[string][]string
}

// Validate executes authorization and validation against the HTTP request.
func Validate(r *http.Request, form FormRequest) (map[string][]error, bool) {
	if !form.Authorize() {
		return nil, false // Unauthorized
	}

	r.ParseForm()
	data := make(map[string]any)
	for k, v := range r.Form {
		if len(v) > 0 {
			data[k] = v[0]
		}
	}

	validator := validation.NewValidator(data, form.Rules())
	errs := validator.Validate()
	
	if len(errs) > 0 {
		return errs, true // Validation failed, but authorized
	}

	return nil, true // Success
}
