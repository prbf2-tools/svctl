package local

import (
	"context"
	"io"
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

	err := s.retrieveProcessPID()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) Update(ctx context.Context, outW io.Writer, inR io.Reader, errW io.Writer) error {
	return s.update(ctx, outW, inR, errW)
}
