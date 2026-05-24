package view

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

func TestCompilerDirectives(t *testing.T) {
	compiler := NewCompiler("./cache")

	t.Run("while directive", func(t *testing.T) {
		raw := `{{ $i := 0 }}
@while(lt $i 5)
{{ $i }}
{{ $i = add $i 1 }}
@endwhile`
		compiled := compiler.CompileString(raw)

		if !strings.Contains(compiled, "{{ range while }}") {
			t.Errorf("Expected compiled while to contain 'range while', got:\n%s", compiled)
		}
		if !strings.Contains(compiled, "{{ if not (lt $i 5) }}{{ break }}{{ end }}") {
			t.Errorf("Expected compiled while to contain break condition, got:\n%s", compiled)
		}

		// Let's actually execute it to be sure
		funcMap := template.FuncMap{
			"while": func() []struct{} { return make([]struct{}, 10) },
			"add":   func(a, b int) int { return a + b },
			"lt":    func(a, b int) bool { return a < b },
		}

		tmpl, err := template.New("test").Funcs(funcMap).Parse(compiled)
		if err != nil {
			t.Fatalf("Failed to parse compiled template: %v\nCompiled:\n%s", err, compiled)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, nil)
		if err != nil {
			t.Fatalf("Failed to execute: %v", err)
		}

		out := strings.ReplaceAll(buf.String(), "\n", "")
		out = strings.ReplaceAll(out, " ", "")
		if out != "01234" {
			t.Errorf("Expected '01234', got '%s'", out)
		}
	})

	t.Run("once directive", func(t *testing.T) {
		raw := `@once
First
@endonce
@once
Second
@endonce`

		compiled := compiler.CompileString(raw)
		if !strings.Contains(compiled, "once \"") {
			t.Fatalf("Expected once func call, got:\n%s", compiled)
		}
	})
}

