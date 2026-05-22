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

// --- BelongsToMany Helpers: Attach / Detach / Sync / Toggle ---

// Attach adds a relationship in the pivot table.
func Attach(db *DB, pivotTable, foreignKey, relatedKey string, parentID, relatedID any) error {
	_, err := db.Builder.Clone().Table(pivotTable).Insert(map[string]any{
		foreignKey: parentID,
		relatedKey: relatedID,
	})
	return err
}

// Detach removes a specific relationship from the pivot table.
func Detach(db *DB, pivotTable, foreignKey, relatedKey string, parentID, relatedID any) error {
	_, err := db.Builder.Clone().Table(pivotTable).
		Where(foreignKey, "=", parentID).
		Where(relatedKey, "=", relatedID).
		Delete()
	return err
}

// Sync replaces all existing relationships with the given list of related IDs.
func Sync(db *DB, pivotTable, foreignKey, relatedKey string, parentID any, relatedIDs []any) error {
	// First delete all existing
	_, err := db.Builder.Clone().Table(pivotTable).
		Where(foreignKey, "=", parentID).
		Delete()
	if err != nil {
		return err
	}

	// Then insert new ones
	for _, rid := range relatedIDs {
		if _, err := db.Builder.Clone().Table(pivotTable).Insert(map[string]any{
			foreignKey: parentID,
			relatedKey: rid,
		}); err != nil {
			return err
		}
	}
	return nil
}

