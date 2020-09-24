package plugin

import (
	"github.com/rumpl/tyro/protos"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client protos.TyroClient }

func (m *GRPCClient) Run(args map[string]string) error {
	_, err := m.client.Run(context.Background(), &protos.RunRequest{
		Args: args,
	})
	return err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Plugin
}

func (m *GRPCServer) Run(
	ctx context.Context,
	req *protos.RunRequest) (*protos.RunResponse, error) {
	return &protos.RunResponse{}, m.Impl.Run(req.Args)
}
