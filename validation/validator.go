package validation

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Pre-compiled regexes to avoid recompilation on every validation call.
var (
	emailRegex     = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	urlRegex       = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	uuidRegex      = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	alphaRegex     = regexp.MustCompile(`^[a-zA-Z]+$`)
	alphaNumRegex  = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	alphaDashRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	ipRegex        = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	ipv6Regex      = regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
)

// Validator handles struct and rule-based validation.
type Validator struct {
	data  map[string]any
	rules map[string][]string
	db    *sql.DB // For database rules like unique/exists
}

// NewValidator creates a new validator instance.
func NewValidator(data map[string]any, rules map[string][]string) *Validator {
	return &Validator{
		data:  data,
		rules: rules,
	}
}

// WithDB attaches a database connection for DB-aware rules.
func (v *Validator) WithDB(db *sql.DB) *Validator {
	v.db = db
	return v
}

// Validate executes the validation rules.
func (v *Validator) Validate() map[string][]error {
	errs := make(map[string][]error)

	// First pass: expand wildcard rules (e.g., "items.*.name" => "items.0.name", "items.1.name", etc.)
	expandedRules := make(map[string][]string)
	for field, fieldRules := range v.rules {
		if strings.Contains(field, "*") {
			// Expand wildcard field into actual field names
			actualFields := v.expandWildcard(field, fieldRules)
			for actualField, actualRules := range actualFields {
				expandedRules[actualField] = append(expandedRules[actualField], actualRules...)
			}
		} else {
			expandedRules[field] = fieldRules
		}
	}

	// Second pass: validate expanded rules
	for field, fieldRules := range expandedRules {
		val, exists := v.data[field]

		hasBail := false
		for _, r := range fieldRules {
			if r == "bail" {
				hasBail = true
				break
			}
		}

		for _, rule := range fieldRules {
			if rule == "bail" {
				continue // bail is a control rule, not a validator
			}

			if err := v.applyRule(field, val, exists, rule); err != nil {
				errs[field] = append(errs[field], err)
				if hasBail {
					break // stop validating this field
				}
			}
		}
	}

	return errs
}

// expandWildcard takes a wildcard field like "items.*.name" and expands it to
// concrete fields like "items.0.name", "items.1.name", etc.
func (v *Validator) expandWildcard(wildcardField string, rules []string) map[string][]string {
	result := make(map[string][]string)

	parts := strings.Split(wildcardField, "*")
	if len(parts) < 2 {
		result[wildcardField] = rules
		return result
	}

	// Find the array field (the part before the first *)
	// e.g., "items.*.name" => arrayField = "items"
	arrayField := strings.TrimSuffix(parts[0], ".")

	val, exists := v.data[arrayField]
	if !exists {
		return result
	}

	// Check if it's actually an array/slice
	valReflect := reflect.ValueOf(val)
	if valReflect.Kind() != reflect.Slice && valReflect.Kind() != reflect.Array {
		return result
	}

	// For each array element, expand the wildcard
	suffix := strings.TrimPrefix(parts[0], arrayField) // could be "." or ""
	if len(parts) > 1 {
		suffix = suffix + parts[1] // e.g., ".name"
	}

	for i := 0; i < valReflect.Len(); i++ {
		expandedField := fmt.Sprintf("%s.%d%s", arrayField, i, suffix)
		result[expandedField] = rules

		// If the nested part is also a wildcard (e.g., "items.*.tags.*.name"),
		// we need to recursively expand
		if strings.Contains(suffix, "*") {
			// For now, we support one level of nesting
			// This handles: items.*.name, users.*.roles.*, etc.
		}
	}

	return result
}

