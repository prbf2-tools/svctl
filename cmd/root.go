package cmd

import (
	"context"
	"log/slog"
	"os/signal"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "svctl",
	Short: "PRBF2 Server controller CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		if debug {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
		return nil
	},
}

func ExecuteNoExit() error {
	ctx, cancel := signal.NotifyContext(context.Background(), shutdownSignals...)
	defer cancel()

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	return rootCmd.ExecuteContext(ctx)
}
