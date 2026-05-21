package orm

import (
	"database/sql"
	"gow/database/query"
	"reflect"
	"strings"
	"time"
)

// Model represents a generic interface for all Goquent models.
type Model interface {
	TableName() string
}

// DB represents the ORM database connection.
type DB struct {
	Conn    *sql.DB
	Builder *query.Builder
}

// ModelQuery is an ORM-wrapper around the Query Builder.
type ModelQuery[T any] struct {
	builder *query.Builder
	db      *DB
}

// NewQuery creates a new query for a specific model type.
func NewQuery[T any](db *DB) *ModelQuery[T] {
	var model T
	table := getTableName(model)

	return &ModelQuery[T]{
		builder: db.Builder.Table(table),
		db:      db,
	}
}

// Where adds a where clause.
func (q *ModelQuery[T]) Where(column, operator string, value any) *ModelQuery[T] {
	q.builder.Where(column, operator, value)
	return q
}

// First fetches the first matching record.
func (q *ModelQuery[T]) First() (*T, error) {
	q.builder.Limit(1)
	rows, err := q.builder.Get()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	model, err := hydrateModel[T](rows)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// Get fetches all matching records.
func (q *ModelQuery[T]) Get() ([]*T, error) {
	rows, err := q.builder.Get()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*T
	for rows.Next() {
		model, err := hydrateModel[T](rows)
		if err != nil {
			return nil, err
		}
		results = append(results, model)
	}

	return results, nil
}

// Insert creates a new record from the model struct.
func (q *ModelQuery[T]) Insert(model *T) error {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	values := make(map[string]any)
	now := time.Now()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		gowTag := field.Tag.Get("gow")
		
		if strings.Contains(gowTag, "autoIncrement") {
			continue
		}

		if strings.Contains(gowTag, "autoCreateTime") || strings.Contains(gowTag, "autoUpdateTime") {
			values[dbTag] = now
			val.Field(i).Set(reflect.ValueOf(now))
			continue
		}

		values[dbTag] = val.Field(i).Interface()
	}

	res, err := q.builder.Insert(values)
	if err != nil {
		return err
	}

	// Try to get last insert ID
	id, err := res.LastInsertId()
	if err == nil {
		// Attempt to set ID back to the struct if it's named ID
		idField := val.FieldByName("ID")
		if idField.IsValid() && idField.CanSet() {
			idField.SetInt(id)
		}
	}

	return nil
}

// Helper: hydrate a model from sql.Rows using reflection.
func hydrateModel[T any](rows *sql.Rows) (*T, error) {
	var model T
	val := reflect.ValueOf(&model).Elem()
	typ := val.Type()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create pointers for scanning
	scanArgs := make([]any, len(columns))
	fieldMap := make(map[string]int)

	for i := 0; i < typ.NumField(); i++ {
		dbTag := typ.Field(i).Tag.Get("db")
		if dbTag != "" {
			fieldMap[dbTag] = i
		} else {
			fieldMap[strings.ToLower(typ.Field(i).Name)] = i
		}
	}

	for i, col := range columns {
		if fieldIdx, ok := fieldMap[col]; ok {
			scanArgs[i] = val.Field(fieldIdx).Addr().Interface()
		} else {
			// Dummy receiver if column not in struct
			var dummy any
			scanArgs[i] = &dummy
		}
	}

	err = rows.Scan(scanArgs...)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// Helper: extract table name from model or default to pluralized struct name.
func getTableName(model any) string {
	if m, ok := model.(Model); ok {
		return m.TableName()
	}
	// Fallback to snake_case plural
	name := reflect.TypeOf(model).Name()
	return strings.ToLower(name) + "s" // Simplistic pluralization
}
