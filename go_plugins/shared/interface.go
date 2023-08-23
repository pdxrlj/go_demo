package shared

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"go_plugin_demo/proto/kv_store"
)

var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

type KvStoreInterface interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}

var _ plugin.GRPCPlugin = (*KvStorePlugin)(nil)

type KvStorePlugin struct {
	plugin.Plugin
	Impl KvStoreInterface
}

func (k *KvStorePlugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	kv_store.RegisterKvStoreServer(server, &GRPCServer{
		Impl: k.Impl,
	})
	return nil
}

func (k *KvStorePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: kv_store.NewKvStoreClient(conn),
	}, nil
}
