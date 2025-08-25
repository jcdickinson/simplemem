package mcp

import (
	"testing"
)

func TestStripMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "This is plain text",
			expected: "This is plain text",
		},
		{
			name:     "heading",
			input:    "# This is a heading",
			expected: "This is a heading",
		},
		{
			name:     "bold and italic",
			input:    "This is **bold** and *italic* text",
			expected: "This is bold and italic text",
		},
		{
			name:     "list items",
			input:    "- Item 1\n- Item 2",
			expected: "Item 1  Item 2",
		},
		{
			name:     "link",
			input:    "[Link text](https://example.com)",
			expected: "Link text",
		},
		{
			name:     "blockquote",
			input:    "> This is a quote",
			expected: "This is a quote",
		},
		{
			name:     "code block",
			input:    "```go\nfunc main() {\n    fmt.Println(\"hello\")\n}\n```",
			expected: "",
		},
		{
			name: "complex markdown",
			input: `# Heading

This is a paragraph with **bold** and *italic* text.

- List item 1
- List item 2

[Link](https://example.com)

> Quote text`,
			expected: "Heading This is a paragraph with bold and italic text. List item 1  List item 2   Link Quote text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("stripMarkdown() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStripMarkdownLength(t *testing.T) {
	// Test that markdown stripping significantly reduces character count
	markdownContent := `# This is a title

This is **bold text** and *italic text*.

- Item 1
- Item 2

[Link](https://example.com)

> This is a quote

` + "```go\nfunc main() {\n    fmt.Println(\"code\")\n}\n```"

	plainText := stripMarkdown(markdownContent)
	
	// Original should be longer than plain text
	if len(plainText) >= len(markdownContent) {
		t.Errorf("Expected plain text to be shorter than markdown, got plain: %d, markdown: %d", 
			len(plainText), len(markdownContent))
	}
	
	// Plain text should be roughly 40-60% of original
	ratio := float64(len(plainText)) / float64(len(markdownContent))
	if ratio < 0.3 || ratio > 0.8 {
		t.Errorf("Expected compression ratio between 0.3-0.8, got %.2f", ratio)
	}
}