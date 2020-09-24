package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	tyroplugin "github.com/rumpl/tyro/pkg/plugin"
)

type Project struct {
	Name       string            `hcl:"project"`
	Properties map[string]string `hcl:"properties"`

	Targets []*Target `hcl:"target,block"`
}

type Step struct {
	Name       string   `hcl:",label"`
	Properties hcl.Body `hcl:",remain"`
}

type Target struct {
	Name  string  `hcl:",label"`
	Steps []*Step `hcl:"step,block"`
}

type Dir struct {
	Dir string `hcl:"dir"`
}

func main() {
	var conf Project
	err := hclsimple.DecodeFile("./test.hcl", nil, &conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("project =", conf.Name)
	for k, v := range conf.Properties {
		fmt.Println(k, "=", v)
	}
	for _, t := range conf.Targets {
		fmt.Println(t.Name)
		for _, s := range t.Steps {
			fmt.Println("\t", s.Name)
			d := &Dir{}
			_ = gohcl.DecodeBody(s.Properties, nil, d)
			fmt.Println("\t\t", d.Dir)
		}
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Trace,
	})

	for _, t := range conf.Targets {
		fmt.Println(t.Name)
		for _, s := range t.Steps {
			client := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: handshakeConfig,
				Plugins:         pluginMap,
				Cmd:             exec.Command("./pls/" + s.Name),
				Logger:          logger,
			})
			defer client.Kill()

			rpcClient, err := client.Client()
			if err != nil {
				log.Fatal(err)
			}

			raw, err := rpcClient.Dispense("greeter")
			if err != nil {
				log.Fatal(err)
			}

			greeter := raw.(tyroplugin.Greeter)
			args := &map[string]string{}
			diags := gohcl.DecodeBody(s.Properties, nil, args)
			if diags.HasErrors() {
				panic(diags)
			}

			fmt.Println(greeter.Run(*args))
		}
	}
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

var pluginMap = map[string]plugin.Plugin{
	"greeter": &tyroplugin.GreeterPlugin{},
}