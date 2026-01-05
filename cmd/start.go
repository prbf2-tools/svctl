package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type startOpts struct {
	*grpcClientCmdOpts
}

func newStartOpts() *startOpts {
	return &startOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func startCmd() *cobra.Command {
	opts := newStartOpts()

	cmd := &cobra.Command{
		Use:          "start <server-id>",
		Short:        "Starts the server",
		Long:         `Send a start signal to daemon to start the server`,
		SilenceUsage: true,
		PreRunE:      opts.PreRunE,
		Args:         cobra.ExactArgs(1),
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

	r, err := c.Start(ctx, &svctl.ServerOpts{Id: o.id})
	if err != nil {
		return fmt.Errorf("error calling function Start: %v", err)
	}

	cmd.Printf("Server started: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(startCmd())
}
