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

// DummyCircuitBuilder generates random circuits of a given size.
// It uses dummy peer IDs that might not be provisioned in the network.
type DummyCircuitBuilder struct {
	size    int
	timeout time.Duration
}

// NewDummyCircuitBuilder creates a new dummy circuit builder with default
// options.
func NewDummyCircuitBuilder(t *testing.T, opts ...echalotte.CircuitOption) echalotte.CircuitBuilder {
	options := &echalotte.CircuitOptions{
		Size:    echalotte.DefaultCircuitSize,
		Timeout: echalotte.DefaultCircuitTimeout,
	}
	err := options.Apply(opts...)
	require.NoError(t, err)

	return &DummyCircuitBuilder{
		size:    options.Size,
		timeout: options.Timeout,
	}
}

// Build a dummy circuit.
func (dcb DummyCircuitBuilder) Build(_ context.Context, _ ...echalotte.CircuitOption) (echalotte.Circuit, error) {
	circuit := make(echalotte.Circuit, dcb.size)
	for i := 0; i < dcb.size; i++ {
		sk, _, _ := crypto.GenerateEd25519Key(rand.Reader)
		peerID, _ := peer.IDFromPrivateKey(sk)
		circuit[i] = peerID
	}

	return circuit, nil
}

// FailingCircuitBuilder simulates a circuit builder that returns an error.
type FailingCircuitBuilder struct{}

// NewFailingCircuitBuilder creates a new FailingCircuitBuilder.
func NewFailingCircuitBuilder() echalotte.CircuitBuilder {
	return &FailingCircuitBuilder{}
}

// Build a circuit will fail.
func (fcb FailingCircuitBuilder) Build(_ context.Context, _ ...echalotte.CircuitOption) (echalotte.Circuit, error) {
	return nil, errors.New("my entire life is a failure")
}
