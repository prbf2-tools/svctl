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
	Path       string
	processPID *int
}

func Open(path string) (*Server, error) {
	s := &Server{
		Path: path,
	}

	content, err := os.ReadFile(filepath.Join(s.Path, pidFile))
	if err == nil {
		pid, err := strconv.Atoi(string(content))
		if err == nil {
			s.processPID = &pid
		}
	}

	return s, nil
}

func (s *Server) Update(ctx context.Context, outW io.Writer, inR io.Reader, errW io.Writer) error {
	return s.update(ctx, outW, inR, errW)
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
