//go:build linux

package game

import (
	"context"
	"io"
	"os/exec"
	"path/filepath"
)

const (
	exe = "prserverupdater-linux64"
)

func (s *Server) update(ctx context.Context, outW io.Writer, inR io.Reader, errW io.Writer) error {
	binPath := filepath.Join(s.Path, updaterPath)
	exePath := filepath.Join(binPath, exe)

	err := makeFileExecutable(exePath)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "./"+exe)
	cmd.Dir = binPath

	cmd.Stdout = outW
	cmd.Stdin = inR
	cmd.Stderr = errW

	return cmd.Run()
}