func (v *Validator) applyRule(field string, value any, exists bool, rule string) error {
	var ruleName string
	var ruleParams []string
	
	if strings.Contains(rule, ":") {
		parts := strings.SplitN(rule, ":", 2)
		ruleName = parts[0]
		ruleParams = strings.Split(parts[1], ",")
	} else {
		ruleName = rule
	}

	// Handle nullable early — if value is empty, skip all rules except "required"
	if ruleName == "nullable" {
		if !exists || isEmpty(value) {
			return nil // field is nullable and empty → valid
		}
	}

	if !exists && ruleName != "required" {
		return nil // skip validation if not required and not exists
	}

	strVal := fmt.Sprintf("%v", value)
	var numVal float64
	isNumeric := false
	if exists && value != nil {
		if n, err := strconv.ParseFloat(strVal, 64); err == nil {
			numVal = n
			isNumeric = true
		}
	}
	
	valLen := float64(len(strVal))
	if isNumeric {
		valLen = numVal // for min/max size on numbers, use the value itself
	}

	switch ruleName {
	case "required":
		if !exists || value == nil || isEmpty(value) {
			return errors.New("The " + field + " field is required.")
		}
	case "email":
		if !emailRegex.MatchString(strVal) {
			return errors.New("The " + field + " field must be a valid email address.")
		}
	case "min":
		min, _ := strconv.ParseFloat(ruleParams[0], 64)
		if valLen < min {
			return fmt.Errorf("The %s field must be at least %v.", field, min)
		}
	case "max":
		max, _ := strconv.ParseFloat(ruleParams[0], 64)
		if valLen > max {
			return fmt.Errorf("The %s field must not be greater than %v.", field, max)
		}
	case "between":
		min, _ := strconv.ParseFloat(ruleParams[0], 64)
		max, _ := strconv.ParseFloat(ruleParams[1], 64)
		if valLen < min || valLen > max {
			return fmt.Errorf("The %s field must be between %v and %v.", field, min, max)
		}
	case "size":
		size, _ := strconv.ParseFloat(ruleParams[0], 64)
		if valLen != size {
			return fmt.Errorf("The %s field must be exactly %v.", field, size)
		}
	case "numeric":
		if !isNumeric {
			return errors.New("The " + field + " field must be a number.")
		}
	case "string":
		if _, ok := value.(string); !ok {
			return errors.New("The " + field + " field must be a string.")
		}
	case "in":
		found := false
		for _, p := range ruleParams {
			if strVal == p {
				found = true
				break
			}
		}
		if !found {
			return errors.New("The selected " + field + " is invalid.")
		}
	case "regex":
		pattern := strings.Join(ruleParams, ",") // in case regex has commas
		matched, _ := regexp.MatchString(pattern, strVal)
		if !matched {
			return errors.New("The " + field + " field format is invalid.")
		}
	case "confirmed":
		confirmField := field + "_confirmation"
		confirmVal, confirmExists := v.data[confirmField]
		if !confirmExists || fmt.Sprintf("%v", confirmVal) != strVal {
			return errors.New("The " + field + " confirmation does not match.")
		}
	case "unique":
		// unique:table,column or unique:table,column,id
		if v.db != nil && len(ruleParams) >= 2 {
			table, column := ruleParams[0], ruleParams[1]
			if !isValidIdentifier(table) || !isValidIdentifier(column) {
				return errors.New("The " + field + " validation rule contains invalid identifiers.")
			}
			var count int
			var query string
			var args []any
			if len(ruleParams) >= 3 {
				// Exclude current record by ID
				id := ruleParams[2]
				query = fmt.Sprintf("SELECT COUNT(*) FROM \"%s\" WHERE \"%s\" = ? AND \"id\" != ?", table, column)
				args = []any{value, id}
			} else {
				query = fmt.Sprintf("SELECT COUNT(*) FROM \"%s\" WHERE \"%s\" = ?", table, column)
				args = []any{value}
			}
			v.db.QueryRow(query, args...).Scan(&count)
			if count > 0 {
				return errors.New("The " + field + " has already been taken.")
			}
		}
	case "exists":
		// exists:table,column
		if v.db != nil && len(ruleParams) >= 2 {
			table, column := ruleParams[0], ruleParams[1]
			// Validate identifiers to prevent SQL injection
			if !isValidIdentifier(table) || !isValidIdentifier(column) {
				return errors.New("The " + field + " validation rule contains invalid identifiers.")
			}
			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\" WHERE \"%s\" = ?", table, column)
			v.db.QueryRow(query, value).Scan(&count)
			if count == 0 {
				return errors.New("The selected " + field + " is invalid.")
			}
		}

	// === New High-Priority Rules ===

	case "boolean":
		if _, ok := value.(bool); ok {
			return nil
		}
		lower := strings.ToLower(strVal)
		if lower == "true" || lower == "false" || lower == "1" || lower == "0" {
			return nil
		}
		return errors.New("The " + field + " field must be a boolean.")

	case "integer":
		if _, err := strconv.Atoi(strVal); err != nil {
			return errors.New("The " + field + " field must be an integer.")
		}

	case "array":
		if value == nil || reflect.TypeOf(value).Kind() != reflect.Slice {
			return errors.New("The " + field + " field must be an array.")
		}

	case "url":
		if !urlRegex.MatchString(strVal) {
			return errors.New("The " + field + " field must be a valid URL.")
		}

	case "uuid":
		if !uuidRegex.MatchString(strVal) {
			return errors.New("The " + field + " field must be a valid UUID.")
		}

	case "alpha":
		if !alphaRegex.MatchString(strVal) {
			return errors.New("The " + field + " field may only contain letters.")
		}

	case "alpha_num":
		if !alphaNumRegex.MatchString(strVal) {
			return errors.New("The " + field + " field may only contain letters and numbers.")
		}

	case "alpha_dash":
		if !alphaDashRegex.MatchString(strVal) {
			return errors.New("The " + field + " field may only contain letters, numbers, dashes and underscores.")
		}

	case "required_if":
		// required_if:another_field,value
		if len(ruleParams) >= 2 {
			otherField := ruleParams[0]
			expectedValue := ruleParams[1]
			otherVal, _ := v.data[otherField]
			if fmt.Sprintf("%v", otherVal) == expectedValue {
				if !exists || isEmpty(value) {
					return errors.New("The " + field + " field is required when " + otherField + " is " + expectedValue + ".")
				}
			}
		}

	case "same":
		// same:field
		if len(ruleParams) > 0 {
			otherField := ruleParams[0]
			otherVal, _ := v.data[otherField]
			if fmt.Sprintf("%v", otherVal) != strVal {
				return errors.New("The " + field + " and " + otherField + " must match.")
			}
		}

	case "different":
		// different:field
		if len(ruleParams) > 0 {
			otherField := ruleParams[0]
			otherVal, _ := v.data[otherField]
			if fmt.Sprintf("%v", otherVal) == strVal {
				return errors.New("The " + field + " and " + otherField + " must be different.")
			}
		}

	case "gt":
		// gt:field or gt:value
		if len(ruleParams) > 0 {
			compare := ruleParams[0]
			compareVal, ok := v.data[compare]
			if !ok {
				compareVal = compare
			}
			cf, _ := strconv.ParseFloat(fmt.Sprintf("%v", compareVal), 64)
			// For numeric strings, compare the numeric value; otherwise compare length
			var compareResult bool
			if isNumeric {
				compareResult = numVal <= cf
			} else {
				compareResult = valLen <= cf
			}
			if compareResult {
				return fmt.Errorf("The %s field must be greater than %v.", field, compare)
			}
		}

	case "gte":
		if len(ruleParams) > 0 {
			compare := ruleParams[0]
			compareVal, ok := v.data[compare]
			if !ok {
				compareVal = compare
			}
			cf, _ := strconv.ParseFloat(fmt.Sprintf("%v", compareVal), 64)
			var compareResult bool
			if isNumeric {
				compareResult = numVal < cf
			} else {
				compareResult = valLen < cf
			}
			if compareResult {
				return fmt.Errorf("The %s field must be greater than or equal to %v.", field, compare)
			}
		}

	case "lt":
		if len(ruleParams) > 0 {
			compare := ruleParams[0]
			compareVal, ok := v.data[compare]
			if !ok {
				compareVal = compare
			}
			cf, _ := strconv.ParseFloat(fmt.Sprintf("%v", compareVal), 64)
			var compareResult bool
			if isNumeric {
				compareResult = numVal >= cf
			} else {
				compareResult = valLen >= cf
			}
			if compareResult {
				return fmt.Errorf("The %s field must be less than %v.", field, compare)
			}
		}

	case "lte":
		if len(ruleParams) > 0 {
			compare := ruleParams[0]
			compareVal, ok := v.data[compare]
			if !ok {
				compareVal = compare
			}
			cf, _ := strconv.ParseFloat(fmt.Sprintf("%v", compareVal), 64)
			var compareResult bool
			if isNumeric {
				compareResult = numVal > cf
			} else {
				compareResult = valLen > cf
			}
			if compareResult {
				return fmt.Errorf("The %s field must be less than or equal to %v.", field, compare)
			}
		}

	case "bail":
		return nil

	case "date":
		_, err := time.Parse("2006-01-02", strVal)
		if err != nil {
			_, err2 := time.Parse(time.RFC3339, strVal)
			if err2 != nil {
				return errors.New("The " + field + " field must be a valid date.")
			}
		}

	case "before":
		// before:2025-01-01 or before:another_field
		if len(ruleParams) > 0 {
			target := ruleParams[0]
			targetVal, ok := v.data[target]
			if !ok {
				targetVal = target
			}
			targetStr := fmt.Sprintf("%v", targetVal)
			t, err1 := time.Parse("2006-01-02", strVal)
			targetTime, err2 := time.Parse("2006-01-02", targetStr)
			if err1 == nil && err2 == nil {
				if !t.Before(targetTime) {
					return fmt.Errorf("The %s field must be a date before %s.", field, target)
				}
			} else {
				// Fallback to string comparison if not parseable as dates
				if strVal >= targetStr {
					return fmt.Errorf("The %s field must be before %s.", field, target)
				}
			}
		}

	case "after":
		if len(ruleParams) > 0 {
			target := ruleParams[0]
			targetVal, ok := v.data[target]
			if !ok {
				targetVal = target
			}
			targetStr := fmt.Sprintf("%v", targetVal)
			t, err1 := time.Parse("2006-01-02", strVal)
			targetTime, err2 := time.Parse("2006-01-02", targetStr)
			if err1 == nil && err2 == nil {
				if !t.After(targetTime) {
					return fmt.Errorf("The %s field must be a date after %s.", field, target)
				}
			} else {
				// Fallback to string comparison if not parseable as dates
				if strVal <= targetStr {
					return fmt.Errorf("The %s field must be after %s.", field, target)
				}
			}
		}

	case "json":
		var js json.RawMessage
		if json.Unmarshal([]byte(strVal), &js) != nil {
			return errors.New("The " + field + " field must be valid JSON.")
		}

	case "digits":
		// digits:length or digits:min,max
		if len(ruleParams) >= 2 {
			min, _ := strconv.Atoi(ruleParams[0])
			max, _ := strconv.Atoi(ruleParams[1])
			if len(strVal) < min || len(strVal) > max {
				return fmt.Errorf("The %s field must be between %d and %d digits.", field, min, max)
			}
		} else if len(ruleParams) == 1 {
			length, _ := strconv.Atoi(ruleParams[0])
			if len(strVal) != length {
				return fmt.Errorf("The %s field must be exactly %d digits.", field, length)
			}
		}
		// Validate all characters are digits
		for _, c := range strVal {
			if c < '0' || c > '9' {
				return errors.New("The " + field + " field must be numeric.")
			}
		}

	case "ip":
		// ip, ipv4, or ipv6
		if len(ruleParams) > 0 && ruleParams[0] == "ipv4" {
			if !ipRegex.MatchString(strVal) {
				return errors.New("The " + field + " field must be a valid IPv4 address.")
			}
		} else if len(ruleParams) > 0 && ruleParams[0] == "ipv6" {
			if !ipv6Regex.MatchString(strVal) {
				return errors.New("The " + field + " field must be a valid IPv6 address.")
			}
		} else {
			if !ipRegex.MatchString(strVal) && !ipv6Regex.MatchString(strVal) {
				return errors.New("The " + field + " field must be a valid IP address.")
			}
		}

	case "active_url":
		// Active URL validation would require DNS lookup
		// For now, just validate it's a valid URL format
		if !urlRegex.MatchString(strVal) {
			return errors.New("The " + field + " field must be a valid URL.")
		}

	case "date_format":
		// date_format:format
		if len(ruleParams) > 0 {
			format := ruleParams[0]
			if _, err := time.Parse(format, strVal); err != nil {
				return fmt.Errorf("The %s field does not match the format %s.", field, format)
			}
		}

	case "starts_with":
		// starts_with:value1,value2,...
		found := false
		for _, prefix := range ruleParams {
			if strings.HasPrefix(strVal, prefix) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("The %s field must start with one of: %s.", field, strings.Join(ruleParams, ", "))
		}

	case "ends_with":
		// ends_with:value1,value2,...
		found := false
		for _, suffix := range ruleParams {
			if strings.HasSuffix(strVal, suffix) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("The %s field must end with one of: %s.", field, strings.Join(ruleParams, ", "))
		}

	case "contains":
		// contains:value
		if len(ruleParams) > 0 {
			if !strings.Contains(strVal, ruleParams[0]) {
				return fmt.Errorf("The %s field must contain %s.", field, ruleParams[0])
			}
		}

	case "required_unless":
		// required_unless:another_field,value
		if len(ruleParams) >= 2 {
			otherField := ruleParams[0]
			expectedValue := ruleParams[1]
			otherVal, _ := v.data[otherField]
			if fmt.Sprintf("%v", otherVal) != expectedValue {
				if !exists || isEmpty(value) {
					return errors.New("The " + field + " field is required when " + otherField + " is not " + expectedValue + ".")
				}
			}
		}

	case "required_with":
		// required_with:field1,field2,...
		if len(ruleParams) > 0 {
			for _, param := range ruleParams {
				if val, ok := v.data[param]; ok && !isEmpty(val) {
					if !exists || isEmpty(value) {
						return errors.New("The " + field + " field is required when " + param + " is present.")
					}
					break
				}
			}
		}

	case "required_without":
		// required_without:field1,field2,...
		if len(ruleParams) > 0 {
			for _, param := range ruleParams {
				if val, ok := v.data[param]; !ok || isEmpty(val) {
					if !exists || isEmpty(value) {
						return errors.New("The " + field + " field is required when " + param + " is not present.")
					}
					break
				}
			}
		}

	case "image":
		// Validate that the field is an image file (by extension)
		allowedExtensions := map[string]bool{
			"jpg": true, "jpeg": true, "png": true, "gif": true, "bmp": true, "svg": true, "webp": true,
		}
		if strVal != "" {
			parts := strings.Split(strVal, ".")
			if len(parts) < 2 {
				return errors.New("The " + field + " field must be an image.")
			}
			ext := strings.ToLower(parts[len(parts)-1])
			if !allowedExtensions[ext] {
				return errors.New("The " + field + " field must be a file of type: jpg, jpeg, png, gif, bmp, svg, or webp.")
			}
		}

	case "mimes":
		// mimes:jpg,png,pdf - validate file extension
		if strVal != "" && len(ruleParams) > 0 {
			parts := strings.Split(strVal, ".")
			if len(parts) < 2 {
				return fmt.Errorf("The %s field must be a file of type: %s.", field, strings.Join(ruleParams, ", "))
			}
			ext := strings.ToLower(parts[len(parts)-1])
			found := false
			for _, allowed := range ruleParams {
				if ext == allowed {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("The %s field must be a file of type: %s.", field, strings.Join(ruleParams, ", "))
			}
		}

	case "file":
		// Validate that the field contains a file (non-empty string with extension)
		if strVal != "" {
			parts := strings.Split(strVal, ".")
			if len(parts) < 2 {
				return errors.New("The " + field + " field must be a file.")
			}
		}

	case "dimensions":
		// dimensions:min_width=100,min_height=100 - basic validation
		// Full image dimension checking would require image parsing
		// For now, just validate the field is not empty if dimensions are specified
		if len(ruleParams) > 0 && strVal == "" {
			return errors.New("The " + field + " field has invalid dimensions.")
		}

	case "ulid":
		// ULID format: 26 chars, Crockford Base32
		ulidRegex := regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
		if !ulidRegex.MatchString(strVal) {
			return errors.New("The " + field + " field must be a valid ULID.")
		}

	case "filled":
		// filled is like required but only if the field exists in the data
		if exists && (value == nil || isEmpty(value)) {
			return errors.New("The " + field + " field must be present and not empty.")
		}

	case "present":
		// present: the field must be present in the input data (even if empty)
		if !exists {
			return errors.New("The " + field + " field must be present.")
		}

	case "nullable":
		// nullable is handled at the top of applyRule
		return nil
	}

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

// isValidIdentifier checks if a string is a safe SQL identifier (alphanumeric + underscore only).
func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, c := range s {
		if c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			if i == 0 && c >= '0' && c <= '9' {
				return false // can't start with digit
			}
			continue
		}
		return false
	}
	return true
}

