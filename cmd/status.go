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

type statusOpts struct {
	*serverOpts
	*daemonOpts
}

func newStatusOpts() *statusOpts {
	return &statusOpts{
		serverOpts: newServerOpts(),
		daemonOpts: newDaemonOpts(),
	}
}

func statusCmd() *cobra.Command {
	opts := newStatusOpts()

	cmd := &cobra.Command{
		Use:          "status",
		Short:        "Server status",
		Long:         `Display the status of the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *statusOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonOpts.AddFlags(cmd)
}

func (o *statusOpts) Run(cmd *cobra.Command, args []string) error {
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

	r, err := c.Status(ctx, &svctl.ServerOpts{Path: path})
	if err != nil {
		return fmt.Errorf("error calling function Status: %v", err)
	}

	cmd.Printf("Server info: %v\n", r)
	return nil
}

func init() {
	rootCmd.AddCommand(statusCmd())
}
