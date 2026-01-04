package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type startOpts struct {
	*serverOpts
	*daemonConnOpts
}

func newStartOpts() *startOpts {
	return &startOpts{
		serverOpts:     newServerOpts(),
		daemonConnOpts: newDaemonConnOpts(),
	}
}

func startCmd() *cobra.Command {
	opts := newStartOpts()

	cmd := &cobra.Command{
		Use:          "start",
		Short:        "Starts the server",
		Long:         `Send a start signal to daemon to start the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *startOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonConnOpts.AddFlags(cmd)
}

func (o *startOpts) Run(cmd *cobra.Command, args []string) error {
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

	r, err := c.Start(ctx, &svctl.ServerOpts{Id: id})
	if err != nil {
		return fmt.Errorf("error calling function Start: %v", err)
	}

	cmd.Printf("Server started: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(startCmd())
}
