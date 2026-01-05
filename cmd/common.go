package cmd

import (
	"net"
	"os"
	"path/filepath"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultSettingsPath = ".svctl"

	defaultDaemonHost = "127.0.0.1"
	defaultDaemonPort = "50051"
)

type serverOpts struct {
	id           string
	settingsPath string
}

func newServerOpts() *serverOpts {
	return &serverOpts{
		settingsPath: defaultSettingsPath,
	}
}

func (opts *serverOpts) PreRunE(cmd *cobra.Command, args []string) error {
	opts.id = args[0]

	return nil
}

func (opts *serverOpts) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&opts.settingsPath, "settings", opts.settingsPath, "Path to settings directory")

	_ = cmd.MarkFlagDirname("settings")
}

func (opts *serverOpts) SettingsPath() (string, error) {
	if filepath.IsAbs(opts.settingsPath) {
		return opts.settingsPath, nil
	}

	return concatWithWorkingDir(opts.settingsPath)
}

func concatWithWorkingDir(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, path), nil
}

type daemonConnOpts struct {
	host string
	port string
}

func newDaemonConnOpts() *daemonConnOpts {
	return &daemonConnOpts{
		host: defaultDaemonHost,
		port: defaultDaemonPort,
	}
}

func (o *daemonConnOpts) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.host, "host", o.host, "Daemon host")
	cmd.Flags().StringVar(&o.port, "port", o.port, "Daemon port")
}

func (o *daemonConnOpts) address() string {
	return net.JoinHostPort(o.host, o.port)
}

func (o *daemonConnOpts) Client() (svctl.ServersClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(o.address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := svctl.NewServersClient(conn)

	return client, conn, nil
}

type grpcClientCmdOpts struct {
	*daemonConnOpts
	*serverOpts
}

func newGrpcClientCmdOpts() *grpcClientCmdOpts {
	return &grpcClientCmdOpts{
		daemonConnOpts: newDaemonConnOpts(),
		serverOpts:     newServerOpts(),
	}
}

func (o *grpcClientCmdOpts) AddFlags(cmd *cobra.Command) {
	o.daemonConnOpts.AddFlags(cmd)
	o.serverOpts.AddFlags(cmd)
}
