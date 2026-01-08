package cmd

import (
	"context"
	"testing"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		mockFunc       func(ctx context.Context, opts *svctl.RenderOpts) (*svctl.ServerInfo, error)
		expectedOutput string
		errOutput      string
	}{
		{
			name: "successful render",
			args: []string{"test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.RenderOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				assert.False(t, opts.ReloadableOnly)
				return &svctl.ServerInfo{
					Id: "test-server",
				}, nil
			},
			expectedOutput: "Server templates rendered: test-server\n",
		},
		{
			name: "successful render with reloadable-only flag",
			args: []string{"--reloadable-only", "test-server"},
			mockFunc: func(ctx context.Context, opts *svctl.RenderOpts) (*svctl.ServerInfo, error) {
				assert.Equal(t, "test-server", opts.Id)
				assert.True(t, opts.ReloadableOnly)
				return &svctl.ServerInfo{
					Id: "test-server",
				}, nil
			},
			expectedOutput: "Server templates rendered: test-server\n",
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
			testServer.Mock.RenderFunc = tt.mockFunc

			// Create command
			cmd := renderCmd()

			// Execute command with server address
			output, err := ExecuteCommandWithServer(t, cmd, testServer, tt.args)

			// Check results
			if tt.errOutput != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errOutput)
			} else {
				require.NoError(t, err)
				assert.Contains(t, output, tt.expectedOutput)
			}
		})
	}
}

