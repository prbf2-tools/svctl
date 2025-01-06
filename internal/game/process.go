package game

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

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

	pid := proc.Pid

	err = proc.Release()
	if err != nil {
		return err
	}

	return s.storeProcessPID(pid)
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
		err := s.retrieveProcessPID()
		if err != nil {
			fmt.Println(err)
		}
		if s.processPID == nil {
			return false
		}
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
		_ = s.clearProcessPID()
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

	if s.processExe() != filepath.Join(s.Path, binaryDir, processExe) {
		return false, nil
	}

	status, err := proc.Status()
	if err != nil || slices.Contains(status, "zombie") {
		return false, err
	}

	isRunning, err := proc.IsRunning()
	if err != nil {
		return false, err
	}

	return isRunning, nil
}

func (s *Server) clearProcessPID() error {
	s.processPID = nil

	return os.Remove(filepath.Join(s.Path, pidFile))
}

func (s *Server) retrieveProcessPID() error {
	content, err := os.ReadFile(filepath.Join(s.Path, pidFile))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return err
	}

	s.processPID = &pid
	return nil
}

func (s *Server) storeProcessPID(pid int) error {
	err := os.WriteFile(filepath.Join(s.Path, pidFile), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return err
	}

	s.processPID = &pid
	return nil
}
