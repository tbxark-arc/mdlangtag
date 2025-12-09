package parser

import (
	"strings"

	"github.com/yuin/goldmark/ast"
)

// WalkCodeBlocks walks all fenced code blocks in document order.
func WalkCodeBlocks(doc *Document, fn func(cb *CodeBlock) error) error {
	fenceIdx := 0
	err := ast.Walk(doc.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		block, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		cb := &CodeBlock{
			Node:       block,
			Code:       block.Lines().Value(doc.Source),
			Fence:      "",
			Info:       "",
			StartPos:   -1,
			EndPos:     -1,
			InfoStart:  -1,
			InfoEnd:    -1,
			Indent:     "",
			LineEnding: "",
		}
		if block.Info != nil {
			cb.Info = string(block.Info.Text(doc.Source))
		}

		if fenceIdx < len(doc.fences) {
			meta := doc.fences[fenceIdx]
			cb.Fence = strings.Repeat(string(meta.fenceChar), meta.fenceLength)
			cb.StartPos = meta.openStart
			cb.EndPos = meta.openEnd
			cb.InfoStart = meta.infoStart
			cb.InfoEnd = meta.infoEnd
			cb.Indent = string(meta.indent)
			cb.LineEnding = string(meta.lineEnding)
		}
		fenceIdx++

		if err := fn(cb); err != nil {
			return ast.WalkStop, err
		}

		return ast.WalkSkipChildren, nil
	})

	return err
}
