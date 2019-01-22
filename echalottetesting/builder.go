package echalottetesting

import (
	"context"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/t-bast/go-libp2p-echalotte"

	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

const (
	// ErrBuildCircuit is returned by circuit builders when they start failing.
	ErrBuildCircuit = "my entire life is a failure"
)

// DummyCircuitBuilder generates random circuits of a given size.
// It uses dummy peer IDs that might not be provisioned in the network.
type DummyCircuitBuilder struct {
	size    int
	timeout time.Duration
	fail    bool
	peers   []peer.ID
}

// NewDummyCircuitBuilder creates a new dummy circuit builder with default
// options.
func NewDummyCircuitBuilder(t *testing.T, opts ...echalotte.CircuitOption) *DummyCircuitBuilder {
	var peers []peer.ID
	for i := 0; i < 10; i++ {
		sk, _, _ := crypto.GenerateEd25519Key(rand.Reader)
		peerID, _ := peer.IDFromPrivateKey(sk)
		peers = append(peers, peerID)
	}

	return NewDummyCircuitBuilderFromNetwork(t, peers, opts...)
}

// NewDummyCircuitBuilderFromNetwork creates a new dummy circuit builder that
// picks peers from a given list.
func NewDummyCircuitBuilderFromNetwork(t *testing.T, peers []peer.ID, opts ...echalotte.CircuitOption) *DummyCircuitBuilder {
	options := &echalotte.CircuitOptions{
		Size:    echalotte.DefaultCircuitSize,
		Timeout: echalotte.DefaultCircuitTimeout,
	}
	err := options.Apply(opts...)
	require.NoError(t, err)

	return &DummyCircuitBuilder{
		size:    options.Size,
		timeout: options.Timeout,
		peers:   peers,
	}
}

// Build a dummy circuit.
func (dcb DummyCircuitBuilder) Build(_ context.Context, _ ...echalotte.CircuitOption) (echalotte.Circuit, error) {
	if dcb.fail {
		return nil, errors.New(ErrBuildCircuit)
	}

	circuit := make(echalotte.Circuit, dcb.size)
	for i := 0; i < dcb.size; i++ {
		circuit[i] = dcb.peers[i]
	}

	return circuit, nil
}

// StartFailing tells the circuit builder to start returning errors.
func (dcb *DummyCircuitBuilder) StartFailing() *DummyCircuitBuilder {
	dcb.fail = true
	return dcb
}

// StopFailing tells the circuit builder to stop returning errors.
func (dcb *DummyCircuitBuilder) StopFailing() *DummyCircuitBuilder {
	dcb.fail = false
	return dcb
}
