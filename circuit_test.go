package echalotte_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t-bast/go-libp2p-echalotte"

	"gx/ipfs/QmY5Grm8pJdiSSVsYxx4uNRgweY72EmYwuSDbRnbFok3iY/go-libp2p-peer"
)

func TestCircuit(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		var circuit echalotte.Circuit
		circuit = []peer.ID{peer.ID("1"), peer.ID("2"), peer.ID("3")}
		// Note: note 1 -> 2 -> 3 because we prettify (b58).
		assert.Equal(t, "r -> s -> t", circuit.String())
	})
}
