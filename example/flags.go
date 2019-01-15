package main

import (
	"flag"

	"gx/ipfs/QmVmDhyTTUcQXFD1rRQ64fGLMSAoaQvNH3hwuaCFAPq2hy/errors"
)

// Config built from command-line flags.
type Config struct {
	BootstrapPeers  addrList
	ListenAddresses addrList
}

// ParseFlags parses configuration flags.
func ParseFlags() (Config, error) {
	config := Config{}
	flag.Var(&config.BootstrapPeers, "peer", "Adds a peer multiaddress to the bootstrap list")
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.Parse()

	if len(config.ListenAddresses) == 0 {
		return config, errors.New("you need to provide at least one listening address")
	}

	return config, nil
}
