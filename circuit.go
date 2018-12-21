package echalotte

import (
	"strings"

	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

// Circuit can be used to create an onion-routed message.
type Circuit []peer.ID

// String representation of the circuit.
func (c Circuit) String() string {
	var relays []string
	for _, relay := range c {
		relays = append(relays, relay.Pretty())
	}

	return strings.Join(relays, " -> ")
}
