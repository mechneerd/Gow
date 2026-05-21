package schema

import (
	"database/sql"
	"fmt"
	"gow/database/dialect"
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

// Create executes a CREATE TABLE statement.
func (s *Builder) Create(table string, callback func(*Blueprint)) error {
	blueprint := NewBlueprint(table)
	callback(blueprint)

	// Translate Blueprint to SQL using Dialect (Simplified for Phase 2)
	// In a full implementation, the Dialect interface would have a CompileCreate method.
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
			// Extremely naive SQLite primary key for now
			colType = "INTEGER PRIMARY KEY AUTOINCREMENT"
		}

		sqlStr.WriteString(fmt.Sprintf("    %s %s", s.dialect.QuoteIdentifier(col.Name), colType))
		
		if !col.Nullable && !col.Primary {
			sqlStr.WriteString(" NOT NULL")
		}
		if col.Unique {
			sqlStr.WriteString(" UNIQUE")
		}
	}
	sqlStr.WriteString("\n)")

	_, err := s.conn.Exec(sqlStr.String())
	return err
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
