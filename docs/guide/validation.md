# Validation

> **Status**: Implemented

GoW provides a fluent validation system for validating incoming request data.

## Basic Usage

```go
import "github.com/mechneerd/gow/validation"

errors := validation.Validate(data, rules)
```

## Available Rules

| Rule | Description |
|------|-------------|
| `required` | Field must be present and not empty |
| `email` | Must be a valid email address |
| `min:N` | Minimum length (strings) or value (numbers) |
| `max:N` | Maximum length (strings) or value (numbers) |
| `gt:N` | Greater than N |
| `gte:N` | Greater than or equal to N |
| `lt:N` | Less than N |
| `lte:N` | Less than or equal to N |
| `in:a,b,c` | Must be one of the listed values |
| `not_in:a,b,c` | Must not be one of the listed values |
| `alpha` | Must contain only letters |
| `alpha_num` | Must contain only letters and numbers |
| `regex:pattern` | Must match the given regex pattern |
| `url` | Must be a valid URL |
| `uuid` | Must be a valid UUID |
| `boolean` | Must be a valid boolean value |
| `date` | Must be a valid date (YYYY-MM-DD) |
| `before:date` | Must be before the given date |
| `after:date` | Must be after the given date |
| `array` | Must be a valid array/slice |
| `unique:table,column` | Must be unique in the database table |
| `exists:table,column` | Must exist in the database table |
| `nullable` | Field can be null/empty (skip other rules if empty) |

## Custom Messages

```go
rules := validation.Map{
    "name":  "required|min:3",
    "email": "required|email|unique:users,email",
}

messages := validation.Messages{
    "name.required": "Please provide your name",
    "name.min":      "Name must be at least 3 characters",
}

errors := validation.Validate(data, rules, messages)
```

## Custom Attributes

```go
attributes := validation.Attributes{
    "name":  "full name",
    "email": "email address",
}

errors := validation.Validate(data, rules, messages, attributes)
```

## Form Requests

For more complex validation, use FormRequest structs:

```go
type CreateUserRequest struct{}

func (r *CreateUserRequest) Authorize(r *http.Request) bool {
    return true
}

func (r *CreateUserRequest) Rules(r *http.Request) validation.Map {
    return validation.Map{
        "name":  "required|min:3",
        "email": "required|email|unique:users,email",
        "password": "required|min:8|confirmed",
    }
}

func (r *CreateUserRequest) Messages() validation.Messages {
    return validation.Messages{
        "password.confirmed": "Password confirmation does not match",
    }
}
```

## Error Handling

```go
errors := validation.Validate(data, rules)

if errors.HasErrors() {
    // Get first error for a field
    firstError := errors.Get("email")

    // Get all errors for a field
    allErrors := errors.Get("email") // returns []string

    // Get all errors as a map
    all := errors.All() // map[string][]string
}
```
