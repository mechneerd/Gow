package orm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Built-in cast types (Laravel-like)
const (
	CastDatetime = "datetime"
	CastJSON     = "json"
	CastBool     = "bool"
	CastInt      = "int"
	CastFloat    = "float"
	CastString   = "string"
)

// Castable models declare their attribute casts.
type Castable interface {
	Casts() map[string]string
}

// getCast returns a Cast implementation for the given type name.
func getCast(castType string) Cast {
	switch castType {
	case CastDatetime:
		return &DateTimeCast{}
	case CastJSON:
		return &JSONCast{}
	case CastBool:
		return &BoolCast{}
	case CastInt:
		return &IntCast{}
	case CastFloat:
		return &FloatCast{}
	case CastString:
		return &StringCast{}
	default:
		return nil
	}
}

// ================== Built-in Cast Implementations ==================

type DateTimeCast struct{}

func (c *DateTimeCast) Get(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", v)
		}
		return t, err
	default:
		return value, nil
	}
}

func (c *DateTimeCast) Set(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	if t, ok := value.(time.Time); ok {
		return t.Format(time.RFC3339), nil
	}
	return value, nil
}

type JSONCast struct{}

func (c *JSONCast) Get(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	if b, ok := value.([]byte); ok {
		var out any
		if err := json.Unmarshal(b, &out); err != nil {
			return nil, err
		}
		return out, nil
	}
	if s, ok := value.(string); ok {
		var out any
		if err := json.Unmarshal([]byte(s), &out); err != nil {
			return nil, err
		}
		return out, nil
	}
	return value, nil
}

func (c *JSONCast) Set(value any) (any, error) {
	if value == nil {
		return nil, nil
	}
	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

type BoolCast struct{}

func (c *BoolCast) Get(value any) (any, error) {
	if value == nil {
		return false, nil
	}
	switch v := value.(type) {
	case bool:
		return v, nil
	case int, int64, float64:
		return reflect.ValueOf(v).Int() != 0, nil
	case string:
		return v == "1" || v == "true", nil
	default:
		return false, nil
	}
}

func (c *BoolCast) Set(value any) (any, error) {
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return value, nil
}

type IntCast struct{}

func (c *IntCast) Get(value any) (any, error) {
	if value == nil {
		return 0, nil
	}
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		var i int
		fmt.Sscanf(v, "%d", &i)
		return i, nil
	default:
		return value, nil
	}
}

func (c *IntCast) Set(value any) (any, error) {
	return value, nil
}

type FloatCast struct{}

func (c *FloatCast) Get(value any) (any, error) {
	if value == nil {
		return 0.0, nil
	}
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f, nil
	default:
		return value, nil
	}
}

func (c *FloatCast) Set(value any) (any, error) {
	return value, nil
}

type StringCast struct{}

func (c *StringCast) Get(value any) (any, error) {
	if value == nil {
		return "", nil
	}
	return fmt.Sprintf("%v", value), nil
}

func (c *StringCast) Set(value any) (any, error) {
	return value, nil
}

// applyCastsAfterHydrate runs registered casts on a hydrated model.
func applyCastsAfterHydrate(model any) {
	if c, ok := model.(Castable); ok {
		casts := c.Casts()
		val := reflect.ValueOf(model).Elem()
		for col, castName := range casts {
			caster := getCast(castName)
			if caster == nil {
				continue
			}
			field := val.FieldByNameFunc(func(n string) bool {
				return strings.EqualFold(n, col) || strings.ToLower(n) == strings.ToLower(col)
			})
			if !field.IsValid() || !field.CanSet() {
				continue
			}
			casted, err := caster.Get(field.Interface())
			if err == nil && casted != nil {
				// Try to set back (handle type differences gracefully)
				rv := reflect.ValueOf(casted)
				if rv.Type().ConvertibleTo(field.Type()) {
					field.Set(rv.Convert(field.Type()))
				} else if field.Kind() == reflect.Interface {
					field.Set(rv)
				}
			}
		}
	}
}

// applyCastsBeforeSave runs Set() on fields before insert/update.
func applyCastsBeforeSave(model any) {
	if c, ok := model.(Castable); ok {
		casts := c.Casts()
		val := reflect.ValueOf(model).Elem()
		for col, castName := range casts {
			caster := getCast(castName)
			if caster == nil {
				continue
			}
			field := val.FieldByNameFunc(func(n string) bool {
				return strings.EqualFold(n, col) || strings.ToLower(n) == strings.ToLower(col)
			})
			if !field.IsValid() || !field.CanSet() {
				continue
			}
			casted, err := caster.Set(field.Interface())
			if err == nil && casted != nil {
				rv := reflect.ValueOf(casted)
				if rv.Type().ConvertibleTo(field.Type()) {
					field.Set(rv.Convert(field.Type()))
				} else if field.Kind() == reflect.Interface {
					field.Set(rv)
				}
			}
		}
	}
}

