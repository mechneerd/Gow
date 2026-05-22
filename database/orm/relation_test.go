package orm

import (
	"database/sql"
	"gow/database/dialect"
	"gow/database/query"
	"testing"

	_ "modernc.org/sqlite"
)

type TestUser struct {
	ID    int    `db:"id" gow:"primaryKey"`
	Name  string `db:"name"`
	Posts []TestPost `db:"-" gow:"hasMany"`
}

func (TestUser) TableName() string { return "users" }

type TestPost struct {
	ID     int    `db:"id" gow:"primaryKey"`
	UserID int    `db:"user_id"`
	Title  string `db:"title"`
	User   *TestUser `db:"-" gow:"belongsTo"`
}

func (TestPost) TableName() string { return "posts" }

func TestEagerLoading(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Exec(`
		CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT);
		CREATE TABLE posts (id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT);
		INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob');
		INSERT INTO posts (id, user_id, title) VALUES (1, 1, 'Post 1'), (2, 1, 'Post 2'), (3, 2, 'Post 3');
	`)
	if err != nil {
		t.Fatal(err)
	}

	db := &DB{
		Conn:    conn,
		Builder: query.NewBuilder(conn, &dialect.SQLiteDialect{}),
	}

	// Test HasMany
	users, err := NewQuery[TestUser](db).With("Posts").Get()
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("Expected 2 users, got %d", len(users))
	}

	if len(users[0].Posts) != 2 {
		t.Errorf("Expected 2 posts for Alice, got %d", len(users[0].Posts))
	} else if users[0].Posts[0].Title != "Post 1" {
		t.Errorf("Expected first post to be 'Post 1', got '%s'", users[0].Posts[0].Title)
	}

	if len(users[1].Posts) != 1 {
		t.Errorf("Expected 1 post for Bob, got %d", len(users[1].Posts))
	}

	// Test BelongsTo
	posts, err := NewQuery[TestPost](db).With("User").Get()
	if err != nil {
		t.Fatalf("Failed to query posts: %v", err)
	}

	if len(posts) != 3 {
		t.Fatalf("Expected 3 posts, got %d", len(posts))
	}

	if posts[0].User == nil {
		t.Error("Expected User to be loaded on first post, got nil")
	} else if posts[0].User.Name != "Alice" {
		t.Errorf("Expected post user to be Alice, got %s", posts[0].User.Name)
	}
}
