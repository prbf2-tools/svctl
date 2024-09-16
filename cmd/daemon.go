package cmd

import (
	"log"
	"net"

	"github.com/sboon-gg/svctl/internal/api"
	"github.com/sboon-gg/svctl/internal/daemon"
	"github.com/sboon-gg/svctl/svctl"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	defaultDaemonHost = "127.0.0.1"
	defaultDaemonPort = "50051"
)

type daemonOpts struct {
	host string
	port string
}

func newDaemonOpts() *daemonOpts {
	return &daemonOpts{
		host: defaultDaemonHost,
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
	cmd.Flags().StringVar(&o.host, "host", o.host, "Host to listen on")
	cmd.Flags().StringVar(&o.port, "port", o.port, "Port to listen on")
}

func (o *daemonOpts) Run(cmd *cobra.Command, args []string) error {
	d, err := daemon.Recover()
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

func (o *daemonOpts) address() string {
	return net.JoinHostPort(o.host, o.port)
}

func init() {
	rootCmd.AddCommand(daemonCmd())
}
