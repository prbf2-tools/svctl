package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sboon-gg/svctl/internal/game"
)

var _ game.GameServer = &Container{}

type Container struct {
	docker  *client.Client
	name    string
	uid     int
	gid     int
	workDir string
}

func Open(c *client.Client, containerName string) (*Container, error) {
	inspect, err := c.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return nil, err
	}

	workDir := inspect.Config.WorkingDir

	readCloser, _, err := c.CopyFromContainer(context.Background(), containerName, path.Join(workDir, "mods/pr/mod.desc"))
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	tarReader := tar.NewReader(readCloser)
	header, err := tarReader.Next()
	if err != nil {
		return nil, err
	}

	return &Container{
		name:    containerName,
		docker:  c,
		workDir: workDir,
		uid:     header.Uid,
		gid:     header.Gid,
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

func (c *Container) WriteFile(filePath string, data []byte) error {
	return c.WriteFileFromReader(filePath, bytes.NewReader(data), int64(len(data)))
}

func (c *Container) WriteFileFromReader(filePath string, reader io.Reader, size int64) error {
	ctx := context.Background()

	fullPath := path.Join(c.workDir, filePath)
	mode := int64(0644)

	stat, err := c.docker.ContainerStatPath(ctx, c.name, fullPath)
	if err == nil {
		mode = int64(stat.Mode)
	}

	f, err := os.CreateTemp("", "svctl-tar-*.tar")
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	tarWriter := tar.NewWriter(f)
	err = tarWriter.WriteHeader(&tar.Header{
		Name: path.Base(filePath),
		Size: size,
		Mode: mode,
		Uid:  c.uid,
		Gid:  c.gid,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, reader)
	if err != nil {
		return err
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return c.docker.CopyToContainer(ctx, c.name, path.Dir(fullPath), f, container.CopyToContainerOptions{})
}

func (c *Container) ReadFile(filePath string) ([]byte, error) {
	ctx := context.Background()

	fullPath := path.Join(c.workDir, filePath)
	println(fullPath)

	readCloser, stat, err := c.docker.CopyFromContainer(ctx, c.name, fullPath)
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	if stat.Mode.IsDir() {
		return nil, game.ErrIsDir
	}

	tarReader := tar.NewReader(readCloser)
	header, err := tarReader.Next()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, header.Size)
	read, err := tarReader.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return buf[:read], nil
}
