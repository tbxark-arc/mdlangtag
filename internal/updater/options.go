package updater

// Options controls how fenced code blocks are updated.
type Options struct {
	Force        bool
	DefaultLang  string
	MinCodeLines int
	Verbose      bool
}

// Stats captures counts from a single update run.
type Stats struct {
	TotalBlocks     int
	UpdatedBlocks   int
	SkippedExisting int
	SkippedTooShort int
	DetectFailed    int
}
