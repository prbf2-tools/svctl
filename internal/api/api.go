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
	// TODO: Fetch settings pats from opts
	err := s.daemon.Register(opts.GetPath(), opts.GetSettingsPath())
	if err != nil {
		return nil, err
	}

	return &svctl.ServerInfo{
		Path:   opts.GetPath(),
		Status: svctl.Status_REGISTERED,
	}, nil
}

func (s *daemonServer) Start(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	err := s.daemon.Start(opts.GetPath())
	if err != nil {
		return nil, err
	}

	return &svctl.ServerInfo{
		Path:   opts.GetPath(),
		Status: svctl.Status_STARTED,
	}, nil
}

func (s *daemonServer) Stop(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	err := s.daemon.Stop(opts.GetPath())
	if err != nil {
		return nil, err
	}

	return &svctl.ServerInfo{
		Path:   opts.GetPath(),
		Status: svctl.Status_STOPPED,
	}, nil
}
