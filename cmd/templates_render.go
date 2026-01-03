package cmd

import (
	"os"
	"path/filepath"

	"github.com/sboon-gg/svctl/pkg/templates"
	"github.com/spf13/cobra"
)

type templatesRenderOpts struct {
	*templatesOpts
	outputPath string
	dryRun     bool
}

func newTemplatesRenderOpts() *templatesRenderOpts {
	return &templatesRenderOpts{
		templatesOpts: newTemplatesOpts(),
	}
}

func templatesRenderCmd() *cobra.Command {
	opts := newTemplatesRenderOpts()

	cmd := &cobra.Command{
		Use:   "render <TEMPLATES_PATH>",
		Short: "Render templates with provided values",
		RunE:  opts.Run,
	}

	cmd.Args = cobra.ExactArgs(1)

	opts.AddFlags(cmd)

	return cmd
}

func (opts *templatesRenderOpts) AddFlags(cmd *cobra.Command) {
	opts.templatesOpts.AddFlags(cmd)

	cmd.Flags().StringVarP(&opts.outputPath, "output", "o", "", "Path to output rendered files")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Print out rendered files")

	cmd.MarkFlagDirname("output")
	cmd.MarkFlagsOneRequired("output", "dry-run")
	cmd.MarkFlagsMutuallyExclusive("output", "dry-run")
}

func (opts *templatesRenderOpts) Run(cmd *cobra.Command, args []string) error {
	renderer, err := templates.NewFromPath(args[0])
	if err != nil {
		return err
	}

	mergedValues, err := opts.MergedValues()
	if err != nil {
		return err
	}

	outputs, err := renderer.Render(nil, mergedValues)
	if err != nil {
		return err
	}

	if opts.dryRun {
		for _, output := range outputs {
			cmd.Println("=== " + output.Destination + " ===")
			cmd.Println(string(output.Content))
		}
		return nil
	}

	for _, output := range outputs {
		fullDir := filepath.Join(opts.outputPath, filepath.Dir(output.Destination))
		err = os.MkdirAll(fullDir, 0o755)
		if err != nil {
			return err
		}

		fullPath := filepath.Join(opts.outputPath, output.Destination)
		err = os.WriteFile(fullPath, output.Content, 0o644)
		if err != nil {
			return err
		}

		cmd.Printf("Wrote %s\n", fullPath)
	}

	return nil
}

func init() {
	templatesCmd.AddCommand(templatesRenderCmd())
}
