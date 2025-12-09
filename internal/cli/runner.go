package cli

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/TBXark/mdlangtag/internal/detector"
	"github.com/TBXark/mdlangtag/internal/parser"
	"github.com/TBXark/mdlangtag/internal/updater"
)

// Runner owns the CLI execution pipeline.
type Runner struct {
	Config   Config
	Detector detector.Detector
}

func NewRunner(cfg Config, det detector.Detector) *Runner {
	return &Runner{Config: cfg, Detector: det}
}

type fileResult struct {
	path  string
	stats updater.Stats
	err   error
}

// Run processes all configured files.
func (r *Runner) Run() error {
	files, err := ListMarkdownFiles(r.Config.Paths)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		if r.Config.Verbose {
			fmt.Println("no markdown files found")
		}
		return nil
	}

	concurrency := r.Config.Concurrency
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}
	if r.Config.Stdout || !r.Config.Write {
		concurrency = 1
	}

	opts := updater.Options{
		Force:        r.Config.Force,
		DefaultLang:  r.Config.DefaultLang,
		MinCodeLines: r.Config.MinLines,
		Verbose:      r.Config.Verbose,
		Candidates:   r.Config.GetCandidates(),
	}

	results := make(chan fileResult, len(files))
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	for _, path := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()
			stats, err := r.processFile(p, opts)
			results <- fileResult{path: p, stats: stats, err: err}
		}(path)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var total updater.Stats
	var firstErr error

	for res := range results {
		if res.err != nil && firstErr == nil {
			firstErr = res.err
		}
		total.TotalBlocks += res.stats.TotalBlocks
		total.UpdatedBlocks += res.stats.UpdatedBlocks
		total.SkippedExisting += res.stats.SkippedExisting
		total.SkippedTooShort += res.stats.SkippedTooShort
		total.DetectFailed += res.stats.DetectFailed

		if r.Config.Verbose {
			fmt.Printf("%s: %+v\n", res.path, res.stats)
		}
	}

	if r.Config.Verbose && len(files) > 1 {
		fmt.Printf("Total: %+v\n", total)
	}

	if firstErr != nil {
		return firstErr
	}
	return nil
}

func (r *Runner) processFile(path string, opts updater.Options) (updater.Stats, error) {
	var stats updater.Stats

	doc, err := parser.ParseFile(path)
	if err != nil {
		return stats, fmt.Errorf("parse %s: %w", path, err)
	}

	stats, err = updater.UpdateDocument(doc, r.Detector, opts)
	if err != nil {
		return stats, fmt.Errorf("update %s: %w", path, err)
	}

	output, err := parser.Render(doc)
	if err != nil {
		return stats, fmt.Errorf("render %s: %w", path, err)
	}

	changed := stats.UpdatedBlocks > 0
	writeOut := r.Config.Write && changed
	printOut := r.Config.Stdout || (!r.Config.Write)

	if writeOut {
		info, statErr := os.Stat(path)
		mode := os.FileMode(0644)
		if statErr == nil {
			mode = info.Mode()
		}
		if err := os.WriteFile(path, output, mode); err != nil {
			return stats, fmt.Errorf("write %s: %w", path, err)
		}
	}

	if printOut {
		fmt.Print(string(output))
		if len(output) == 0 || output[len(output)-1] != '\n' {
			fmt.Println()
		}
	}

	return stats, nil
}
