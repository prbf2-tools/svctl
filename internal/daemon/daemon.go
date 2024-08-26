package daemon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sboon-gg/svctl/internal/fsm"
	"github.com/sboon-gg/svctl/internal/server"
)

const (
	svctlDir  = "svctl"
	stateFile = "state.yaml"
)

type Daemon struct {
	cacheDir string
	Servers  map[string]*fsm.FSM
	ServerManager
}

func New() (*Daemon, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	svctlCacheDir := filepath.Join(cacheDir, svctlDir)

	err = os.MkdirAll(svctlCacheDir, 0755)
	if err != nil {
		return nil, err
	}

	serverManager, err := NewServerManager(filepath.Join(svctlCacheDir, stateFile))
	if err != nil {
		return nil, err
	}

	return &Daemon{
		Servers:       make(map[string]*fsm.FSM),
		cacheDir:      svctlCacheDir,
		ServerManager: *serverManager,
	}, nil
}

func Recover() (*Daemon, error) {
	d, err := New()
	if err != nil {
		return nil, err
	}

	for svPath, sv := range d.ServerManager.Servers {
		s, err := server.Open(
			svPath,
			filepath.Join(svPath, sv.SettingsPath),
		)
		if err != nil {
			return nil, err
		}

		var initState fsm.State
		initState = fsm.NewStateStopped()
		if sv.CurrentState == running {
			initState = fsm.NewStateRunning(fsm.NewRestartCounter(fsm.MaxRestarts))
		}

		d.Servers[svPath] = fsm.New(s, initState)
	}

	return d, nil
}

func (s *Daemon) Register(serverPath, settingsPath string) error {
	err := s.ServerManager.AddServer(serverPath, settingsPath)
	if err != nil {
		return err
	}

	sv, err := server.Open(serverPath, settingsPath)
	if err != nil {
		return err
	}

	s.Servers[serverPath] = fsm.New(sv, fsm.NewStateStopped())

	return nil
}

func (s *Daemon) Start(path string) error {
	sv, err := s.findServer(path)
	if err != nil {
		return err
	}

	err = sv.Event(fsm.EventStart)
	if err != nil {
		return err
	}

	return s.ServerManager.ChangeState(path, running)
}

func (s *Daemon) Stop(path string) error {
	sv, err := s.findServer(path)
	if err != nil {
		return err
	}

	err = sv.Event(fsm.EventStart)
	if err != nil {
		return err
	}

	return s.ServerManager.ChangeState(path, stopped)
}

func (d *Daemon) findServer(path string) (*fsm.FSM, error) {
	s, ok := d.Servers[path]
	if !ok {
		return nil, fmt.Errorf("server %q not found", path)
	}

	return s, nil
}
