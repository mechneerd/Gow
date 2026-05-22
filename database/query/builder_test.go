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

func TestBuilderJoins(t *testing.T) {
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(nil, d)

	b.Table("users").
		Join("posts", "users.id", "=", "posts.user_id").
		LeftJoin("profiles", "users.id", "=", "profiles.user_id")

	sql, _ := b.ToSQL()
	expected := `SELECT * FROM "users" INNER JOIN "posts" ON "users.id" = "posts.user_id" LEFT JOIN "profiles" ON "users.id" = "profiles.user_id"`
	
	if sql != expected {
		t.Errorf("Expected %s, got %s", expected, sql)
	}
}

func TestBuilderAdvancedWheres(t *testing.T) {
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(nil, d)

	b.Table("users").
		WhereIn("id", []any{1, 2, 3}).
		WhereNull("deleted_at").
		WhereNotNull("email").
		WhereBetween("age", []any{18, 65})

	sql, args := b.ToSQL()
	expected := `SELECT * FROM "users" WHERE "id" IN (?, ?, ?) AND "deleted_at" IS NULL AND "email" IS NOT NULL AND "age" BETWEEN ? AND ?`

	if sql != expected {
		t.Errorf("Expected %s, got %s", expected, sql)
	}
	
	if len(args) != 5 {
		t.Errorf("Expected 5 args, got %d", len(args))
	}
}

func TestBuilderConditional(t *testing.T) {
	d := &dialect.SQLiteDialect{}
	b := NewBuilder(nil, d)

	b.Table("users").
		When(true, func(q *Builder) {
			q.Where("active", "=", 1)
		}).
		When(false, func(q *Builder) {
			q.Where("admin", "=", 1)
		})

	sql, args := b.ToSQL()
	expected := `SELECT * FROM "users" WHERE "active" = ?`
	
	if sql != expected {
		t.Errorf("Expected SQL: %s, got %s", expected, sql)
	}
	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}
