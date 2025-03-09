package daemon

import (
	"fmt"
	"log/slog"
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
	Servers map[string]*fsm.FSM
	ServerManager
	config *Config
}

func New(configFile string) (*Daemon, error) {
	if configFile == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}

		svctlConfigDir := filepath.Join(configDir, svctlDir)
		err = os.MkdirAll(svctlConfigDir, 0755)
		if err != nil {
			return nil, err
		}

		configFile = filepath.Join(svctlConfigDir, "config.yaml")
	}

	config, err := NewConfig(configFile)
	if err != nil {
		return nil, err
	}

	if config.CacheFile == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}

		svctlCacheDir := filepath.Join(cacheDir, svctlDir)

		err = os.MkdirAll(svctlCacheDir, 0755)
		if err != nil {
			return nil, err
		}

		config.CacheFile = filepath.Join(svctlCacheDir, stateFile)
	}

	serverManager, err := NewServerManager(config.CacheFile)
	if err != nil {
		return nil, err
	}

	return &Daemon{
		Servers:       make(map[string]*fsm.FSM),
		ServerManager: *serverManager,
		config:        config,
	}, nil
}

func Recover(configFile string) (*Daemon, error) {
	d, err := New(configFile)
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
			slog.Error("Unable to open server", "svPath", svPath, "settingsPath", settingsPath, "err", err)
			continue
		}

		var machine *fsm.FSM
		isRunning, err := s.IsRunning()
		if err != nil {
			s.Log.Error("Unable to check if server is running", "err", err)
		}

		if isRunning {
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
	GameStatus   *server.Status
}

func (s *Daemon) Status(path string) (*ServerStatus, error) {
	sv, err := s.findServer(path)
	if err != nil {
		return nil, err
	}

	gameStatus, err := sv.Server().Status()
	if err != nil {
		sv.Log.Error("Unable to get Gamespy 3 query status", "err", err)
	}

	info, ok := s.ServersInfo[path]
	if !ok {
		return nil, fmt.Errorf("server %q not found", path)
	}

	status := ServerStatus{
		DesiredState: info.DesiredState,
		CurrentState: Stopped,
		SettingsPath: info.SettingsPath,
		GameStatus:   gameStatus,
	}

	if sv.Server() != nil {
		if isRunning, err := sv.Server().IsRunning(); err == nil && isRunning {
			status.CurrentState = Running
		}
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
