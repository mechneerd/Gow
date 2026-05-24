package orm

import (
	"reflect"
	"strings"
	"sync"
)

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
		TableName: strings.ToLower(typ.Name()) + "s", // Default simplistic pluralization
		Fields:    make([]FieldMeta, 0, typ.NumField()),
	}

	// Try to instantiate to check for TableName() interface
	// This requires reflection to create a dummy value.
	// We handle TableName overriding in getTableName, but we can cache the default here.
	dummyPtr := reflect.New(typ).Interface()
	if m, ok := dummyPtr.(Model); ok {
		meta.TableName = m.TableName()
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

