package orm

import (
	"reflect"
	"strings"
	"sync"
)

var softDeleteCache sync.Map // map[reflect.Type]string

// getSoftDeleteColumn checks if a model type has a field tagged with gow:"softDelete".
// It caches the result for future lookups. Returns the column name (from "db" tag) or empty string.
func getSoftDeleteColumn(modelType reflect.Type) string {
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if col, ok := softDeleteCache.Load(modelType); ok {
		return col.(string)
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		gowTag := field.Tag.Get("gow")
		if strings.Contains(gowTag, "softDelete") {
			dbTag := field.Tag.Get("db")
			if dbTag == "" || dbTag == "-" {
				dbTag = strings.ToLower(field.Name) // Fallback if no db tag is specified
			}
			softDeleteCache.Store(modelType, dbTag)
			return dbTag
		}
	}

	softDeleteCache.Store(modelType, "")
	return ""
}
