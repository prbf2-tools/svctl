package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateConfig(t *testing.T) {
	path := "./testdata"

	generatedConfig, err := generateConfig(path)
	require.NoError(t, err)

	expectedConfig := map[ArtifactType]artifactConfig{}
	require.Equal(t, expectedConfig, generatedConfig)
}
