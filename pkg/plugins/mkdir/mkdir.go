package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	tyroplugin "github.com/rumpl/tyro/pkg/plugin"
)

type Mkdir struct {
	logger hclog.Logger
}

func (g *Mkdir) Run(args map[string]string) error {
	f := args["dir"]
	if err := os.Mkdir(f, 0755); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "tyro",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	plug := &Mkdir{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"tyro": &tyroplugin.TyroPlugin{Impl: plug},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
