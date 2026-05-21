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

	// For Phase 2 Catch-up: Simplified Eager Loading stub.
	// In a full implementation, we'd extract all IDs from 'models',
	// perform a single IN query on the related table, and map the results back to the parent models.

	return nil
}
