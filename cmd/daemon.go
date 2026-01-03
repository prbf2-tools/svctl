package cmd

import (
	"log"
	"net"

	"github.com/sboon-gg/svctl/internal/api"
	"github.com/sboon-gg/svctl/internal/daemon"
	"github.com/sboon-gg/svctl/svctl/v1"
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
		log.Fatalf("failed to listen on address %s: %v", o.address(), err)
	}

	s := grpc.NewServer()
	svctl.RegisterServersServer(s, api.NewDaemonServer(d))
	log.Printf("gRPC server listening at %v", lis.Addr())

	go func() {
		<-cmd.Context().Done()
		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(daemonCmd())
}
