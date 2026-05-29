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

	for field, fieldRules := range v.rules {
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
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
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
		// unique:table,column
		if v.db != nil && len(ruleParams) >= 2 {
			table, column := ruleParams[0], ruleParams[1]
			// Validate identifiers to prevent SQL injection
			if !isValidIdentifier(table) || !isValidIdentifier(column) {
				return errors.New("The " + field + " validation rule contains invalid identifiers.")
			}
			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\" WHERE \"%s\" = ?", table, column)
			v.db.QueryRow(query, value).Scan(&count)
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
		urlRegex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
		if !urlRegex.MatchString(strVal) {
			return errors.New("The " + field + " field must be a valid URL.")
		}

	case "uuid":
		uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
		if !uuidRegex.MatchString(strVal) {
			return errors.New("The " + field + " field must be a valid UUID.")
		}

	case "alpha":
		alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
		if !alphaRegex.MatchString(strVal) {
			return errors.New("The " + field + " field may only contain letters.")
		}

	case "alpha_num":
		alphaNumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
		if !alphaNumRegex.MatchString(strVal) {
			return errors.New("The " + field + " field may only contain letters and numbers.")
		}

	case "alpha_dash":
		alphaDashRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
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
			if valLen <= cf {
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
			if valLen < cf {
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
			if valLen >= cf {
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
			if valLen > cf {
				return fmt.Errorf("The %s field must be less than or equal to %v.", field, compare)
			}
		}

	case "nullable":
		// If value is empty/null, skip all further rules for this field
		if !exists || isEmpty(value) {
			return nil // skip remaining rules
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
			if strVal >= targetStr {
				return fmt.Errorf("The %s field must be a date before %s.", field, target)
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
			if strVal <= targetStr {
				return fmt.Errorf("The %s field must be a date after %s.", field, target)
			}
		}

	case "json":
		var js json.RawMessage
		if json.Unmarshal([]byte(strVal), &js) != nil {
			return errors.New("The " + field + " field must be valid JSON.")
		}
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

