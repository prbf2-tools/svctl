package cmd

import (
	"os"
	"path/filepath"

	"github.com/sboon-gg/svctl/internal/server"
	"github.com/spf13/cobra"
)

const (
	defaultSettingsPath = ".svctl"
)

type serverOpts struct {
	serverPath string
	svctlPath  string
}

func newServerOpts() *serverOpts {
	return &serverOpts{
		serverPath: ".",
		svctlPath:  defaultSettingsPath,
	}
}

func (opts *serverOpts) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&opts.serverPath, "path", "p", opts.serverPath, "Path to server directory")
	cmd.Flags().StringVar(&opts.svctlPath, "settings", opts.svctlPath, "Path to settings directory")
}

func (opts *serverOpts) Path() (string, error) {
	if filepath.IsAbs(opts.serverPath) {
		return opts.serverPath, nil
	}

	return concatWithWorkingDir(opts.serverPath)
}

func (opts *serverOpts) SvctlPath() (string, error) {
	if filepath.IsAbs(opts.svctlPath) {
		return opts.svctlPath, nil
	}

	return concatWithWorkingDir(opts.svctlPath)
}

func concatWithWorkingDir(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, path), nil
}

func (opts *serverOpts) Server() (*server.Server, error) {
	path, err := opts.Path()
	if err != nil {
		return nil, err
	}

	svctlPath, err := opts.SvctlPath()
	if err != nil {
		return nil, err
	}

	return server.Open(path, svctlPath)
}
