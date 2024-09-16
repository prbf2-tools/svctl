//go:build windows

package game

import (
	"context"
	"io"
)

func (s *Server) update(ctx context.Context, outW io.Writer, inR io.Reader, errW io.Writer) error {
	return nil
}
