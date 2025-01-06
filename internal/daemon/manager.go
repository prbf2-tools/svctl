package daemon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ServerState string

const (
	Running ServerState = "running"
	Stopped ServerState = "stopped"
)

type ServerInfo struct {
	ServerPath   string      `yaml:"serverPath"`
	SettingsPath string      `yaml:"settingsPath"`
	DesiredState ServerState `yaml:"desiredState"`
}

type ServerManager struct {
	ServersInfo map[string]*ServerInfo
	cachePath   string
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
		ServersInfo: servers,
		cachePath:   cachePath,
	}, nil
}

func (m *ServerManager) AddServer(serverPath, settingsPath string) error {
	if _, ok := m.ServersInfo[serverPath]; ok {
		return fmt.Errorf("server %q already exists", serverPath)
	}

	m.ServersInfo[serverPath] = &ServerInfo{
		ServerPath:   serverPath,
		SettingsPath: settingsPath,
		DesiredState: Stopped,
	}

	return m.Flush()
}

func (m *ServerManager) ChangeState(serverPath string, state ServerState) error {
	s, ok := m.ServersInfo[serverPath]
	if !ok {
		return fmt.Errorf("server %q not found", serverPath)
	}

	s.DesiredState = state
	return m.Flush()
}

func (m *ServerManager) Flush() error {
	content, err := yaml.Marshal(m.ServersInfo)
	if err != nil {
		return err
	}

	return os.WriteFile(m.cachePath, content, 0644)
}

func (d *Daemon) cachePath(path string) string {
	return filepath.Join(d.cacheDir, path)
}
