package dialect

import (
	"fmt"
	"strings"
)

// SQLiteDialect implements the Dialect interface for SQLite.
type SQLiteDialect struct{}

func (d *SQLiteDialect) QuoteIdentifier(name string) string {
	if name == "*" {
		return name
	}
	// Basic quoting, assuming no internal double quotes for simplicity
	return `"` + name + `"`
}

func (d *SQLiteDialect) Placeholder(index int) string {
	return "?"
}

func (d *SQLiteDialect) CompileSelect(query SelectQuery) (string, []any) {
	var sql strings.Builder
	var args []any

	// SELECT
	sql.WriteString("SELECT ")
	if len(query.Columns) == 0 {
		sql.WriteString("*")
	} else {
		for i, col := range query.Columns {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(d.QuoteIdentifier(col))
		}
	}

	// FROM
	sql.WriteString(" FROM ")
	sql.WriteString(d.QuoteIdentifier(query.Table))

	// WHERE
	if len(query.Wheres) > 0 {
		sql.WriteString(" WHERE ")
		for i, w := range query.Wheres {
			if i > 0 {
				sql.WriteString(" " + w.Boolean + " ")
			}
			sql.WriteString(d.QuoteIdentifier(w.Column))
			sql.WriteString(" " + w.Operator + " ")
			sql.WriteString(d.Placeholder(len(args) + 1))
			args = append(args, w.Value)
		}
	}

	// ORDER BY
	if len(query.OrderBys) > 0 {
		sql.WriteString(" ORDER BY ")
		for i, o := range query.OrderBys {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(d.QuoteIdentifier(o.Column) + " " + o.Direction)
		}
	}

	// LIMIT / OFFSET
	if query.Limit != nil {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", *query.Limit))
	}
	if query.Offset != nil {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", *query.Offset))
	}

	return sql.String(), args
}

func (d *SQLiteDialect) CompileInsert(table string, columns []string, values [][]any) (string, []any) {
	var sql strings.Builder
	var args []any

	sql.WriteString("INSERT INTO ")
	sql.WriteString(d.QuoteIdentifier(table))
	sql.WriteString(" (")

	for i, col := range columns {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(d.QuoteIdentifier(col))
	}
	sql.WriteString(") VALUES ")

	for i, row := range values {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString("(")
		for j, val := range row {
			if j > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(d.Placeholder(len(args) + 1))
			args = append(args, val)
		}
		sql.WriteString(")")
	}

	return sql.String(), args
}

func (d *SQLiteDialect) CompileUpdate(table string, values map[string]any, wheres []WhereClause) (string, []any) {
	var sql strings.Builder
	var args []any

	sql.WriteString("UPDATE ")
	sql.WriteString(d.QuoteIdentifier(table))
	sql.WriteString(" SET ")

	first := true
	for col, val := range values {
		if !first {
			sql.WriteString(", ")
		}
		sql.WriteString(d.QuoteIdentifier(col) + " = " + d.Placeholder(len(args)+1))
		args = append(args, val)
		first = false
	}

	if len(wheres) > 0 {
		sql.WriteString(" WHERE ")
		for i, w := range wheres {
			if i > 0 {
				sql.WriteString(" " + w.Boolean + " ")
			}
			sql.WriteString(d.QuoteIdentifier(w.Column))
			sql.WriteString(" " + w.Operator + " ")
			sql.WriteString(d.Placeholder(len(args) + 1))
			args = append(args, w.Value)
		}
	}

	return sql.String(), args
}

func (d *SQLiteDialect) CompileDelete(table string, wheres []WhereClause) (string, []any) {
	var sql strings.Builder
	var args []any

	sql.WriteString("DELETE FROM ")
	sql.WriteString(d.QuoteIdentifier(table))

	if len(wheres) > 0 {
		sql.WriteString(" WHERE ")
		for i, w := range wheres {
			if i > 0 {
				sql.WriteString(" " + w.Boolean + " ")
			}
			sql.WriteString(d.QuoteIdentifier(w.Column))
			sql.WriteString(" " + w.Operator + " ")
			sql.WriteString(d.Placeholder(len(args) + 1))
			args = append(args, w.Value)
		}
	}

	return sql.String(), args
}
