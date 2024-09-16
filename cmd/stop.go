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

type stopOpts struct {
	*serverOpts
	*daemonOpts
}

func newStopOpts() *stopOpts {
	return &stopOpts{
		serverOpts: newServerOpts(),
		daemonOpts: newDaemonOpts(),
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
	o.daemonOpts.AddFlags(cmd)
}

func (o *stopOpts) Run(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial(o.daemonOpts.address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %v", o.daemonOpts.address(), err)
	}
	defer conn.Close()
	c := svctl.NewServersClient(conn)

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
	defer cancel()

	path, err := o.Path()
	if err != nil {
		return err
	}

	r, err := c.Stop(ctx, &svctl.ServerOpts{Path: path})
	if err != nil {
		return fmt.Errorf("error calling function Stop: %v", err)
	}

	cmd.Printf("Server status: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(stopCmd())
}
