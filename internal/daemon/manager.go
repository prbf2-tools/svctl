package daemon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type serverState string

const (
	running serverState = "running"
	stopped serverState = "stopped"
)

type ServerInfo struct {
	ServerPath   string      `yaml:"serverPath"`
	SettingsPath string      `yaml:"settingsPath"`
	CurrentState serverState `yaml:"currentState"`
}

type ServerManager struct {
	Servers   map[string]*ServerInfo
	cachePath string
}

func NewServerManager(cachePath string) (*ServerManager, error) {
	content, err := os.ReadFile(cachePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	servers := make(map[string]*ServerInfo)
	err = yaml.Unmarshal(content, &servers)
	if err != nil {
		return nil, err
	}

	return &ServerManager{
		Servers:   servers,
		cachePath: cachePath,
	}, nil
}

func (m *ServerManager) AddServer(serverPath, settingsPath string) error {
	if _, ok := m.Servers[serverPath]; ok {
		return fmt.Errorf("server %q already exists", serverPath)
	}

	m.Servers[serverPath] = &ServerInfo{
		ServerPath:   serverPath,
		SettingsPath: settingsPath,
		CurrentState: stopped,
	}

	return m.Flush()
}

func (m *ServerManager) ChangeState(serverPath string, state serverState) error {
	s, ok := m.Servers[serverPath]
	if !ok {
		return fmt.Errorf("server %q not found", serverPath)
	}

	s.CurrentState = state
	return m.Flush()
}

func (m *ServerManager) Flush() error {
	content, err := yaml.Marshal(m.Servers)
	if err != nil {
		return err
	}

	return os.WriteFile(m.cachePath, content, 0644)
}

func (d *Daemon) cachePath(path string) string {
	return filepath.Join(d.cacheDir, path)
}
