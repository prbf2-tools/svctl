package settings

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ValuesSource struct {
	File   string         `yaml:"file"`
	Values map[string]any `yaml:"values"`
}

type GameConfig struct {
	ServerIP    string `yaml:"serverIP"`
	ServerPort  string `yaml:"serverPort"`
	GamespyPort string `yaml:"gamespyPort"`
	ExternalIP  string `yaml:"externalIP"`
}

type Config struct {
	Values        []ValuesSource `yaml:"values"`
	Loggers       []LoggerConfig `yaml:"loggers"`
	TemplatesPath string         `yaml:"templates"`
	Game          GameConfig     `yaml:"game"`
}

func (s *Settings) Config() (*Config, error) {
	var config Config

	content, err := os.ReadFile(filepath.Join(s.path, configFileName))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func writeConfig(path string, conf *Config) error {
	content, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path, configFileName), content, 0644)
}
