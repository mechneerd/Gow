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

	// 1. Echo statements
	// {!! $var !!} -> unescaped raw output (requires "raw" template func registered in Engine)
	compiled = regexp.MustCompile(`\{\!\!\s*(.+?)\s*\!\!\}`).ReplaceAllString(compiled, `{{ raw $1 }}`)

	// 2. Control Structures
	compiled = strings.ReplaceAll(compiled, "@else", "{{ else }}")
	compiled = strings.ReplaceAll(compiled, "@endif", "{{ end }}")
	compiled = strings.ReplaceAll(compiled, "@endforeach", "{{ end }}")
	compiled = strings.ReplaceAll(compiled, "@endfor", "{{ end }}")

	// @if(condition) -> {{ if condition }}
	ifRe := regexp.MustCompile(`@if\s*\((.*?)\)`)
	compiled = ifRe.ReplaceAllString(compiled, "{{ if $1 }}")

	// @elseif(condition) -> {{ else if condition }}
	elseIfRe := regexp.MustCompile(`@elseif\s*\((.*?)\)`)
	compiled = elseIfRe.ReplaceAllString(compiled, "{{ else if $1 }}")

	// @foreach($items as $item) with full $loop support
	// We compile it to use index form + mkloop so $loop is available inside the block.
	foreachRe := regexp.MustCompile(`@foreach\s*\(\s*(.+?)\s+as\s+(.+?)\s*\)`)
	compiled = foreachRe.ReplaceAllString(compiled, `{{ range $index, $2 := $1 }}{{ $loop := mkloop $index (len $1) }}`)

	// @while(condition) - We implement this using a custom template func 'while'
	// that returns a large slice, and we break out of it when the condition is false.
	whileRe := regexp.MustCompile(`@while\s*\((.*?)\)`)
	compiled = whileRe.ReplaceAllString(compiled, "{{ range while }}{{ if not ($1) }}{{ break }}{{ end }}")
	compiled = strings.ReplaceAll(compiled, "@endwhile", "{{ end }}")

	// @for(init; condition; post) - Go templates don't natively support C-style for loops inside templates easily.
	// We will compile this into a generic action to warn users or simulate it. 
	// For simplicity, we assume users will use @foreach in Go context.
	forRe := regexp.MustCompile(`@for\s*\((.*?)\)`)
	compiled = forRe.ReplaceAllString(compiled, "{{ range $1 }}")

	// @switch / @case
	compiled = strings.ReplaceAll(compiled, "@endswitch", "{{ end }}")
	switchRe := regexp.MustCompile(`@switch\s*\((.*?)\)`)
	compiled = switchRe.ReplaceAllString(compiled, `{{ $switch_var := $1 }}`)
	
	caseRe := regexp.MustCompile(`@case\s*\((.*?)\)`)
	// We use a simple if/else if chain for switch in Go templates
	compiled = caseRe.ReplaceAllStringFunc(compiled, func(match string) string {
		val := caseRe.FindStringSubmatch(match)[1]
		return fmt.Sprintf(`{{ if eq $switch_var %s }}`, val)
	})
	compiled = strings.ReplaceAll(compiled, "@break", "") // Go template 'if' blocks don't need breaks
	
	// @default
	compiled = strings.ReplaceAll(compiled, "@default", "{{ else }}")

	// 3. Layouts & Sections
	// @extends('layouts.app') -> We leave it in the compiled code as a special tag for the Engine to extract dependencies
	// Actually, in Go, the layout file is parsed and executed. The layout file executes {{ block "content" . }}{{ end }}.
	// The child file defines {{ define "content" }}...{{ end }}.
	// We just need the Engine to parse BOTH files, and ExecuteTemplate(w, "layout.html", data).
	extendsRe := regexp.MustCompile(`@extends\s*\(['"](.*?)['"]\)`)
	compiled = extendsRe.ReplaceAllString(compiled, `{{/* extends "$1" */}}`)

	// @yield('content') -> {{ block "content" . }}{{ end }}
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

	// @lang('welcome.message') -> {{ __ "welcome.message" }}
	langRe := regexp.MustCompile(`@lang\s*\(\s*['"](.+?)['"]\s*\)`)
	compiled = langRe.ReplaceAllString(compiled, `{{ __ "$1" }}`)

	// 6. Authorization Directives (Full Implementation)

	// @auth -> {{ if auth }}
	compiled = strings.ReplaceAll(compiled, "@auth", `{{ if auth }}`)
	compiled = strings.ReplaceAll(compiled, "@endauth", "{{ end }}")

	// @guest -> {{ if guest }}
	compiled = strings.ReplaceAll(compiled, "@guest", `{{ if guest }}`)
	compiled = strings.ReplaceAll(compiled, "@endguest", "{{ end }}")

	// @can('update', $post)
	canRe := regexp.MustCompile(`@can\s*\((.*?)\)`)
	compiled = canRe.ReplaceAllStringFunc(compiled, func(match string) string {
		val := canRe.FindStringSubmatch(match)[1]
		val = strings.ReplaceAll(val, ",", " ")
		val = strings.ReplaceAll(val, "'", "\"")
		return fmt.Sprintf(`{{ if can %s }}`, val)
	})
	compiled = strings.ReplaceAll(compiled, "@endcan", "{{ end }}")

	// @cannot('update', $post)
	cannotRe := regexp.MustCompile(`@cannot\s*\((.*?)\)`)
	compiled = cannotRe.ReplaceAllStringFunc(compiled, func(match string) string {
		val := cannotRe.FindStringSubmatch(match)[1]
		val = strings.ReplaceAll(val, ",", " ")
		val = strings.ReplaceAll(val, "'", "\"")
		return fmt.Sprintf(`{{ if not (can %s) }}`, val)
	})
	compiled = strings.ReplaceAll(compiled, "@endcannot", "{{ end }}")

	// @canany('update', 'delete', $post)
	canAnyRe := regexp.MustCompile(`@canany\s*\((.*?)\)`)
	compiled = canAnyRe.ReplaceAllStringFunc(compiled, func(match string) string {
		val := canAnyRe.FindStringSubmatch(match)[1]
		val = strings.ReplaceAll(val, ",", " ")
		val = strings.ReplaceAll(val, "'", "\"")
		return fmt.Sprintf(`{{ if canany %s }}`, val)
	})
	compiled = strings.ReplaceAll(compiled, "@endcanany", "{{ end }}")

	// 7. Advanced Directives
	// @class(['p-4', 'font-bold' => true]) -> we would map to a func `class(...)`
	classRe := regexp.MustCompile(`@class\s*\((.*?)\)`)
	compiled = classRe.ReplaceAllString(compiled, `class=$1`) // Simplified for demo

	// @checked(true) -> {{ if true }}checked="checked"{{ end }}
	checkedRe := regexp.MustCompile(`@checked\s*\((.*?)\)`)
	compiled = checkedRe.ReplaceAllString(compiled, `{{ if $1 }}checked="checked"{{ end }}`)

	// @method('DELETE') -> <input type="hidden" name="_method" value="DELETE">
	methodRe := regexp.MustCompile(`@method\s*\(\s*['"](.+?)['"]\s*\)`)
	compiled = methodRe.ReplaceAllString(compiled, `<input type="hidden" name="_method" value="$1">`)

	// @error('email') ... @enderror -> {{ if error email }}...{{ end }}
	errorRe := regexp.MustCompile(`@error\s*\(\s*['"](.+?)['"]\s*\)`)
	compiled = errorRe.ReplaceAllString(compiled, `{{ if error "$1" }}`)
	compiled = strings.ReplaceAll(compiled, "@enderror", "{{ end }}")

	// @props declaration for components
	compiled = strings.ReplaceAll(compiled, "@props", "{{/* @props")

	// @aware - access parent component data
	awareRe := regexp.MustCompile(`@aware\s*\((.*?)\)`)
	compiled = awareRe.ReplaceAllString(compiled, `{{ /* aware: $1 */ }}`)

	// @json($data) -> JSON encode data
	jsonRe := regexp.MustCompile(`@json\s*\((.*?)\)`)
	compiled = jsonRe.ReplaceAllString(compiled, `{{ json $1 }}`)

	// Generate a short hash of the raw template to ensure @once IDs are unique per file
	fileHash := fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))[:8]

	// @once
	onceCount := 0
	onceRe := regexp.MustCompile(`@once`)
	compiled = onceRe.ReplaceAllStringFunc(compiled, func(match string) string {
		onceCount++
		return fmt.Sprintf(`{{ if once "%s_%d" }}`, fileHash, onceCount)
	})
	compiled = strings.ReplaceAll(compiled, "@endonce", `{{ end }}`)

	// ==================== Components & Slots ====================

	// Self-closing: <x-alert type="error" />
	selfClose := regexp.MustCompile(`<x-([a-zA-Z0-9_-]+)([^>]*?)\s*/>`)
	compiled = selfClose.ReplaceAllStringFunc(compiled, func(m string) string {
		matches := selfClose.FindStringSubmatch(m)
		name := matches[1]
		attrs := matches[2]
		return fmt.Sprintf(`{{ component "%s" . "%s" "" }}`, name, attrs)
	})

	// With content (slots) - handle nesting via iterative innermost-first replacement
	compiled = replaceNestedComponents(compiled)

	return compiled
}

