package schema

import (
	"database/sql"
	"fmt"
	"github.com/mechneerd/gow/database/dialect"
	"strings"
)

// Builder handles DDL execution against the database.
type Builder struct {
	conn    *sql.DB
	dialect dialect.Dialect
}

// NewBuilder creates a new Schema Builder.
func NewBuilder(conn *sql.DB, d dialect.Dialect) *Builder {
	return &Builder{
		conn:    conn,
		dialect: d,
	}
}

// Create executes a CREATE TABLE statement, then creates any indexes separately.
func (s *Builder) Create(table string, callback func(*Blueprint)) error {
	blueprint := NewBlueprint(table)
	callback(blueprint)

	var sqlStr strings.Builder
	sqlStr.WriteString("CREATE TABLE ")
	sqlStr.WriteString(s.dialect.QuoteIdentifier(table))
	sqlStr.WriteString(" (\n")

	for i, col := range blueprint.Columns() {
		if i > 0 {
			sqlStr.WriteString(",\n")
		}
		
		// Map generic types to dialect specific ones
		colType := col.Type
		if col.Type == "varchar" {
			colType = fmt.Sprintf("VARCHAR(%d)", col.Length)
		} else if col.Type == "bigint" && col.AutoIncrement {
			colType = s.dialect.AutoIncrementSQL()
		}

		sqlStr.WriteString(fmt.Sprintf("    %s %s", s.dialect.QuoteIdentifier(col.Name), colType))
		
		if !col.Nullable && !col.Primary {
			sqlStr.WriteString(" NOT NULL")
		}
		if col.Unique {
			sqlStr.WriteString(" UNIQUE")
		}
	}

	// Composite primary key
	if len(blueprint.PrimaryKeyColumns()) > 0 {
		sqlStr.WriteString(",\n    PRIMARY KEY (")
		for i, col := range blueprint.PrimaryKeyColumns() {
			if i > 0 {
				sqlStr.WriteString(", ")
			}
			sqlStr.WriteString(s.dialect.QuoteIdentifier(col))
		}
		sqlStr.WriteString(")")
	}

	// Table-level unique constraints (inline)
	for _, idx := range blueprint.Indexes() {
		if strings.HasPrefix(idx, "unique:") {
			cols := strings.TrimPrefix(idx, "unique:")
			sqlStr.WriteString(",\n    UNIQUE (")
			parts := strings.Split(cols, ",")
			for i, c := range parts {
				if i > 0 {
					sqlStr.WriteString(", ")
				}
				sqlStr.WriteString(s.dialect.QuoteIdentifier(strings.TrimSpace(c)))
			}
			sqlStr.WriteString(")")
		}
	}

	// Foreign keys (inline)
	for _, fk := range blueprint.ForeignKeys() {
		sqlStr.WriteString(",\n    FOREIGN KEY (")
		sqlStr.WriteString(s.dialect.QuoteIdentifier(fk.Column))
		sqlStr.WriteString(") REFERENCES ")
		sqlStr.WriteString(s.dialect.QuoteIdentifier(fk.ReferencedTable))
		sqlStr.WriteString("(")
		sqlStr.WriteString(s.dialect.QuoteIdentifier(fk.ReferencedColumn))
		sqlStr.WriteString(")")
		if fk.OnDelete != "" {
			sqlStr.WriteString(" ON DELETE ")
			sqlStr.WriteString(strings.ToUpper(fk.OnDelete))
		}
		if fk.OnUpdate != "" {
			sqlStr.WriteString(" ON UPDATE ")
			sqlStr.WriteString(strings.ToUpper(fk.OnUpdate))
		}
	}

	sqlStr.WriteString("\n)")

	if _, err := s.conn.Exec(sqlStr.String()); err != nil {
		return err
	}

	// Create non-unique indexes separately (not supported inline in SQLite)
	for _, idx := range blueprint.Indexes() {
		if !strings.HasPrefix(idx, "unique:") {
			cols := strings.Split(idx, ",")
			colNames := make([]string, len(cols))
			for i, c := range cols {
				colNames[i] = s.dialect.QuoteIdentifier(strings.TrimSpace(c))
			}
			idxName := fmt.Sprintf("idx_%s_%s", table, strings.Join(cols, "_"))
			createIdx := fmt.Sprintf("CREATE INDEX %s ON %s (%s)",
				s.dialect.QuoteIdentifier(idxName),
				s.dialect.QuoteIdentifier(table),
				strings.Join(colNames, ", "),
			)
			if _, err := s.conn.Exec(createIdx); err != nil {
				return err
			}
		}
	}

	return nil
}

// Drop executes a DROP TABLE statement.
func (s *Builder) Drop(table string) error {
	query := fmt.Sprintf("DROP TABLE %s", s.dialect.QuoteIdentifier(table))
	_, err := s.conn.Exec(query)
	return err
}

// DropIfExists executes a DROP TABLE IF EXISTS statement.
func (s *Builder) DropIfExists(table string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", s.dialect.QuoteIdentifier(table))
	_, err := s.conn.Exec(query)
	return err
}

