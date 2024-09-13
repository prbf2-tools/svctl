package settings

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sboon-gg/svctl/pkg/templates"
)

const (
	configFileName      = "config.yaml"
	defaultValuesFile   = "values.yaml"
	defaultTemplatesDir = "templates"
)

type Settings struct {
	path      string
	Templates *templates.Renderer
	Log       *slog.Logger
}

func Open(path string) (*Settings, error) {
	s := &Settings{
		path: path,
	}

	config, err := s.Config()
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(path, config.Loggers)
	if err != nil {
		return nil, err
	}

	s.Log = logger

	if config.TemplatesPath != "" {
		templatesPath := config.TemplatesPath
		if !filepath.IsAbs(templatesPath) {
			templatesPath = filepath.Join(path, config.TemplatesPath)
		}

		_, err = os.Stat(templatesPath)
		if err != nil {
			return nil, fmt.Errorf("templates path %q not found", templatesPath)
		}

		t, err := templates.NewFromPath(templatesPath)
		if err != nil {
			return nil, err
		}

		s.Templates = t
	}

	return s, nil
}

type Opts struct {
	TemplatesRepo string
	Token         string
}

func Initialize(path string, opts *Opts) (*Settings, error) {
	if opts == nil {
		opts = &Opts{}
	}

	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		err := os.Mkdir(path, 0755)
		if err != nil {
			return nil, err
		}
	}

	config := &Config{
		Loggers: []LoggerConfig{
			{
				Level: slog.LevelDebug,
				Stdout: &StdoutLogger{
					Type: textLogger,
				},
			},
			{
				Level: slog.LevelInfo,
				File: &FileLogger{
					Path: "svctl.log",
					Type: jsonLogger,
				},
			},
		},
	}

	if opts.TemplatesRepo != "" {
		templatesPath := filepath.Join(path, defaultTemplatesDir)
		err = cloneTemplates(templatesPath, opts.TemplatesRepo, opts.Token)
		if err != nil {
			return nil, err
		}

		t, err := templates.NewFromPath(templatesPath)
		if err != nil {
			return nil, err
		}

		defaultsContent, err := t.DefaultsContent()
		if err != nil {
			return nil, err
		}

		err = writeValues(path, defaultsContent)
		if err != nil {
			return nil, err
		}

		config.Values = append(config.Values, ValuesSource{
			File: defaultValuesFile,
		})

		config.TemplatesPath = templatesPath
	}

	err = writeConfig(path, config)
	if err != nil {
		return nil, err
	}

	return Open(path)
}
