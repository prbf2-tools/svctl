package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateConfig(t *testing.T) {
	path := "./testdata"

	generatedConfig, err := generateConfig(path)
	require.NoError(t, err)

	expectedConfig := map[ArtifactType]artifactConfig{
		ArtifactTypeChatLog: {
			Path: "admin/logs",
			File: "chatlog_%Y-%m-%d_%H%M.txt",
		},
		ArtifactTypeCoincidentIPsLog: {
			Path: "mods/pr/settings/",
			File: "IPcoincidences.log",
		},
		ArtifactTypeAdminLog: {
			Path: "admin/logs",
			File: "ra_adminlog.txt",
		},
		ArtifactTypeBanLog: {
			Path: "mods/pr/settings/",
			File: "banlist_info.log",
		},
		ArtifactTypeTicketsLog: {
			Path: "admin/logs",
			File: "tickets.log",
		},
		ArtifactTypeJoinLog: {
			Path: "admin/logs",
			File: "joinlog.log",
		},
		ArtifactTypePlayerProfilesLog: {
			Path: "admin/logs",
			File: "playerprofiles.log",
		},
		ArtifactPlayerDataErrorsLog: {
			Path: "admin/logs",
			File: "playerdataerrors.log",
		},
		ArtifactTypePythonErrorLog: {
			Path: "mods/pr/settings",
			File: "python_errors_v([0-9]+.[0-9]+.[0-9]+.[0-9]+).log",
		},
		ArtifactTypePythonLaunchErrorLog: {
			Path: "mods/pr/settings",
			File: "python_launch_error.log",
		},
		ArtifactTypeBF2Demo: {
			Path: "mods/pr/demos",
			File: "demo_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer",
		},
		ArtifactTypePRDemo: {
			Path: "demos",
			File: "tracker_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer",
		},
		ArtifactTypePRDemoPrivate: {
			Path: "demos_private",
			File: "tracker_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer",
		},
	}
	require.Equal(t, expectedConfig, generatedConfig)
}
