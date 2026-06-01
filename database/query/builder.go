package query

import (
	"context"
	"database/sql"
	"sync"
	"time"
	
	"github.com/mechneerd/gow/database/dialect"
)

// QueryEvent represents an executed database query.
type QueryEvent struct {
	SQL      string
	Bindings []any
	Duration time.Duration
	// Caller stack can be added here if needed
}

// QueryListener is a callback function that listens for query events.
type QueryListener func(QueryEvent)

var (
	queryListeners []QueryListener
	listenerMu     sync.RWMutex
)

// Listen registers a new query listener.
func Listen(listener QueryListener) {
	listenerMu.Lock()
	defer listenerMu.Unlock()
	queryListeners = append(queryListeners, listener)
}

func dispatchQueryEvent(sql string, bindings []any, duration time.Duration) {
	listenerMu.RLock()
	listeners := make([]QueryListener, len(queryListeners))
	copy(listeners, queryListeners)
	listenerMu.RUnlock()

	if len(listeners) == 0 {
		return
	}
	event := QueryEvent{
		SQL:      sql,
		Bindings: bindings,
		Duration: duration,
	}
	for _, listener := range listeners {
		listener(event)
	}
}

// QueryExecer abstracts *sql.DB and *sql.Tx
type QueryExecer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Builder provides a fluent API for building SQL queries.
type Builder struct {
	conn      QueryExecer
	dialect   dialect.Dialect
	ctx       context.Context
	
	query     dialect.SelectQuery
}

// NewBuilder creates a new query builder.
func NewBuilder(conn QueryExecer, d dialect.Dialect) *Builder {
	return &Builder{
		conn:    conn,
		dialect: d,
		ctx:     context.Background(),
		query:   dialect.SelectQuery{},
	}
}

// Clone creates a deep copy of the builder.
func (b *Builder) Clone() *Builder {
	q := b.query
	q.Wheres = append([]dialect.WhereClause{}, b.query.Wheres...)
	q.OrderBys = append([]dialect.OrderByClause{}, b.query.OrderBys...)
	q.Joins = append([]dialect.JoinClause{}, b.query.Joins...)
	q.Columns = append([]string{}, b.query.Columns...)
	q.GroupBys = append([]string{}, b.query.GroupBys...)
	q.Havings = append([]dialect.WhereClause{}, b.query.Havings...)
	if b.query.Aggregate != nil {
		agg := *b.query.Aggregate
		q.Aggregate = &agg
	}
	if b.query.Limit != nil {
		limit := *b.query.Limit
		q.Limit = &limit
	}
	if b.query.Offset != nil {
		offset := *b.query.Offset
		q.Offset = &offset
	}
	return &Builder{
		conn:    b.conn,
		dialect: b.dialect,
		ctx:     b.ctx,
		query:   q,
	}
}

// WithConn returns a new builder with the given connection (useful for transactions).
func (b *Builder) WithConn(conn QueryExecer) *Builder {
	clone := b.Clone()
	clone.conn = conn
	return clone
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

// Distinct adds a DISTINCT clause to the query.
func (b *Builder) Distinct() *Builder {
	b.query.Distinct = true
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

// GroupBy adds GROUP BY columns.
func (b *Builder) GroupBy(columns ...string) *Builder {
	b.query.GroupBys = append(b.query.GroupBys, columns...)
	return b
}

// Having adds a HAVING clause (after GROUP BY).
func (b *Builder) Having(column, operator string, value any) *Builder {
	b.query.Havings = append(b.query.Havings, dialect.WhereClause{
		Type:     "Basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  "AND",
	})
	return b
}

// OrHaving adds an OR HAVING clause.
func (b *Builder) OrHaving(column, operator string, value any) *Builder {
	b.query.Havings = append(b.query.Havings, dialect.WhereClause{
		Type:     "Basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  "OR",
	})
	return b
}

// SelectRaw allows raw SELECT expressions (e.g. "COUNT(*) as total").
func (b *Builder) SelectRaw(sql string, args ...any) *Builder {
	b.query.RawColumns = append(b.query.RawColumns, sql)
	// Note: Raw select args are passed through but may require dialect-specific handling.
	return b
}

// WhereRaw adds a raw WHERE condition.
func (b *Builder) WhereRaw(sql string, args ...any) *Builder {
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:    "Raw",
		RawSQL:  sql,
		RawArgs: args,
		Boolean: "AND",
	})
	return b
}

// OrWhereRaw adds a raw OR WHERE condition.
func (b *Builder) OrWhereRaw(sql string, args ...any) *Builder {
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:    "Raw",
		RawSQL:  sql,
		RawArgs: args,
		Boolean: "OR",
	})
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
	start := time.Now()
	rows, err := b.conn.QueryContext(b.ctx, sqlQuery, args...)
	dispatchQueryEvent(sqlQuery, args, time.Since(start))
	return rows, err
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
	start := time.Now()
	res, err := b.conn.ExecContext(b.ctx, sqlQuery, args...)
	dispatchQueryEvent(sqlQuery, args, time.Since(start))
	return res, err
}

