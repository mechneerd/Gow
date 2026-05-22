# Views & Goblade

> **Status**: ?? In Progress (Core working, tests/edge cases pending)


GoW uses a powerful server-side rendering engine called **Goblade**, which allows you to seamlessly mix HTML with Go control structures using an expressive, Blade-like syntax.

## Basic Usage

You can return a view from your route or controller using the `view` package.

```go
func Index(w http.ResponseWriter, r *http.Request) {
    view.Make(w, "welcome", map[string]any{
        "name": "Taylor",
    })
}
```

## Control Structures

Goblade provides convenient shortcuts for common Go template control structures.

### If Statements

```html
@if(name == "Taylor")
    I have a name!
@elseif(name == "Victoria")
    I have a different name!
@else
    I don't have a name!
@endif
```

### Loops

You can iterate over arrays and slices easily. Goblade automatically injects a `$loop` variable that gives you valuable information about the current iteration.

```html
@foreach(users)
    @if($loop.first)
        This is the first iteration.
    @endif

    <p>This is user {{ .id }}</p>

    @if($loop.last)
        This is the last iteration.
    @endif
@endforeach
```

## Advanced Directives

### Conditional Classes (`@class`)

The `@class` directive conditionally compiles a CSS class string based on a map of booleans.

```html
<span @class(['font-bold': true, 'bg-red': hasError])></span>
```

### Conditional Checked (`@checked`)

```html
<input type="checkbox" name="active" value="active" @checked(user.isActive)>
```

### One-Time Rendering (`@once`)

The `@once` directive allows you to define a portion of the template that will only be evaluated and rendered once per request cycle. This is incredibly useful for pushing a specific piece of JavaScript to a stack from within a component.

```html
@once
    <script>
        // Custom JavaScript that will only be included once
    </script>
@endonce
```

## Layouts & Components

Goblade supports extending base layouts and defining reusable components.

### Defining a Layout

```html
<!-- resources/views/layouts/app.gohtml -->
<html>
    <head>
        <title>App Name - @yield('title')</title>
    </head>
    <body>
        <div class="container">
            @yield('content')
        </div>
    </body>
</html>
```

### Extending a Layout

```html
<!-- resources/views/child.gohtml -->
@extends('layouts.app')

@section('title', 'Page Title')

@section('content')
    <p>This is my body content.</p>
@endsection
```

### Components

You can define reusable components in the `components/` directory and use them seamlessly.

```html
<!-- resources/views/components/alert.gohtml -->
<div class="alert alert-danger">
    {{ .slot }}
</div>
```

```html
<!-- Using the component -->
<x-alert>
    <strong>Whoops!</strong> Something went wrong!
</x-alert>
```
