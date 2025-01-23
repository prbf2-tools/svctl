package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestsConfig(t *testing.T) {
	path := "./testdata"

	sv, err := Open(path)
	require.NoError(t, err)

	generatedConfig, err := sv.LogsConfig()
	require.NoError(t, err)

	expectedConfig := map[LogType]ArtifactConfig{
		LogTypeChat: {
			Path: "admin/logs",
			File: "chatlog_%Y-%m-%d_%H%M.txt",
		},
		LogTypeCoincidentIPs: {
			Path: "mods/pr/settings/",
			File: "IPcoincidences.log",
		},
		LogTypeAdmin: {
			Path: "admin/logs",
			File: "ra_adminlog.txt",
		},
		LogTypeBan: {
			Path: "mods/pr/settings/",
			File: "banlist_info.log",
		},
		LogTypeTickets: {
			Path: "admin/logs",
			File: "tickets.log",
		},
		LogTypeJoin: {
			Path: "admin/logs",
			File: "joinlog.log",
		},
		LogTypePlayerProfiles: {
			Path: "admin/logs",
			File: "playerprofiles.log",
		},
		LogTypePlayerDataErrors: {
			Path: "admin/logs",
			File: "playerdataerrors.log",
		},
		LogTypePythonErrors: {
			Path: "mods/pr/settings",
			File: "python_errors_v([0-9]+.[0-9]+.[0-9]+.[0-9]+).log",
		},
		LogTypePythonLaunchErrors: {
			Path: "mods/pr/settings",
			File: "python_launch_error.log",
		},
	}
	require.Equal(t, expectedConfig, generatedConfig)
}