// replaceNestedComponents handles <x-name>...</x-name> with proper nesting support.
// It iteratively replaces innermost components first.
func replaceNestedComponents(s string) string {
	tagRe := regexp.MustCompile(`<x-([a-zA-Z0-9_-]+)([^>]*?)>`)
	changed := true
	for changed {
		changed = false
		idx := tagRe.FindStringIndex(s)
		if idx == nil {
			break
		}
		// Find the matching closing tag by counting nesting depth
		start := idx[0]
		tagStart := s[idx[0]:idx[1]]
		matches := tagRe.FindStringSubmatch(tagStart)
		name := matches[1]

		depth := 1
		pos := idx[1]
		closePattern := "</x-" + name + ">"
		openPattern := "<x-" + name

		for depth > 0 && pos < len(s) {
			nextOpen := strings.Index(s[pos:], openPattern)
			nextClose := strings.Index(s[pos:], closePattern)

			if nextClose == -1 {
				break
			}

			if nextOpen != -1 && nextOpen < nextClose {
				depth++
				pos += nextOpen + len(openPattern)
			} else {
				depth--
				if depth == 0 {
					closeEnd := pos + nextClose + len(closePattern)
					attrs := matches[2]
					slotContent := strings.TrimSpace(s[idx[1]:pos+nextClose])
					replacement := fmt.Sprintf(`{{ component "%s" . "%s" "%s" }}`, name, attrs, slotContent)
					s = s[:start] + replacement + s[closeEnd:]
					changed = true
				} else {
					pos += nextClose + len(closePattern)
				}
			}
		}
	}
	return s
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