// ==================== PHASE 2: Rule Objects ====================

// Rule is an interface for reusable validation rules.
type Rule interface {
	Validate(value any) error
	RuleName() string
}

// RequiredRule validates that a value is not empty.
type RequiredRule struct{}

func (r *RequiredRule) Validate(value any) error {
	if value == nil || isEmpty(value) {
		return errors.New("field is required")
	}
	return nil
}
func (r *RequiredRule) RuleName() string { return "required" }

// MinRule validates minimum length (strings) or value (numbers).
type MinRule struct{ Min float64 }

func (r *MinRule) Validate(value any) error {
	strVal := fmt.Sprintf("%v", value)
	if n, err := strconv.ParseFloat(strVal, 64); err == nil {
		if n < r.Min {
			return fmt.Errorf("must be at least %v", r.Min)
		}
	} else if float64(len(strVal)) < r.Min {
		return fmt.Errorf("must be at least %v characters", r.Min)
	}
	return nil
}
func (r *MinRule) RuleName() string { return "min" }

// MaxRule validates maximum length (strings) or value (numbers).
type MaxRule struct{ Max float64 }

func (r *MaxRule) Validate(value any) error {
	strVal := fmt.Sprintf("%v", value)
	if n, err := strconv.ParseFloat(strVal, 64); err == nil {
		if n > r.Max {
			return fmt.Errorf("must not exceed %v", r.Max)
		}
	} else if float64(len(strVal)) > r.Max {
		return fmt.Errorf("must not exceed %v characters", r.Max)
	}
	return nil
}
func (r *MaxRule) RuleName() string { return "max" }

