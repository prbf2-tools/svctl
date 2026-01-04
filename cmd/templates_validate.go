package cmd

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/prbf2-tools/svctl/pkg/templates"
	"github.com/spf13/cobra"
)

type templatesValidateOpts struct {
	*templatesOpts
}

func newTemplatesValidateOpts() *templatesValidateOpts {
	return &templatesValidateOpts{
		templatesOpts: newTemplatesOpts(),
	}
}

func templatesValidateCmd() *cobra.Command {
	opts := newTemplatesValidateOpts()

	cmd := &cobra.Command{
		Use:   "validate <TEMPLATES_PATH>",
		Short: "Validate templates in the specified path",
		RunE:  opts.Run,
	}

	cmd.Args = cobra.ExactArgs(1)

	opts.AddFlags(cmd)

	return cmd
}

func (o *templatesValidateOpts) AddFlags(cmd *cobra.Command) {
	// Add flags here if needed in the future
}

func (o *templatesValidateOpts) Run(cmd *cobra.Command, args []string) error {
	templatesPath := args[0]

	renderer, err := templates.NewFromPath(templatesPath)
	if err != nil {
		return err
	}

	schema, err := renderer.Schema()
	if err != nil {
		return err
	}

	cmd.Println("Validating defaults...")

	defaultsContent, err := renderer.DefaultsContent()
	if err != nil {
		return err
	}

	var defaults map[string]any
	err = yaml.Unmarshal(defaultsContent, &defaults)
	if err != nil {
		return err
	}

	err = schema.Validate(defaults)
	if err != nil {
		return fmt.Errorf("defaults validation failed: %w", err)
	}

	cmd.Println("Defaults validation succeeded")
	cmd.Println("Validating values merged with defaults...")

	mergedValues, err := o.MergedValues()
	if err != nil {
		return err
	}

	mergedValues, err = templates.MergeValues(defaults, mergedValues)
	if err != nil {
		return err
	}
	m := map[string]any(mergedValues)

	err = schema.Validate(m)
	if err != nil {
		return fmt.Errorf("merged values validation failed: %w", err)
	}

	cmd.Println("Merged values validation succeeded")
	return nil
}

func init() {
	templatesCmd.AddCommand(templatesValidateCmd())
}
