package server

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/sboon-gg/svctl/internal/game"
	"github.com/sboon-gg/svctl/internal/game/docker"
	"github.com/sboon-gg/svctl/internal/game/local"
	"github.com/sboon-gg/svctl/internal/settings"
	"github.com/sboon-gg/svctl/pkg/templates"
)

type Server struct {
	game.GameServer
	settings.Settings
}

func Open(server game.GameServer, settingsPath string) (*Server, error) {
	s, err := settings.Open(settingsPath)
	if err != nil {
		return nil, err
	}

	sv := &Server{
		GameServer: server,
		Settings:   *s,
	}

	sv.Log = sv.Log.With("server", server.ID())

	return sv, nil
}

func OpenLocal(serverPath, settingsPath string) (*Server, error) {
	g, err := local.Open(serverPath)
	if err != nil {
		return nil, err
	}

	return Open(g, settingsPath)
}

func OpenDocker(containerName, settingsPath string) (*Server, error) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	g, err := docker.Open(c, containerName)
	if err != nil {
		return nil, err
	}

	return Open(g, settingsPath)
}

func (s *Server) Render(reloadableOnly bool) error {
	if s.Settings.Templates == nil {
		return nil
	}

	values, gameConfig, err := s.Settings.TemplateData()
	if err != nil {
		return err
	}

	outputs, err := s.Settings.Templates.Render(gameConfig, values)
	if err != nil {
		return err
	}

	for _, output := range outputs {
		if reloadableOnly && !output.Reloadable {
			continue
		}

		err = s.WriteFile(output.Destination, output.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) DryRender() ([]templates.RenderOutput, error) {
	if s.Settings.Templates == nil {
		return nil, nil
	}

	values, gameConfig, err := s.Settings.TemplateData()
	if err != nil {
		return nil, err
	}

	return s.Settings.Templates.Render(gameConfig, values)
}

func (s *Server) ApplyPatches() error {
	cfg, err := s.Settings.Config()
	if err != nil {
		return err
	}

	for _, patch := range cfg.Patches {
		if patch.Source == "" || patch.Destination == "" {
			continue
		}

		f, err := os.Open(filepath.Join(s.Settings.Path, patch.Source))
		if err != nil {
			return err
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return err
		}

		err = s.WriteFileFromReader(patch.Destination, f, stat.Size())
		if err != nil {
			return err
		}
	}

	return nil
}