// Upsert inserts or updates on conflict using dialect-specific syntax.
func (b *Builder) Upsert(values map[string]any, updateColumns []string, conflictColumns []string) (sql.Result, error) {
	columns := make([]string, 0, len(values))
	row := make([]any, 0, len(values))
	for k, v := range values {
		columns = append(columns, k)
		row = append(row, v)
	}

	sqlQ, args := b.dialect.CompileUpsert(b.query.Table, columns, [][]any{row}, conflictColumns, updateColumns)

	start := time.Now()
	res, err := b.conn.ExecContext(b.ctx, sqlQ, args...)
	dispatchQueryEvent(sqlQ, args, time.Since(start))
	return res, err
}

// LockForUpdate adds FOR UPDATE pessimistic lock (useful in transactions).
func (b *Builder) LockForUpdate() *Builder {
	b.query.Lock = "FOR UPDATE"
	return b
}

// SharedLock adds FOR SHARE / FOR READ ONLY lock.
func (b *Builder) SharedLock() *Builder {
	b.query.Lock = "FOR SHARE"
	return b
}

// Update executes an update statement.
func (b *Builder) Update(values map[string]any) (sql.Result, error) {
	sqlQuery, args := b.dialect.CompileUpdate(b.query.Table, values, b.query.Wheres)
	start := time.Now()
	res, err := b.conn.ExecContext(b.ctx, sqlQuery, args...)
	dispatchQueryEvent(sqlQuery, args, time.Since(start))
	return res, err
}

// Delete executes a delete statement.
func (b *Builder) Delete() (sql.Result, error) {
	sqlQuery, args := b.dialect.CompileDelete(b.query.Table, b.query.Wheres)
	start := time.Now()
	res, err := b.conn.ExecContext(b.ctx, sqlQuery, args...)
	dispatchQueryEvent(sqlQuery, args, time.Since(start))
	return res, err
}

// --- JOIN CLAUSES ---

func (b *Builder) join(joinType, table, first, operator, second string) *Builder {
	b.query.Joins = append(b.query.Joins, dialect.JoinClause{
		Type:     joinType,
		Table:    table,
		First:    first,
		Operator: operator,
		Second:   second,
	})
	return b
}

func (b *Builder) Join(table, first, operator, second string) *Builder {
	return b.join("INNER", table, first, operator, second)
}

func (b *Builder) LeftJoin(table, first, operator, second string) *Builder {
	return b.join("LEFT", table, first, operator, second)
}

func (b *Builder) RightJoin(table, first, operator, second string) *Builder {
	return b.join("RIGHT", table, first, operator, second)
}

func (b *Builder) CrossJoin(table string) *Builder {
	return b.join("CROSS", table, "", "", "")
}

// --- ADVANCED WHERE CLAUSES ---

func (b *Builder) WhereIn(column string, values []any) *Builder {
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:    "In",
		Column:  column,
		Values:  values,
		Boolean: "AND",
	})
	return b
}

func (b *Builder) WhereNull(column string) *Builder {
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:    "Null",
		Column:  column,
		Boolean: "AND",
	})
	return b
}

func (b *Builder) WhereNotNull(column string) *Builder {
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:    "NotNull",
		Column:  column,
		Boolean: "AND",
	})
	return b
}

func (b *Builder) WhereBetween(column string, values []any) *Builder {
	if len(values) != 2 {
		// in a real framework we'd return an error or panic gracefully
		return b
	}
	b.query.Wheres = append(b.query.Wheres, dialect.WhereClause{
		Type:    "Between",
		Column:  column,
		Values:  values,
		Boolean: "AND",
	})
	return b
}

// --- AGGREGATES ---

func (b *Builder) aggregate(function, column string) (int, error) {
	// Clone builder to avoid mutating the original
	b2 := b.Clone()
	b2.query.Aggregate = &dialect.AggregateClause{
		Function: function,
		Column:   column,
	}
	
	sqlQuery, args := b2.dialect.CompileSelect(b2.query)
	var result int
	start := time.Now()
	err := b2.conn.QueryRowContext(b2.ctx, sqlQuery, args...).Scan(&result)
	dispatchQueryEvent(sqlQuery, args, time.Since(start))
	return result, err
}

func (b *Builder) Count(column string) (int, error) {
	if column == "" {
		column = "*"
	}
	return b.aggregate("COUNT", column)
}

func (b *Builder) Max(column string) (int, error) {
	return b.aggregate("MAX", column)
}

func (b *Builder) Min(column string) (int, error) {
	return b.aggregate("MIN", column)
}

func (b *Builder) Avg(column string) (int, error) {
	return b.aggregate("AVG", column)
}

func (b *Builder) Sum(column string) (int, error) {
	return b.aggregate("SUM", column)
}

// --- CONDITIONAL ---

// When executes the given callback if the condition is true.
func (b *Builder) When(condition bool, callback func(*Builder)) *Builder {
	if condition {
		callback(b)
	}
	return b
}

