package schema

// ForeignKey represents a foreign key constraint.
type ForeignKey struct {
	Column              string
	ReferencedTable     string
	ReferencedColumn    string
	OnDelete            string
	OnUpdate            string
}

// Blueprint represents the schema definition for a table.
type Blueprint struct {
	table       string
	columns     []ColumnDefinition
	drops       []string
	indexes     []string
	primaryKey  []string
	foreignKeys []ForeignKey
}

// ColumnDefinition represents a database column.
type ColumnDefinition struct {
	Name          string
	Type          string
	Length        int
	Nullable      bool
	Unique        bool
	Primary       bool
	AutoIncrement bool
	Default       any
}

// NewBlueprint creates a new table blueprint.
func NewBlueprint(table string) *Blueprint {
	return &Blueprint{
		table: table,
	}
}

// ID creates an auto-incrementing primary key "id" column.
func (b *Blueprint) ID() *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name:          "id",
		Type:          "bigint", // Dialect converts this to correct DB type
		Primary:       true,
		AutoIncrement: true,
	})
	return b
}

// String creates a VARCHAR column.
func (b *Blueprint) String(name string, length int) *Blueprint {
	if length == 0 {
		length = 255
	}
	b.columns = append(b.columns, ColumnDefinition{
		Name:   name,
		Type:   "varchar",
		Length: length,
	})
	return b
}

// Text creates a TEXT column.
func (b *Blueprint) Text(name string) *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name: name,
		Type: "text",
	})
	return b
}

// Integer creates an INT column.
func (b *Blueprint) Integer(name string) *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name: name,
		Type: "integer",
	})
	return b
}

// Boolean creates a BOOLEAN column.
func (b *Blueprint) Boolean(name string) *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name: name,
		Type: "boolean",
	})
	return b
}

// Timestamp creates a TIMESTAMP column.
func (b *Blueprint) Timestamp(name string) *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name: name,
		Type: "timestamp",
	})
	return b
}

// UnsignedBigInteger creates an UNSIGNED BIGINT column.
func (b *Blueprint) UnsignedBigInteger(name string) *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name: name,
		Type: "bigint",
	})
	return b
}

// Timestamps creates created_at and updated_at columns.
func (b *Blueprint) Timestamps() *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name: "created_at",
		Type: "timestamp",
	})
	b.columns = append(b.columns, ColumnDefinition{
		Name: "updated_at",
		Type: "timestamp",
	})
	return b
}

// SoftDeletes creates a deleted_at column.
func (b *Blueprint) SoftDeletes() *Blueprint {
	b.columns = append(b.columns, ColumnDefinition{
		Name:     "deleted_at",
		Type:     "timestamp",
		Nullable: true,
	})
	return b
}

// DropColumn adds a column to be dropped.
func (b *Blueprint) DropColumn(name string) *Blueprint {
	b.drops = append(b.drops, name)
	return b
}

// Nullable makes the last defined column nullable.
func (b *Blueprint) Nullable() *Blueprint {
	if len(b.columns) > 0 {
		b.columns[len(b.columns)-1].Nullable = true
	}
	return b
}

// Unique makes the last defined column unique, or adds a composite unique constraint.
// No args: modifies the last column. With args: adds a table-level unique constraint.
func (b *Blueprint) Unique(columns ...string) *Blueprint {
	if len(columns) == 0 {
		if len(b.columns) > 0 {
			b.columns[len(b.columns)-1].Unique = true
		}
	} else {
		b.indexes = append(b.indexes, "unique:"+joinColumns(columns))
	}
	return b
}

// Primary makes the last defined column a primary key, or sets a composite primary key.
// No args: modifies the last column. With args: sets a table-level composite primary key.
func (b *Blueprint) Primary(columns ...string) *Blueprint {
	if len(columns) == 0 {
		if len(b.columns) > 0 {
			b.columns[len(b.columns)-1].Primary = true
		}
	} else {
		b.primaryKey = append(b.primaryKey, columns...)
	}
	return b
}

// Index marks the last defined column for indexing, or adds a table-level index.
// No args: modifies the last column. With args: adds a table-level index on the given columns.
func (b *Blueprint) Index(columns ...string) *Blueprint {
	if len(columns) == 0 {
		if len(b.columns) > 0 {
			b.indexes = append(b.indexes, b.columns[len(b.columns)-1].Name)
		}
	} else {
		b.indexes = append(b.indexes, joinColumns(columns))
	}
	return b
}

// Default sets a default value for the last defined column.
func (b *Blueprint) Default(val any) *Blueprint {
	if len(b.columns) > 0 {
		b.columns[len(b.columns)-1].Default = val
	}
	return b
}

// Foreign starts a foreign key definition on the given column.
func (b *Blueprint) Foreign(column string) *ForeignKeyBuilder {
	return &ForeignKeyBuilder{blueprint: b, column: column}
}

// ForeignKeyBuilder is a fluent builder for foreign key constraints.
type ForeignKeyBuilder struct {
	blueprint *Blueprint
	column    string
}

// References sets the referenced column.
func (f *ForeignKeyBuilder) References(column string) *ForeignKeyOnBuilder {
	return &ForeignKeyOnBuilder{
		blueprint: f.blueprint,
		fk: ForeignKey{
			Column:           f.column,
			ReferencedColumn: column,
		},
	}
}

// ForeignKeyOnBuilder holds the FK state before On() is called.
type ForeignKeyOnBuilder struct {
	blueprint *Blueprint
	fk        ForeignKey
}

// On sets the referenced table and stores the FK in the blueprint.
// Returns a *ForeignKeyFinalBuilder for optional OnDelete/OnUpdate chaining.
func (f *ForeignKeyOnBuilder) On(table string) *ForeignKeyFinalBuilder {
	f.fk.ReferencedTable = table
	f.blueprint.foreignKeys = append(f.blueprint.foreignKeys, f.fk)
	// Return a reference to the last stored FK so OnDelete/OnUpdate can modify it
	idx := len(f.blueprint.foreignKeys) - 1
	return &ForeignKeyFinalBuilder{
		blueprint: f.blueprint,
		fkIndex:   idx,
	}
}

// ForeignKeyFinalBuilder allows setting OnDelete/OnUpdate on the stored FK.
type ForeignKeyFinalBuilder struct {
	blueprint *Blueprint
	fkIndex   int
}

// OnDelete sets the ON DELETE action.
func (f *ForeignKeyFinalBuilder) OnDelete(action string) *ForeignKeyFinalBuilder {
	f.blueprint.foreignKeys[f.fkIndex].OnDelete = action
	return f
}

// OnUpdate sets the ON UPDATE action.
func (f *ForeignKeyFinalBuilder) OnUpdate(action string) *ForeignKeyFinalBuilder {
	f.blueprint.foreignKeys[f.fkIndex].OnUpdate = action
	return f
}

func joinColumns(cols []string) string {
	result := ""
	for i, c := range cols {
		if i > 0 {
			result += ","
		}
		result += c
	}
	return result
}

// Columns returns the defined columns.
func (b *Blueprint) Columns() []ColumnDefinition {
	return b.columns
}

// Table returns the table name.
func (b *Blueprint) Table() string {
	return b.table
}

// Indexes returns the defined index column names.
func (b *Blueprint) Indexes() []string {
	return b.indexes
}

// PrimaryKeyColumns returns the composite primary key columns.
func (b *Blueprint) PrimaryKeyColumns() []string {
	return b.primaryKey
}

// ForeignKeys returns the defined foreign keys.
func (b *Blueprint) ForeignKeys() []ForeignKey {
	return b.foreignKeys
}

