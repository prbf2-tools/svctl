//go:build windows

package game

import (
	"context"
	"fmt"
	"io"
)

func (s *Server) update(ctx context.Context, outW io.Writer, inR io.Reader, errW io.Writer) error {
	fmt.Println("Updating game on Windows is not supported yet")
	return nil
}
