package dialect

import (
	"testing"
)

func TestPostgresDialect_QuoteIdentifier(t *testing.T) {
	d := &PostgresDialect{}
	if d.QuoteIdentifier("users") != `"users"` {
		t.Errorf("expected \"users\", got %s", d.QuoteIdentifier("users"))
	}
	if d.QuoteIdentifier("*") != "*" {
		t.Error("expected * for star")
	}
}

func TestPostgresDialect_Placeholder(t *testing.T) {
	d := &PostgresDialect{}
	if d.Placeholder(1) != "$1" || d.Placeholder(3) != "$3" {
		t.Error("Postgres placeholders should be $n")
	}
}

func TestPostgresDialect_CompileSelect(t *testing.T) {
	d := &PostgresDialect{}
	q := SelectQuery{
		Table:   "users",
		Columns: []string{"id", "name"},
		Wheres: []WhereClause{
			{Type: "Basic", Column: "age", Operator: ">", Value: 18, Boolean: "AND"},
		},
		Limit:  func(i int) *int { return &i }(10),
	}

	sql, args := d.CompileSelect(q)
	expected := `SELECT "id", "name" FROM "users" WHERE "age" > $1 LIMIT 10`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 1 || args[0] != 18 {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestPostgresDialect_CompileInsert(t *testing.T) {
	d := &PostgresDialect{}
	sql, args := d.CompileInsert("users", []string{"name", "email"}, [][]any{{"Alice", "a@b.com"}})
	expected := `INSERT INTO "users" ("name", "email") VALUES ($1, $2)`
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}