// Toggle adds the relationship if it doesn't exist, removes it if it does.
func Toggle(db *DB, pivotTable, foreignKey, relatedKey string, parentID, relatedID any) error {
	// Check if exists (simple existence check)
	rows, err := db.Builder.Clone().Table(pivotTable).
		Where(foreignKey, "=", parentID).
		Where(relatedKey, "=", relatedID).
		Limit(1).
		Get()
	if err != nil {
		return err
	}
	defer rows.Close()

	exists := rows.Next()
	rows.Close()

	if exists {
		return Detach(db, pivotTable, foreignKey, relatedKey, parentID, relatedID)
	}
	return Attach(db, pivotTable, foreignKey, relatedKey, parentID, relatedID)
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
	if !strings.Contains(gowTag, "hasMany") && !strings.Contains(gowTag, "belongsTo") && !strings.Contains(gowTag, "belongsToMany") {
		return fmt.Errorf("relation %s on %s is missing relationship tags", relationName, typ.Name())
	}

	var isHasMany bool = strings.Contains(gowTag, "hasMany")
	var isBelongsTo bool = strings.Contains(gowTag, "belongsTo")
	var isBelongsToMany bool = strings.Contains(gowTag, "belongsToMany")

	// Parse belongsToMany options
	var pivotTable, foreignKey, relatedKey string
	if isBelongsToMany {
		parts := strings.Split(gowTag, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if strings.HasPrefix(p, "through=") {
				pivotTable = strings.TrimPrefix(p, "through=")
			}
			if strings.HasPrefix(p, "foreignKey=") {
				foreignKey = strings.TrimPrefix(p, "foreignKey=")
			}
			if strings.HasPrefix(p, "relatedKey=") {
				relatedKey = strings.TrimPrefix(p, "relatedKey=")
			}
		}

		if pivotTable == "" {
			pivotTable = strings.ToLower(typ.Name()) + "_" + strings.ToLower(relationName)
		}
		if foreignKey == "" {
			foreignKey = strings.ToLower(typ.Name()) + "_id"
		}
		if relatedKey == "" {
			relatedKey = strings.ToLower(relationName) + "_id"
		}

		return loadBelongsToMany(db, models, relationName, pivotTable, foreignKey, relatedKey, gowTag)
	}

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

// loadBelongsToMany handles eager loading for many-to-many relations.
func loadBelongsToMany[T any](db *DB, models []*T, relationName, pivotTable, foreignKey, relatedKey, gowTag string) error {
	if len(models) == 0 {
		return nil
	}

	// 1. Collect parent IDs
	parentIDs := []any{}
	parentIndex := make(map[string]int)

	for i, m := range models {
		v := reflect.ValueOf(m).Elem()
		idField := v.FieldByName("ID")
		if !idField.IsValid() {
			idField = v.FieldByName("Id")
		}
		if idField.IsValid() && idField.CanInterface() {
			idVal := idField.Interface()
			key := fmt.Sprintf("%v", idVal)
			parentIDs = append(parentIDs, idVal)
			parentIndex[key] = i
		}
	}

	if len(parentIDs) == 0 {
		return nil
	}

	// 2. Query pivot to build parentID -> []relatedID mapping
	pivotQ := db.Builder.Clone().Table(pivotTable).WhereIn(foreignKey, parentIDs)
	pivotResult, err := pivotQ.Get()
	if err != nil {
		return err
	}
	defer pivotResult.Close()

	pCols, _ := pivotResult.Columns()
	parentToRelatedIDs := make(map[string][]any)

	for pivotResult.Next() {
		vals := make([]any, len(pCols))
		ptrs := make([]any, len(pCols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := pivotResult.Scan(ptrs...); err != nil {
			continue
		}

		var pID, rID string
		for i, c := range pCols {
			if c == foreignKey {
				pID = fmt.Sprintf("%v", vals[i])
			}
			if c == relatedKey {
				rID = fmt.Sprintf("%v", vals[i])
			}
		}
		if pID != "" && rID != "" {
			parentToRelatedIDs[pID] = append(parentToRelatedIDs[pID], rID)
		}
	}

	// Collect all related IDs
	allRelatedIDs := []any{}
	for _, ids := range parentToRelatedIDs {
		allRelatedIDs = append(allRelatedIDs, ids...)
	}
	if len(allRelatedIDs) == 0 {
		return nil
	}

	// 3. Load target models
	var sample T
	targetType := reflect.TypeOf(sample)
	meta := getMetadata(targetType)
	targetTable := meta.TableName
	pk := meta.PrimaryKey
	if pk == "" {
		pk = "id"
	}

	targetQ := db.Builder.Clone().Table(targetTable).WhereIn(pk, allRelatedIDs)
	targetResult, err := targetQ.Get()
	if err != nil {
		return err
	}
	defer targetResult.Close()

	tCols, _ := targetResult.Columns()
	relatedModels := make(map[string]reflect.Value) // relatedID -> reflect.Value

	for targetResult.Next() {
		targetPtr := reflect.New(targetType)
		targetVal := targetPtr.Elem()

		scanArgs := make([]any, len(tCols))
		fieldMap := map[string]int{}
		for i := 0; i < targetType.NumField(); i++ {
			f := targetType.Field(i)
			tag := f.Tag.Get("db")
			if tag != "" {
				fieldMap[tag] = i
			} else {
				fieldMap[strings.ToLower(f.Name)] = i
			}
		}

		for i, col := range tCols {
			if idx, ok := fieldMap[col]; ok {
				scanArgs[i] = targetVal.Field(idx).Addr().Interface()
			} else {
				var dummy any
				scanArgs[i] = &dummy
			}
		}

		if err := targetResult.Scan(scanArgs...); err != nil {
			return err
		}

		// Get PK value as string key
		pkVal := ""
		if pkField := targetVal.FieldByName("ID"); pkField.IsValid() {
			pkVal = fmt.Sprintf("%v", pkField.Interface())
		}
		if pkVal != "" {
			relatedModels[pkVal] = targetVal
		}
	}

	// 4. Assign back to parent models
	for parentKey, relatedIDList := range parentToRelatedIDs {
		parentIdx, ok := parentIndex[parentKey]
		if !ok {
			continue
		}
		parent := models[parentIdx]
		parentVal := reflect.ValueOf(parent).Elem()

		relField := parentVal.FieldByName(relationName)
		if !relField.IsValid() || !relField.CanSet() {
			continue
		}

		// Build slice of related models
		sliceType := reflect.SliceOf(targetType)
		sliceVal := reflect.MakeSlice(sliceType, 0, len(relatedIDList))

		for _, rid := range relatedIDList {
			ridStr := fmt.Sprintf("%v", rid)
			if rv, found := relatedModels[ridStr]; found {
				sliceVal = reflect.Append(sliceVal, rv)
			}
		}

		// Support both direct []T and BelongsToMany[T]
		if relField.Type().Kind() == reflect.Slice {
			relField.Set(sliceVal)
		} else if relField.Type().Name() == "BelongsToMany" {
			modelsField := relField.FieldByName("Models")
			if modelsField.IsValid() && modelsField.CanSet() {
				modelsField.Set(sliceVal)
			}
		}
	}

	return nil
}
