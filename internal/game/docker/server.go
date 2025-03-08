package docker

import (
	"bytes"
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sboon-gg/svctl/internal/game"
)

var _ game.GameServer = &Container{}

type Container struct {
	docker *client.Client
	name   string
}

func Open(c *client.Client, containerName string) (*Container, error) {
	_, err := c.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return nil, err
	}

	return &Container{
		name:   containerName,
		docker: c,
	}, nil
}

func (c *Container) Start() error {
	return c.docker.ContainerStart(context.Background(), c.name, container.StartOptions{})
}

func (c *Container) Stop() error {
	return c.docker.ContainerStop(context.Background(), c.name, container.StopOptions{})
}

func (c *Container) IsRunning() (bool, error) {
	inspect, err := c.docker.ContainerInspect(context.Background(), c.name)
	if err != nil {
		return false, err
	}

	return inspect.State.Running, nil
}

func (c *Container) WriteFile(path string, data []byte) error {
	ctx := context.Background()

	reader := bytes.NewReader(data)

	return c.docker.CopyToContainer(ctx, c.name, path, reader, container.CopyToContainerOptions{})
}

func (c *Container) ReadFile(path string) ([]byte, error) {
	ctx := context.Background()

	readCloser, stat, err := c.docker.CopyFromContainer(ctx, c.name, path)
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	if stat.Mode.IsDir() {
		return nil, game.ErrIsDir
	}

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(readCloser)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
