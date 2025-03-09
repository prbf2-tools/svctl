package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/sboon-gg/svctl/svctl"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type resetOpts struct {
	*serverOpts
	*daemonOpts
}

func newResetOpts() *resetOpts {
	return &resetOpts{
		serverOpts: newServerOpts(),
		daemonOpts: newDaemonOpts(),
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
	o.daemonOpts.AddFlags(cmd)
}

func (o *resetOpts) Run(cmd *cobra.Command, args []string) error {
	conn, err := grpc.NewClient(o.daemonOpts.address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %v", o.daemonOpts.address(), err)
	}
	defer conn.Close()
	c := svctl.NewServersClient(conn)

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	path, err := o.Path()
	if err != nil {
		return err
	}

	r, err := c.Reset(ctx, &svctl.ServerOpts{Path: path})
	if err != nil {
		return fmt.Errorf("error calling function Reset: %v", err)
	}

	cmd.Printf("Reseted completed: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(resetCmd())
}
