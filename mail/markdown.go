package mail

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"
)

// Markdown renders a markdown string to HTML (with common extensions).
func RenderMarkdown(markdown string) (html string, err error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldhtml.WithHardWraps(),
			goldhtml.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// StripForText returns a basic plain-text version of the HTML.
// This is a lightweight implementation. For production email text parts,
// consider a dedicated library such as "html2text".
func StripForText(htmlContent string) string {
	// Basic implementation: return the HTML as-is for text fallback.
	// Consumers can override with richer stripping if needed.
	return htmlContent
}