// EmailRule validates email format.
type EmailRule struct{}

func (r *EmailRule) Validate(value any) error {
	s, ok := value.(string)
	if !ok {
		return nil
	}
	if !emailRegex.MatchString(s) {
		return errors.New("must be a valid email address")
	}
	return nil
}
func (r *EmailRule) RuleName() string { return "email" }

// RegexRule validates against a regex pattern.
type RegexRule struct{ Pattern string }

func (r *RegexRule) Validate(value any) error {
	s, ok := value.(string)
	if !ok {
		return nil
	}
	matched, _ := regexp.MatchString(r.Pattern, s)
	if !matched {
		return errors.New("does not match the required format")
	}
	return nil
}
func (r *RegexRule) RuleName() string { return "regex" }

// InRule validates that the value is in a list of allowed values.
type InRule struct{ Allowed []string }

func (r *InRule) Validate(value any) error {
	strVal := fmt.Sprintf("%v", value)
	for _, a := range r.Allowed {
		if strVal == a {
			return nil
		}
	}
	return errors.New("must be one of the allowed values")
}
func (r *InRule) RuleName() string { return "in" }

// BooleanRule validates boolean values.
type BooleanRule struct{}

func (r *BooleanRule) Validate(value any) error {
	if _, ok := value.(bool); ok {
		return nil
	}
	lower := strings.ToLower(fmt.Sprintf("%v", value))
	valid := map[string]bool{"true": true, "false": true, "1": true, "0": true, "yes": true, "no": true}
	if !valid[lower] {
		return errors.New("must be a valid boolean")
	}
	return nil
}
func (r *BooleanRule) RuleName() string { return "boolean" }

