package server

import (
	"os"
	"path/filepath"

	"github.com/sboon-gg/svctl/internal/game"
	"github.com/sboon-gg/svctl/internal/settings"
	"github.com/sboon-gg/svctl/pkg/templates"
)

type Server struct {
	game.Server
	settings.Settings
}

func Open(serverPath, settingsPath string) (*Server, error) {
	s, err := settings.Open(settingsPath)
	if err != nil {
		return nil, err
	}

	g, err := game.Open(serverPath)
	if err != nil {
		return nil, err
	}

	sv := &Server{
		Server:   *g,
		Settings: *s,
	}

	sv.Log = sv.Log.With("server", filepath.Base(serverPath))

	return sv, nil
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

		dst := filepath.Join(s.Path, output.Destination)
		err = os.MkdirAll(filepath.Dir(dst), 0755)
		if err != nil {
			return err
		}

		err = os.WriteFile(dst, output.Content, 0644)
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
