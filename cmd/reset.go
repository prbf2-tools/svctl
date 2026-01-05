package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type resetOpts struct {
	*grpcClientCmdOpts
}

func newResetOpts() *resetOpts {
	return &resetOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func resetCmd() *cobra.Command {
	opts := newResetOpts()

	cmd := &cobra.Command{
		Use:          "reset <server-id>",
		Short:        "Resets the server",
		Long:         `Send a reset signal to daemon to reset the server`,
		SilenceUsage: true,
		PreRunE:      opts.PreRunE,
		Args:         cobra.ExactArgs(1),
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *resetOpts) Run(cmd *cobra.Command, args []string) error {
	c, conn, err := o.Client()
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %v", o.address(), err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			cmd.PrintErrf("error closing connection: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	r, err := c.Reset(ctx, &svctl.ServerOpts{Id: o.id})
	if err != nil {
		return fmt.Errorf("error calling function Reset: %v", err)
	}

	cmd.Printf("Reseted completed: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(resetCmd())
}
