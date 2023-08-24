package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"python_go_proto/go_plugin/greeting"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Error,
		Output: os.Stdout,
	})

	client := plugin.NewClient(
		&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "PLUGIN_MAGIC_COOKIE",
				MagicCookieValue: "1234567890",
			},
			Plugins: map[string]plugin.Plugin{
				"greeter": &greeting.GreeterPlugin{},
			},

			Cmd:    exec.Command("sh", "-c", os.Getenv("KV_PLUGIN")),
			Logger: logger,
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolGRPC,
			},
		},
	)

	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		logger.Error("Error creating RPC client", "error", err)
		os.Exit(1)
	}

	dispense, err := rpcClient.Dispense("greeter")
	if err != nil {
		logger.Error("Error creating RPC client", "error", err)
		os.Exit(1)
	}

	greeter := dispense.(greeting.Greeter)
	greet, err := greeter.Greet("Golang")
	if err != nil {
		logger.Error("Error creating RPC client", "error", err)
		os.Exit(1)
	}
	fmt.Printf("greet: %s\n", greet)
}
