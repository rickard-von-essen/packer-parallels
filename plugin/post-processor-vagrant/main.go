package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/rickard-von-essen/packer-parallels/post-processor/vagrant"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(vagrant.PostProcessor))
	server.Serve()
}
