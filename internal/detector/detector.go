package detector

// Hint provides optional context for language detection.
type Hint struct {
	LineCount  int
	Candidates []string
}

// Detector is a generic language detector.
type Detector interface {
	Detect(code []byte, hint Hint) (lang string, ok bool, err error)
}
