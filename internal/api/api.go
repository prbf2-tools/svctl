package api

import (
	"context"

	"github.com/sboon-gg/svctl/internal/daemon"
	"github.com/sboon-gg/svctl/svctl"
)

type daemonServer struct {
	svctl.UnimplementedServersServer
	daemon *daemon.Daemon
}

func NewDaemonServer(daemon *daemon.Daemon) svctl.ServersServer {
	return &daemonServer{
		daemon: daemon,
	}
}

func (s *daemonServer) Register(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	err := s.daemon.Register(opts.GetPath(), opts.GetSettingsPath())
	if err != nil {
		return nil, err
	}

	info, err := s.fetchServerInfo(opts.GetPath())
	if err != nil {
		return nil, err
	}

	info.Status = svctl.Status_REGISTERED

	return info, nil
}

func (s *daemonServer) Start(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	err := s.daemon.Start(opts.GetPath())
	if err != nil {
		return nil, err
	}

	info, err := s.fetchServerInfo(opts.GetPath())
	if err != nil {
		return nil, err
	}

	info.Status = svctl.Status_STARTING

	return info, nil
}

func (s *daemonServer) Stop(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	err := s.daemon.Stop(opts.GetPath())
	if err != nil {
		return nil, err
	}

	info, err := s.fetchServerInfo(opts.GetPath())
	if err != nil {
		return nil, err
	}

	info.Status = svctl.Status_STOPPING

	return info, nil
}

func (s *daemonServer) Status(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	info, err := s.fetchServerInfo(opts.GetPath())
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *daemonServer) Reset(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	err := s.daemon.Reset(opts.GetPath())
	if err != nil {
		return nil, err
	}

	return s.fetchServerInfo(opts.GetPath())
}

func (s *daemonServer) fetchServerInfo(path string) (*svctl.ServerInfo, error) {
	status, err := s.daemon.Status(path)
	if err != nil {
		return nil, err
	}

	info := &svctl.ServerInfo{
		Path:         path,
		SettingsPath: status.SettingsPath,
		DesiredState: serverStateToProto(status.DesiredState),
		CurrentState: serverStateToProto(status.CurrentState),
	}

	if status.GameStatus != nil {
		info.GameStatus = &svctl.GameStatus{
			Hostname:   status.GameStatus.Hostname,
			Port:       status.GameStatus.Port,
			Gamemode:   status.GameStatus.GameMode,
			MapName:    status.GameStatus.MapName,
			MapSize:    int64(status.GameStatus.MapSize),
			MaxPlayers: int64(status.GameStatus.MaxPlayers),
			NumPlayers: int64(status.GameStatus.NumPlayers),
		}
	}

	return info, nil
}

func serverStateToProto(state daemon.ServerState) svctl.State {
	switch state {
	case daemon.Running:
		return svctl.State_RUNNING
	case daemon.Stopped:
		return svctl.State_STOPPED
	default:
		return svctl.State_STOPPED
	}
}
