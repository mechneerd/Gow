package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gow/database/pagination"
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
	tx      *sql.Tx // non-nil when this instance represents an active transaction
}

// Transaction executes a function within a database transaction.
func (db *DB) Transaction(ctx context.Context, fn func(txDB *DB) error) error {
	txDB, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			txDB.Rollback()
			panic(p)
		}
	}()

	if err := fn(txDB); err != nil {
		txDB.Rollback()
		return err
	}

	return txDB.Commit()
}

// TransactionFn is a convenience method that runs the given function inside a transaction
// using context.Background(). Useful for simple cases.
func (db *DB) TransactionFn(fn func(txDB *DB) error) error {
	return db.Transaction(context.Background(), fn)
}

// Begin starts a new database transaction and returns a new *DB instance
// that operates within that transaction.
func (db *DB) Begin(ctx context.Context) (*DB, error) {
	if db.tx != nil {
		return nil, errors.New("cannot begin nested transaction (use savepoints for nesting)")
	}

	sqlDB, ok := db.Conn.(*sql.DB)
	if !ok {
		return nil, errors.New("cannot start transaction: connection is not a root *sql.DB")
	}

	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	txDB := &DB{
		Conn:    tx,
		Builder: db.Builder.WithConn(tx),
		tx:      tx,
	}
	return txDB, nil
}

// BeginTx is a convenience method equivalent to Begin(context.Background()).
func (db *DB) BeginTx() (*DB, error) {
	return db.Begin(context.Background())
}

// Commit commits the current transaction.
// It only works if this DB instance was created via Begin() or Transaction().
func (db *DB) Commit() error {
	if db.tx == nil {
		return errors.New("Commit called on non-transactional DB")
	}
	return db.tx.Commit()
}

// Rollback rolls back the current transaction.
func (db *DB) Rollback() error {
	if db.tx == nil {
		return errors.New("Rollback called on non-transactional DB")
	}
	return db.tx.Rollback()
}

