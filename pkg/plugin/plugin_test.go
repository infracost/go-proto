package plugin

import (
	"context"
	"os"
	"os/exec"
	"testing"

	goplugin "github.com/hashicorp/go-plugin"
	pluginpb "github.com/infracost/proto/gen/go/infracost/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// exitEnv makes the re-executed test binary exit before serving, standing in
// for a binary that never completes the handshake.
const exitEnv = "PLUGIN_TEST_EXIT"

// TestMain lets the test binary act as a plugin: Dial re-executes it with the
// handshake cookie set, and it serves instead of running the tests.
func TestMain(m *testing.M) {
	if os.Getenv(Handshake.MagicCookieKey) == Handshake.MagicCookieValue {
		if os.Getenv(exitEnv) != "" {
			os.Exit(1)
		}
		name := "unset"
		Serve(func(s *grpc.Server) {
			srv := &testServer{name: &name}
			pluginpb.RegisterPluginServiceServer(s, srv)
			pluginpb.RegisterParserServiceServer(s, srv)
		}, BeforeServe(func() { name = "test-plugin" }))
		os.Exit(0)
	}
	os.Exit(m.Run())
}

type testServer struct {
	pluginpb.UnimplementedPluginServiceServer
	pluginpb.UnimplementedParserServiceServer
	name *string
}

func (s *testServer) GetPluginInfo(context.Context, *pluginpb.GetPluginInfoRequest) (*pluginpb.GetPluginInfoResponse, error) {
	return &pluginpb.GetPluginInfoResponse{Name: *s.name, Type: pluginpb.PluginType_PARSER}, nil
}

func (s *testServer) Parse(context.Context, *pluginpb.ParseRequest) (*pluginpb.ParseResponse, error) {
	return &pluginpb.ParseResponse{}, nil
}

func dialSelf(t *testing.T, cfg goplugin.ClientConfig) *grpc.ClientConn {
	t.Helper()
	conn, client, err := Dial(os.Args[0], cfg)
	require.NoError(t, err)
	t.Cleanup(client.Kill)
	return conn
}

func TestDialAndServe(t *testing.T) {
	conn := dialSelf(t, goplugin.ClientConfig{})

	info, err := pluginpb.NewPluginServiceClient(conn).GetPluginInfo(t.Context(), &pluginpb.GetPluginInfoRequest{})
	require.NoError(t, err)
	assert.Equal(t, "test-plugin", info.GetName(), "BeforeServe hook should run before serving")
	assert.Equal(t, pluginpb.PluginType_PARSER, info.GetType())
}

func TestServeAcceptsMessagesAboveGRPCDefault(t *testing.T) {
	conn := dialSelf(t, goplugin.ClientConfig{})

	// gRPC servers reject requests over 4 MB unless the limit is raised.
	_, err := pluginpb.NewParserServiceClient(conn).Parse(t.Context(), &pluginpb.ParseRequest{
		RawOptions: make([]byte, 8*1024*1024),
	})
	require.NoError(t, err)
}

func TestDialKeepsCallerCmd(t *testing.T) {
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), exitEnv+"=1")

	conn, client, err := Dial(os.Args[0], goplugin.ClientConfig{Cmd: cmd})
	require.Error(t, err)
	assert.Nil(t, conn)
	assert.Nil(t, client)
	assert.NotNil(t, cmd.ProcessState, "Dial should have started the caller's command")
}

func TestDialRejectsBadPath(t *testing.T) {
	_, _, err := Dial(t.TempDir(), goplugin.ClientConfig{})
	require.ErrorContains(t, err, "is a directory")

	_, _, err = Dial("/nonexistent/plugin", goplugin.ClientConfig{})
	require.ErrorIs(t, err, os.ErrNotExist)
}
