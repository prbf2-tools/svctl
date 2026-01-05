/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type registerOpts struct {
	*grpcClientCmdOpts
	localPath     string
	containerName string
	serviceName   string
}

func newRegisterOpts() *registerOpts {
	return &registerOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func registerCmd() *cobra.Command {
	opts := newRegisterOpts()

	cmd := &cobra.Command{
		Use:     "register <server-id>",
		Short:   "register a new server",
		PreRunE: opts.PreRunE,
		Args:    cobra.ExactArgs(1),
		RunE:    opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (opts *registerOpts) AddFlags(cmd *cobra.Command) {
	opts.grpcClientCmdOpts.AddFlags(cmd)

	cmd.Flags().StringVarP(&opts.localPath, "path", "p", "", "Path to server directory")
	cmd.Flags().StringVarP(&opts.containerName, "container", "c", "", "Name of the Docker container")
	cmd.Flags().StringVarP(&opts.serviceName, "service", "s", "", "Name of the systemd service")

	_ = cmd.MarkFlagDirname("path")
	cmd.MarkFlagsMutuallyExclusive("path", "container", "service")
	cmd.MarkFlagsOneRequired("path", "container", "service")
}

func (opts *registerOpts) Run(cmd *cobra.Command, args []string) error {
	c, conn, err := opts.Client()
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %v", opts.address(), err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			cmd.PrintErrf("error closing connection: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	location, err := opts.Location()
	if err != nil {
		return err
	}

	settingsPath, err := opts.SettingsPath()
	if err != nil {
		return err
	}

	r, err := c.Register(ctx, &svctl.RegisterServerOpts{
		Id:           opts.id,
		SettingsPath: settingsPath,
		Location:     location,
		Type:         opts.ServerType(),
	})
	if err != nil {
		return fmt.Errorf("error calling function Register: %v", err)
	}

	cmd.Printf("Server status: %v\n", r.GetStatus().String())
	return nil
}

func (opts *registerOpts) Location() (string, error) {
	switch {
	case opts.serviceName != "":
		return opts.serviceName, nil
	case opts.containerName != "":
		return opts.containerName, nil
	case opts.localPath != "":
		return concatWithWorkingDir(opts.localPath)
	default:
		return "", fmt.Errorf("no server identifier provided")
	}
}

func (opts *registerOpts) ServerType() svctl.ServerType {
	switch {
	case opts.serviceName != "":
		return svctl.ServerType_SERVER_TYPE_SYSTEMD
	case opts.containerName != "":
		return svctl.ServerType_SERVER_TYPE_DOCKER
	case opts.localPath != "":
		return svctl.ServerType_SERVER_TYPE_LOCAL
	default:
		return svctl.ServerType_SERVER_TYPE_UNSPECIFIED
	}
}

func init() {
	rootCmd.AddCommand(registerCmd())
}
