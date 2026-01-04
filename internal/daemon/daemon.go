package daemon

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/prbf2-tools/svctl/internal/fsm"
	"github.com/prbf2-tools/svctl/internal/server"
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

	for serverID, sv := range d.ServersInfo {
		settingsPath := sv.SettingsPath
		if !filepath.IsAbs(settingsPath) {
			settingsPath = filepath.Join(serverID, sv.SettingsPath)
		}

		var s *server.Server

		switch sv.Type {
		case LocalServer:
			s, err = server.OpenLocal(
				serverID,
				settingsPath,
			)
			if err != nil {
				slog.Error("Unable to open local server", "serverID", serverID, "settingsPath", settingsPath, "err", err)
				continue
			}
		case DockerServer:
			s, err = server.OpenDocker(
				serverID,
				settingsPath,
			)
			if err != nil {
				slog.Error("Unable to open docker server", "serverID", serverID, "settingsPath", settingsPath, "err", err)
				continue
			}
		case SystemdServer:
			s, err = server.OpenSystemd(
				serverID,
				settingsPath,
			)
			if err != nil {
				slog.Error("Unable to open systemd server", "serverID", serverID, "settingsPath", settingsPath, "err", err)
				continue
			}
		default:
			slog.Error("Unknown server type", "serverID", serverID, "type", sv.Type)
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

		d.Servers[serverID] = machine
	}

	return d, nil
}

func (s *Daemon) Register(serverID, settingsPath string, typ ServerType) error {
	err := s.AddServer(serverID, settingsPath, typ)
	if err != nil {
		return err
	}

	var sv *server.Server
	switch typ {
	case LocalServer:
		sv, err = server.OpenLocal(serverID, settingsPath)
		if err != nil {
			return err
		}
	case DockerServer:
		sv, err = server.OpenDocker(serverID, settingsPath)
		if err != nil {
			return err
		}
	case SystemdServer:
		sv, err = server.OpenSystemd(serverID, settingsPath)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown server type %q", typ)
	}

	s.Servers[serverID] = fsm.New(sv, sv.Log, fsm.NewStateStopped())

	return nil
}

func (s *Daemon) Start(path string) error {
	sv, err := s.findServer(path)
	if err != nil {
		return err
	}

	err = s.ChangeState(path, Running)
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

	err = s.ChangeState(path, Stopped)
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
	err = s.ChangeState(path, Stopped)
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

func (s *Daemon) Status(id string) (*ServerStatus, error) {
	sv, err := s.findServer(id)
	if err != nil {
		return nil, err
	}

	// gameStatus, err := sv.Server().Status()
	// if err != nil {
	// 	sv.Log.Error("Unable to get Gamespy 3 query status", "err", err)
	// }

	info, ok := s.ServersInfo[id]
	if !ok {
		return nil, fmt.Errorf("server %q not found", id)
	}

	status := ServerStatus{
		DesiredState: info.DesiredState,
		CurrentState: Stopped,
		SettingsPath: info.SettingsPath,
		// GameStatus:   gameStatus,
	}

	if sv.Server() != nil {
		if isRunning, err := sv.Server().IsRunning(); err == nil && isRunning {
			status.CurrentState = Running
		}
	}

	return &status, nil
}

func (d *Daemon) findServer(id string) (*fsm.FSM, error) {
	s, ok := d.Servers[id]
	if !ok {
		return nil, fmt.Errorf("server %q not found", id)
	}

	return s, nil
}
