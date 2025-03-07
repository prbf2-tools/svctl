package game

import "errors"

type GameServer interface {
	Start() error
	IsRunning() (bool, error)
	Stop() error
	WriteFile(path string, data []byte) error
	ReadFile(path string) ([]byte, error)
}

var (
	ErrIsDir = errors.New("path is a directory")
)
