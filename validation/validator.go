package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Validator handles struct and rule-based validation.
type Validator struct {
	data  map[string]any
	rules map[string][]string
}

// NewValidator creates a new validator instance.
func NewValidator(data map[string]any, rules map[string][]string) *Validator {
	return &Validator{
		data:  data,
		rules: rules,
	}
}

// Validate executes the validation rules.
func (v *Validator) Validate() map[string][]error {
	errs := make(map[string][]error)

	for field, fieldRules := range v.rules {
		val, exists := v.data[field]
		for _, rule := range fieldRules {
			if err := v.applyRule(field, val, exists, rule); err != nil {
				errs[field] = append(errs[field], err)
			}
		}
	}

	return errs
}

func (v *Validator) applyRule(field string, value any, exists bool, rule string) error {
	if rule == "required" {
		if !exists || value == nil || isEmpty(value) {
			return errors.New("The " + field + " field is required.")
		}
	}
	if rule == "email" {
		if exists && !strings.Contains(fmt.Sprintf("%v", value), "@") {
			return errors.New("The " + field + " field must be a valid email address.")
		}
	}
	// Note: Min, Max, and other rules would be implemented here in a full framework.
	return nil
}

func isEmpty(val any) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	return false
}
