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

	for svPath, sv := range d.ServerManager.ServersInfo {
		settingsPath := sv.SettingsPath
		if !filepath.IsAbs(settingsPath) {
			settingsPath = filepath.Join(svPath, sv.SettingsPath)
		}

		s, err := server.Open(
			svPath,
			settingsPath,
		)
		if err != nil {
			s.Log.Error("Unable to open server", "svPath", svPath, "settingsPath", settingsPath, "err", err)
			continue
		}

		var machine *fsm.FSM
		if s.IsRunning() {
			machine = fsm.New(s, s.Log, fsm.NewStateRunning(nil))
			if sv.DesiredState == Stopped {
				err := machine.Event(fsm.EventStop)
				if err != nil {
					s.Log.Error("Unable to reach desired state on recovery", "desiredState", sv.DesiredState, "err", err)
					continue
				}
			}
		} else {
			machine = fsm.New(s, s.Log, fsm.NewStateStopped())
			if sv.DesiredState == Running {
				err := machine.Event(fsm.EventStart)
				if err != nil {
					s.Log.Error("Unable to reach desired state on recovery", "desiredState", sv.DesiredState, "err", err)
					continue
				}
			}
		}

		d.Servers[svPath] = machine
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

	s.Servers[serverPath] = fsm.New(sv, sv.Log, fsm.NewStateStopped())

	return nil
}

func (s *Daemon) Start(path string) error {
	sv, err := s.findServer(path)
	if err != nil {
		return err
	}

	err = s.ServerManager.ChangeState(path, Running)
	if err != nil {
		return err
	}

	err = sv.Server().Render(false)
	if err != nil {
		return err
	}

	err = sv.Event(fsm.EventStart)
	if err != nil {
		return err
	}

	sv.Log.Info("Server started", "op", "Daemon.Start")
	return nil
}

func (s *Daemon) Stop(path string) error {
	sv, err := s.findServer(path)
	if err != nil {
		return err
	}

	err = s.ServerManager.ChangeState(path, Stopped)
	if err != nil {
		return err
	}

	err = sv.Event(fsm.EventStop)
	if err != nil {
		return err
	}

	sv.Log.Info("Server stopped", "op", "Daemon.Stop")
	return nil
}

func (s *Daemon) Reset(path string) error {
	sv, err := s.findServer(path)
	if err != nil {
		return err
	}

	sv.Log.Info("Reseting server", "op", "Daemon.Reset")
	err = s.ServerManager.ChangeState(path, Stopped)
	if err != nil {
		return err
	}

	return sv.Event(fsm.EventReset)
}

type ServerStatus struct {
	DesiredState ServerState
	CurrentState ServerState
	SettingsPath string
}

func (s *Daemon) Status(path string) (*ServerStatus, error) {
	sv, err := s.findServer(path)
	if err != nil {
		return nil, err
	}

	info, ok := s.ServersInfo[path]
	if !ok {
		return nil, fmt.Errorf("server %q not found", path)
	}

	status := ServerStatus{
		DesiredState: info.DesiredState,
		CurrentState: Stopped,
		SettingsPath: info.SettingsPath,
	}

	if sv.Server() != nil && sv.Server().IsRunning() {
		status.CurrentState = Running
	}

	return &status, nil
}

func (d *Daemon) findServer(path string) (*fsm.FSM, error) {
	s, ok := d.Servers[path]
	if !ok {
		return nil, fmt.Errorf("server %q not found", path)
	}

	return s, nil
}
