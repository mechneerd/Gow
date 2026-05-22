package routing

import (
	"gow/database/orm"
	"net/http"
	"reflect"
	"strings"
)

// Dispatcher handles reflective execution of controller methods and performs Route Model Binding.
type Dispatcher struct {
	db *orm.DB
}

// NewDispatcher creates a new dispatcher with database access for model binding.
func NewDispatcher(db *orm.DB) *Dispatcher {
	return &Dispatcher{db: db}
}

// Wrap takes any function and wraps it into a standard http.HandlerFunc.
// It reflects on the function parameters and attempts to resolve them:
// 1. http.ResponseWriter
// 2. *http.Request
// 3. Route parameters via Route Model Binding if it's an ORM Model
func (d *Dispatcher) Wrap(handler any) HandlerFunc {
	val := reflect.ValueOf(handler)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		panic("Dispatcher.Wrap: handler must be a function")
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		in := make([]reflect.Value, typ.NumIn())
		params := r.Context().Value(ParamsKey).(map[string]string)

		for i := 0; i < typ.NumIn(); i++ {
			paramType := typ.In(i)

			// 1. Inject ResponseWriter
			if paramType.Implements(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()) {
				in[i] = reflect.ValueOf(w)
				continue
			}

			// 2. Inject *http.Request
			if paramType == reflect.TypeOf((*http.Request)(nil)) {
				in[i] = reflect.ValueOf(r)
				continue
			}

			// 3. Route Model Binding
			// We check if it's a struct pointer that might be an ORM model
			if paramType.Kind() == reflect.Ptr && paramType.Elem().Kind() == reflect.Struct {
				// Try to find a matching parameter from the route
				// For a parameter `*models.User`, we look for `{user}` in the route params
				modelName := strings.ToLower(paramType.Elem().Name())
				if idStr, ok := params[modelName]; ok {
					// We have a match! We need to query the database.
					// Since Go doesn't have easy generic method invocation at runtime without the exact type parameter,
					// we have to construct a query dynamically using the builder directly.
					
					// Instantiate the struct
					modelPtr := reflect.New(paramType.Elem())
					
					// Determine table name
					tableName := modelName + "s" // Simplistic pluralization
					if m, ok := modelPtr.Interface().(orm.Model); ok {
						tableName = m.TableName()
					}
					
					// Query builder
					builder := d.db.Builder.Table(tableName)
					builder.Where("id", "=", idStr)
					builder.Limit(1)
					
					rows, err := builder.Get()
					if err == nil {
						defer rows.Close()
						if rows.Next() {
							// We need to hydrate. We can use the exported hydrate logic or simple reflection.
							columns, _ := rows.Columns()
							scanArgs := make([]any, len(columns))
							fieldMap := make(map[string]int)
							
							for j := 0; j < paramType.Elem().NumField(); j++ {
								dbTag := paramType.Elem().Field(j).Tag.Get("db")
								if dbTag != "" {
									fieldMap[dbTag] = j
								} else {
									fieldMap[strings.ToLower(paramType.Elem().Field(j).Name)] = j
								}
							}
							
							for j, col := range columns {
								if fieldIdx, ok := fieldMap[col]; ok {
									scanArgs[j] = modelPtr.Elem().Field(fieldIdx).Addr().Interface()
								} else {
									var dummy any
									scanArgs[j] = &dummy
								}
							}
							
							_ = rows.Scan(scanArgs...)
						}
					}
					
					in[i] = modelPtr
					continue
				}
			}

			// Fallback: zero value
			in[i] = reflect.Zero(paramType)
		}

		out := val.Call(in)
		
		if len(out) > 0 {
			if err, ok := out[0].Interface().(error); ok && err != nil {
				return err
			}
		}
		return nil
	}
}
