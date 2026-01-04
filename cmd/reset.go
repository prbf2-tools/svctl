package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type resetOpts struct {
	*serverOpts
	*daemonConnOpts
}

func newResetOpts() *resetOpts {
	return &resetOpts{
		serverOpts:     newServerOpts(),
		daemonConnOpts: newDaemonConnOpts(),
	}
}

func resetCmd() *cobra.Command {
	opts := newResetOpts()

	cmd := &cobra.Command{
		Use:          "reset",
		Short:        "Resets the server",
		Long:         `Send a reset signal to daemon to reset the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *resetOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonConnOpts.AddFlags(cmd)
}

func (o *resetOpts) Run(cmd *cobra.Command, args []string) error {
	c, conn, err := o.Client()
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %v", o.address(), err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	id, err := o.ID()
	if err != nil {
		return err
	}

	r, err := c.Reset(ctx, &svctl.ServerOpts{Id: id})
	if err != nil {
		return fmt.Errorf("error calling function Reset: %v", err)
	}

	cmd.Printf("Reseted completed: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(resetCmd())
}
