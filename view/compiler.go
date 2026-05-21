package view

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Compiler handles transpiling Goblade templates into standard Go html/template files.
type Compiler struct {
	CachePath string
}

// NewCompiler creates a new Goblade compiler.
func NewCompiler(cachePath string) *Compiler {
	return &Compiler{
		CachePath: cachePath,
	}
}

// Compile takes a raw Goblade string and returns the compiled html/template string.
// For Phase 3, we implement a lightweight parser/replacer as requested.
func (c *Compiler) CompileString(raw string) string {
	compiled := raw

	// 1. Echo statements: {{ $var }} -> {{ $var }} (Standard Go handles it)
	// But let's support Blade unescaped {!! $var !!} by translating to {{ .var | raw }} if we register a func.
	// For now, we'll map standard blade @directives.

	// 2. Control Structures
	compiled = strings.ReplaceAll(compiled, "@else", "{{ else }}")
	compiled = strings.ReplaceAll(compiled, "@endif", "{{ end }}")
	compiled = strings.ReplaceAll(compiled, "@endforeach", "{{ end }}")

	// @if(condition) -> {{ if condition }}
	ifRe := regexp.MustCompile(`@if\s*\((.*?)\)`)
	compiled = ifRe.ReplaceAllString(compiled, "{{ if $1 }}")

	// @elseif(condition) -> {{ else if condition }}
	elseIfRe := regexp.MustCompile(`@elseif\s*\((.*?)\)`)
	compiled = elseIfRe.ReplaceAllString(compiled, "{{ else if $1 }}")

	// @foreach(range) -> {{ range range }}
	foreachRe := regexp.MustCompile(`@foreach\s*\((.*?)\)`)
	compiled = foreachRe.ReplaceAllString(compiled, "{{ range $1 }}")

	// 3. Layouts & Sections
	// @extends('layouts.app') -> ignored directly in compilation string, handled by engine aggregating files
	// However, we can map @yield('content') -> {{ block "content" . }}{{ end }}
	yieldRe := regexp.MustCompile(`@yield\s*\(['"](.*?)['"]\)`)
	compiled = yieldRe.ReplaceAllString(compiled, `{{ block "$1" . }}{{ end }}`)

	// @section('content') -> {{ define "content" }}
	sectionRe := regexp.MustCompile(`@section\s*\(['"](.*?)['"]\)`)
	compiled = sectionRe.ReplaceAllString(compiled, `{{ define "$1" }}`)
	compiled = strings.ReplaceAll(compiled, "@endsection", "{{ end }}")
	compiled = strings.ReplaceAll(compiled, "@show", "{{ end }}")

	// 4. Includes
	// @include('partials.header') -> {{ template "partials.header" . }}
	includeRe := regexp.MustCompile(`@include\s*\(['"](.*?)['"]\)`)
	compiled = includeRe.ReplaceAllString(compiled, `{{ template "$1" . }}`)

	// 5. CSRF
	compiled = strings.ReplaceAll(compiled, "@csrf", `<input type="hidden" name="_token" value="{{ ._csrf_token }}">`)

	return compiled
}

// CompileFile reads a file, compiles it, and caches the output, returning the path to the cached file.
func (c *Compiler) CompileFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(content))
	cachedPath := filepath.Join(c.CachePath, hash+".gohtml")

	// Check if already compiled
	if _, err := os.Stat(cachedPath); err == nil {
		return cachedPath, nil
	}

	compiled := c.CompileString(string(content))

	// Ensure cache dir exists
	os.MkdirAll(c.CachePath, 0755)

	err = os.WriteFile(cachedPath, []byte(compiled), 0644)
	if err != nil {
		return "", err
	}

	return cachedPath, nil
}
