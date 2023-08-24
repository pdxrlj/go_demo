package greeting

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Greeter interface {
	Greet(string) (string, error)
}

var _ plugin.GRPCPlugin = (*GreeterPlugin)(nil)

type GreeterPlugin struct {
	Impl Greeter
	plugin.Plugin
}

func (g *GreeterPlugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	RegisterGreetingServiceServer(server, &GRPCServer{
		Impl: g.Impl,
	})
	return nil
}

func (g *GreeterPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: NewGreetingServiceClient(conn),
	}, nil
}
