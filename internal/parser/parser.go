package parser

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// ParseBytes parses markdown bytes into a Document.
func ParseBytes(src []byte, path string) (*Document, error) {
	md := goldmark.New()
	root := md.Parser().Parse(text.NewReader(src))
	return &Document{
		Source: src,
		Root:   root,
		MD:     md,
		Path:   path,
		fences: scanFencedCodeBlocks(src),
	}, nil
}

// ParseFile reads the given file and parses it into a Document.
func ParseFile(path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return ParseBytes(data, path)
}

// Render returns the Markdown bytes for the document.
func Render(doc *Document) ([]byte, error) {
	return append([]byte(nil), doc.Source...), nil
}

// Refresh reparses the current document content to refresh AST and fence metadata.
func Refresh(doc *Document) error {
	parsed, err := ParseBytes(doc.Source, doc.Path)
	if err != nil {
		return err
	}
	doc.Root = parsed.Root
	doc.MD = parsed.MD
	doc.fences = parsed.fences
	return nil
}

type fenceBlock struct {
	openStart   int
	openEnd     int
	infoStart   int
	infoEnd     int
	fenceChar   byte
	fenceLength int
	indent      []byte
	lineEnding  []byte
}

func scanFencedCodeBlocks(src []byte) []fenceBlock {
	var blocks []fenceBlock
	var pending *fenceBlock

	offset := 0
	for offset <= len(src) {
		lineStart := offset
		lineEnd := bytes.IndexByte(src[offset:], '\n')
		if lineEnd < 0 {
			lineEnd = len(src)
		} else {
			lineEnd = offset + lineEnd + 1
		}
		line := src[lineStart:lineEnd]
		content, ending := splitLineEnding(line)
		pos, indentWidth := leadingIndent(content)

		if pending == nil {
			if fb, ok := parseFenceLine(content, lineStart, pos, indentWidth, ending); ok {
				pending = &fb
			}
		} else {
			if isClosingFence(content, pending.fenceChar, pending.fenceLength, indentWidth, pos) {
				blocks = append(blocks, *pending)
				pending = nil
			}
		}

		if lineEnd >= len(src) {
			break
		}
		offset = lineEnd
	}
	if pending != nil {
		blocks = append(blocks, *pending)
	}
	return blocks
}

func splitLineEnding(line []byte) ([]byte, []byte) {
	if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
		return line[:len(line)-2], line[len(line)-2:]
	}
	if len(line) >= 1 && line[len(line)-1] == '\n' {
		return line[:len(line)-1], line[len(line)-1:]
	}
	return line, nil
}

func leadingIndent(line []byte) (idx int, width int) {
	for idx < len(line) {
		switch line[idx] {
		case ' ':
			width++
			idx++
		case '\t':
			width += 4
			idx++
		default:
			return
		}
	}
	return
}

func parseFenceLine(content []byte, lineStart, pos, indentWidth int, ending []byte) (fenceBlock, bool) {
	var fb fenceBlock
	if indentWidth >= 4 || pos >= len(content) {
		return fb, false
	}

	ch := content[pos]
	if ch != '`' && ch != '~' {
		return fb, false
	}

	i := pos
	for i < len(content) && content[i] == ch {
		i++
	}
	length := i - pos
	if length < 3 {
		return fb, false
	}

	rest := content[i:]
	left := util.TrimLeftSpaceLength(rest)
	right := util.TrimRightSpaceLength(rest)
	infoStart, infoEnd := -1, -1
	if left < len(rest)-right {
		value := rest[left : len(rest)-right]
		if ch == '`' && bytes.IndexByte(value, '`') > -1 {
			return fb, false
		}
		infoStart = lineStart + i + left
		infoEnd = lineStart + len(content) - right
	}

	indent := make([]byte, pos)
	copy(indent, content[:pos])

	fb = fenceBlock{
		openStart:   lineStart,
		openEnd:     lineStart + len(content) + len(ending),
		infoStart:   infoStart,
		infoEnd:     infoEnd,
		fenceChar:   ch,
		fenceLength: length,
		indent:      indent,
		lineEnding:  ending,
	}
	return fb, true
}

func isClosingFence(content []byte, ch byte, minLen, indentWidth, pos int) bool {
	if indentWidth >= 4 || pos >= len(content) || content[pos] != ch {
		return false
	}
	i := pos
	for i < len(content) && content[i] == ch {
		i++
	}
	if i-pos < minLen {
		return false
	}
	return util.IsBlank(content[i:])
}
