package view

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	// Convert "auth.login" to "auth/login.blade.php" or "auth/login.html"
	relPath := strings.ReplaceAll(name, ".", "/")
	var absPath string

	for _, vp := range e.ViewPaths {
		// Look for .blade.php or .goblade or .html
		possiblePaths := []string{
			filepath.Join(vp, relPath+".blade.php"),
			filepath.Join(vp, relPath+".goblade"),
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
				possiblePaths := []string{
					filepath.Join(vp, layoutRelPath+".blade.php"),
					filepath.Join(vp, layoutRelPath+".goblade"),
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
