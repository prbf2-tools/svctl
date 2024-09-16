package cmd

import (
	"github.com/sboon-gg/svctl/internal/game"
	"github.com/spf13/cobra"
)

type updateOpts struct {
	path string
}

func updateCmd() *cobra.Command {
	opts := &updateOpts{}

	cmd := &cobra.Command{
		Use:   "update <path>",
		Short: "Manually update the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.path = args[0]
			return opts.Run(cmd)
		},
	}

	return cmd
}

func (o *updateOpts) Run(cmd *cobra.Command) error {
	gameServer, err := game.Open(o.path)
	if err != nil {
		return err
	}

	return gameServer.Update(cmd.Context(), cmd.OutOrStdout(), cmd.InOrStdin(), cmd.OutOrStderr())
}

func init() {
	rootCmd.AddCommand(updateCmd())
}
