package plugin

import (
	"github.com/hashicorp/go-plugin"
	"github.com/rumpl/tyro/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Plugin interface {
	Run(args map[string]string) error
}

type TyroPlugin struct {
	plugin.Plugin
	Impl Plugin
}

func (p *TyroPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	protos.RegisterTyroServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (TyroPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: protos.NewTyroClient(c)}, nil
}
