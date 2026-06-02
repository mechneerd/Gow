package schema

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mechneerd/gow/database/query"
)

// SchemaIndex represents an index on a table.
type SchemaIndex struct {
	Name    string
	Columns []string
	Unique  bool
	Type    string
}

// SchemaForeignKey represents a foreign key constraint.
type SchemaForeignKey struct {
	Name              string
	Columns           []string
	ReferencedTable   string
	ReferencedColumns []string
	OnDelete          string
	OnUpdate          string
}

// SchemaTable represents a database table schema.
type SchemaTable struct {
	Name        string
	Columns     []ColumnDefinition
	Indexes     []SchemaIndex
	ForeignKeys []SchemaForeignKey
	Comment     string
}

// GetColumn returns a column by name.
func (t *SchemaTable) GetColumn(name string) *ColumnDefinition {
	for i := range t.Columns {
		if t.Columns[i].Name == name {
			return &t.Columns[i]
		}
	}
	return nil
}

// HasColumn checks if a column exists.
func (t *SchemaTable) HasColumn(name string) bool {
	return t.GetColumn(name) != nil
}

// ColumnNames returns all column names.
func (t *SchemaTable) ColumnNames() []string {
	var names []string
	for _, col := range t.Columns {
		names = append(names, col.Name)
	}
	return names
}

// ColumnCount returns the number of columns.
func (t *SchemaTable) ColumnCount() int {
	return len(t.Columns)
}

// HasIndex checks if an index exists.
func (t *SchemaTable) HasIndex(name string) bool {
	for _, idx := range t.Indexes {
		if idx.Name == name {
			return true
		}
	}
	return false
}

// HasForeignKey checks if a foreign key exists.
func (t *SchemaTable) HasForeignKey(name string) bool {
	for _, fk := range t.ForeignKeys {
		if fk.Name == name {
			return true
		}
	}
	return false
}

// AlterTable represents an ALTER TABLE operation.
type AlterTable struct {
	tableName      string
	addColumns     []ColumnDefinition
	dropColumns    []string
	renameColumns  []RenameColumnOp
	modifyColumns  []ModifyColumnOp
	addIndexes     []SchemaIndex
	dropIndexes    []string
	addForeignKeys []SchemaForeignKey
	dropForeignKeys []string
}

// RenameColumnOp represents a rename column operation.
type RenameColumnOp struct {
	From string
	To   string
}

// ModifyColumnOp represents a modify column operation.
type ModifyColumnOp struct {
	Name string
	Type string
}

// NewAlterTable creates a new AlterTable.
func NewAlterTable(tableName string) *AlterTable {
	return &AlterTable{tableName: tableName}
}

// AddColumn adds a new column.
func (a *AlterTable) AddColumn(col ColumnDefinition) *AlterTable {
	a.addColumns = append(a.addColumns, col)
	return a
}

// DropColumn drops a column.
func (a *AlterTable) DropColumn(name string) *AlterTable {
	a.dropColumns = append(a.dropColumns, name)
	return a
}

// RenameColumn renames a column.
func (a *AlterTable) RenameColumn(from, to string) *AlterTable {
	a.renameColumns = append(a.renameColumns, RenameColumnOp{From: from, To: to})
	return a
}

// ModifyColumn modifies a column type.
func (a *AlterTable) ModifyColumn(name, newType string) *AlterTable {
	a.modifyColumns = append(a.modifyColumns, ModifyColumnOp{Name: name, Type: newType})
	return a
}

// AddIndex adds a new index.
func (a *AlterTable) AddIndex(idx SchemaIndex) *AlterTable {
	a.addIndexes = append(a.addIndexes, idx)
	return a
}

// DropIndex drops an index.
func (a *AlterTable) DropIndex(name string) *AlterTable {
	a.dropIndexes = append(a.dropIndexes, name)
	return a
}

// AddForeignKey adds a foreign key constraint.
func (a *AlterTable) AddForeignKey(fk SchemaForeignKey) *AlterTable {
	a.addForeignKeys = append(a.addForeignKeys, fk)
	return a
}

// DropForeignKey drops a foreign key constraint.
func (a *AlterTable) DropForeignKey(name string) *AlterTable {
	a.dropForeignKeys = append(a.dropForeignKeys, name)
	return a
}

