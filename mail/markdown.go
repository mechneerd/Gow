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

// StripMarkdownForText is a very basic way to get a plain text version.
// For production, consider using a proper HTML-to-text converter.
func StripForText(htmlContent string) string {
	// Very naive stripping for now. In real apps you'd use a proper lib like "html2text".
	// This is acceptable as a starting point.
	return htmlContent // TODO: improve with proper stripper
}
