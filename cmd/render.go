package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/cobra"
)

type renderOpts struct {
	*serverOpts
	dryRun bool
	watch  bool
	values []string
}

func newRenderOpts() *renderOpts {
	return &renderOpts{
		serverOpts: newServerOpts(),
	}
}

func init() {
	rootCmd.AddCommand(renderCmd())
}

func renderCmd() *cobra.Command {
	opts := newRenderOpts()

	cmd := &cobra.Command{
		Use:          "render",
		Short:        "Render templates",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(cmd)
		},
	}

	opts.serverOpts.AddFlags(cmd)

	// cmd.Flags().BoolVar(&opts.defaults, "defaults", true, "Use default values")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Print out rendered files")
	cmd.Flags().BoolVar(&opts.watch, "watch", false, "Watch all values and config files")
	cmd.Flags().StringSliceVar(&opts.values, "values", []string{}, "Additional values files - relative to execution working directory")

	return cmd
}

func (opts *renderOpts) Run(cmd *cobra.Command) error {
	watchedFiles, err := opts.render()
	if err != nil {
		return err
	}

	if !opts.watch {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for _, file := range watchedFiles {
		err = watcher.Add(file)
		if err != nil {
			return err
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Has(fsnotify.Write | fsnotify.Create) {
				watchedFiles, err = opts.render()
				if err != nil {
					log.Printf("Error rendering files: %s", err)
				} else {
					for _, file := range watchedFiles {
						_ = watcher.Add(file)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println("error:", err)
		}
	}
}

func (opts *renderOpts) render() ([]string, error) {
	var files []string

	si, err := opts.Server()
	if err != nil {
		return files, errors.New("Script has not been initialized, run `init` first.")
	}

	if opts.dryRun {
		outputs, err := si.DryRender()
		if err != nil {
			return nil, err
		}
		for _, out := range outputs {
			fmt.Printf("File: %s\n---\n%s", out.Destination, string(out.Content))
		}
	} else {
		err := si.Render()
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}
