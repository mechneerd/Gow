# Authorization

In addition to providing built-in authentication services, GoW also provides a simple way to authorize user actions against a given resource.

## Gates

Gates are closures that determine if a user is authorized to perform a given action. You typically define gates in a service provider.

```go
gate := authorization.NewGate()

gate.Define("update-post", func(user any, args ...any) bool {
    post := args[0].(*models.Post)
    u := user.(*models.User)
    
    return u.ID == post.AuthorID
})
```

You can then authorize actions using the `Allows` or `Denies` methods:

```go
if gate.Allows(currentUser, "update-post", post) {
    // The user can update the post...
}
```

## Policies

Policies are structs that organize authorization logic around a particular model or resource.

```go
type PostPolicy struct {}

func (p *PostPolicy) Update(user *models.User, post *models.Post) bool {
    return user.ID == post.AuthorID
}
```

You can register policies against their respective models:

```go
gate.Policy("Post", &PostPolicy{})
```

## Blade Directives

When writing Goblade templates, you often want to conditionally render elements based on the user's permissions. GoW provides the `@can` and `@cannot` directives for this exact purpose.

```html
@can('update', $post)
    <button>Edit Post</button>
@else
    <span>You do not have permission to edit this post.</span>
@endcan
```
