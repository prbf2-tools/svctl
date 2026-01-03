package cmd

import (
	"net"
	"os"
	"path/filepath"

	"github.com/sboon-gg/svctl/internal/server"
	"github.com/sboon-gg/svctl/svctl/v1"
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
	serverPath    string
	containerName string
	settingsPath  string
}

func newServerOpts() *serverOpts {
	return &serverOpts{
		settingsPath: defaultSettingsPath,
	}
}

func (opts *serverOpts) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.serverPath, "path", "p", "", "Path to server directory")
	cmd.Flags().StringVarP(&opts.containerName, "container", "c", "", "Name of the Docker container")
	cmd.Flags().StringVar(&opts.settingsPath, "settings", opts.settingsPath, "Path to settings directory")

	cmd.MarkFlagDirname("path")
	cmd.MarkFlagDirname("settings")
	cmd.MarkFlagsMutuallyExclusive("path", "container")
	cmd.MarkFlagsOneRequired("path", "container")
}

func (opts *serverOpts) ID() (string, error) {
	if opts.containerName != "" {
		return opts.containerName, nil
	}

	if filepath.IsAbs(opts.serverPath) {
		return opts.serverPath, nil
	}

	return concatWithWorkingDir(opts.serverPath)
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

func (opts *serverOpts) Server() (*server.Server, error) {
	id, err := opts.ID()
	if err != nil {
		return nil, err
	}

	svctlPath, err := opts.SettingsPath()
	if err != nil {
		return nil, err
	}

	if opts.containerName != "" {
		return server.OpenDocker(id, svctlPath)
	}

	return server.OpenLocal(id, svctlPath)
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
