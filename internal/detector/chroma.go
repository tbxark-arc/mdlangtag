package detector

import (
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
)

// ChromaDetector uses chroma's analyzers for language detection.
type ChromaDetector struct {
	MinBytes int
}

// NewChromaDetector returns a chroma-backed detector with sane defaults.
func NewChromaDetector() *ChromaDetector {
	return &ChromaDetector{MinBytes: 8}
}

// Detect attempts to identify the language of the given code snippet.
func (d *ChromaDetector) Detect(code []byte) (string, bool, error) {
	if d != nil && d.MinBytes > 0 && len(code) < d.MinBytes {
		return "", false, nil
	}

	lexer := lexers.Analyse(string(code))
	if lexer == nil {
		return "", false, nil
	}

	name := strings.ToLower(strings.TrimSpace(lexer.Config().Name))
	if name == "" {
		return "", false, nil
	}
	return name, true, nil
}
