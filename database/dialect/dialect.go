package dialect

import "fmt"

// ValidOperators is the whitelist of allowed SQL comparison operators.
var ValidOperators = map[string]bool{
	"=": true, "==": true, "!=": true, "<>": true,
	">": true, "<": true, ">=": true, "<=": true,
	"LIKE": true, "NOT LIKE": true, "ILIKE": true,
	"IN": true, "NOT IN": true,
	"IS": true, "IS NOT": true,
	"BETWEEN": true, "NOT BETWEEN": true,
}

// ValidateOperator checks if an operator is in the whitelist.
// Returns the operator if valid, or an error message if not.
func ValidateOperator(op string) (string, error) {
	if ValidOperators[op] {
		return op, nil
	}
	return "", fmt.Errorf("invalid SQL operator: %s", op)
}

// QuoteIdentifier quotes a table or column name (e.g. `col` for MySQL, "col" for Postgres)
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

	// CompileUpsert compiles an INSERT ... ON CONFLICT (or equivalent) statement.
	CompileUpsert(table string, columns []string, values [][]any, conflictCols []string, updateCols []string) (string, []any)

	// AutoIncrementSQL returns the SQL fragment for an auto-increment primary key column.
	AutoIncrementSQL() string
}

// SelectQuery represents the components of a SELECT query for compilation.
type SelectQuery struct {
	Table      string
	Columns    []string
	RawColumns []string // Raw SQL expressions (e.g. "COUNT(*) as total") — not quoted
	Wheres     []WhereClause
	OrderBys   []OrderByClause
	Joins      []JoinClause
	Aggregate  *AggregateClause
	Limit      *int
	Offset     *int
	GroupBys   []string
	Havings    []WhereClause
	Lock       string // e.g. "FOR UPDATE", "FOR SHARE" for pessimistic locking
}

// JoinClause represents a table JOIN condition.
type JoinClause struct {
	Type     string // INNER, LEFT, RIGHT, CROSS
	Table    string
	First    string
	Operator string
	Second   string
}

// AggregateClause represents an aggregate function to apply instead of selecting columns.
type AggregateClause struct {
	Function string // COUNT, MAX, MIN, SUM, AVG
	Column   string // usually * or a column name
}

// WhereClause represents a WHERE condition.
type WhereClause struct {
	Type     string // Basic, In, Null, NotNull, Between, Subquery, Raw
	Column   string
	Operator string
	Value    any
	Values   []any // For IN, BETWEEN
	Boolean  string // AND / OR
	RawSQL   string // For Raw queries
	RawArgs  []any
}

// OrderByClause represents an ORDER BY condition.
type OrderByClause struct {
	Column    string
	Direction string
}

