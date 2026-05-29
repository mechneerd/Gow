package routing

import (
	"fmt"
	"github.com/mechneerd/gow/database/orm"
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
func (d *Dispatcher) Wrap(handler any) HandlerFunc {
	val := reflect.ValueOf(handler)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		return func(w http.ResponseWriter, r *http.Request) error {
			return fmt.Errorf("Dispatcher.Wrap: handler must be a function, got %T", handler)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		in := make([]reflect.Value, typ.NumIn())
		params, _ := r.Context().Value(ParamsKey).(map[string]string)
		if params == nil {
			params = make(map[string]string)
		}

		for i := 0; i < typ.NumIn(); i++ {
			paramType := typ.In(i)

			if paramType.Implements(reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()) {
				in[i] = reflect.ValueOf(w)
				continue
			}

			if paramType == reflect.TypeOf((*http.Request)(nil)) {
				in[i] = reflect.ValueOf(r)
				continue
			}

			if paramType.Kind() == reflect.Ptr && paramType.Elem().Kind() == reflect.Struct {
				modelName := strings.ToLower(paramType.Elem().Name())
				if idStr, ok := params[modelName]; ok {
					modelPtr := reflect.New(paramType.Elem())

					tableName := modelName + "s"
					if m, ok := modelPtr.Interface().(orm.Model); ok {
						tableName = m.TableName()
					}

					builder := d.db.Builder.Table(tableName)
					builder.Where("id", "=", idStr)
					builder.Limit(1)

					rows, err := builder.Get()
					if err == nil {
						defer rows.Close()
						if rows.Next() {
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

							if err := rows.Scan(scanArgs...); err != nil {
								return fmt.Errorf("failed to hydrate model %s: %w", paramType.Elem().Name(), err)
							}
						}
					}

					in[i] = modelPtr
					continue
				}
			}

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

