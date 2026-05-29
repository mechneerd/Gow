package orm

import (
	"reflect"
	"strings"
	"sync"
)

// pluralize converts a model name to a table name by lowercasing and appending "s".
// For production, consider using a proper inflection library.
func pluralize(name string) string {
	lower := strings.ToLower(name)
	// Basic English pluralization rules
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") ||
		strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return lower + "es"
	}
	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		prev := lower[len(lower)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return lower[:len(lower)-1] + "ies"
		}
	}
	return lower + "s"
}

// FieldMeta represents cached metadata for a struct field.
type FieldMeta struct {
	Name       string // struct field name
	Column     string // db tag value
	Index      int    // reflect field index for fast access
	IsPrimary  bool
	IsAuto     bool // auto-increment
}

// ModelMetadata holds cached reflection data for an ORM model.
type ModelMetadata struct {
	TableName     string
	PrimaryKey    string
	SoftDeleteCol string // empty string if no soft delete field
	Fields        []FieldMeta
}

// metaCache stores the metadata for each reflect.Type to avoid reflecting on every query.
var metaCache sync.Map // map[reflect.Type]*ModelMetadata

// getMetadata parses a model type and caches its structural metadata.
func getMetadata(typ reflect.Type) *ModelMetadata {
	// Ensure we're working with the struct type, not a pointer
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if cached, ok := metaCache.Load(typ); ok {
		return cached.(*ModelMetadata)
	}

	meta := &ModelMetadata{
		TableName: pluralize(typ.Name()),
		Fields:    make([]FieldMeta, 0, typ.NumField()),
	}

	// Check if type implements Model interface (TableName() method) without instantiating
	modelType := reflect.TypeOf((*Model)(nil)).Elem()
	ptrType := reflect.PointerTo(typ)
	if ptrType.Implements(modelType) || typ.Implements(modelType) {
		// Need a dummy instance to call TableName() — only once per type
		dummyPtr := reflect.New(typ).Interface()
		if m, ok := dummyPtr.(Model); ok {
			meta.TableName = m.TableName()
		}
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		
		dbTag := field.Tag.Get("db")
		gowTag := field.Tag.Get("gow")

		column := dbTag
		if column == "" {
			column = strings.ToLower(field.Name)
		}

		isPrimary := strings.Contains(gowTag, "primaryKey") || dbTag == "id"
		isAuto := isPrimary // simplistic assumption, can be expanded based on tags

		if isPrimary {
			meta.PrimaryKey = column
		}

		if strings.Contains(gowTag, "softDelete") || column == "deleted_at" {
			meta.SoftDeleteCol = column
		}

		meta.Fields = append(meta.Fields, FieldMeta{
			Name:      field.Name,
			Column:    column,
			Index:     i,
			IsPrimary: isPrimary,
			IsAuto:    isAuto,
		})
	}

	if meta.PrimaryKey == "" {
		meta.PrimaryKey = "id" // Fallback default
	}

	metaCache.Store(typ, meta)
	return meta
}

// GetTableName resolves the table name for a model instance using metadata.
func GetTableName(model any) string {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	meta := getMetadata(typ)
	return meta.TableName
}

// ClearCache removes cached metadata for the given model type.
// Call this after altering struct tags or table mappings at runtime.
func ClearCache(modelType any) {
	typ := reflect.TypeOf(modelType)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	metaCache.Delete(typ)
}

