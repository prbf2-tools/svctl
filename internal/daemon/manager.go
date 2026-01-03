package daemon

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type ServerState string

const (
	Running ServerState = "running"
	Stopped ServerState = "stopped"
)

type ServerType string

const (
	LocalServer   ServerType = "local"
	DockerServer  ServerType = "docker"
	SystemdServer ServerType = "systemd"
)

type ServerInfo struct {
	ServerID     string      `yaml:"id"`
	Type         ServerType  `yaml:"type"`
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

func (m *ServerManager) AddServer(serverID, settingsPath string, typ ServerType) error {
	if _, ok := m.ServersInfo[serverID]; ok {
		return fmt.Errorf("server %q already exists", serverID)
	}

	m.ServersInfo[serverID] = &ServerInfo{
		ServerID:     serverID,
		Type:         typ,
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
