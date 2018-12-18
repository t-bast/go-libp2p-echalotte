package echalotte

import (
	"context"

	"gx/ipfs/QmYJtCabf3prS3HKQUGgqDLVxvbT9iDx5mfeVfhtCcJxxE/go-libp2p-discovery"
)

const (
	// OnionRelay is the name of the namespace advertized by onion relays.
	OnionRelay = "/libp2p/onion"
)

// CircuitBuilder lets you build random circuits for onion routing.
type CircuitBuilder struct{}

// NewCircuitBuilder creates a new circuit builder that leverages the given
// discovery component to find other peers that provide onion relays.
func NewCircuitBuilder(ctx context.Context, dscvr discovery.Discovery) (*CircuitBuilder, error) {
	// TODO: set appropriate TTL option and auto-refresh in go routine.
	// This will prevent peers that aren't responsive from being chosen in an
	// onion circuit.
	_, err := dscvr.Advertise(ctx, OnionRelay)
	if err != nil {
		return nil, err
	}

	return &CircuitBuilder{}, nil
}
