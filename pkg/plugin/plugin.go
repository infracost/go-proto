// Package plugin holds the go-plugin contract shared by Infracost plugin
// binaries and the hosts that launch them: the handshake, the gRPC server
// harness a plugin calls from main, and the client side that launches a plugin
// binary and returns a gRPC connection to it.
package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Name is the key a plugin is registered and dispensed under.
const Name = "plugin"

// maxMessageSize caps gRPC messages in both directions on the plugin server.
const maxMessageSize = 64 * 1024 * 1024

// Handshake must match between a host and the plugins it launches. Changing it
// breaks compatibility with every released plugin.
var Handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "INFRACOST_PLUGIN",
	MagicCookieValue: "de8c7e96-497c-4168-80c4-fc875c8ce764",
}

var (
	_ goplugin.Plugin     = grpcPlugin{}
	_ goplugin.GRPCPlugin = grpcPlugin{}
)

// grpcPlugin is the single go-plugin entry both sides use. On the server it
// registers the plugin's services; on the client it hands back the raw
// connection so callers build whichever service clients they need.
type grpcPlugin struct {
	goplugin.NetRPCUnsupportedPlugin
	register func(*grpc.Server)
}

func (p grpcPlugin) GRPCServer(_ *goplugin.GRPCBroker, s *grpc.Server) error {
	if p.register == nil {
		return errors.New("plugin has no services to register")
	}
	p.register(s)
	return nil
}

func (grpcPlugin) GRPCClient(_ context.Context, _ *goplugin.GRPCBroker, conn *grpc.ClientConn) (any, error) {
	return conn, nil
}

type serveConfig struct {
	before []func()
}

// ServeOption configures Serve.
type ServeOption func(*serveConfig)

// BeforeServe runs fn before the plugin starts serving, for process-wide setup
// such as tuning the runtime's memory limit.
func BeforeServe(fn func()) ServeOption {
	return func(c *serveConfig) { c.before = append(c.before, fn) }
}

// Serve runs the plugin server and blocks until the host disconnects. register
// adds the plugin's gRPC services, for example the PluginService and a
// ParserService or ProviderService implementation. Call it from the plugin
// binary's main.
func Serve(register func(*grpc.Server), opts ...ServeOption) {
	var cfg serveConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	for _, fn := range cfg.before {
		fn()
	}

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         goplugin.PluginSet{Name: grpcPlugin{register: register}},
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts,
				grpc.MaxRecvMsgSize(maxMessageSize),
				grpc.MaxSendMsgSize(maxMessageSize),
			)
			return grpc.NewServer(opts...)
		},
	})
}

// Dial launches the plugin binary at path and returns a gRPC connection to it,
// along with the go-plugin client that owns the subprocess. The caller must call
// Kill on the client once done; on error the subprocess has already been
// killed.
//
// Dial sets the handshake, plugin set and protocol on cfg. When cfg.Cmd is nil
// it runs path with no arguments, and when cfg.Logger is nil go-plugin's own
// logging is discarded. Every other field, such as SyncStdout, SyncStderr or a
// Cmd with a custom environment, is passed through unchanged.
func Dial(path string, cfg goplugin.ClientConfig) (*grpc.ClientConn, *goplugin.Client, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}
	if stat.IsDir() {
		return nil, nil, fmt.Errorf("%s is a directory", path)
	}

	cfg.HandshakeConfig = Handshake
	cfg.Plugins = goplugin.PluginSet{Name: grpcPlugin{}}
	cfg.AllowedProtocols = []goplugin.Protocol{goplugin.ProtocolGRPC}
	if cfg.Cmd == nil {
		cfg.Cmd = exec.Command(path)
	}
	if cfg.Logger == nil {
		cfg.Logger = hclog.NewNullLogger()
	}

	client := goplugin.NewClient(&cfg)
	conn, err := dispense(client)
	if err != nil {
		client.Kill()
		return nil, nil, err
	}
	return conn, client, nil
}

func dispense(client *goplugin.Client) (*grpc.ClientConn, error) {
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}
	raw, err := rpcClient.Dispense(Name)
	if err != nil {
		return nil, err
	}
	conn, ok := raw.(*grpc.ClientConn)
	if !ok {
		return nil, fmt.Errorf("unexpected dispensed type %T", raw)
	}
	return conn, nil
}
