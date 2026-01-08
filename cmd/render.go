package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type renderOpts struct {
	*grpcClientCmdOpts
	reloadableOnly bool
}

func newRenderOpts() *renderOpts {
	return &renderOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func renderCmd() *cobra.Command {
	opts := newRenderOpts()

	cmd := &cobra.Command{
		Use:          "render <server-id>",
		Short:        "Render server templates",
		Long:         `Send a render signal to daemon to render server templates`,
		SilenceUsage: true,
		PreRunE:      opts.PreRunE,
		Args:         cobra.ExactArgs(1),
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *renderOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonConnOpts.AddFlags(cmd)

	cmd.Flags().BoolVar(&o.reloadableOnly, "reloadable-only", false, "Only render reloadable templates")
}

func (o *renderOpts) Run(cmd *cobra.Command, args []string) error {
	c, conn, err := o.Client()
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			cmd.PrintErrf("error closing connection: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	r, err := c.Render(ctx, &svctl.RenderOpts{
		Id:             o.id,
		ReloadableOnly: o.reloadableOnly,
	})
	if err != nil {
		return fmt.Errorf("error calling function Render: %v", err)
	}

	cmd.Printf("Server templates rendered: %v\n", r.GetId())
	return nil
}

func init() {
	rootCmd.AddCommand(renderCmd())
}

