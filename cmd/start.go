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

type startOpts struct {
	*serverOpts
	*daemonOpts
}

func newStartOpts() *startOpts {
	return &startOpts{
		serverOpts: newServerOpts(),
		daemonOpts: newDaemonOpts(),
	}
}

func startCmd() *cobra.Command {
	opts := newStartOpts()

	cmd := &cobra.Command{
		Use:          "start",
		Short:        "Starts the server",
		Long:         `Starts the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *startOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonOpts.AddFlags(cmd)
}

func (o *startOpts) Run(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial(o.daemonOpts.address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	r, err := c.Start(ctx, &svctl.ServerOpts{Path: path})
	if err != nil {
		return fmt.Errorf("error calling function Start: %v", err)
	}

	cmd.Printf("Server started: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(startCmd())
}
