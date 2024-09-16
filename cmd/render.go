package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type renderOpts struct {
	*serverOpts
	dryRun bool
}

func newRenderOpts() *renderOpts {
	return &renderOpts{
		serverOpts: newServerOpts(),
	}
}

func init() {
	rootCmd.AddCommand(renderCmd())
}

func renderCmd() *cobra.Command {
	opts := newRenderOpts()

	cmd := &cobra.Command{
		Use:          "render",
		Short:        "Render templates",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(cmd)
		},
	}

	opts.serverOpts.AddFlags(cmd)

	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Print out rendered files")

	return cmd
}

func (opts *renderOpts) Run(cmd *cobra.Command) error {
	sv, err := opts.Server()
	if err != nil {
		return err
	}

	if opts.dryRun {
		outputs, err := sv.DryRender()
		if err != nil {
			return err
		}
		for _, out := range outputs {
			fmt.Printf("File: %s\n---\n%s", out.Destination, string(out.Content))
		}
	} else {
		err := sv.Render()
		if err != nil {
			return err
		}
	}

	return nil
}
