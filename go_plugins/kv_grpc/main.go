package main

import (
	"os"

	"github.com/hashicorp/go-plugin"

	"go_plugin_demo/shared"
)

var _ shared.KvStoreInterface = (*KvStore)(nil)

type KvStore struct {
}

func (k *KvStore) Put(key string, value []byte) error {
	return os.WriteFile("kv_"+key, value, 0644)
}

func (k *KvStore) Get(key string) ([]byte, error) {
	return os.ReadFile("kv_" + key)
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.KvStorePlugin{Impl: &KvStore{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
