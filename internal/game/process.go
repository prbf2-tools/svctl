package game

import (
	"errors"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v3/process"
)

var (
	ErrProcessAlreadyRunning = errors.New("process already running")
)

var commonProcessArgs = []string{
	"+modPath", "mods/pr",
	"+noStatusMonitor", "1",
	"+multi", "1",
	"+dedicated", "1",
}

func (s *Server) Start() error {
	if s.IsRunning() {
		return ErrProcessAlreadyRunning
	}

	proc, err := s.startProcess()
	if err != nil {
		return err
	}

	return s.storeProcessPID(proc.Pid)
}

func (s *Server) Stop() error {
	if !s.IsRunning() {
		return nil
	}

	proc, err := process.NewProcess(int32(*s.processPID))
	if err != nil {
		return err
	}

	err = proc.Kill()
	if err != nil {
		return err
	}

	return s.clearProcessPID()
}

func (s *Server) IsRunning() bool {
	if s.processPID == nil {
		return false
	}

	// On Windows we need to check if the process isn't hanging on an error dialog.
	if runtime.GOOS == "windows" {
		health, err := processHealth(*s.processPID)
		if err == nil && !health {
			s.clearProcessPID()
			return false
		}
	}

	isRunning, err := s.isRunning()
	if err != nil || !isRunning {
		s.clearProcessPID()
		return false
	}

	return isRunning
}

func (s *Server) isRunning() (bool, error) {
	if s.processPID == nil {
		return false, nil
	}

	proc, err := process.NewProcess(int32(*s.processPID))
	if err != nil {
		return false, err
	}

	if s.processExe() != processExe {
		return false, nil
	}

	isRunning, err := proc.IsRunning()
	if err != nil {
		return false, err
	}

	return isRunning, nil
}

func (s *Server) clearProcessPID() error {
	s.processPID = nil

	return s.RemoveFile(pidFile)
}

func (s *Server) storeProcessPID(pid int) error {
	err := s.WriteFile(pidFile, []byte(strconv.Itoa(pid)))
	if err != nil {
		return err
	}

	s.processPID = &pid
	return nil
}
