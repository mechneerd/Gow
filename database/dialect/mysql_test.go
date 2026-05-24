package dialect

import (
	"testing"
)

func TestMySQLDialect_QuoteIdentifier(t *testing.T) {
	d := &MySQLDialect{}
	if d.QuoteIdentifier("users") != "`users`" {
		t.Errorf("expected `users`, got %s", d.QuoteIdentifier("users"))
	}
	if d.QuoteIdentifier("*") != "*" {
		t.Error("expected * for star")
	}
}

func TestMySQLDialect_Placeholder(t *testing.T) {
	d := &MySQLDialect{}
	if d.Placeholder(1) != "?" || d.Placeholder(5) != "?" {
		t.Error("MySQL placeholders should always be ?")
	}
}

func TestMySQLDialect_CompileSelect(t *testing.T) {
	d := &MySQLDialect{}
	q := SelectQuery{
		Table:   "users",
		Columns: []string{"id", "name"},
		Wheres: []WhereClause{
			{Type: "Basic", Column: "age", Operator: ">", Value: 18, Boolean: "AND"},
		},
		Limit:  func(i int) *int { return &i }(10),
	}

	sql, args := d.CompileSelect(q)
	expected := "SELECT `id`, `name` FROM `users` WHERE `age` > ? LIMIT 10"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 1 || args[0] != 18 {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestMySQLDialect_CompileInsert(t *testing.T) {
	d := &MySQLDialect{}
	sql, args := d.CompileInsert("users", []string{"name", "email"}, [][]any{{"Alice", "a@b.com"}})
	expected := "INSERT INTO `users` (`name`, `email`) VALUES (?, ?)"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %d", len(args))
	}
}

