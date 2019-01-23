package echalottetesting

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"gx/ipfs/QmNiJiXwWE3kRhZrC5ej3kSjWHm337pYfhjLGSCDNKJP2s/go-libp2p-crypto"
	"gx/ipfs/QmaoXrM4Z41PD48JY36YqQGKQpLGjyLA2cKcLsES7YddAq/go-libp2p-host"
	"gx/ipfs/QmdJdFQc5U3RAKgJQGmWR7SSM7TLuER5FWz5Wq6Tzs2CnS/go-libp2p"
)

// RandomHost creates a host on a random port.
func RandomHost(ctx context.Context, t *testing.T) host.Host {
	sk, _, err := crypto.GenerateEd25519Key(crand.Reader)
	require.NoError(t, err)

	return HostWithIdentity(ctx, t, sk)
}

// HostWithIdentity creates a random host with the given private key.
func HostWithIdentity(ctx context.Context, t *testing.T, sk crypto.PrivKey) host.Host {
	port := 9000 + rand.Intn(1000)
	addr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)
	h, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings(addr),
		libp2p.Identity(sk))
	require.NoError(t, err)

	return h
}
