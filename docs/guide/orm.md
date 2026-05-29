> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Goquent ORM

> **Status**: ✅ Implemented (Very Mature)

Goquent is a powerful, Laravel Eloquent-inspired ORM built on top of a fluent query builder. It supports relationships, eager loading, soft deletes, scopes, observers, casting, accessors, mutators, and advanced querying features.

## Defining Models

```go
type User struct {
    ID        int    `db:"id"`
    Name      string `db:"name"`
    Email     string `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

## Mass Assignment Protection

```go
func (u *User) Fillable() []string {
    return []string{"name", "email"}
}

func (u *User) Guarded() []string {
    return []string{"id"}
}
```

## Relationships

Goquent supports all common relationship types including **Polymorphic** and **Through** relations.

### Relationship Tags Reference

All relationships are defined using the `gow` struct tag. The following tags are supported:

| Tag | Description | Example |
|---|---|---|
| `gow:"hasMany"` | One-to-many relationship | `Comments []Comment \`gow:"hasMany"\`` |
| `gow:"belongsTo"` | Inverse one-to-many (foreign key on this model) | `User orm.BelongsTo[User] \`gow:"belongsTo"\`` |
| `gow:"belongsToMany"` | Many-to-many via pivot table | `Roles orm.BelongsToMany[Role] \`gow:"belongsToMany"\`` |
| `gow:"hasOne"` | One-to-one relationship | `Profile orm.HasOne[Profile] \`gow:"hasOne"\`` |
| `gow:"morphMany"` | Polymorphic one-to-many | `Comments orm.MorphMany[Comment] \`gow:"morphMany"\`` |
| `gow:"morphOne"` | Polymorphic one-to-one | `Image orm.MorphOne[Image] \`gow:"morphOne"\`` |
| `gow:"morphTo"` | Inverse polymorphic | `Commentable orm.MorphTo \`gow:"morphTo"\`` |
| `gow:"hasOneThrough"` | One-to-one through intermediate | `Country orm.HasOneThrough[Country] \`gow:"hasOneThrough"\`` |
| `gow:"hasManyThrough"` | One-to-many through intermediate | `Country orm.HasManyThrough[Country] \`gow:"hasManyThrough"\`` |

### Tag Options

| Option | Description | Example |
|---|---|---|
| `foreignKey=<column>` | Custom foreign key column | `gow:"hasMany,foreignKey=post_id"` |
| `relatedKey=<column>` | Custom related key column | `gow:"hasMany,relatedKey=user_id"` |
| `through=<table>` | Intermediate table for Through relations | `gow:"hasManyThrough,through=users"` |
| `morphType=<column>` | Polymorphic type column | `gow:"morphMany,morphType=commentable_type"` |
| `morphId=<column>` | Polymorphic ID column | `gow:"morphMany,morphId=commentable_id"` |
| `type=<Name>` | Concrete type name for morph relations | `gow:"morphMany,type=Post"` |
| `pivot=<table>` | Pivot table name (default: `table1_table2`) | `gow:"belongsToMany,pivot=role_user"` |

### Standard Relations

```go
type Post struct {
    ID        int    `db:"id"`
    UserID    int    `db:"user_id"`
    User      orm.BelongsTo[User] `gow:"belongsTo"`
    Comments  []Comment           `gow:"hasMany"`
    Tags      orm.BelongsToMany[Tag] `gow:"belongsToMany,pivot=post_tag"`
}
```

The foreign key is automatically derived from the relationship name (e.g., `User` → `user_id`). Override with `foreignKey=` if needed.

### BelongsToMany (Many-to-Many)

```go
type Post struct {
    Tags orm.BelongsToMany[Tag] `gow:"belongsToMany,pivot=post_tag,foreignKey=post_id,relatedKey=tag_id"`
}
```

Default pivot table name: `<pluralized_model>_<pluralized_relation>` (e.g., `post_tags`). Override with `pivot=`.

### Polymorphic Relations

```go
type Comment struct {
    CommentableType string `db:"commentable_type"`
    CommentableID   int    `db:"commentable_id"`

    // From the owning side
    Post orm.MorphMany[Post] `gow:"morphMany,morphType=commentable_type,morphId=commentable_id,type=Post"`
}

// On the inverse side
type Post struct {
    Comments orm.MorphMany[Comment] `gow:"morphMany,type=Post"`
}
```

