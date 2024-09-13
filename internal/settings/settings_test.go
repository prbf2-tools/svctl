package settings_test

import (
	"path/filepath"
	"testing"

	"github.com/sboon-gg/svctl/internal/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettingsInit(t *testing.T) {
	path := t.TempDir()

	// TODO: Test cloning repo
	s, err := settings.Initialize(path, nil)
	require.NoError(t, err)
	require.NotNil(t, s)
	assert.Nil(t, s.Templates)
	assert.NotNil(t, s.Log)
}

func TestSettingsOpen(t *testing.T) {
	path := t.TempDir()

	_, err := settings.Initialize(path, nil)
	require.NoError(t, err)

	s, err := settings.Open(path)
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestSettingsDefaultConfig(t *testing.T) {
	path := t.TempDir()

	_, err := settings.Initialize(path, nil)
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(path, "config.yaml"))
}
