package cmd

import (
	"github.com/prbf2-tools/svctl/internal/settings"
	"github.com/spf13/cobra"
)

type settingsInitOpts struct {
	serverOpts
	templatesRepo string
	token         string
}

func newSettingsInitOpts() *settingsInitOpts {
	return &settingsInitOpts{
		serverOpts: *newServerOpts(),
	}
}

func settingsInitCmd() *cobra.Command {
	opts := newSettingsInitOpts()

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize svctl dir (.svctl) in PRBF2 installation directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	opts.AddFlags(cmd)

	return cmd
}

func (opts *settingsInitOpts) AddFlags(cmd *cobra.Command) {
	opts.serverOpts.AddFlags(cmd)

	cmd.Flags().StringVar(&opts.templatesRepo, "templates-repo", "", "Repository with templates")
	cmd.Flags().StringVar(&opts.token, "token", "", "Token to use when cloning templates repo")
}

func (opts *settingsInitOpts) Run() error {
	svctlPath, err := opts.SettingsPath()
	if err != nil {
		return err
	}

	_, err = settings.Initialize(svctlPath, &settings.Opts{
		TemplatesRepo: opts.templatesRepo,
		Token:         opts.token,
	})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	settingsCmd.AddCommand(settingsInitCmd())
}
