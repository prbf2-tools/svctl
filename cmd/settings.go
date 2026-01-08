package cmd

import (
	"github.com/spf13/cobra"
)

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Manage settings",
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}

