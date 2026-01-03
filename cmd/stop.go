package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/sboon-gg/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type stopOpts struct {
	*serverOpts
	*daemonConnOpts
}

func newStopOpts() *stopOpts {
	return &stopOpts{
		serverOpts:     newServerOpts(),
		daemonConnOpts: newDaemonConnOpts(),
	}
}

func stopCmd() *cobra.Command {
	opts := newStopOpts()

	cmd := &cobra.Command{
		Use:          "stop",
		Short:        "Stops the server",
		Long:         `Send a stop signal to daemon to stop the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *stopOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonConnOpts.AddFlags(cmd)
}

func (o *stopOpts) Run(cmd *cobra.Command, args []string) error {
	c, conn, err := o.Client()
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %v", o.address(), err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
	defer cancel()

	id, err := o.ID()
	if err != nil {
		return err
	}

	r, err := c.Stop(ctx, &svctl.ServerOpts{Id: id})
	if err != nil {
		return fmt.Errorf("error calling function Stop: %v", err)
	}

	cmd.Printf("Server status: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(stopCmd())
}
