package cmd

import (
	"context"
	"testing"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResetCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockFunc       func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error)
		expectedOutput string
		expectError    bool
	}{
		{
			name: "successful reset",
			args: []string{"test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				return &svctl.ServerInfo{
					Id:           "test-server",
					Status:       svctl.Status_STATUS_REGISTERED,
					DesiredState: svctl.State_STATE_STOPPED,
					CurrentState: svctl.State_STATE_STOPPED,
				}, nil
			},
			expectedOutput: "Reseted completed: STATUS_REGISTERED\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			testServer := NewTestGRPCServer(t)
			defer testServer.Close()

			// Configure mock
			testServer.Mock.ResetFunc = tt.mockFunc

			// Create command
			cmd := resetCmd()

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

func TestResetCmdMissingArgs(t *testing.T) {
	testServer := NewTestGRPCServer(t)
	defer testServer.Close()

	cmd := resetCmd()

	_, err := ExecuteCommandWithServer(t, cmd, testServer, []string{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg(s), received 0")
}

