package view

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEngineExtends(t *testing.T) {
	// Create a temporary directory for views and cache
	tempDir, err := os.MkdirTemp("", "gow_views")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	viewsDir := filepath.Join(tempDir, "views")
	cacheDir := filepath.Join(tempDir, "cache")
	os.MkdirAll(viewsDir, 0755)
	os.MkdirAll(cacheDir, 0755)

	// Create master layout
	masterContent := `<html>
<head><title>@yield('title')</title></head>
<body>
@yield('content')
</body>
</html>`
	os.WriteFile(filepath.Join(viewsDir, "master.goblade"), []byte(masterContent), 0644)

	// Create parent layout that extends master
	parentContent := `@extends('master')
@section('title')
My App
@endsection
@section('content')
<div class="container">
@yield('main')
</div>
@endsection`
	os.WriteFile(filepath.Join(viewsDir, "parent.goblade"), []byte(parentContent), 0644)

	// Create child view that extends parent
	childContent := `@extends('parent')
@section('main')
<h1>Hello, {{ .Name }}</h1>
@endsection`
	os.WriteFile(filepath.Join(viewsDir, "child.goblade"), []byte(childContent), 0644)

	engine := NewEngine([]string{viewsDir}, cacheDir)
	
	out, err := engine.Make("child", map[string]any{"Name": "World"})
	if err != nil {
		t.Fatalf("Make failed: %v", err)
	}

	// Output should have master, parent, and child content
	expectedStrings := []string{
		"<html>",
		"<title>", "My App", "</title>",
		"<div class=\"container\">",
		"<h1>Hello, World</h1>",
		"</div>",
		"</html>",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(out, s) {
			t.Errorf("Expected output to contain '%s'\nOutput:\n%s", s, out)
		}
	}
}

func TestEngineCycleDetection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gow_views_cycle")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	viewsDir := filepath.Join(tempDir, "views")
	cacheDir := filepath.Join(tempDir, "cache")
	os.MkdirAll(viewsDir, 0755)
	os.MkdirAll(cacheDir, 0755)

	// Create a cycle: a -> b -> a
	os.WriteFile(filepath.Join(viewsDir, "a.goblade"), []byte(`@extends('b')`), 0644)
	os.WriteFile(filepath.Join(viewsDir, "b.goblade"), []byte(`@extends('a')`), 0644)

	engine := NewEngine([]string{viewsDir}, cacheDir)
	
	_, err = engine.Make("a", nil)
	if err == nil {
		t.Fatal("Expected error due to cycle, got nil")
	}

	if !strings.Contains(err.Error(), "cycle detected in layout chain") {
		t.Errorf("Expected cycle detection error, got: %v", err)
	}
}

