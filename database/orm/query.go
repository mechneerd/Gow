package orm

import (
	"context"
	"database/sql"
	"errors"
	"gow/database/query"
	"reflect"
	"strings"
	"time"
)

var ErrCancelled = errors.New("operation cancelled by model observer")

// Model represents a generic interface for all Goquent models.
type Model interface {
	TableName() string
}

// DB represents the ORM database connection.
type DB struct {
	Conn    query.QueryExecer
	Builder *query.Builder
}

// Transaction executes a function within a database transaction.
func (db *DB) Transaction(ctx context.Context, fn func(txDB *DB) error) error {
	sqlDB, ok := db.Conn.(*sql.DB)
	if !ok {
		return errors.New("cannot start transaction: connection is not a root *sql.DB")
	}

	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	txDB := &DB{
		Conn:    tx,
		Builder: db.Builder.WithConn(tx),
	}

	if err := fn(txDB); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ModelQuery is an ORM-wrapper around the Query Builder.
type ModelQuery[T any] struct {
	builder       *query.Builder
	db            *DB
	with          []string
	softDeleteCol string
	withTrashed   bool
	onlyTrashed   bool
}

// NewQuery creates a new query for a specific model type.
func NewQuery[T any](db *DB) *ModelQuery[T] {
	var model T
	table := getTableName(model)

	builder := db.Builder.Clone().Table(table)
	
	// Apply global scopes
	modelName := reflect.TypeOf((*T)(nil)).Elem().Name()
	if scopes, ok := globalScopes[modelName]; ok {
		for _, scope := range scopes {
			builder = scope.Apply(builder)
		}
	}

	q := &ModelQuery[T]{
		builder:       builder,
		db:            db,
		with:          make([]string, 0),
		softDeleteCol: getSoftDeleteColumn(reflect.TypeOf((*T)(nil))),
	}
	return q
}

// WithTrashed includes soft-deleted records in the query.
func (q *ModelQuery[T]) WithTrashed() *ModelQuery[T] {
	q.withTrashed = true
	return q
}

// OnlyTrashed limits the query to only soft-deleted records.
func (q *ModelQuery[T]) OnlyTrashed() *ModelQuery[T] {
	q.onlyTrashed = true
	return q
}

func (q *ModelQuery[T]) applySoftDeletes() {
	if q.softDeleteCol != "" {
		if q.onlyTrashed {
			q.builder.WhereNotNull(q.softDeleteCol)
		} else if !q.withTrashed {
			q.builder.WhereNull(q.softDeleteCol)
		}
	}
}

// With adds a relationship to be eager-loaded.
func (q *ModelQuery[T]) With(relation string) *ModelQuery[T] {
	q.with = append(q.with, relation)
	return q
}

// Where adds a where clause.
func (q *ModelQuery[T]) Where(column, operator string, value any) *ModelQuery[T] {
	q.builder.Where(column, operator, value)
	return q
}

// First fetches the first matching record.
func (q *ModelQuery[T]) First() (*T, error) {
	q.applySoftDeletes()
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
	q.applySoftDeletes()
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

	if len(results) > 0 && len(q.with) > 0 {
		err = eagerLoadRelationships(q.db, results, q.with)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// Insert creates a new record from the model struct.
func (q *ModelQuery[T]) Insert(model *T) error {
	if !DispatchModelEvent(model, "creating") {
		return ErrCancelled
	}

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
		
		// Skip primary key if it's zero
		if strings.Contains(gowTag, "primaryKey") || dbTag == "id" {
			switch val.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if val.Field(i).Int() == 0 {
					continue
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if val.Field(i).Uint() == 0 {
					continue
				}
			case reflect.String:
				if val.Field(i).String() == "" {
					continue
				}
			}
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

	DispatchModelEvent(model, "created")
	return nil
}

// Find fetches a model by its primary key.
func (q *ModelQuery[T]) Find(id any) (*T, error) {
	// Assume primary key is always 'id' for now, later we can check struct tags
	q.builder.Where("id", "=", id)
	return q.First()
}

// Update updates an existing model.
func (q *ModelQuery[T]) Update(model *T) error {
	if !DispatchModelEvent(model, "updating") {
		return ErrCancelled
	}

	val := reflect.ValueOf(model).Elem()
	typ := val.Type()

	values := make(map[string]any)
	now := time.Now()
	
	var pkValue any
	var pkColumn string = "id" // default

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		gowTag := field.Tag.Get("gow")
		
		if strings.Contains(gowTag, "primaryKey") {
			pkColumn = dbTag
			pkValue = val.Field(i).Interface()
			continue
		}
		if dbTag == "id" && pkValue == nil {
			pkValue = val.Field(i).Interface()
			continue
		}

		if strings.Contains(gowTag, "autoUpdateTime") {
			values[dbTag] = now
			val.Field(i).Set(reflect.ValueOf(now))
			continue
		}

		values[dbTag] = val.Field(i).Interface()
	}
	
	if pkValue == nil {
		return sql.ErrNoRows // Cannot update without a primary key
	}

	q.applySoftDeletes()
	q.builder.Where(pkColumn, "=", pkValue)
	_, err := q.builder.Update(values)
	if err == nil {
		DispatchModelEvent(model, "updated")
	}
	return err
}

// Delete deletes the given model (or rows matching the query if model is nil).
func (q *ModelQuery[T]) Delete(model *T) error {
	if model != nil {
		if !DispatchModelEvent(model, "deleting") {
			return ErrCancelled
		}
		val := reflect.ValueOf(model).Elem()
		typ := val.Type()
		var pkValue any
		var pkColumn string = "id"
		
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			dbTag := field.Tag.Get("db")
			gowTag := field.Tag.Get("gow")
			
			if strings.Contains(gowTag, "primaryKey") || dbTag == "id" {
				pkColumn = dbTag
				pkValue = val.Field(i).Interface()
				break
			}
		}
		
		if pkValue == nil {
			return sql.ErrNoRows
		}
		
		q.builder.Where(pkColumn, "=", pkValue)
	}
	
	q.applySoftDeletes()

	if q.softDeleteCol != "" {
		// Soft delete via Update
		now := time.Now()
		_, err := q.builder.Update(map[string]any{
			q.softDeleteCol: now,
		})
		if err == nil && model != nil {
			DispatchModelEvent(model, "deleted")
			// Update the model struct's deleted_at field if possible
			val := reflect.ValueOf(model).Elem()
			for i := 0; i < val.Type().NumField(); i++ {
				field := val.Type().Field(i)
				dbTag := field.Tag.Get("db")
				if dbTag == "" {
					dbTag = strings.ToLower(field.Name)
				}
				if dbTag == q.softDeleteCol {
					// Handle sql.NullTime or time.Time pointer
					f := val.Field(i)
					if f.Type() == reflect.TypeOf(time.Time{}) {
						f.Set(reflect.ValueOf(now))
					} else if f.Type() == reflect.TypeOf(&time.Time{}) {
						f.Set(reflect.ValueOf(&now))
					} else if f.Type() == reflect.TypeOf(sql.NullTime{}) {
						f.Set(reflect.ValueOf(sql.NullTime{Time: now, Valid: true}))
					}
					break
				}
			}
		}
		return err
	}
	
	_, err := q.builder.Delete()
	if err == nil && model != nil {
		DispatchModelEvent(model, "deleted")
	}
	return err
}

// ForceDelete permanently deletes the model.
func (q *ModelQuery[T]) ForceDelete(model *T) error {
	// Temporarily clear soft delete so it bypasses scope and does a hard delete
	originalSoftDelete := q.softDeleteCol
	q.softDeleteCol = ""
	err := q.Delete(model)
	q.softDeleteCol = originalSoftDelete
	return err
}

// Restore restores a soft-deleted model.
func (q *ModelQuery[T]) Restore(model *T) error {
	if q.softDeleteCol == "" {
		return errors.New("model does not use soft deletes")
	}

	if !DispatchModelEvent(model, "restoring") {
		return ErrCancelled
	}

	val := reflect.ValueOf(model).Elem()
	var pkValue any
	var pkColumn string = "id"
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		dbTag := field.Tag.Get("db")
		gowTag := field.Tag.Get("gow")
		if strings.Contains(gowTag, "primaryKey") || dbTag == "id" {
			pkColumn = dbTag
			pkValue = val.Field(i).Interface()
			break
		}
	}

	if pkValue == nil {
		return sql.ErrNoRows
	}

	q.builder.Where(pkColumn, "=", pkValue)
	
	// We want to restore regardless of current trash state, so don't applySoftDeletes constraints here
	_, err := q.builder.Update(map[string]any{
		q.softDeleteCol: nil,
	})
	
	if err == nil {
		DispatchModelEvent(model, "restored")
		// Update the model struct
		for i := 0; i < val.Type().NumField(); i++ {
			field := val.Type().Field(i)
			dbTag := field.Tag.Get("db")
			if dbTag == "" {
				dbTag = strings.ToLower(field.Name)
			}
			if dbTag == q.softDeleteCol {
				f := val.Field(i)
				if f.Type() == reflect.TypeOf(time.Time{}) {
					f.Set(reflect.ValueOf(time.Time{}))
				} else if f.Type() == reflect.TypeOf(&time.Time{}) {
					f.Set(reflect.Zero(f.Type()))
				} else if f.Type() == reflect.TypeOf(sql.NullTime{}) {
					f.Set(reflect.ValueOf(sql.NullTime{Valid: false}))
				}
				break
			}
		}
	}
	return err
}

// Save inserts a new model or updates an existing one based on primary key presence.
func (q *ModelQuery[T]) Save(model *T) error {
	val := reflect.ValueOf(model).Elem()
	typ := val.Type()
	
	var isNew bool = true
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		gowTag := field.Tag.Get("gow")
		
		if strings.Contains(gowTag, "primaryKey") || dbTag == "id" {
			v := val.Field(i).Interface()
			// Basic check for empty zero values
			switch val.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v.(int64) != 0 {
					isNew = false
				}
			case reflect.String:
				if v.(string) != "" {
					isNew = false
				}
			}
			break
		}
	}
	
	if isNew {
		return q.Insert(model)
	}
	return q.Update(model)
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
