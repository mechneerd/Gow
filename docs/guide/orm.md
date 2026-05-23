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

### Standard Relations

```go
type Post struct {
    // ...
    UserID int
    User   orm.BelongsTo[User] `gow:"belongsTo"`
    Comments []Comment         `gow:"hasMany"`
}
```

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
