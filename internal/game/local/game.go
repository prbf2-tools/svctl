package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sboon-gg/svctl/internal/game"
)

const (
	pidFile     = "prbf2.pid"
	updaterPath = "mods/pr/bin"
)

type runningProcess interface {
	IsRunning() (bool, error)
	Kill() error
	PID() int
}

var _ game.GameServer = &Server{}

type Server struct {
	path    string
	process runningProcess
}

func Open(path string) (*Server, error) {
	s := &Server{
		path: path,
	}

	err := s.recoverProcess()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) ID() string {
	return s.path
}

func (s *Server) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.path, path))
}

func (s *Server) WriteFile(path string, data []byte) error {
	fullPath := filepath.Join(s.path, path)

	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		// Ignore error, the write will fail if the directory doesn't exist.
		_ = os.MkdirAll(filepath.Dir(fullPath), 0755)
	}

	return os.WriteFile(fullPath, data, 0644)
}

func (s *Server) WriteFileFromReader(path string, r io.Reader, _ int64) error {
	fullPath := filepath.Join(s.path, path)

	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		// Ignore error, the write will fail if the directory doesn't exist.
		_ = os.MkdirAll(filepath.Dir(fullPath), 0755)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}

func (s *Server) Update(ctx context.Context, outW io.Writer, inR io.Reader, errW io.Writer) error {
	return s.update(ctx, outW, inR, errW)
}

func (s *Server) Start() error {
	isRunning, err := s.IsRunning()
	if err != nil {
		return err
	}

	if isRunning {
		return ErrProcessAlreadyRunning
	}

	proc, err := s.startProcess()
	if err != nil {
		return err
	}

	s.process = watchProcess(proc)

	return s.storeProcessPID(proc.Pid)
}

func (s *Server) Stop() error {
	isRunning, err := s.IsRunning()
	if err != nil {
		return err
	}

	if !isRunning {
		return nil
	}
	err = s.process.Kill()
	if err != nil {
		return err
	}

	return s.clearProcess()
}

func (s *Server) IsRunning() (bool, error) {
	if s.process == nil {
		err := s.recoverProcess()
		if err != nil {
			return false, err
		}
		if s.process == nil {
			return false, nil
		}
	}

	// On Windows we need to check if the process isn't hanging on an error dialog.
	if runtime.GOOS == "windows" {
		health, err := processHealth(s.process.PID())
		if err == nil && !health {
			_ = s.clearProcess()
			return false, nil
		}
	}

	isRunning, err := s.process.IsRunning()
	if err != nil || !isRunning {
		_ = s.clearProcess()
		return false, err
	}

	return isRunning, nil
}

func (s *Server) clearProcess() error {
	s.process = nil

	return os.Remove(filepath.Join(s.path, pidFile))
}

func (s *Server) recoverProcess() error {
	content, err := s.ReadFile(pidFile)
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

	s.process = recoverProcess(pid, s.processExe())

	return nil
}

func (s *Server) storeProcessPID(pid int) error {
	return s.WriteFile(pidFile, []byte(strconv.Itoa(pid)))
}
