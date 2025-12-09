package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

// Config mirrors CLI flags.
type Config struct {
	Paths       []string
	Write       bool
	Stdout      bool
	Force       bool
	DefaultLang string
	MinLines    int
	Verbose     bool
	Concurrency int
}

// BindFlags registers CLI flags on the provided Cobra command.
func BindFlags(cmd *cobra.Command, cfg *Config) {
	cmd.Flags().BoolVarP(&cfg.Write, "write", "w", false, "write result back to files")
	cmd.Flags().BoolVar(&cfg.Stdout, "stdout", false, "print output to stdout")
	cmd.Flags().BoolVar(&cfg.Force, "force", false, "overwrite existing language info")
	cmd.Flags().StringVar(&cfg.DefaultLang, "default", "", "fallback language when detection fails")
	cmd.Flags().IntVar(&cfg.MinLines, "min-lines", 0, "skip blocks with fewer than this many lines")
	cmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "enable verbose logging")
	cmd.Flags().IntVarP(&cfg.Concurrency, "concurrency", "j", 1, "number of files to process concurrently")
}

// FinalizeConfig applies defaults and positional args after flag parsing.
func FinalizeConfig(cfg *Config, args []string) error {
	cfg.Paths = args
	if len(cfg.Paths) == 0 {
		return errors.New("no input paths provided")
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 1
	}
	return nil
}
