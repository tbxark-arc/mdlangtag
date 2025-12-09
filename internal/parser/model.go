package parser

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
)

// Document represents a parsed Markdown document.
type Document struct {
	Source []byte
	Root   ast.Node
	MD     goldmark.Markdown
	Path   string

	fences []fenceBlock
}

// CodeBlock wraps a fenced code block node with handy metadata.
type CodeBlock struct {
	Node       *ast.FencedCodeBlock
	Code       []byte
	Info       string
	Fence      string
	StartPos   int
	EndPos     int
	InfoStart  int
	InfoEnd    int
	Indent     string
	LineEnding string
}
