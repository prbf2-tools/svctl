package cmd

import (
	"context"
	"testing"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockFunc       func(ctx context.Context, opts *svctl.RegisterServerOpts) (*svctl.ServerInfo, error)
		expectedOutput string
		expectError    bool
	}{
		{
			name: "successful register with local path",
			args: []string{"--path", "/test/path", "test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.RegisterServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				assert.Equal(t, svctl.ServerType_SERVER_TYPE_LOCAL, opts.Type)
				assert.Equal(t, "/test/path", opts.Location)
				return &svctl.ServerInfo{
					Id:     "test-server",
					Status: svctl.Status_STATUS_REGISTERED,
				}, nil
			},
			expectedOutput: "Server status: STATUS_REGISTERED\n",
		},
		{
			name: "successful register with container",
			args: []string{"--container", "test-container", "test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.RegisterServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				assert.Equal(t, "test-container", opts.Location)
				assert.Equal(t, svctl.ServerType_SERVER_TYPE_DOCKER, opts.Type)
				return &svctl.ServerInfo{
					Id:     "test-server",
					Status: svctl.Status_STATUS_REGISTERED,
				}, nil
			},
			expectedOutput: "Server status: STATUS_REGISTERED\n",
		},
		{
			name: "successful register with service",
			args: []string{"--service", "test-service", "test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.RegisterServerOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				assert.Equal(t, "test-service", opts.Location)
				assert.Equal(t, svctl.ServerType_SERVER_TYPE_SYSTEMD, opts.Type)
				return &svctl.ServerInfo{
					Id:     "test-server",
					Status: svctl.Status_STATUS_REGISTERED,
				}, nil
			},
			expectedOutput: "Server status: STATUS_REGISTERED\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			testServer := NewTestGRPCServer(t)
			defer testServer.Close()

			// Configure mock
			testServer.Mock.RegisterFunc = tt.mockFunc

			// Create command
			cmd := registerCmd()

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

func TestRegisterCmdMissingArgs(t *testing.T) {
	testServer := NewTestGRPCServer(t)
	defer testServer.Close()

	cmd := registerCmd()

	_, err := ExecuteCommandWithServer(t, cmd, testServer, []string{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "accepts 1 arg(s), received 0")
}

func TestRegisterCmdMissingLocation(t *testing.T) {
	testServer := NewTestGRPCServer(t)
	defer testServer.Close()

	cmd := registerCmd()

	_, err := ExecuteCommandWithServer(t, cmd, testServer, []string{"test-server"})

	require.Error(t, err)
	require.Contains(t, err.Error(), "flags in the group")
}

