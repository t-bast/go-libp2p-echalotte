package echalottetesting

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"gx/ipfs/QmaoXrM4Z41PD48JY36YqQGKQpLGjyLA2cKcLsES7YddAq/go-libp2p-host"
	"gx/ipfs/QmdJdFQc5U3RAKgJQGmWR7SSM7TLuER5FWz5Wq6Tzs2CnS/go-libp2p"
)

// RandomHost creates a host on a random port.
func RandomHost(ctx context.Context, t *testing.T) host.Host {
	port := 9000 + rand.Intn(1000)
	addr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)
	h, err := libp2p.New(ctx, libp2p.ListenAddrStrings(addr))
	require.NoError(t, err)

	return h
}
