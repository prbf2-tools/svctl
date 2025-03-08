package daemon

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Docker    DockerConfig `yaml:"docker"`
	CacheFile string       `yaml:"cache"`
}

type DockerConfig struct {
}

func NewConfig(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
