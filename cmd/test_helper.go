package cmd

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/prbf2-tools/svctl/svctl/v1"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockServersServer implements the svctl.ServersServer interface for testing
type MockServersServer struct {
	svctl.UnimplementedServersServer

	// Test functions that can be overridden
	RegisterFunc func(context.Context, *svctl.RegisterServerOpts) (*svctl.ServerInfo, error)
	StartFunc    func(context.Context, *svctl.ServerOpts) (*svctl.ServerInfo, error)
	StopFunc     func(context.Context, *svctl.ServerOpts) (*svctl.ServerInfo, error)
	StatusFunc   func(context.Context, *svctl.ServerOpts) (*svctl.ServerInfo, error)
	ResetFunc    func(context.Context, *svctl.ServerOpts) (*svctl.ServerInfo, error)
	RenderFunc   func(context.Context, *svctl.RenderOpts) (*svctl.ServerInfo, error)
}

func (m *MockServersServer) Register(ctx context.Context, opts *svctl.RegisterServerOpts) (*svctl.ServerInfo, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, opts)
	}
	return &svctl.ServerInfo{
		Id:           opts.Id,
		Location:     opts.Location,
		SettingsPath: opts.SettingsPath,
		Status:       svctl.Status_STATUS_REGISTERED,
		DesiredState: svctl.State_STATE_STOPPED,
		CurrentState: svctl.State_STATE_STOPPED,
	}, nil
}

func (m *MockServersServer) Start(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	if m.StartFunc != nil {
		return m.StartFunc(ctx, opts)
	}
	return &svctl.ServerInfo{
		Id:           opts.Id,
		Status:       svctl.Status_STATUS_STARTING,
		DesiredState: svctl.State_STATE_RUNNING,
		CurrentState: svctl.State_STATE_RUNNING,
	}, nil
}

func (m *MockServersServer) Stop(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	if m.StopFunc != nil {
		return m.StopFunc(ctx, opts)
	}
	return &svctl.ServerInfo{
		Id:           opts.Id,
		Status:       svctl.Status_STATUS_STOPPING,
		DesiredState: svctl.State_STATE_STOPPED,
		CurrentState: svctl.State_STATE_STOPPED,
	}, nil
}

func (m *MockServersServer) Status(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	if m.StatusFunc != nil {
		return m.StatusFunc(ctx, opts)
	}
	return &svctl.ServerInfo{
		Id:           opts.Id,
		Location:     "/test/location",
		SettingsPath: "/test/settings",
		Status:       svctl.Status_STATUS_REGISTERED,
		DesiredState: svctl.State_STATE_STOPPED,
		CurrentState: svctl.State_STATE_STOPPED,
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
}

func (m *MockServersServer) Reset(ctx context.Context, opts *svctl.ServerOpts) (*svctl.ServerInfo, error) {
	if m.ResetFunc != nil {
		return m.ResetFunc(ctx, opts)
	}
	return &svctl.ServerInfo{
		Id:           opts.Id,
		Status:       svctl.Status_STATUS_REGISTERED,
		DesiredState: svctl.State_STATE_STOPPED,
		CurrentState: svctl.State_STATE_STOPPED,
	}, nil
}

func (m *MockServersServer) Render(ctx context.Context, opts *svctl.RenderOpts) (*svctl.ServerInfo, error) {
	if m.RenderFunc != nil {
		return m.RenderFunc(ctx, opts)
	}
	return &svctl.ServerInfo{
		Id:           opts.Id,
		Status:       svctl.Status_STATUS_REGISTERED,
		DesiredState: svctl.State_STATE_STOPPED,
		CurrentState: svctl.State_STATE_STOPPED,
	}, nil
}

// TestGRPCServer sets up a real gRPC server with mock implementation for testing
type TestGRPCServer struct {
	Server   *grpc.Server
	Listener net.Listener
	Mock     *MockServersServer
	Address  string
}

// NewTestGRPCServer creates a new test gRPC server with a mock implementation
func NewTestGRPCServer(t *testing.T) *TestGRPCServer {
	// Create listener on available port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	// Create gRPC server
	server := grpc.NewServer()

	// Create mock service
	mock := &MockServersServer{}
	svctl.RegisterServersServer(server, mock)

	// Start server in background
	go func() {
		err = server.Serve(lis)
		require.NoError(t, err)
	}()

	return &TestGRPCServer{
		Server:   server,
		Listener: lis,
		Mock:     mock,
		Address:  lis.Addr().String(),
	}
}

// Close shuts down the test server
func (ts *TestGRPCServer) Close() {
	if ts.Server != nil {
		ts.Server.Stop()
	}
	if ts.Listener != nil {
		_ = ts.Listener.Close()
	}
}

// Host returns the host part of the server address
func (ts *TestGRPCServer) Host() string {
	host, _, _ := net.SplitHostPort(ts.Address)
	return host
}

// Port returns the port part of the server address
func (ts *TestGRPCServer) Port() string {
	_, port, _ := net.SplitHostPort(ts.Address)
	return port
}

// ExecuteCommandWithServer executes a cobra command with the given arguments and server address
func ExecuteCommandWithServer(t *testing.T, cmd *cobra.Command, server *TestGRPCServer, args []string) (string, error) {
	// Capture output
	output := &testOutput{}
	cmd.SetOut(output)
	cmd.SetErr(output)

	// Add server address flags and other args
	allArgs := append([]string{"--host", server.Host(), "--port", server.Port()}, args...)
	cmd.SetArgs(allArgs)

	// Execute
	err := cmd.Execute()

	return output.String(), err
}

// testOutput captures command output for testing
type testOutput struct {
	data []byte
}

func (o *testOutput) Write(p []byte) (n int, err error) {
	o.data = append(o.data, p...)
	return len(p), nil
}

func (o *testOutput) String() string {
	return string(o.data)
}

// TestServerInfo returns a standard test ServerInfo for consistent testing
func TestServerInfo(id string) *svctl.ServerInfo {
	return &svctl.ServerInfo{
		Id:           id,
		Location:     fmt.Sprintf("/test/%s", id),
		SettingsPath: fmt.Sprintf("/test/%s/.svctl", id),
		Status:       svctl.Status_STATUS_REGISTERED,
		DesiredState: svctl.State_STATE_STOPPED,
		CurrentState: svctl.State_STATE_STOPPED,
		GameStatus: &svctl.GameStatus{
			Hostname:   fmt.Sprintf("Test Server %s", id),
			Port:       "16567",
			Gamemode:   "gpm_cq",
			MapName:    "strike_at_karkand",
			MapSize:    4,
			NumPlayers: 32,
			MaxPlayers: 100,
		},
	}
}

// ExecuteCommand executes a cobra command with the given arguments for testing (without server)
func ExecuteCommand(t *testing.T, cmd *cobra.Command, args []string) (string, error) {
	// Capture output
	output := &testOutput{}
	cmd.SetOut(output)
	cmd.SetErr(output)

	// Set args
	cmd.SetArgs(args)

	// Execute
	err := cmd.Execute()

	return output.String(), err
}
