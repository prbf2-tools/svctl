package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type unregisterOpts struct {
	*grpcClientCmdOpts
}

func newUnregisterOpts() *unregisterOpts {
	return &unregisterOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func unregisterCmd() *cobra.Command {
	opts := newUnregisterOpts()

	cmd := &cobra.Command{
		Use:          "unregister <server-id>",
		Short:        "Unregisters the server",
		Long:         `Send an unregister signal to daemon to unregister the server`,
		SilenceUsage: true,
		PreRunE:      opts.PreRunE,
		Args:         cobra.ExactArgs(1),
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *unregisterOpts) Run(cmd *cobra.Command, args []string) error {
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

	r, err := c.Unregister(ctx, &svctl.ServerOpts{Id: o.id})
	if err != nil {
		return fmt.Errorf("error calling function Unregister: %v", err)
	}

	cmd.Printf("Server unregistered: %v\n", r.GetId())
	return nil
}

func init() {
	rootCmd.AddCommand(unregisterCmd())
}
