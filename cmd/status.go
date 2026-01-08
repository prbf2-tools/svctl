package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
)

type statusOpts struct {
	*grpcClientCmdOpts
	json bool
}

func newStatusOpts() *statusOpts {
	return &statusOpts{
		grpcClientCmdOpts: newGrpcClientCmdOpts(),
	}
}

func statusCmd() *cobra.Command {
	opts := newStatusOpts()

	cmd := &cobra.Command{
		Use:          "status <server-id>",
		Short:        "Server status",
		Long:         `Display the status of the server`,
		SilenceUsage: true,
		PreRunE:      opts.PreRunE,
		Args:         cobra.ExactArgs(1),
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *statusOpts) AddFlags(cmd *cobra.Command) {
	o.grpcClientCmdOpts.AddFlags(cmd)

	cmd.Flags().BoolVarP(&o.json, "json", "j", false, "Output as JSON")
}

func (o *statusOpts) Run(cmd *cobra.Command, args []string) error {
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

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	status, err := c.Status(ctx, &svctl.ServerOpts{Id: o.id})
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
	cmd.Printf("  Path: %s\n", status.Location)
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
