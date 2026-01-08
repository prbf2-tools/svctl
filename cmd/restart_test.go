package cmd

import (
	"context"
	"testing"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestartCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		stopFunc       func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error)
		startFunc      func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error)
		expectedOutput []string
		errOutput      string
	}{
		{
			name: "successful restart",
			args: []string{"test-server"},
			stopFunc: func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				return &svctl.ServerInfo{
					Id:     "test-server",
					Status: svctl.Status_STATUS_STOPPING,
				}, nil
			},
			startFunc: func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				return &svctl.ServerInfo{
					Id:     "test-server",
					Status: svctl.Status_STATUS_STARTING,
				}, nil
			},
			expectedOutput: []string{
				"Stopping server test-server...",
				"Starting server test-server...",
				"Server restarted: STATUS_STARTING",
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

			// Configure mock functions
			testServer.Mock.StopFunc = tt.stopFunc
			testServer.Mock.StartFunc = tt.startFunc

			// Create command
			cmd := restartCmd()

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