Register morph types at bootstrap:

```go
orm.RegisterMorph("post", models.Post{})
```

### HasOneThrough / HasManyThrough

```go
type Country struct {
    Posts orm.HasManyThrough[Post] `gow:"hasManyThrough,through=users,foreignKey=country_id,relatedKey=user_id"`
}
```

## Attribute Casting

Models can implement `Castable`:

```go
func (u *User) Casts() map[string]string {
    return map[string]string{
        "settings": "json",
        "birthday": "datetime",
        "is_active": "bool",
    }
}
```

Built-in casts: `datetime`, `json`, `bool`, `int`, `float`, `string`.

## Accessors & Mutators

```go
func (u *User) GetFullNameAttribute() string {
    return u.FirstName + " " + u.LastName
}

func (u *User) SetPasswordAttribute(value string) {
    hashed, _ := hashing.Make(value)
    u.Password = hashed
}
```

Usage:

```go
name := orm.GetModelAttribute(user, "full_name")
orm.SetModelAttribute(user, "password", "secret123")
```

## Local & Global Scopes

Global scopes:

```go
orm.AddGlobalScope("User", &ActiveUsersScope{})
```

Local scopes (auto-discovered):

```go
func (u *User) ScopeActive(q *query.Builder) *query.Builder {
    return q.Where("active", "=", true)
}

// Usage
users := orm.NewQuery[User](db).Where("age", ">", 18).Active().Get()
```

## Pessimistic Locking

```go
user, _ := orm.NewQuery[User](db).
    LockForUpdate().
    Find(1)
```

Also supports `.SharedLock()`.

## Upsert

```go
builder.Upsert(map[string]any{
    "email": "john@example.com",
    "name":  "John",
}, []string{"name"})
```

## Chunking Large Datasets

```go
orm.NewQuery[User](db).Chunk(100, func(users []User) error {
    // Process batch
    return nil
})
```

## Relation Touching

```go
func (p *Post) Touches() []string {
    return []string{"author"}
}
```

When a Post is updated, the related Author's `updated_at` is automatically touched.

## Observers & Lifecycle Events

```go
orm.Observe("Post", &PostObserver{})
```

Or use interface hooks (`BeforeCreate`, `AfterSave`, etc.).

## Eager Loading

```go
posts := orm.NewQuery[Post](db).
    With("Comments", "Author.Country").
    Get()
```

## Soft Deletes, Pagination, Transactions

All fully supported (see older sections or `Current_Capabilities.md` for details).

## Best Practices

- Use `MassAssignable` on all models.
- Prefer casting + accessors over manual transformation.
- Use `Chunk` instead of loading thousands of records at once.
- Register polymorphic types early in bootstrap.

---

## Deep-Dive Examples

### Polymorphic Comments

```go
type Comment struct {
    ID               int    `db:"id"`
    Body             string `db:"body"`
    CommentableType  string `db:"commentable_type"`
    CommentableID    int    `db:"commentable_id"`
}

type Post struct {
    ID       int
    Title    string
    Comments orm.MorphMany[Comment] `gow:"morphMany,morphType=commentable_type,morphId=commentable_id,type=Post"`
}

type Video struct {
    ID       int
    Title    string
    Comments orm.MorphMany[Comment] `gow:"morphMany,morphType=commentable_type,morphId=commentable_id,type=Video"`
}
```

Register once:

```go
orm.RegisterMorph("post", models.Post{})
orm.RegisterMorph("video", models.Video{})
```

### Attribute Casting + Accessors Together

```go
type User struct {
    Settings map[string]any `db:"settings"`
    Birthday time.Time      `db:"birthday"`
}

func (u *User) Casts() map[string]string {
    return map[string]string{
        "settings": "json",
        "birthday": "datetime",
    }
}

func (u *User) GetAgeAttribute() int {
    return time.Now().Year() - u.Birthday.Year()
}
```

### HasManyThrough Example

```go
type Country struct {
    ID    int
    Name  string
    Posts orm.HasManyThrough[Post] `gow:"hasManyThrough,through=users,foreignKey=country_id,relatedKey=user_id"`
}
```

This lets you do `country.Posts` directly.
