package settings

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ValuesSource struct {
	File string `yaml:"file"`
}

type Config struct {
	Values        []ValuesSource `yaml:"values"`
	Loggers       []LoggerConfig `yaml:"loggers"`
	TemplatesPath string         `yaml:"templates"`
}

func (s *Settings) Config() (*Config, error) {
	var config Config

	content, err := os.ReadFile(filepath.Join(s.path, ConfigFile))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (s *Settings) WriteConfig(conf *Config) error {
	return writeConfig(s.path, conf)
}

func writeConfig(path string, conf *Config) error {
	content, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path, ConfigFile), content, 0644)
}
