package validation

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
		for _, rule := range fieldRules {
			if err := v.applyRule(field, val, exists, rule); err != nil {
				errs[field] = append(errs[field], err)
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
			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
			v.db.QueryRow(query, value).Scan(&count)
			if count > 0 {
				return errors.New("The " + field + " has already been taken.")
			}
		}
	case "exists":
		// exists:table,column
		if v.db != nil && len(ruleParams) >= 2 {
			table, column := ruleParams[0], ruleParams[1]
			var count int
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", table, column)
			v.db.QueryRow(query, value).Scan(&count)
			if count == 0 {
				return errors.New("The selected " + field + " is invalid.")
			}
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
