package local

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	GameDir string
}

func (s *ServerSuite) SetupTest() {
	if gameDir, ok := os.LookupEnv("PRBF2_GAME_DIR"); ok {
		s.GameDir = gameDir
	}
}

func (s *ServerSuite) TestServer_Open() {
	dir := s.T().TempDir()

	sv, err := Open(dir)
	s.Require().NoError(err)
	s.Nil(sv.process)
	s.Equal(dir, sv.path)
}

func (s *ServerSuite) TestServer_Open_WithPIDFile() {
	dir := s.T().TempDir()

	err := os.WriteFile(filepath.Join(dir, pidFile), []byte("123"), 0644)
	s.Require().NoError(err)

	sv, err := Open(dir)
	s.Require().NoError(err)
	s.Equal(123, sv.process.PID())
	s.Equal(dir, sv.path)
}

func (s *ServerSuite) TestServerLifecycle() {
	if s.GameDir == "" {
		s.T().Skip("PRBF2_GAME_DIR not set")
	}

	sv, err := Open(s.GameDir)
	s.Require().NoError(err)

	err = sv.Start()
	s.Require().NoError(err)
	s.NotNil(sv.process)

	s.Eventually(func() bool {
		isRunning, err := sv.IsRunning()
		s.Assert().NoError(err)
		return isRunning
	}, 10*time.Second, 100*time.Millisecond)

	err = sv.Stop()
	s.Require().NoError(err)

	s.Never(func() bool {
		isRunning, err := sv.IsRunning()
		s.Assert().NoError(err)
		return isRunning
	}, 10*time.Second, 100*time.Millisecond)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
