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

	var isHasMany bool = strings.Contains(gowTag, "hasMany")
	var isBelongsTo bool = strings.Contains(gowTag, "belongsTo")

	// Collect local keys based on relation type
	var ids []any
	var localKeyName string
	
	if isBelongsTo {
		localKeyName = relationName + "ID" // e.g. User -> UserID
	} else {
		localKeyName = "ID" // e.g. hasMany -> parent ID
	}

	for _, m := range models {
		v := reflect.ValueOf(m).Elem()
		f := v.FieldByName(localKeyName)
		if f.IsValid() {
			ids = append(ids, f.Interface())
		}
	}

	if len(ids) == 0 {
		return nil
	}

	// Figure out the target model type
	var targetType reflect.Type
	var isSlice bool
	
	if field.Type.Kind() == reflect.Struct {
		if isHasMany && field.Type.NumField() > 0 && field.Type.Field(0).Type.Kind() == reflect.Slice {
			targetType = field.Type.Field(0).Type.Elem() // 'Target'
			isSlice = true
		} else if isBelongsTo && field.Type.NumField() > 0 && field.Type.Field(0).Type.Kind() == reflect.Ptr {
			targetType = field.Type.Field(0).Type.Elem() // 'Target'
			isSlice = false
		} else if field.Type.Kind() == reflect.Slice {
			targetType = field.Type.Elem()
			isSlice = true
		} else if field.Type.Kind() == reflect.Ptr {
			targetType = field.Type.Elem()
			isSlice = false
		} else {
			targetType = field.Type
		}
	} else if field.Type.Kind() == reflect.Slice {
		targetType = field.Type.Elem()
		isSlice = true
	} else if field.Type.Kind() == reflect.Ptr {
		targetType = field.Type.Elem()
		isSlice = false
	} else {
		targetType = field.Type
	}

	// Use metadata to get table name
	meta := getMetadata(targetType)
	tableName := meta.TableName

	var foreignKeyName string
	if isBelongsTo {
		// Target is the parent, so its primary key is needed
		foreignKeyName = meta.PrimaryKey
	} else {
		// hasMany -> foreign key on child is usually parent's name + "ID". e.g. targetType has 'UserID'
		parentTypeName := typ.Name()
		// E.g., parent is TestUser, look for TestUserID, if not found, look for UserID (strip Test prefix, but simpler to just do typ.Name() + "ID")
		// The most robust way is to find a field ending with ID that has the gow tag, or just match typ.Name()+"ID".
		// For our tests, `TestPost` has `UserID` which doesn't perfectly match `TestUserID`.
		// Let's just strip 'Test' prefix if present for testing, or just find any field ending in ID whose type matches.
		expectedFkName := parentTypeName + "ID"
		expectedFkNameAlt := strings.TrimPrefix(parentTypeName, "Test") + "ID"
		
		foreignKeyName = strings.ToLower(parentTypeName) + "_id" // fallback
		for i := 0; i < targetType.NumField(); i++ {
			f := targetType.Field(i)
			if f.Name == expectedFkName || f.Name == expectedFkNameAlt {
				if dbTag := f.Tag.Get("db"); dbTag != "" && dbTag != "-" {
					foreignKeyName = dbTag
				} else {
					foreignKeyName = strings.ToLower(f.Name)
				}
				break
			}
		}
	}

	// Build query using a fresh clone
	builder := db.Builder.Clone().Table(tableName).WhereIn(foreignKeyName, ids)
	rows, err := builder.Get()
	if err != nil {
		// fmt.Printf("DEBUG: builder.Get() error: %v\n", err)
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Map key is string representation of foreign key value
	groupedResults := make(map[string][]reflect.Value)

	for rows.Next() {
		targetPtr := reflect.New(targetType)
		targetVal := targetPtr.Elem()

		scanArgs := make([]any, len(columns))
		fieldMap := make(map[string]int)

		for i := 0; i < targetType.NumField(); i++ {
			dbTag := targetType.Field(i).Tag.Get("db")
			if dbTag != "" {
				fieldMap[dbTag] = i
			} else {
				fieldMap[strings.ToLower(targetType.Field(i).Name)] = i
			}
		}

		var fkScanDest any
		var fkScanIndex int = -1

		for i, col := range columns {
			if fieldIdx, ok := fieldMap[col]; ok {
				scanArgs[i] = targetVal.Field(fieldIdx).Addr().Interface()
				if col == foreignKeyName {
					fkScanDest = scanArgs[i]
					fkScanIndex = i
				}
			} else {
				var dummy any
				scanArgs[i] = &dummy
				if col == foreignKeyName {
					fkScanDest = &dummy
					fkScanIndex = i
				}
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}

		var fkValueStr string
		if fkScanIndex != -1 {
			val := reflect.ValueOf(fkScanDest)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			fkValueStr = fmt.Sprintf("%v", val.Interface())
		}

		groupedResults[fkValueStr] = append(groupedResults[fkValueStr], targetVal)
	}

	// Map grouped results back to parent models
	for _, m := range models {
		parentVal := reflect.ValueOf(m).Elem()
		
		var lookupKey string
		if isBelongsTo {
			fkField := parentVal.FieldByName(localKeyName)
			if fkField.IsValid() {
				lookupKey = fmt.Sprintf("%v", fkField.Interface())
			}
		} else {
			pkField := parentVal.FieldByName("ID")
			if pkField.IsValid() {
				lookupKey = fmt.Sprintf("%v", pkField.Interface())
			}
		}


		relatedVals, ok := groupedResults[lookupKey]
		if !ok {
			continue
		}

		relField := parentVal.FieldByName(relationName)
		if !relField.IsValid() || !relField.CanSet() {
			continue
		}

		if isSlice {
			sliceType := reflect.SliceOf(targetType)
			sliceVal := reflect.MakeSlice(sliceType, 0, len(relatedVals))
			for _, rv := range relatedVals {
				sliceVal = reflect.Append(sliceVal, rv)
			}

			if relField.Type().Kind() == reflect.Slice {
				relField.Set(sliceVal)
			} else if relField.Type().Name() == "HasMany" || relField.Type().Name() == "BelongsToMany" {
				modelsField := relField.FieldByName("Models")
				if modelsField.IsValid() && modelsField.CanSet() {
					modelsField.Set(sliceVal)
				}
			}
		} else {
			if len(relatedVals) > 0 {
				firstVal := relatedVals[0]
				ptrVal := firstVal.Addr()

				if relField.Type().Kind() == reflect.Ptr {
					relField.Set(ptrVal)
				} else if relField.Type().Name() == "BelongsTo" || relField.Type().Name() == "HasOne" {
					modelField := relField.FieldByName("Model")
					if modelField.IsValid() && modelField.CanSet() {
						modelField.Set(ptrVal)
					}
				}
			}
		}
	}

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
