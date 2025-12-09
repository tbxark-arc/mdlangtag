package detector

import (
	"strings"

	"github.com/go-enry/go-enry/v2"
)

// Top 20 most popular programming languages for accurate detection
var topLanguages = []string{
	"Python",
	"JavaScript",
	"TypeScript",
	"Java",
	"C",
	"C++",
	"C#",
	"PHP",
	"Ruby",
	"Go",
	"Rust",
	"Kotlin",
	"Swift",
	"Objective-C",
	"Shell",
	"HTML",
	"CSS",
	"SQL",
	"JSON",
	"XML",
}

// EnryDetector uses go-enry for language detection.
type EnryDetector struct {
	MinBytes int
}

// NewEnryDetector returns an Enry-backed detector with sane defaults.
func NewEnryDetector() *EnryDetector {
	return &EnryDetector{MinBytes: 8}
}

// Detect attempts to identify the language of the given code snippet.
func (d *EnryDetector) Detect(code []byte, hint Hint) (string, bool, error) {
	if d != nil && d.MinBytes > 0 && len(code) < d.MinBytes {
		return "", false, nil
	}

	// Use provided candidates if available, otherwise use default top languages
	var candidateList []string
	if len(hint.Candidates) > 0 {
		candidateList = hint.Candidates
	} else {
		candidateList = topLanguages
	}

	candidates := make([]string, 0, len(candidateList))
	candidates = append(candidates, candidateList...)

	lang, _ := enry.GetLanguageByClassifier(code, candidates)
	if lang == "" {
		return "", false, nil
	}
	return strings.ToLower(lang), true, nil
}
