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
	ServerPort  int    `yaml:"serverPort"`
	GamespyPort int    `yaml:"gamespyPort"`
	ExternalIP  string `yaml:"externalIP"`
}

type DockerConfig struct {
}

type PatchConfig struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

type Config struct {
	Values        []ValuesSource `yaml:"values"`
	Loggers       []LoggerConfig `yaml:"loggers"`
	TemplatesPath string         `yaml:"templates"`
	Patches       []PatchConfig  `yaml:"patches"`
	Game          GameConfig     `yaml:"game"`
	Docker        *DockerConfig  `yaml:"docker,omtempty"`
}

func (s *Settings) Config() (*Config, error) {
	var config Config

	content, err := os.ReadFile(filepath.Join(s.Path, configFileName))
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