// InTransaction returns true if this DB instance is operating inside a transaction.
func (db *DB) InTransaction() bool {
	return db.tx != nil
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

func NewQuery[T any](db *DB) *ModelQuery[T] {
	meta := getMetadata(reflect.TypeOf((*T)(nil)))
	
	builder := db.Builder.Clone().Table(meta.TableName)
	
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
		softDeleteCol: meta.SoftDeleteCol,
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

// Paginate returns a length-aware paginator (with total count).
func (q *ModelQuery[T]) Paginate(perPage, page int) *pagination.LengthAwarePaginator[T] {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	// Clone for count
	countQ := q.builder.Clone()
	total, _ := countQ.Count("*")

	// Apply pagination
	offset := (page - 1) * perPage
	q.builder.Offset(offset).Limit(perPage)

	items, _ := q.Get() // reuse Get which handles soft deletes + eager loading

	// Convert []*T to []T for the paginator (or keep as is)
	var plainItems []T
	for _, item := range items {
		if item != nil {
			plainItems = append(plainItems, *item)
		}
	}

	return pagination.NewLengthAwarePaginator(plainItems, total, perPage, page, "")
}

// SimplePaginate returns a simple paginator (no total count).
func (q *ModelQuery[T]) SimplePaginate(perPage, page int) *pagination.SimplePaginator[T] {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	offset := (page - 1) * perPage
	q.builder.Offset(offset).Limit(perPage + 1) // fetch one extra to detect if there's a next page

	items, _ := q.Get()

	hasMore := len(items) > perPage
	if hasMore {
		items = items[:perPage]
	}

	var plainItems []T
	for _, item := range items {
		if item != nil {
			plainItems = append(plainItems, *item)
		}
	}

	return pagination.NewSimplePaginator(plainItems, perPage, page, "", hasMore)
}

// CursorPaginate is a basic cursor-based paginator (keyset pagination).
// For production use, you should pass a proper cursor column (usually `id` or `created_at`).
func (q *ModelQuery[T]) CursorPaginate(perPage int, after any) *pagination.CursorPaginator[T] {
	if perPage < 1 {
		perPage = 15
	}

	if after != nil {
		q.builder.Where("id", ">", after) // simplistic cursor on 'id'
	}

	q.builder.Limit(perPage + 1)

	items, _ := q.Get()

	hasMore := len(items) > perPage
	if hasMore {
		items = items[:perPage]
	}

	var plainItems []T
	for _, item := range items {
		if item != nil {
			plainItems = append(plainItems, *item)
		}
	}

	var nextCursor *string
	if hasMore && len(plainItems) > 0 {
		last := plainItems[len(plainItems)-1]
		v := reflect.ValueOf(last)
		if v.Kind() == reflect.Struct {
			if idField := v.FieldByName("ID"); idField.IsValid() {
				s := fmt.Sprintf("%v", idField.Interface())
				nextCursor = &s
			}
		}
	}

	return &pagination.CursorPaginator[T]{
		Items:      plainItems,
		NextCursor: nextCursor,
		PerPage:    perPage,
	}
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

	meta := getMetadata(typ)

	for _, field := range meta.Fields {
		if field.Column == "" || field.Column == "-" {
			continue
		}

		if field.IsAuto {
			continue
		}
		
		// Skip primary key if it's zero
		if field.IsPrimary {
			v := val.Field(field.Index)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v.Int() == 0 {
					continue
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if v.Uint() == 0 {
					continue
				}
			case reflect.String:
				if v.String() == "" {
					continue
				}
			}
		}

		// Since we don't store full tags in FieldMeta to save space, we can either
		// re-fetch the tag or we could have added IsCreateTime to FieldMeta.
		// For simplicity, we re-fetch the gow tag from reflection here.
		gowTag := typ.Field(field.Index).Tag.Get("gow")
		if strings.Contains(gowTag, "autoCreateTime") || strings.Contains(gowTag, "autoUpdateTime") {
			values[field.Column] = now
			val.Field(field.Index).Set(reflect.ValueOf(now))
			continue
		}

		values[field.Column] = val.Field(field.Index).Interface()
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
	
	meta := getMetadata(typ)
	
	var pkValue any
	pkColumn := meta.PrimaryKey

	for _, field := range meta.Fields {
		if field.Column == "" || field.Column == "-" {
			continue
		}

		if field.IsPrimary {
			pkValue = val.Field(field.Index).Interface()
			continue
		}

		gowTag := typ.Field(field.Index).Tag.Get("gow")

		if strings.Contains(gowTag, "autoUpdateTime") {
			values[field.Column] = now
			val.Field(field.Index).Set(reflect.ValueOf(now))
			continue
		}

		values[field.Column] = val.Field(field.Index).Interface()
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
		meta := getMetadata(typ)
		var pkValue any
		pkColumn := meta.PrimaryKey
		
		for _, field := range meta.Fields {
			if field.IsPrimary {
				pkValue = val.Field(field.Index).Interface()
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
			typ := val.Type()
			meta := getMetadata(typ)
			for _, field := range meta.Fields {
				if field.Column == q.softDeleteCol {
					// Handle sql.NullTime or time.Time pointer
					f := val.Field(field.Index)
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
	meta := getMetadata(val.Type())
	var pkValue any
	pkColumn := meta.PrimaryKey
	for _, field := range meta.Fields {
		if field.IsPrimary {
			pkValue = val.Field(field.Index).Interface()
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
		meta := getMetadata(val.Type())
		for _, field := range meta.Fields {
			if field.Column == q.softDeleteCol {
				f := val.Field(field.Index)
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
	
	meta := getMetadata(typ)
	var isNew bool = true
	
	for _, field := range meta.Fields {
		if field.IsPrimary {
			v := val.Field(field.Index)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v.Int() != 0 {
					isNew = false
				}
			case reflect.String:
				if v.String() != "" {
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
	meta := getMetadata(typ)

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	
	scanArgs := make([]any, len(columns))
	fieldMap := make(map[string]int)

	for _, field := range meta.Fields {
		fieldMap[field.Column] = field.Index
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
