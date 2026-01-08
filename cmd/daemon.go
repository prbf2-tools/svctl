package cmd

import (
	"fmt"
	"net"

	"github.com/prbf2-tools/svctl/internal/api"
	"github.com/prbf2-tools/svctl/internal/daemon"
	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type daemonOpts struct {
	*daemonConnOpts
	configFile string
}

func newDaemonOpts() *daemonOpts {
	return &daemonOpts{
		daemonConnOpts: newDaemonConnOpts(),
	}
}

func daemonCmd() *cobra.Command {
	opts := newDaemonOpts()

	cmd := &cobra.Command{
		Use:  "daemon",
		RunE: opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *daemonOpts) AddFlags(cmd *cobra.Command) {
	o.daemonConnOpts.AddFlags(cmd)
	cmd.Flags().StringVar(&o.configFile, "config", o.configFile, "Path to the daemon config file")
}

func (o *daemonOpts) Run(cmd *cobra.Command, args []string) error {
	d, err := daemon.Recover(o.configFile)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", o.address())
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %w", o.address(), err)
	}

	s := grpc.NewServer()
	svctl.RegisterServersServer(s, api.NewDaemonServer(d))
	cmd.Printf("gRPC server listening at %v", lis.Addr())

	go func() {
		<-cmd.Context().Done()
		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(daemonCmd())
}
