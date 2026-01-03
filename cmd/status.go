package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sboon-gg/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type statusOpts struct {
	*serverOpts
	*daemonConnOpts
	json bool
}

func newStatusOpts() *statusOpts {
	return &statusOpts{
		serverOpts:     newServerOpts(),
		daemonConnOpts: newDaemonConnOpts(),
	}
}

func statusCmd() *cobra.Command {
	opts := newStatusOpts()

	cmd := &cobra.Command{
		Use:          "status",
		Short:        "Server status",
		Long:         `Display the status of the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *statusOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonConnOpts.AddFlags(cmd)

	cmd.Flags().BoolVarP(&o.json, "json", "j", false, "Output as JSON")
}

func (o *statusOpts) Run(cmd *cobra.Command, args []string) error {
	c, conn, err := o.Client()
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at %s: %v", o.address(), err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	id, err := o.ID()
	if err != nil {
		return err
	}

	status, err := c.Status(ctx, &svctl.ServerOpts{Id: id})
	if err != nil {
		return fmt.Errorf("error calling function Status: %v", err)
	}

	if o.json {
		content, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %v", err)
		}

		cmd.Println(string(content))
		return nil
	}

	cmd.Printf("Server Info:\n")
	cmd.Printf("  Path: %s\n", status.Path)
	cmd.Printf("  Settings Path: %s\n", status.SettingsPath)
	cmd.Printf("  Desired State: %v\n", status.DesiredState)
	cmd.Printf("  Current State: %v\n", status.CurrentState)
	if status.GameStatus != nil {
		cmd.Printf("  Game Info:\n")
		cmd.Printf("     Hostname: %s\n", status.GameStatus.Hostname)
		cmd.Printf("     Port: %s\n", status.GameStatus.Port)
		cmd.Printf("     Game Mode: %s\n", status.GameStatus.Gamemode)
		cmd.Printf("     Map Name: %s\n", status.GameStatus.MapName)
		cmd.Printf("     Map Size: %d\n", status.GameStatus.MapSize)
		cmd.Printf("     Number of Players: %d\n", status.GameStatus.NumPlayers)
		cmd.Printf("     Maximum Players: %d\n", status.GameStatus.MaxPlayers)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(statusCmd())
}
