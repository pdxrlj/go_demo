package greeting

import (
	"context"
)

var _ GreetingServiceServer = (*GRPCServer)(nil)

type GRPCServer struct {
	Impl Greeter
}

func (G *GRPCServer) Greet(ctx context.Context, request *GreetingRequest) (*GreetingResponse, error) {
	greet, err := G.Impl.Greet(request.GetName())
	if err != nil {
		return nil, err
	}
	return &GreetingResponse{
		Message: greet,
	}, nil
}

type GRPCClient struct {
	client GreetingServiceClient
}

func (G *GRPCClient) Greet(req string) (string, error) {
	resp, err := G.client.Greet(context.Background(), &GreetingRequest{
		Name: req,
	})
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}
