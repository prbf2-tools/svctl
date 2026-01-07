package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type restartOpts struct {
	*grpcClientCmdOpts
}

func newRestartOpts() *restartOpts {
	return &restartOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func restartCmd() *cobra.Command {
	opts := newRestartOpts()

	cmd := &cobra.Command{
		Use:          "restart <server-id>",
		Short:        "Restarts the server",
		Long:         `Send stop and start signals to daemon to restart the server`,
		SilenceUsage: true,
		PreRunE:      opts.PreRunE,
		Args:         cobra.ExactArgs(1),
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *restartOpts) AddFlags(cmd *cobra.Command) {
	o.grpcClientCmdOpts.AddFlags(cmd)
}

func (o *restartOpts) Run(cmd *cobra.Command, args []string) error {
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

	serverOpts := &svctl.ServerOpts{Id: o.id}

	// Stop the server first
	cmd.Printf("Stopping server %s...\n", o.id)
	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
	defer cancel()

	_, err = c.Stop(ctx, serverOpts)
	if err != nil {
		return fmt.Errorf("error calling function Stop: %v", err)
	}

	// Start the server
	cmd.Printf("Starting server %s...\n", o.id)
	ctx, cancel = context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	r, err := c.Start(ctx, serverOpts)
	if err != nil {
		return fmt.Errorf("error calling function Start: %v", err)
	}

	cmd.Printf("Server restarted: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(restartCmd())
}

