package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"go_plugin_demo/shared"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Error,
	})
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.KvStorePlugin{},
		},
		Cmd: exec.Command("sh", "-c", "./kv-go-grpc"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		Logger: logger,
	})

	defer client.Kill()
	rpcClient, err := client.Client()
	if err != nil {
		panic(err)
	}

	raw, err := rpcClient.Dispense("kv")

	if err != nil {
		panic(err)
	}

	kv := raw.(shared.KvStoreInterface)

	//err = kv.Put("demo", []byte("hello world"))
	//if err != nil {
	//	panic(err)
	//}

	get, err := kv.Get("demo")
	if err != nil {
		return
	}
	fmt.Printf("get: %s\n", string(get))

	err = kv.Put("demo", []byte("hello"))
	if err != nil {
		panic(err)
	}

}
