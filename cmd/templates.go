package cmd

import (
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
