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
	ID           string      `yaml:"id"`
	Type         ServerType  `yaml:"type"`
	Location     string      `yaml:"location"`
	SettingsPath string      `yaml:"settingsPath"`
	DesiredState ServerState `yaml:"desiredState"`
}

type ServerManager struct {
	ServersInfo map[string]*ServerInfo
	cachePath   string
}

func NewServerManager(cachePath string) (*ServerManager, error) {
	manager := &ServerManager{
		ServersInfo: make(map[string]*ServerInfo),
		cachePath:   cachePath,
	}

	content, err := os.ReadFile(cachePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err != nil {
		if os.IsNotExist(err) {
			return manager, nil
		}
		return nil, err
	}

	err = yaml.Unmarshal(content, &manager.ServersInfo)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *ServerManager) AddServer(serverID, location, settingsPath string, typ ServerType) error {
	if _, ok := m.ServersInfo[serverID]; ok {
		return fmt.Errorf("server %q already exists", serverID)
	}

	m.ServersInfo[serverID] = &ServerInfo{
		ID:           serverID,
		Location:     location,
		Type:         typ,
		SettingsPath: settingsPath,
		DesiredState: Stopped,
	}

	return m.Flush()
}

func (m *ServerManager) RemoveServer(serverID string) error {
	if _, ok := m.ServersInfo[serverID]; !ok {
		return fmt.Errorf("server %q not found", serverID)
	}

	delete(m.ServersInfo, serverID)
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
