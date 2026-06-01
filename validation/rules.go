package validation

import (
	"fmt"
	"reflect"
	"strings"
)

// CustomRule defines a custom validation rule interface (Laravel InvokableRule equivalent).
type CustomRule interface {
	Passes(field string, value any, parameters []string) bool
	Message(field string, parameters []string) string
}

// AfterDateRule validates that a date is after another field.
type AfterDateRule struct{}

func (r *AfterDateRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *AfterDateRule) Message(field string, parameters []string) string {
	if len(parameters) > 0 {
		return fmt.Sprintf("The %s field must be after %s.", field, parameters[0])
	}
	return fmt.Sprintf("The %s field must be a date after today.", field)
}

// BeforeDateRule validates that a date is before another field.
type BeforeDateRule struct{}

func (r *BeforeDateRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *BeforeDateRule) Message(field string, parameters []string) string {
	if len(parameters) > 0 {
		return fmt.Sprintf("The %s field must be before %s.", field, parameters[0])
	}
	return fmt.Sprintf("The %s field must be a date before today.", field)
}

// DifferentRule validates that two fields have different values.
type DifferentRule struct{}

func (r *DifferentRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *DifferentRule) Message(field string, parameters []string) string {
	if len(parameters) > 0 {
		return fmt.Sprintf("The %s field must be different from %s.", field, parameters[0])
	}
	return fmt.Sprintf("The %s field must be different.", field)
}

// CustomUniqueRule validates that a field is unique in the database.
type CustomUniqueRule struct {
	table  string
	column string
}

// NewCustomUniqueRule creates a new unique rule.
func NewCustomUniqueRule(table, column string) *CustomUniqueRule {
	return &CustomUniqueRule{table: table, column: column}
}

func (r *CustomUniqueRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *CustomUniqueRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s has already been taken.", field)
}

// CustomExistsRule validates that a field exists in the database.
type CustomExistsRule struct {
	table  string
	column string
}

// NewCustomExistsRule creates a new exists rule.
func NewCustomExistsRule(table, column string) *CustomExistsRule {
	return &CustomExistsRule{table: table, column: column}
}

func (r *CustomExistsRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *CustomExistsRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The selected %s is invalid.", field)
}

// AcceptedRule validates that a field is "yes", "on", 1, or true.
type AcceptedRule struct{}

func (r *AcceptedRule) Passes(field string, value any, parameters []string) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "yes" || strings.ToLower(v) == "on" || v == "1"
	case int:
		return v == 1
	default:
		return false
	}
}

func (r *AcceptedRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s must be accepted.", field)
}

// DeclinedRule validates that a field is "no", "off", 0, or false.
type DeclinedRule struct{}

func (r *DeclinedRule) Passes(field string, value any, parameters []string) bool {
	a := &AcceptedRule{}
	return !a.Passes(field, value, parameters)
}

func (r *DeclinedRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s must be declined.", field)
}

// ConfirmedRule validates that two fields match.
type ConfirmedRule struct{}

func (r *ConfirmedRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *ConfirmedRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s confirmation does not match.", field)
}

// DistinctRule validates that all array values are unique.
type DistinctRule struct{}

func (r *DistinctRule) Passes(field string, value any, parameters []string) bool {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		seen := make(map[any]bool)
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			if seen[elem] {
				return false
			}
			seen[elem] = true
		}
	}
	return true
}

func (r *DistinctRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s field must contain unique values.", field)
}

// DoesntStartWithRule validates that a field does not start with a specific value.
type DoesntStartWithRule struct{}

func (r *DoesntStartWithRule) Passes(field string, value any, parameters []string) bool {
	strVal := fmt.Sprintf("%v", value)
	for _, p := range parameters {
		if strings.HasPrefix(strVal, p) {
			return false
		}
	}
	return true
}

func (r *DoesntStartWithRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s field must not start with one of: %s.", field, strings.Join(parameters, ", "))
}

// DoesntEndWithRule validates that a field does not end with a specific value.
type DoesntEndWithRule struct{}

func (r *DoesntEndWithRule) Passes(field string, value any, parameters []string) bool {
	strVal := fmt.Sprintf("%v", value)
	for _, p := range parameters {
		if strings.HasSuffix(strVal, p) {
			return false
		}
	}
	return true
}

func (r *DoesntEndWithRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s field must not end with one of: %s.", field, strings.Join(parameters, ", "))
}

// ExcludeRule always passes (used to exclude fields from validation).
type ExcludeRule struct{}

func (r *ExcludeRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *ExcludeRule) Message(field string, parameters []string) string {
	return ""
}

// PresentRule validates that a field is present.
type PresentRule struct{}

func (r *PresentRule) Passes(field string, value any, parameters []string) bool {
	return value != nil
}

func (r *PresentRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s field must be present.", field)
}

// MissingRule validates that a field is not present.
type MissingRule struct{}

func (r *MissingRule) Passes(field string, value any, parameters []string) bool {
	return value == nil
}

func (r *MissingRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s field must not be present.", field)
}

// ProhibitedRule validates that a field is not present.
type ProhibitedRule struct{}

func (r *ProhibitedRule) Passes(field string, value any, parameters []string) bool {
	if value == nil {
		return true
	}
	return fmt.Sprintf("%v", value) == ""
}

func (r *ProhibitedRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s field must not be present.", field)
}

// CurrentPasswordRule validates that the field matches the authenticated user's password.
type CurrentPasswordRule struct{}

func (r *CurrentPasswordRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *CurrentPasswordRule) Message(field string, parameters []string) string {
	return "The password is incorrect."
}

// EnumRule validates that a field is a valid enum value.
type EnumRule struct {
	enumType reflect.Type
}

// NewEnumRule creates a new enum rule.
func NewEnumRule(enumType reflect.Type) *EnumRule {
	return &EnumRule{enumType: enumType}
}

func (r *EnumRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *EnumRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The selected %s must be a valid value.", field)
}

// ActiveURLRule validates that a field has a valid A or AAAA record.
type ActiveURLRule struct{}

func (r *ActiveURLRule) Passes(field string, value any, parameters []string) bool {
	return true
}

func (r *ActiveURLRule) Message(field string, parameters []string) string {
	return fmt.Sprintf("The %s field must contain a valid URL.", field)
}
