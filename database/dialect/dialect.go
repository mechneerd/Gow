package dialect

// Dialect represents a SQL dialect for a specific database driver.
type Dialect interface {
	// QuoteIdentifier quotes a table or column name (e.g. `col` for MySQL, "col" for Postgres)
	QuoteIdentifier(name string) string
	
	// Placeholder returns the bind parameter placeholder for a given index (1-based).
	// e.g. $1 for Postgres, ? for MySQL/SQLite
	Placeholder(index int) string
	
	// CompileSelect compiles a SELECT statement
	CompileSelect(query SelectQuery) (string, []any)
	
	// CompileInsert compiles an INSERT statement
	CompileInsert(table string, columns []string, values [][]any) (string, []any)
	
	// CompileUpdate compiles an UPDATE statement
	CompileUpdate(table string, values map[string]any, wheres []WhereClause) (string, []any)
	
	// CompileDelete compiles a DELETE statement
	CompileDelete(table string, wheres []WhereClause) (string, []any)
}

// SelectQuery represents the components of a SELECT query for compilation.
type SelectQuery struct {
	Table    string
	Columns  []string
	Wheres   []WhereClause
	OrderBys []OrderByClause
	Limit    *int
	Offset   *int
}

// WhereClause represents a WHERE condition.
type WhereClause struct {
	Type     string // Basic, In, Null, etc.
	Column   string
	Operator string
	Value    any
	Boolean  string // AND / OR
}

// OrderByClause represents an ORDER BY condition.
type OrderByClause struct {
	Column    string
	Direction string
}
