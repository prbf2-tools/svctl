package cmd

import (
	"github.com/spf13/cobra"
)

type updateOpts struct {
	path string
}

func updateCmd() *cobra.Command {
	opts := &updateOpts{}

	cmd := &cobra.Command{
		Use: "update",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.path = args[0]
			return opts.Run(cmd)
		},
	}

	return cmd
}

func (o *updateOpts) Run(cmd *cobra.Command) error {
	// cache := prbf2update.NewCache(filepath.Join(o.path, ".update-cache"))
	//
	// u := prbf2update.New(o.path, cache)
	//
	// result, err := u.Update()
	// if err != nil {
	// 	return err
	// }
	//
	// cmd.Println("Updated from", result.OldVersion, "to", result.NewVersion)
	// TODO: implement this as FSM event
	cmd.Println("Not implemented")

	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd())
}
