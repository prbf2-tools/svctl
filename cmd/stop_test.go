package cmd

import (
	"context"
	"testing"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStopCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockFunc       func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error)
		expectedOutput string
		expectError    bool
	}{
		{
			name: "successful stop",
			args: []string{"test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				return &svctl.ServerInfo{
					Id:     "test-server",
					Status: svctl.Status_STATUS_STOPPING,
				}, nil
			},
			expectedOutput: "Server status: STATUS_STOPPING\n",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			testServer := NewTestGRPCServer(t)
			defer testServer.Close()

			// Configure mock
			testServer.Mock.StopFunc = tt.mockFunc

			// Create command
			cmd := stopCmd()

			// Execute command with server address
			output, err := ExecuteCommandWithServer(t, cmd, testServer, tt.args)

			// Check results
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, output, tt.expectedOutput)
			}
		})
	}
}

func TestStopCmdMissingArgs(t *testing.T) {
	testServer := NewTestGRPCServer(t)
	defer testServer.Close()

	cmd := stopCmd()
	
	_, err := ExecuteCommandWithServer(t, cmd, testServer, []string{})
	
	require.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg(s), received 0")
}