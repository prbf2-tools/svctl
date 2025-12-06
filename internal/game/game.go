package game

import (
	"errors"
	"io"
)

type GameServer interface {
	Start() error
	IsRunning() (bool, error)
	Stop() error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	WriteFileFromReader(path string, reader io.Reader, size int64) error
}

var (
	ErrIsDir = errors.New("path is a directory")
)
