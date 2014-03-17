package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/rickard-von-essen/packer-parallels/pvm"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(pvm.Builder))
	server.Serve()
}
