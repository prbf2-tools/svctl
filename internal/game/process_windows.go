//go:build windows

package game

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

const (
	processExe = "prbf2_w32ded.exe"
)

func (s *Server) processExe() string {
	return filepath.Join(s.Path, processExe)
}

func (s *Server) startProcess() (*os.Process, error) {
	allArgs := append([]string{processExe}, commonProcessArgs...)

	proc, err := os.StartProcess(processExe, allArgs, &os.ProcAttr{
		Dir: s.Path,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	})
	if err != nil {
		return nil, err
	}

	setHighPriority(proc.Pid)

	return proc, nil
}

const processAllAccess = windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0xffff

// setHighPriority is generally used by PR sysadmins
// TODO: make this configurable
func setHighPriority(pid int) error {
	handle, err := windows.OpenProcess(processAllAccess, false, uint32(pid))
	if err != nil {
		return err
	}

	err = windows.SetPriorityClass(handle, windows.HIGH_PRIORITY_CLASS)
	if err != nil {
		return err
	}

	return windows.CloseHandle(handle)
}
