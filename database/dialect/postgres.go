package dialect

import (
	"fmt"
	"strings"
)

// PostgresDialect implements the Dialect interface for PostgreSQL.
type PostgresDialect struct{}

func (d *PostgresDialect) QuoteIdentifier(name string) string {
	if name == "*" {
		return name
	}
	// PostgreSQL uses double quotes for identifier quoting
	return `"` + name + `"`
}

func (d *PostgresDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

func (d *PostgresDialect) compileWheres(wheres []WhereClause, args *[]any) string {
	if len(wheres) == 0 {
		return ""
	}

	var sql strings.Builder
	sql.WriteString(" WHERE ")

	for i, w := range wheres {
		if i > 0 {
			sql.WriteString(" " + w.Boolean + " ")
		}

		switch w.Type {
		case "Basic":
			sql.WriteString(d.QuoteIdentifier(w.Column))
			op, _ := ValidateOperator(w.Operator)
			sql.WriteString(" " + op + " ")
			sql.WriteString(d.Placeholder(len(*args) + 1))
			*args = append(*args, w.Value)
		case "In":
			sql.WriteString(d.QuoteIdentifier(w.Column) + " IN (")
			for j, val := range w.Values {
				if j > 0 {
					sql.WriteString(", ")
				}
				sql.WriteString(d.Placeholder(len(*args) + 1))
				*args = append(*args, val)
			}
			sql.WriteString(")")
		case "Null":
			sql.WriteString(d.QuoteIdentifier(w.Column) + " IS NULL")
		case "NotNull":
			sql.WriteString(d.QuoteIdentifier(w.Column) + " IS NOT NULL")
		case "Between":
			sql.WriteString(d.QuoteIdentifier(w.Column) + " BETWEEN ")
			sql.WriteString(d.Placeholder(len(*args) + 1) + " AND ")
			*args = append(*args, w.Values[0])
			sql.WriteString(d.Placeholder(len(*args) + 1))
			*args = append(*args, w.Values[1])
		case "Raw":
			sql.WriteString(w.RawSQL)
			*args = append(*args, w.RawArgs...)
		}
	}
	return sql.String()
}

func (d *PostgresDialect) CompileSelect(query SelectQuery) (string, []any) {
	var sql strings.Builder
	var args []any

	// SELECT
	sql.WriteString("SELECT ")
	if query.Aggregate != nil {
		sql.WriteString(fmt.Sprintf("%s(%s)", query.Aggregate.Function, d.QuoteIdentifier(query.Aggregate.Column)))
	} else if len(query.Columns) == 0 {
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

	// JOINS
	for _, j := range query.Joins {
		sql.WriteString(" " + j.Type + " JOIN ")
		sql.WriteString(d.QuoteIdentifier(j.Table) + " ON ")
		sql.WriteString(d.QuoteIdentifier(j.First) + " " + j.Operator + " " + d.QuoteIdentifier(j.Second))
	}

	// WHERE
	sql.WriteString(d.compileWheres(query.Wheres, &args))

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

	// GROUP BY
	if len(query.GroupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		for i, g := range query.GroupBys {
			if i > 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(d.QuoteIdentifier(g))
		}
	}

	// HAVING
	if len(query.Havings) > 0 {
		sql.WriteString(" HAVING ")
		havingArgs := args
		for i, h := range query.Havings {
			if i > 0 {
				sql.WriteString(" " + h.Boolean + " ")
			}
			sql.WriteString(d.QuoteIdentifier(h.Column) + " " + h.Operator + " ")
			sql.WriteString(d.Placeholder(len(havingArgs) + 1))
			havingArgs = append(havingArgs, h.Value)
		}
		args = havingArgs
	}

	// LIMIT / OFFSET
	if query.Limit != nil {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", *query.Limit))
	}
	if query.Offset != nil {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", *query.Offset))
	}

	// Locking (FOR UPDATE / FOR SHARE)
	if query.Lock != "" {
		sql.WriteString(" " + query.Lock)
	}

	return sql.String(), args
}

func (d *PostgresDialect) CompileInsert(table string, columns []string, values [][]any) (string, []any) {
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

func (d *PostgresDialect) CompileUpdate(table string, values map[string]any, wheres []WhereClause) (string, []any) {
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

	sql.WriteString(d.compileWheres(wheres, &args))

	return sql.String(), args
}

func (d *PostgresDialect) CompileDelete(table string, wheres []WhereClause) (string, []any) {
	var sql strings.Builder
	var args []any

	sql.WriteString("DELETE FROM ")
	sql.WriteString(d.QuoteIdentifier(table))

	sql.WriteString(d.compileWheres(wheres, &args))

	return sql.String(), args
}

