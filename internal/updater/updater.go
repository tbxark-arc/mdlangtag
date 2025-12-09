package updater

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/TBXark/mdlangtag/internal/detector"
	"github.com/TBXark/mdlangtag/internal/parser"
)

type replacement struct {
	start int
	end   int
	value []byte
}

// UpdateDocument updates fenced code block info strings using the provided detector.
func UpdateDocument(doc *parser.Document, det detector.Detector, opts Options) (Stats, error) {
	var stats Stats
	var reps []replacement

	if det == nil {
		return stats, fmt.Errorf("detector is required")
	}

	err := parser.WalkCodeBlocks(doc, func(cb *parser.CodeBlock) error {
		stats.TotalBlocks++

		if cb.Info != "" && !opts.Force {
			stats.SkippedExisting++
			return nil
		}

		lines := countLines(cb.Code)
		if opts.MinCodeLines > 0 && lines < opts.MinCodeLines {
			stats.SkippedTooShort++
			return nil
		}

		hint := detector.Hint{
			LineCount:  lines,
			Candidates: opts.Candidates,
		}

		lang, ok, err := det.Detect(cb.Code, hint)
		if err != nil {
			return err
		}

		if !ok {
			stats.DetectFailed++
			if opts.DefaultLang == "" {
				return nil
			}
			lang = opts.DefaultLang
		}

		lang = strings.TrimSpace(lang)
		if lang == "" {
			stats.DetectFailed++
			return nil
		}

		rep, err := buildReplacement(doc, cb, lang)
		if err != nil {
			return err
		}
		reps = append(reps, rep)
		stats.UpdatedBlocks++
		return nil
	})
	if err != nil {
		return stats, err
	}

	if len(reps) == 0 {
		return stats, nil
	}

	sort.Slice(reps, func(i, j int) bool {
		return reps[i].start < reps[j].start
	})

	newSource, err := applyReplacements(doc.Source, reps)
	if err != nil {
		return stats, err
	}

	doc.Source = newSource
	if err := parser.Refresh(doc); err != nil {
		return stats, err
	}

	return stats, nil
}

func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte{'\n'}) + 1
}

func buildReplacement(doc *parser.Document, cb *parser.CodeBlock, lang string) (replacement, error) {
	if cb.StartPos < 0 || cb.EndPos <= cb.StartPos {
		return replacement{}, fmt.Errorf("missing fence positions for code block")
	}

	line := doc.Source[cb.StartPos:cb.EndPos]

	indent, fence, ending := cb.Indent, cb.Fence, cb.LineEnding
	if indent == "" || fence == "" || ending == "" {
		ind, fnc, end := parseFenceLineFromSource(line)
		if indent == "" {
			indent = ind
		}
		if fence == "" {
			fence = fnc
		}
		if ending == "" {
			ending = end
		}
	}

	if fence == "" {
		return replacement{}, fmt.Errorf("unable to determine fence characters")
	}

	spacing := 1
	if cb.InfoStart >= 0 {
		base := cb.StartPos + len(indent) + len(fence)
		if cb.InfoStart >= base {
			spacing = cb.InfoStart - base
		}
		if spacing < 0 {
			spacing = 1
		}
	}

	var b strings.Builder
	b.WriteString(indent)
	b.WriteString(fence)
	if lang != "" {
		if spacing < 1 {
			spacing = 1
		}
		b.WriteString(strings.Repeat(" ", spacing))
		b.WriteString(lang)
	}
	b.WriteString(ending)

	return replacement{
		start: cb.StartPos,
		end:   cb.EndPos,
		value: []byte(b.String()),
	}, nil
}

func parseFenceLineFromSource(line []byte) (indent string, fence string, ending string) {
	content := line
	if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
		content = line[:len(line)-2]
		ending = "\r\n"
	} else if len(line) >= 1 && line[len(line)-1] == '\n' {
		content = line[:len(line)-1]
		ending = "\n"
	}

	idx := 0
	for idx < len(content) && (content[idx] == ' ' || content[idx] == '\t') {
		idx++
	}
	indent = string(content[:idx])

	start := idx
	var fenceChar byte
	for idx < len(content) && (content[idx] == '`' || content[idx] == '~') {
		if fenceChar == 0 {
			fenceChar = content[idx]
		} else if content[idx] != fenceChar {
			break
		}
		idx++
	}
	fence = string(content[start:idx])
	return
}

func applyReplacements(src []byte, reps []replacement) ([]byte, error) {
	var buf bytes.Buffer
	cursor := 0
	for _, rep := range reps {
		if rep.start > len(src) {
			return nil, fmt.Errorf("replacement start %d beyond source length %d", rep.start, len(src))
		}
		if rep.start < cursor {
			return nil, fmt.Errorf("overlapping replacement at %d", rep.start)
		}
		if rep.end > len(src) {
			rep.end = len(src)
		}
		buf.Write(src[cursor:rep.start])
		buf.Write(rep.value)
		cursor = rep.end
	}
	buf.Write(src[cursor:])
	return buf.Bytes(), nil
}
