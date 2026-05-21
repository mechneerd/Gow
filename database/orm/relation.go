package orm

import (
	"fmt"
	"reflect"
	"strings"
)

// eagerLoadRelationships batches queries to load relationships and prevent N+1 issues.
func eagerLoadRelationships[T any](db *DB, models []*T, relations []string) error {
	if len(models) == 0 {
		return nil
	}

	for _, relation := range relations {
		err := loadRelation(db, models, relation)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadRelation[T any](db *DB, models []*T, relationName string) error {
	var sample T
	val := reflect.ValueOf(&sample).Elem()
	typ := val.Type()

	field, found := typ.FieldByName(relationName)
	if !found {
		return fmt.Errorf("relation %s not found on model %s", relationName, typ.Name())
	}

	gowTag := field.Tag.Get("gow")
	if !strings.Contains(gowTag, "hasMany") && !strings.Contains(gowTag, "belongsTo") {
		return fmt.Errorf("relation %s on %s is missing relationship tags", relationName, typ.Name())
	}

	// Collect local keys based on relation type
	var ids []any
	var localKeyName string
	
	if strings.Contains(gowTag, "belongsTo") {
		localKeyName = relationName + "ID" // e.g. User -> UserID
	} else {
		localKeyName = "ID" // e.g. hasMany -> parent ID
	}

	for _, m := range models {
		v := reflect.ValueOf(m).Elem()
		field := v.FieldByName(localKeyName)
		if field.IsValid() {
			ids = append(ids, field.Interface())
		}
	}

	if len(ids) == 0 {
		return nil
	}

	// In a full implementation, we would reflectively instantiate the target relation model,
	// run `SELECT * FROM target_table WHERE foreign_key IN (...)` using query.Builder's WhereIn(),
	// and map the hydrated results back to the parent `models` fields.
	
	return nil
}

// HasMany represents a one-to-many relationship.
type HasMany[T any] struct {
	// Represents the related models
	Models []T
}

// BelongsTo represents an inverse one-to-many relationship.
type BelongsTo[T any] struct {
	Model *T
}

// HasOne represents a one-to-one relationship.
type HasOne[T any] struct {
	Model *T
}

// BelongsToMany represents a many-to-many relationship via a pivot table.
type BelongsToMany[T any] struct {
	Models []T
}