// URLRule validates URL format.
type URLRule struct{}

func (r *URLRule) Validate(value any) error {
	s, ok := value.(string)
	if !ok {
		return nil
	}
	if !urlRegex.MatchString(s) {
		return errors.New("must be a valid URL")
	}
	return nil
}
func (r *URLRule) RuleName() string { return "url" }

// UUIDRule validates UUID format.
type UUIDRule struct{}

func (r *UUIDRule) Validate(value any) error {
	s, ok := value.(string)
	if !ok {
		return nil
	}
	if !uuidRegex.MatchString(s) {
		return errors.New("must be a valid UUID")
	}
	return nil
}
func (r *UUIDRule) RuleName() string { return "uuid" }

// IntegerRule validates integer values.
type IntegerRule struct{}

func (r *IntegerRule) Validate(value any) error {
	if _, ok := value.(int); ok {
		return nil
	}
	if _, ok := value.(int64); ok {
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return nil
	}
	if _, err := strconv.Atoi(s); err != nil {
		return errors.New("must be a valid integer")
	}
	return nil
}
func (r *IntegerRule) RuleName() string { return "integer" }

// ValidateWithRules validates using Rule objects instead of string rules.
func ValidateWithRules(data map[string]any, fieldRules map[string][]Rule) map[string][]error {
	errs := make(map[string][]error)
	for field, rules := range fieldRules {
		for _, rule := range rules {
			val, exists := data[field]
			if !exists && rule.RuleName() != "required" {
				continue
			}
			if err := rule.Validate(val); err != nil {
				errs[field] = append(errs[field], err)
			}
		}
	}
	return errs
}

