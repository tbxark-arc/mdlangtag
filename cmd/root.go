package cmd

import (
	"github.com/spf13/cobra"

	"github.com/TBXark/mdlangtag/internal/cli"
	"github.com/TBXark/mdlangtag/internal/detector"
)

// Execute runs the root command.
func Execute() error {
	return NewRootCommand().Execute()
}

// NewRootCommand constructs the CLI root command.
func NewRootCommand() *cobra.Command {
	var cfg cli.Config

	cmd := &cobra.Command{
		Use:           "mdlangtag [paths...]",
		Short:         "Detect and fill fenced code block languages in Markdown files",
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cli.FinalizeConfig(&cfg, args); err != nil {
				return err
			}
			runner := cli.NewRunner(cfg, detector.NewChromaDetector())
			return runner.Run()
		},
	}

	cli.BindFlags(cmd, &cfg)
	return cmd
}
