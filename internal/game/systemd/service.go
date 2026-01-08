package systemd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/unit"
	"github.com/prbf2-tools/svctl/internal/game"
)

var _ game.GameServer = &Service{}

type Service struct {
	name string
	path string
}

func Open(serviceName string) (*Service, error) {
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	unitPathProp, err := conn.GetUnitPropertyContext(ctx, fmt.Sprintf("%s.service", serviceName), "FragmentPath")
	if err != nil {
		return nil, err
	}

	unitPath := unitPathProp.Value.Value().(string)

	unitFile, err := os.Open(unitPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = unitFile.Close()
	}()

	unitOptions, err := unit.DeserializeOptions(unitFile)
	if err != nil {
		return nil, err
	}

	var path string
	for _, option := range unitOptions {
		if option.Section == "Service" && option.Name == "WorkingDirectory" {
			path = option.Value
			break
		}
	}

	if path == "" {
		return nil, fmt.Errorf("working directory not found in unit file")
	}

	return &Service{
		name: serviceName,
		path: path,
	}, nil
}

func (c *Service) ID() string {
	return c.name
}

func (c *Service) Start() error {
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.StartUnitContext(ctx, fmt.Sprintf("%s.service", c.name), "fail", nil)
	return err
}

func (c *Service) Stop() error {
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.StopUnitContext(ctx, fmt.Sprintf("%s.service", c.name), "fail", nil)
	return err
}

func (c *Service) IsRunning() (bool, error) {
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	unitStatus, err := conn.GetUnitPropertyContext(ctx, fmt.Sprintf("%s.service", c.name), "ActiveState")
	if err != nil {
		return false, err
	}

	activeState := unitStatus.Value.Value().(string)
	return activeState == "active", nil
}

func (c *Service) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(c.path, path))
}

func (c *Service) WriteFile(path string, data []byte) error {
	fullPath := filepath.Join(c.path, path)

	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		// Ignore error, the write will fail if the directory doesn't exist.
		_ = os.MkdirAll(filepath.Dir(fullPath), 0755)
	}

	return os.WriteFile(fullPath, data, 0644)
}

func (c *Service) WriteFileFromReader(path string, r io.Reader, _ int64) error {
	fullPath := filepath.Join(c.path, path)

	if _, err := os.Stat(filepath.Dir(fullPath)); os.IsNotExist(err) {
		// Ignore error, the write will fail if the directory doesn't exist.
		_ = os.MkdirAll(filepath.Dir(fullPath), 0755)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Sync()
		_ = f.Close()
	}()

	_, err = io.Copy(f, r)
	return err
}
