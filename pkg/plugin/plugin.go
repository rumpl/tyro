package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Plugin is the interface that we're exposing as a plugin.
type Plugin interface {
	Run(args map[string]string) error
}

// Here is an implementation that talks over RPC
type PluginRPC struct{ client *rpc.Client }

func (g *PluginRPC) Run(args map[string]string) error {
	var resp error
	err := g.client.Call("Plugin.Run", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		// panic(err)
		return err
	}

	return resp
}

type PluginRPCServer struct {
	// This is the real implementation
	Impl Plugin
}

func (s *PluginRPCServer) Run(args map[string]string, resp *error) error {
	*resp = s.Impl.Run(args)
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a PluginRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return PluginRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type TyroPlugin struct {
	// Impl Injection
	Impl Plugin
}

func (p *TyroPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PluginRPCServer{Impl: p.Impl}, nil
}

func (TyroPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PluginRPC{client: c}, nil
}
