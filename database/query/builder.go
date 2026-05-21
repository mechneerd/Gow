package query

import (
	"context"
	"database/sql"
	"gow/database/dialect"
)

// Builder provides a fluent API for building SQL queries.
type Builder struct {
	conn      *sql.DB
	dialect   dialect.Dialect
	ctx       context.Context
	
	query     dialect.SelectQuery
}

// NewBuilder creates a new query builder.
func NewBuilder(conn *sql.DB, d dialect.Dialect) *Builder {
	return &Builder{
		conn:    conn,
		dialect: d,
		ctx:     context.Background(),
		query:   dialect.SelectQuery{},
	}
}

// Table sets the target table.
func (b *Builder) Table(table string) *Builder {
	b.query.Table = table
	return b
}

// Select specifies the columns to select.
func (b *Builder) Select(columns ...string) *Builder {
	b.query.Columns = append(b.query.Columns, columns...)
	return b
}

// Where adds a basic where clause.
func (b *Builder) Where(column, operator string, value any) *Builder {
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:     "Basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  "AND",
	})
	return b
}

// OrWhere adds an OR where clause.
func (b *Builder) OrWhere(column, operator string, value any) *Builder {
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:     "Basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  "OR",
	})
	return b
}

// OrderBy adds an order by clause.
func (b *Builder) OrderBy(column, direction string) *Builder {
	b.query.OrderBys = append(b.query.OrderBys, dialect.OrderByClause{
		Column:    column,
		Direction: direction,
	})
	return b
}

// Limit sets the limit.
func (b *Builder) Limit(limit int) *Builder {
	b.query.Limit = &limit
	return b
}

// Offset sets the offset.
func (b *Builder) Offset(offset int) *Builder {
	b.query.Offset = &offset
	return b
}

// ToSQL compiles the query to SQL.
func (b *Builder) ToSQL() (string, []any) {
	return b.dialect.CompileSelect(b.query)
}

// Get executes the query and returns the raw rows.
// In Goquent (ORM), we will wrap this to hydrate models.
func (b *Builder) Get() (*sql.Rows, error) {
	sqlQuery, args := b.ToSQL()
	return b.conn.QueryContext(b.ctx, sqlQuery, args...)
}

// Insert executes an insert statement.
func (b *Builder) Insert(values map[string]any) (sql.Result, error) {
	columns := make([]string, 0, len(values))
	row := make([]any, 0, len(values))
	
	for k, v := range values {
		columns = append(columns, k)
		row = append(row, v)
	}
	
	sqlQuery, args := b.dialect.CompileInsert(b.query.Table, columns, [][]any{row})
	return b.conn.ExecContext(b.ctx, sqlQuery, args...)
}

// Update executes an update statement.
func (b *Builder) Update(values map[string]any) (sql.Result, error) {
	sqlQuery, args := b.dialect.CompileUpdate(b.query.Table, values, b.query.Wheres)
	return b.conn.ExecContext(b.ctx, sqlQuery, args...)
}

// Delete executes a delete statement.
func (b *Builder) Delete() (sql.Result, error) {
	sqlQuery, args := b.dialect.CompileDelete(b.query.Table, b.query.Wheres)
	return b.conn.ExecContext(b.ctx, sqlQuery, args...)
}
