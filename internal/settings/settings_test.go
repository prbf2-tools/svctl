package settings_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/sboon-gg/svctl/internal/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettingsInit(t *testing.T) {
	path := t.TempDir()

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

func TestSettingsDoubleInit(t *testing.T) {
	path := t.TempDir()

	_, err := settings.Initialize(path, nil)
	require.NoError(t, err)

	_, err = settings.Initialize(path, nil)
	require.Error(t, err)
}

func TestSettingsWithTemplatesRepo(t *testing.T) {
	settingsPath := t.TempDir()
	templatesRepoPath := t.TempDir()

	createRepoFromTemplatesDir(t, "../../pkg/templates/testdata/example/", templatesRepoPath)

	s, err := settings.Initialize(settingsPath, &settings.Opts{
		TemplatesRepo: templatesRepoPath,
	})
	require.NoError(t, err)
	assert.NotNil(t, s.Templates)

	assert.FileExists(t, filepath.Join(settingsPath, "config.yaml"))
	assert.FileExists(t, filepath.Join(settingsPath, "values.yaml"))
	assert.DirExists(t, filepath.Join(settingsPath, "templates"))
}

func createRepoFromTemplatesDir(t *testing.T, templatesDir, repoDir string) {
	err := filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		src := filepath.Join(wd, path)
		dst := filepath.Join(repoDir, d.Name())

		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}

		defer func() {
			cerr := out.Close()
			if err == nil {
				err = cerr
			}
		}()
		if _, err = io.Copy(out, in); err != nil {
			return err
		}
		err = out.Sync()
		return nil
	})
	require.NoError(t, err)

	repo, err := git.PlainInit(repoDir, false)
	require.NoError(t, err)

	tree, err := repo.Worktree()
	require.NoError(t, err)

	_, err = tree.Add(".")
	require.NoError(t, err)

	_, err = tree.Commit("Initial commit", &git.CommitOptions{})
	require.NoError(t, err)
}
