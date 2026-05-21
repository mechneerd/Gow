# Goquent ORM

The Goquent ORM provides a beautiful, simple Active Record implementation for working with your database. Each database table has a corresponding "Model" struct that is used to interact with that table.

## Defining Models

Models are standard Go structs heavily enriched by struct tags. By default, Goquent assumes your table name is the plural, lowercase version of your struct name.

```go
package models

type Post struct {
    ID        int    `db:"id"`
    Title     string `db:"title"`
    Body      string `db:"body"`
    AuthorID  int    `db:"author_id"`
    CreatedAt string `db:"created_at"`
    UpdatedAt string `db:"updated_at"`
}
```

## Fluent Query Builder

Goquent ships with a fluent query builder to construct dynamic SQL safely.

```go
import "gow/database/orm"

// Retrieve all active posts created by author 1
posts := orm.Table("posts").
    Where("author_id", "=", 1).
    Where("status", "=", "active").
    Get()
```

## Relationships

Goquent makes managing and querying relationships trivial using Eager Loading (`With`).

```go
// Retrieve posts and eagerly load their associated authors to prevent N+1 queries
posts := orm.Table("posts").
    With("Author").
    Get()
```
