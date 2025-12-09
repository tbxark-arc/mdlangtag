package detector

// Detector is a generic language detector.
type Detector interface {
	Detect(code []byte) (lang string, ok bool, err error)
}