// ToSQL generates the ALTER TABLE SQL.
func (a *AlterTable) ToSQL(dialect string) string {
	var stmts []string

	for _, col := range a.addColumns {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;",
			a.tableName, col.Name, col.Type))
	}

	for _, name := range a.dropColumns {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", a.tableName, name))
	}

	for _, rc := range a.renameColumns {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;",
			a.tableName, rc.From, rc.To))
	}

	for _, mc := range a.modifyColumns {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s;",
			a.tableName, mc.Name, mc.Type))
	}

	for _, idx := range a.addIndexes {
		unique := ""
		if idx.Unique {
			unique = "UNIQUE "
		}
		stmts = append(stmts, fmt.Sprintf("CREATE %sINDEX %s ON %s (%s);",
			unique, idx.Name, a.tableName, strings.Join(idx.Columns, ", ")))
	}

	for _, name := range a.dropIndexes {
		stmts = append(stmts, fmt.Sprintf("DROP INDEX %s ON %s;", name, a.tableName))
	}

	for _, fk := range a.addForeignKeys {
		stmts = append(stmts, fmt.Sprintf("ALTER TABLE %s ADD FOREIGN KEY (%s) REFERENCES %s(%s);",
			a.tableName,
			strings.Join(fk.Columns, ", "),
			fk.ReferencedTable,
			strings.Join(fk.ReferencedColumns, ", ")))
	}

	return strings.Join(stmts, "\n")
}

// ColumnInspector provides methods to inspect existing table columns.
type ColumnInspector struct {
	db      query.QueryExecer
	dialect string
}

// NewColumnInspector creates a new ColumnInspector.
func NewColumnInspector(db query.QueryExecer, dialect string) *ColumnInspector {
	return &ColumnInspector{db: db, dialect: dialect}
}

// GetColumns returns the columns of a table.
func (ci *ColumnInspector) GetColumns(tableName string) ([]ColumnDefinition, error) {
	var queryStr string
	switch ci.dialect {
	case "mysql":
		queryStr = fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)
	case "postgres":
		queryStr = fmt.Sprintf(`SELECT column_name, data_type, is_nullable, column_default 
			FROM information_schema.columns WHERE table_name = '%s'`, tableName)
	case "sqlite":
		queryStr = fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", ci.dialect)
	}

	rows, err := ci.db.QueryContext(context.Background(), queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnDefinition
	for rows.Next() {
		var col ColumnDefinition
		var nullable string
		var defaultVal *string
		err := rows.Scan(&col.Name, &col.Type, &nullable, &defaultVal)
		if err != nil {
			continue
		}
		col.Nullable = nullable == "YES" || nullable == "1"
		if defaultVal != nil {
			col.Default = *defaultVal
		}
		columns = append(columns, col)
	}

	return columns, nil
}

// HasTable checks if a table exists.
func (ci *ColumnInspector) HasTable(tableName string) bool {
	var queryStr string
	switch ci.dialect {
	case "mysql":
		queryStr = "SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?"
	case "postgres":
		queryStr = fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = '%s'", tableName)
	case "sqlite":
		queryStr = "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
	default:
		return false
	}

	var count int
	ci.db.QueryRowContext(context.Background(), queryStr).Scan(&count)
	return count > 0
}

// HasColumn checks if a column exists in a table.
func (ci *ColumnInspector) HasColumn(tableName, columnName string) bool {
	cols, err := ci.GetColumns(tableName)
	if err != nil {
		return false
	}
	for _, col := range cols {
		if col.Name == columnName {
			return true
		}
	}
	return false
}

// GetIndexes returns the indexes of a table.
func (ci *ColumnInspector) GetIndexes(tableName string) ([]SchemaIndex, error) {
	var queryStr string
	switch ci.dialect {
	case "mysql":
		queryStr = fmt.Sprintf("SHOW INDEX FROM %s", tableName)
	case "postgres":
		queryStr = fmt.Sprintf(`SELECT indexname, indexdef FROM pg_indexes WHERE tablename = '%s'`, tableName)
	case "sqlite":
		queryStr = fmt.Sprintf("PRAGMA index_list(%s)", tableName)
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", ci.dialect)
	}

	rows, err := ci.db.QueryContext(context.Background(), queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []SchemaIndex
	for rows.Next() {
		var idx SchemaIndex
		rows.Scan(&idx.Name, &idx.Columns)
		indexes = append(indexes, idx)
	}

	return indexes, nil
}

// GetCurrentTimestamp returns the current database timestamp.
func (ci *ColumnInspector) GetCurrentTimestamp() time.Time {
	var queryStr string
	switch ci.dialect {
	case "mysql", "postgres":
		queryStr = "SELECT NOW()"
	case "sqlite":
		queryStr = "SELECT datetime('now')"
	default:
		return time.Now()
	}

	var t time.Time
	ci.db.QueryRowContext(context.Background(), queryStr).Scan(&t)
	return t
}
