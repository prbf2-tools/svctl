//go:build linux

package game

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

const (
	processExe = "prbf2_l64ded"
	binaryDir  = "bin/amd-64"
)

func (s *Server) processExe() string {
	return filepath.Join(s.Path, binaryDir, processExe)
}

func (s *Server) startProcess() (*os.Process, error) {
	binDir := filepath.Join(s.Path, binaryDir)
	fullExe := filepath.Join(binDir, processExe)

	err := makeFileExecutable(fullExe)
	if err != nil {
		return nil, err
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("LD_LIBRARY_PATH=%s", binDir))

	allArgs := append([]string{fullExe}, commonProcessArgs...)

	return os.StartProcess(fullExe, allArgs, &os.ProcAttr{
		Dir: s.Path,
		Env: env,
		Sys: &syscall.SysProcAttr{
			Setpgid: true,
		},
	})
}

func processHealth(_ int) (bool, error) {
	return true, nil
}
