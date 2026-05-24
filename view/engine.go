package view

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	authaccess "github.com/mechneerd/gow/auth/access"
	"github.com/mechneerd/gow/localization"
)

// Engine is the Goblade view rendering engine.
type Engine struct {
	ViewPaths []string
	Compiler  *Compiler
}

// NewEngine creates a new view engine.
func NewEngine(viewPaths []string, cachePath string) *Engine {
	return &Engine{
		ViewPaths: viewPaths,
		Compiler:  NewCompiler(cachePath),
	}
}

// Make compiles and renders a view by name (e.g., "auth.login")
func (e *Engine) Make(name string, data map[string]any) (string, error) {
	// Convert "auth.login" to "auth/login.goblade" (preferred), .blade.php (compat), or .html
	relPath := strings.ReplaceAll(name, ".", "/")
	var absPath string

	for _, vp := range e.ViewPaths {
		// Prefer .goblade (Go-native), fall back to .blade.php for compatibility, then .html
		possiblePaths := []string{
			filepath.Join(vp, relPath+".goblade"),
			filepath.Join(vp, relPath+".blade.php"),
			filepath.Join(vp, relPath+".html"),
		}

		for _, p := range possiblePaths {
			if _, err := os.Stat(p); err == nil {
				absPath = p
				break
			}
		}
		if absPath != "" {
			break
		}
	}

	if absPath == "" {
		return "", os.ErrNotExist
	}

	// For simplicity in Phase 3, we'll compile the main file.
	// If the file uses @extends, we'd normally trace and compile the layout too.
	// For this iteration, we compile the file, parse it into Go templates, and execute.
	
	var rootLayoutBaseName string
	currentPath := absPath
	visited := make(map[string]bool)
	var filesToParse []string

	for {
		if visited[currentPath] {
			return "", fmt.Errorf("cycle detected in layout chain involving: %s", currentPath)
		}
		visited[currentPath] = true

		cachedPath, err := e.Compiler.CompileFile(currentPath)
		if err != nil {
			return "", err
		}
		filesToParse = append(filesToParse, cachedPath)

		compiledContent, _ := os.ReadFile(cachedPath)
		contentStr := string(compiledContent)

		// Look for extends
		// e.g. {{/* extends "layouts.app" */}}
		extendsRe := regexp.MustCompile(`\{\{/\*\s*extends\s*["'](.*?)["']\s*\*/\}\}`)
		matches := extendsRe.FindStringSubmatch(contentStr)

		if len(matches) > 1 {
			layoutName := matches[1]
			layoutRelPath := strings.ReplaceAll(layoutName, ".", "/")
			var layoutAbsPath string
			for _, vp := range e.ViewPaths {
				// Prefer .goblade for layouts too
				possiblePaths := []string{
					filepath.Join(vp, layoutRelPath+".goblade"),
					filepath.Join(vp, layoutRelPath+".blade.php"),
					filepath.Join(vp, layoutRelPath+".html"),
				}
				for _, p := range possiblePaths {
					if _, err := os.Stat(p); err == nil {
						layoutAbsPath = p
						break
					}
				}
				if layoutAbsPath != "" {
					break
				}
			}
			
			if layoutAbsPath == "" {
				return "", fmt.Errorf("layout not found: %s", layoutName)
			}
			
			currentPath = layoutAbsPath
		} else {
			// No more extends. This is the root layout (or just the normal view if no extends)
			rootLayoutBaseName = filepath.Base(cachedPath)
			break
		}
	}

	onceMap := make(map[string]bool)

	// Helper to extract the current authenticated user from common data keys
	getAuthUser := func(d map[string]any) any {
		candidates := []string{"Auth", "User", "auth", "user", "CurrentUser", "AuthenticatedUser"}
		for _, key := range candidates {
			if v, ok := d[key]; ok && v != nil {
				return v
			}
		}
		return nil
	}

	getCanFunc := func(d map[string]any) func(string, ...any) bool {
		return func(ability string, args ...any) bool {
			user := getAuthUser(d)
			if user == nil {
				return false
			}

			if gateIface, ok := d["Gate"]; ok {
				if gate, ok := gateIface.(*authaccess.Gate); ok {
					return gate.Allows(user, ability, args...)
				}
			}
			return false
		}
	}

	funcMap := template.FuncMap{
		"once": func(id string) bool {
			if onceMap[id] {
				return false
			}
			onceMap[id] = true
			return true
		},
		"while": func() []struct{} {
			// Returns a very large slice to simulate a while loop.
			// Go templates do not natively support infinite loops.
			return make([]struct{}, 100000)
		},
		// raw allows unescaped output: {!! $html !!}
		"raw": func(s any) template.HTML {
			switch v := s.(type) {
			case string:
				return template.HTML(v)
			case template.HTML:
				return v
			default:
				return template.HTML(fmt.Sprintf("%v", v))
			}
		},
		// $loop support for @foreach
		"mkloop": func(index, total int) *Loop {
			return newLoop(index, total)
		},
		// Component helper: merges data + attributes + slot content
		// component helper used by the compiler for <x-*> tags
		"component": func(name string, data map[string]any, attrStr, slotContent string) (string, error) {
			attrs := parseAttributes(attrStr)

			// Build context for the component
			componentData := make(map[string]any)
			for k, v := range data {
				componentData[k] = v
			}
			componentData["attributes"] = attrs
			componentData["slot"] = slotContent
			componentData["__component"] = name

			// Render the component view (e.g. "components.alert")
			return e.Make("components."+name, componentData)
		},

		// ==================== Auth Directives Support ====================

		// auth() returns the current authenticated user (if any)
		"auth": func() any {
			// Common places where controllers put the user
			candidates := []string{"Auth", "User", "auth", "user", "CurrentUser"}
			for _, key := range candidates {
				if v, ok := data[key]; ok && v != nil {
					return v
				}
			}
			return nil
		},

		// guest() returns true if no user is authenticated
		"guest": func() bool {
			return getAuthUser(data) == nil
		},

		// can(ability, args...)
		"can": getCanFunc(data),

		// cannot(ability, args...)
		"cannot": func(ability string, args ...any) bool {
			return !getCanFunc(data)(ability, args...)
		},

		// canany(ability1, ability2, ...)
		"canany": func(abilities ...any) bool {
			canFunc := getCanFunc(data)
			for _, a := range abilities {
				if s, ok := a.(string); ok {
					if canFunc(s) {
						return true
					}
				}
			}
			return false
		},

		// Localization helpers
		"__":     localization.Translate,
		"trans":  localization.Trans,
		"lang":   localization.Translate,
	}

	// baseName is arbitrary for the template set, but ParseFiles will use the base name of each file for its templates
	tmpl, err := template.New(rootLayoutBaseName).Funcs(funcMap).ParseFiles(filesToParse...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, rootLayoutBaseName, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

