package cmd

import (
	"context"
	"testing"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockFunc       func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error)
		expectedOutput []string
		errOutput      string
	}{
		{
			name: "successful status",
			args: []string{"test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				return &svctl.ServerInfo{
					Id:           "test-server",
					Location:     "/test/location",
					SettingsPath: "/test/settings",
					Status:       svctl.Status_STATUS_REGISTERED,
					DesiredState: svctl.State_STATE_RUNNING,
					CurrentState: svctl.State_STATE_RUNNING,
					GameStatus: &svctl.GameStatus{
						Hostname:   "Test Server",
						Port:       "16567",
						Gamemode:   "gpm_cq",
						MapName:    "strike_at_karkand",
						MapSize:    64,
						NumPlayers: 32,
						MaxPlayers: 64,
					},
				}, nil
			},
			expectedOutput: []string{
				"Server Info:",
				"Path: /test/location",
				"Settings Path: /test/settings",
				"Desired State: STATE_RUNNING",
				"Current State: STATE_RUNNING",
				"Hostname: Test Server",
				"Port: 16567",
				"Game Mode: gpm_cq",
				"Map Name: strike_at_karkand",
			},
		},
		{
			name: "json output",
			args: []string{"--json", "test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
				return TestServerInfo("test-server"), nil
			},
			expectedOutput: []string{
				`"id": "test-server"`,
				`"location": "/test/test-server"`,
			},
		},
		{
			name:      "missing arguments",
			args:      []string{},
			errOutput: "accepts 1 arg(s), received 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			testServer := NewTestGRPCServer(t)
			defer testServer.Close()

			// Configure mock
			testServer.Mock.StatusFunc = tt.mockFunc

			// Create command
			cmd := statusCmd()

			// Execute command with server address
			output, err := ExecuteCommandWithServer(t, cmd, testServer, tt.args)

			// Check results
			if tt.errOutput != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errOutput)
			} else {
				require.NoError(t, err)
				for _, expected := range tt.expectedOutput {
					assert.Contains(t, output, expected)
				}
			}
		})
	}
}

