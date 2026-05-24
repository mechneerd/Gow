package schema

// Blueprint represents the schema definition for a table.
type Blueprint struct {
	table   string
	columns []ColumnDefinition
	drops   []string
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

// Unique makes the last defined column unique.
func (b *Blueprint) Unique() *Blueprint {
	if len(b.columns) > 0 {
		b.columns[len(b.columns)-1].Unique = true
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

// Columns returns the defined columns.
func (b *Blueprint) Columns() []ColumnDefinition {
	return b.columns
}

// Table returns the table name.
func (b *Blueprint) Table() string {
	return b.table
}

