package daemon

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sboon-gg/svctl/internal/daemon/fsm"
	"github.com/sboon-gg/svctl/internal/server"
	"github.com/sboon-gg/svctl/internal/settings"
	"github.com/sboon-gg/svctl/pkg/prbf2update"
)

type RunningServer struct {
	server.Server

	fsm *fsm.FSM
	log *slog.Logger
}

func OpenServer(svPath string, updaterCache *prbf2update.Cache) (*fsm.FSM, error) {
	settingsPath := filepath.Join(svPath, settings.SvctlDir)
	s, err := server.Open(svPath, settingsPath)
	if err != nil {
		return nil, err
	}

	svFSM := fsm.New(s, updaterCache)

	cache, err := s.Settings.Cache()
	if err != nil {
		return nil, err
	}

	if cache.PID != -1 {
		proc, err := os.FindProcess(cache.PID)
		if err == nil {
			err = svFSM.Adopt(proc)
			if err != nil {
				// Process can be already dead
				cache.PID = -1
				s.Settings.WriteCache(cache)
			}
		}
	}

	return svFSM, nil
}

func (rs *RunningServer) Start() error {
	rs.log.Info("Rendering templates")

	err := rs.fullRender()
	if err != nil {
		return err
	}

	rs.log.Info("Starting server")

	return rs.fsm.Start()
}

func (rs *RunningServer) Stop() error {
	rs.log.Info("Stopping server")

	return rs.fsm.Stop()
}

func (rs *RunningServer) Restart() error {
	rs.log.Info("Restarting server")

	return rs.fsm.Restart()
}

func (rs *RunningServer) setProcess(newState fsm.StateT) {
	if newState == fsm.StateTRestarting {
		rs.log.Info("Process restarting")
	}

	if newState == fsm.StateTStopped {
		rs.log.Info("Process stopped")
	}

	rs.log.Debug("Setting PID in cache")

	pid := -1
	if rs.fsm != nil {
		pid = rs.fsm.Pid()
	}

	rs.log = initLog(rs.Path).With(slog.Int("pid", pid))

	err := rs.Settings.StorePID(pid)
	if err != nil {
		rs.log.Error("Failed to store PID", slog.String("err", err.Error()))
	}
}

func (rs *RunningServer) fullRender() error {
	return rs.Render()
}

func initLog(svPath string) *slog.Logger {
	return slog.Default().With(slog.String("sv", svPath))
}
