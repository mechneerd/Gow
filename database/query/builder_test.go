package query

import (
	"gow/database/dialect"
	"testing"
)

func TestBuilderSelect(t *testing.T) {
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(nil, d)

	b.Table("users").
		Select("id", "name", "email").
		Where("status", "=", "active").
		OrWhere("age", ">", 18).
		OrderBy("created_at", "DESC").
		Limit(10).
		Offset(20)

	sql, args := b.ToSQL()

	expectedSQL := `SELECT "id", "name", "email" FROM "users" WHERE "status" = ? OR "age" > ? ORDER BY "created_at" DESC LIMIT 10 OFFSET 20`
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 2 || args[0] != "active" || args[1] != 18 {
		t.Errorf("Unexpected args: %v", args)
	}
}

func TestBuilderInsert(t *testing.T) {
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(nil, d)
	b.Table("users")

	sql, args := d.CompileInsert(b.query.Table, []string{"name", "email"}, [][]any{{"John Doe", "john@example.com"}})

	expectedSQL := `INSERT INTO "users" ("name", "email") VALUES (?, ?)`
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
	if len(args) != 2 {
		t.Errorf("Unexpected args: %v", args)
	}
}
