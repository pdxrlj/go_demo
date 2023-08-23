package shared

import (
	"context"

	"go_plugin_demo/proto/kv_store"
)

type GRPCServer struct {
	Impl KvStoreInterface
}

func (g GRPCServer) Get(ctx context.Context, request *kv_store.GetRequest) (*kv_store.GetResponse, error) {
	v, err := g.Impl.Get(request.Key)
	return &kv_store.GetResponse{
		Value: v,
	}, err
}

func (g GRPCServer) Put(ctx context.Context, request *kv_store.PutRequest) (*kv_store.Empty, error) {
	return &kv_store.Empty{}, g.Impl.Put(request.Key, request.Value)
}

// ==============================================

type GRPCClient struct {
	client kv_store.KvStoreClient
}

func (g *GRPCClient) Put(key string, value []byte) error {
	_, err := g.client.Put(context.Background(), &kv_store.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (g *GRPCClient) Get(key string) ([]byte, error) {
	resp, err := g.client.Get(context.Background(), &kv_store.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}
