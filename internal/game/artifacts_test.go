package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateConfig(t *testing.T) {
	path := "./testdata"

	sv, err := Open(path)
	require.NoError(t, err)

	generatedConfig, err := sv.ArtifactsConfig()
	require.NoError(t, err)

	expectedConfig := map[ArtifactType]ArtifactConfig{
		ArtifactTypeBF2Demo: {
			Path: "mods/pr/demos",
			File: "auto_%Y_%m_%d_%H_%M_%S.bf2demo",
		},
		ArtifactTypePRDemo: {
			Path: "demos",
			File: "tracker_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer.PRdemo",
		},
		ArtifactTypePRDemoPrivate: {
			Path: "demos_private",
			File: "tracker_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer.p.PRdemo",
		},
		ArtifactTypeJSONSummary: {
			Path: "json",
			File: "tracker_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer.json",
		},
	}
	require.Equal(t, expectedConfig, generatedConfig)
}
