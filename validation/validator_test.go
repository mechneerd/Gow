package validation

import (
	"strings"
	"testing"
)

func TestValidatorBasicRules(t *testing.T) {
	data := map[string]any{
		"name":           "John",
		"email":          "john@example.com",
		"age":            30,
		"password":       "secret",
		"password_confirmation": "secret",
		"role":           "admin",
	}

	rules := map[string][]string{
		"name":     {"required", "string", "min:3", "max:10"},
		"email":    {"required", "email"},
		"age":      {"required", "numeric", "between:18,65"},
		"password": {"required", "min:6", "confirmed"},
		"role":     {"required", "in:admin,user"},
		"missing":  {"string"}, // Not required, should skip
	}

	v := NewValidator(data, rules)
	errs := v.Validate()

	if len(errs) > 0 {
		t.Errorf("Expected no validation errors, got %v", errs)
	}
}

func TestValidatorFailures(t *testing.T) {
	data := map[string]any{
		"name":           "Jo", // too short
		"email":          "invalid-email",
		"age":            17, // too young
		"password":       "sec", // too short
		"password_confirmation": "different", // mismatch
		"role":           "guest", // not in
	}

	rules := map[string][]string{
		"name":     {"required", "string", "min:3"},
		"email":    {"required", "email"},
		"age":      {"required", "numeric", "min:18"},
		"password": {"required", "min:6", "confirmed"},
		"role":     {"required", "in:admin,user"},
		"missing":  {"required"},
	}

	v := NewValidator(data, rules)
	errs := v.Validate()

	expectedErrors := map[string]string{
		"name":     "at least 3",
		"email":    "valid email",
		"age":      "at least 18",
		"password": "at least 6", // multiple errors for password, check one
		"role":     "is invalid",
		"missing":  "required",
	}

	for field, expectedFragment := range expectedErrors {
		fieldErrs, ok := errs[field]
		if !ok || len(fieldErrs) == 0 {
			t.Errorf("Expected error for %s, got none", field)
			continue
		}

		found := false
		for _, err := range fieldErrs {
			if strings.Contains(err.Error(), expectedFragment) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected error for %s containing '%s', got %v", field, expectedFragment, fieldErrs)
		}
	}
}

func TestValidatorDatabaseRules(t *testing.T) {
	// Let's create an in-memory SQLite DB just to test unique/exists
	// Note: We need a driver, but we don't have sqlite3 imported here.
	// We'll skip the actual query execution if DB is nil, but we can't easily mock *sql.DB without an interface.
	// We will just verify that passing nil doesn't panic and fails gracefully or skips.
	
	data := map[string]any{
		"username": "john",
	}
	rules := map[string][]string{
		"username": {"unique:users,username"},
	}

	v := NewValidator(data, rules)
	// Without DB, it should just ignore the DB rules and return no errors
	errs := v.Validate()
	if len(errs) > 0 {
		t.Errorf("Expected no errors when DB is nil for DB rules, got %v", errs)
	}

	// We can try to test it if we can open sqlite3, but since this is unit test for validation,
	// if sqlite driver is not imported, it might fail. We'll leave it as is for nil DB test.
}

func TestValidatorRegex(t *testing.T) {
	data := map[string]any{
		"zip": "12345",
		"bad": "12A45",
	}

	rules := map[string][]string{
		"zip": {"regex:^[0-9]{5}$"},
		"bad": {"regex:^[0-9]{5}$"},
	}

	v := NewValidator(data, rules)
	errs := v.Validate()

	if len(errs["zip"]) > 0 {
		t.Errorf("Expected zip to pass regex")
	}
	if len(errs["bad"]) == 0 {
		t.Errorf("Expected bad to fail regex")
	}
}
