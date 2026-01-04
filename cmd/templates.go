package cmd

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/prbf2-tools/svctl/pkg/templates"
	"github.com/spf13/cobra"
)

type templatesOpts struct {
	values []string
}

func newTemplatesOpts() *templatesOpts {
	return &templatesOpts{}
}

func (opts *templatesOpts) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&opts.values, "values", "f", []string{}, "Path to values file(s)")
	cmd.MarkFlagFilename("values", "yaml", "yml", "json")
}

func (opts *templatesOpts) MergedValues() (templates.Values, error) {
	valuesSources := make([]templates.Values, len(opts.values))
	for i, path := range opts.values {
		var values templates.Values
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		err = yaml.NewDecoder(file).Decode(&values)
		file.Close()
		if err != nil {
			return nil, err
		}

		valuesSources[i] = values
	}

	return templates.MergeValues(valuesSources...)
}

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
