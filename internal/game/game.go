package game

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const (
	pidFile     = "prbf2.pid"
	updaterPath = "mods/pr/bin"
)

type Server struct {
	path       string
	processPID *int
}

func Open(path string) (*Server, error) {
	//TODO: adopt and store PID in directory
	s := &Server{
		path: path,
	}

	content, err := s.ReadFile(pidFile)
	if err == nil {
		pid, err := strconv.Atoi(string(content))
		if err == nil {
			s.processPID = &pid
		}
	}

	return s, nil
}

func (s *Server) WriteFile(path string, data []byte) error {
	return os.WriteFile(filepath.Join(s.path, path), data, 0644)
}

func (s *Server) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.path, path))
}

func (s *Server) RemoveFile(path string) error {
	return os.Remove(filepath.Join(s.path, path))
}

func (s *Server) Update(ctx context.Context, outW, inW, errW io.Writer) error {
	return s.update(ctx, outW, inW, errW)
}

func makeFileExecutable(exePath string) error {
	info, err := os.Stat(exePath)
	if err != nil {
		return err
	}

	if info.Mode().Perm()&0100 == 0 {
		err = os.Chmod(exePath, info.Mode()|0100)
		if err != nil {
			return err
		}
	}

	return nil
}