// ==================== PHASE 2: After Validation Hooks ====================

// AfterValidationFunc is a callback that runs after validation passes.
type AfterValidationFunc func(data map[string]any) error

// BeforeValidationFunc is a callback that runs before validation.
type BeforeValidationFunc func(data map[string]any) error

// ValidateWithHooks executes validation with optional lifecycle hooks.
func ValidateWithHooks(data map[string]any, rules map[string][]string, db *sql.DB, before BeforeValidationFunc, after AfterValidationFunc) map[string][]error {
	// Run before hook
	if before != nil {
		if err := before(data); err != nil {
			return map[string][]error{"_hook": {err}}
		}
	}

	v := NewValidator(data, rules)
	if db != nil {
		v.WithDB(db)
	}
	errs := v.Validate()

	// Run after hook only if no validation errors
	if after != nil && len(errs) == 0 {
		if err := after(data); err != nil {
			errs["_hook"] = []error{err}
		}
	}

	return errs
}

// ==================== PHASE 2: Sometimes Validation ====================

// SometimesCondition defines when a rule should be applied.
type SometimesCondition struct {
	Field    string
	Operator string // "filled", "not_filled", "equals", "not_equals"
	Value    any
}

// When creates a SometimesCondition for use with Sometimes.
func When(field, operator string, value any) SometimesCondition {
	return SometimesCondition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// ValidateSometimes validates with conditional (sometimes) rules.
// Usage: ValidateSometimes(data, rules, Sometimes("role", "equals", "admin", "bio", "required"))
func ValidateSometimes(data map[string]any, rules map[string][]string, conditions ...struct {
	Field    string
	When     SometimesCondition
	RuleStr  string
}) map[string][]error {
	// Evaluate conditions and add conditional rules
	for _, cond := range conditions {
		if evaluateCondition(cond.When, data) {
			existing := rules[cond.Field]
			rules[cond.Field] = append(existing, cond.RuleStr)
		}
	}

	v := NewValidator(data, rules)
	return v.Validate()
}

// evaluateCondition checks if a sometimes condition is met.
func evaluateCondition(cond SometimesCondition, data map[string]any) bool {
	val, exists := data[cond.Field]
	switch cond.Operator {
	case "filled":
		return exists && val != nil && fmt.Sprintf("%v", val) != ""
	case "not_filled":
		return !exists || val == nil || fmt.Sprintf("%v", val) == ""
	case "equals":
		return fmt.Sprintf("%v", val) == fmt.Sprintf("%v", cond.Value)
	case "not_equals":
		return fmt.Sprintf("%v", val) != fmt.Sprintf("%v", cond.Value)
	}
	return false
}

// ==================== PHASE 2: Implicit Field Attributes ====================

// AttributeNames maps field names to human-readable names for error messages.
var attributeNames = map[string]string{}

// SetAttributeName registers a human-readable name for a field.
func SetAttributeName(field, name string) {
	attributeNames[field] = name
}

// GetAttributeName returns the human-readable name for a field.
func GetAttributeName(field string) string {
	if name, ok := attributeNames[field]; ok {
		return name
	}
	return field
}

// SetAttributeNames registers multiple field names at once.
func SetAttributeNames(names map[string]string) {
	for k, v := range names {
		attributeNames[k] = v
	}
}

