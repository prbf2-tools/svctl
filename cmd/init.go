package cmd

import (
	"github.com/sboon-gg/svctl/internal/settings"
	"github.com/spf13/cobra"
)

type initOpts struct {
	serverOpts
	templatesRepo string
	token         string
}

func newInitOpts() *initOpts {
	return &initOpts{
		serverOpts: *newServerOpts(),
	}
}

func initCmd() *cobra.Command {
	opts := newInitOpts()

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize svctl dir (.svctl) in PRBF2 installation directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run()
		},
	}

	opts.serverOpts.AddFlags(cmd)
	cmd.Flags().StringVar(&opts.templatesRepo, "templates-repo", "", "Repository with templates")
	cmd.Flags().StringVar(&opts.token, "token", "", "Token to use when cloning templates repo")

	return cmd
}

func (opts *initOpts) Run() error {
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
	rootCmd.AddCommand(initCmd())
}
