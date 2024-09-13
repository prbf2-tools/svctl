package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/sboon-gg/svctl/svctl"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type startOpts struct {
	*serverOpts
	*daemonOpts
}

func newStartOpts() *startOpts {
	return &startOpts{
		serverOpts: newServerOpts(),
		daemonOpts: newDaemonOpts(),
	}
}

func startCmd() *cobra.Command {
	opts := newStartOpts()

	cmd := &cobra.Command{
		Use:          "start",
		Short:        "Starts the server",
		Long:         `Starts the server`,
		SilenceUsage: true,
		RunE:         opts.Run,
	}

	opts.AddFlags(cmd)

	return cmd
}

func (o *startOpts) AddFlags(cmd *cobra.Command) {
	o.serverOpts.AddFlags(cmd)
	o.daemonOpts.AddFlags(cmd)
}

func (o *startOpts) Run(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", o.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server at localhost:%s: %v", err, o.port)
	}
	defer conn.Close()
	c := svctl.NewServersClient(conn)

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Second)
	defer cancel()

	path, err := o.Path()
	if err != nil {
		return err
	}

	r, err := c.Start(ctx, &svctl.ServerOpts{Path: path})
	if err != nil {
		return fmt.Errorf("error calling function Start: %v", err)
	}

	cmd.Printf("Server started: %v\n", r.GetStatus().String())
	return nil
}

// func (o *startOpts) Run(cmd *cobra.Command, args []string) error {
// 	sv, err := o.Server()
// 	if err != nil {
// 		return err
// 	}
//
// 	cache, err := sv.Cache()
// 	if err != nil {
// 		return err
// 	}
//
// 	var proc *prbf2.Proc
// 	if cache.PID > 0 {
// 		proc, err = prbf2.OpenProc(cache.PID)
// 		if err != nil {
// 			return err
// 		}
//
// 		err = proc.HealthCheck()
// 		if err != nil {
// 			proc, err = prbf2.NewProc(sv.Path)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			return fmt.Errorf("server is already running")
// 		}
// 	} else {
// 		proc, err = prbf2.NewProc(sv.Path)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	err = proc.HealthCheck()
// 	if err != nil {
// 		return err
// 	}
//
// 	cache.PID = proc.Pid
//
// 	err = sv.WriteCache(cache)
// 	if err != nil {
// 		return err
// 	}
//
// 	cmd.Println("Server started")
// 	return nil
// }

func init() {
	rootCmd.AddCommand(startCmd())
}
