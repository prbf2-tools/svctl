package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type stopOpts struct {
	*grpcClientCmdOpts
}

func newStopOpts() *stopOpts {
	return &stopOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func stopCmd() *cobra.Command {
	opts := newStopOpts()

	cmd := &cobra.Command{
		Use:          "stop <server-id>",
		Short:        "Stops the server",
		Long:         `Send a stop signal to daemon to stop the server`,
		SilenceUsage: true,
		PreRunE:      opts.PreRunE,
		Args:         cobra.ExactArgs(1),
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *stopOpts) Run(cmd *cobra.Command, args []string) error {
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

	ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
	defer cancel()

	r, err := c.Stop(ctx, &svctl.ServerOpts{Id: o.id})
	if err != nil {
		return fmt.Errorf("error calling function Stop: %v", err)
	}

	cmd.Printf("Server status: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(stopCmd())
}
