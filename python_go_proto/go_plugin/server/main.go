package main

import (
	"github.com/hashicorp/go-plugin"

	"python_go_proto/go_plugin/greeting"
)

var _ greeting.Greeter = (*GreeterServer)(nil)

type GreeterServer struct {
}

func (g *GreeterServer) Greet(s string) (string, error) {
	return "Hello " + s, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "PLUGIN_MAGIC_COOKIE",
			MagicCookieValue: "1234567890",
		},

		Plugins: map[string]plugin.Plugin{
			"greeter": &greeting.GreeterPlugin{Impl: &GreeterServer{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
