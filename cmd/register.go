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
	*serverOpts
	*daemonConnOpts
}

func newRegisterOpts() *registerOpts {
	return &registerOpts{
		serverOpts:     newServerOpts(),
		daemonConnOpts: newDaemonConnOpts(),
	}
}

func registerCmd() *cobra.Command {
	opts := newRegisterOpts()

	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a new server",
		RunE:  opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *registerOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonConnOpts.AddFlags(cmd)
}

func (o *registerOpts) Run(cmd *cobra.Command, args []string) error {
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

	id, err := o.ID()
	if err != nil {
		return err
	}

	settingsPath, err := o.SettingsPath()
	if err != nil {
		return err
	}

	r, err := c.Register(ctx, &svctl.RegisterServerOpts{
		Id:           id,
		SettingsPath: settingsPath,
		Type:         o.serverType(),
	})
	if err != nil {
		return fmt.Errorf("error calling function Register: %v", err)
	}

	cmd.Printf("Server status: %v\n", r.GetStatus().String())
	return nil
}

func init() {
	rootCmd.AddCommand(registerCmd())
}
