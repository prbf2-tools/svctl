package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/sboon-gg/svctl/internal/api"
	"github.com/sboon-gg/svctl/internal/daemon"
	"github.com/sboon-gg/svctl/svctl"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	defaultDaemonPort = "50051"
)

type daemonOpts struct {
	port string
}

func newDaemonOpts() *daemonOpts {
	return &daemonOpts{
		port: defaultDaemonPort,
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
	cmd.Flags().StringVar(&o.port, "port", o.port, "Port to listen on")
}

func (o *daemonOpts) Run(cmd *cobra.Command, args []string) error {
	d, err := daemon.Recover()
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", defaultDaemonPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", defaultDaemonPort, err)
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
